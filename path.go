package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

const (
	pathCommandSizeBits = 2
	byteSizeBits        = 8
)

type pathFlag uint8

const (
	pathFlagClosed pathFlag = 1 << (1 + iota)
	pathFlagUsesCommands
	pathFlagNoCurves
)

type pathCommandType uint8

const (
	pathCommandHLine pathCommandType = iota
	pathCommandVLine
	pathCommandLine
	pathCommandCurve
)

type PathElement any // Point | HLine | VLine | Curve

type HLine struct {
	X float32
}

type VLine struct {
	Y float32
}

type Point struct {
	X float32
	Y float32
}

type Curve struct {
	PointIn  Point
	Point    Point
	PointOut Point
}

type Path struct {
	isClosed bool
	Elements []PathElement
}

func readPoint(r io.Reader) (Point, error) {
	var p Point
	x, err := readFloatCoord(r)
	if err != nil {
		return p, fmt.Errorf("reading x coord: %w", err)
	}
	y, err := readFloatCoord(r)
	if err != nil {
		return p, fmt.Errorf("reading y coord: %w", err)
	}
	p.X = x
	p.Y = y

	return p, nil
}

func splitCommandTypes(rawTypes []uint8, count uint8) []pathCommandType {
	pct := make([]pathCommandType, 0, count)
	const pctsPerByte = (byteSizeBits / pathCommandSizeBits)
	for i := range count {
		segment := i / pctsPerByte
		shift := i % pctsPerByte

		commandType := (rawTypes[segment] >> (shift * pathCommandSizeBits)) & 0x3
		pct = append(pct, pathCommandType(commandType))
	}

	return pct
}

func readPath(r io.Reader) (*Path, error) {
	path := &Path{}
	var flag pathFlag
	err := binary.Read(r, binary.LittleEndian, &flag)
	if err != nil {
		return nil, fmt.Errorf("reading flags: %w", err)
	}
	path.isClosed = flag&pathFlagClosed != 0

	switch {
	case flag&pathFlagNoCurves != 0:
		var count uint8
		err := binary.Read(r, binary.LittleEndian, &count)
		if err != nil {
			return nil, fmt.Errorf("reading count for path no curves: %w", err)
		}

		var points []PathElement
		for i := byte(0); i < count; i++ {
			p, err := readPoint(r)
			if err != nil {
				return nil, fmt.Errorf("reading point: %w", err)
			}
			points = append(points, p)
		}
		path.Elements = points
	case flag&pathFlagUsesCommands != 0:
		var count uint8
		err := binary.Read(r, binary.LittleEndian, &count)
		if err != nil {
			return nil, fmt.Errorf("reading count for path with commands: %w", err)
		}

		// Each command is 2 bits, aligned in a byte
		bytesForCommandTypes := uint8(math.Ceil(pathCommandSizeBits * float64(count) / byteSizeBits))
		pathRawCommandTypes := make([]uint8, bytesForCommandTypes)
		err = binary.Read(r, binary.LittleEndian, &pathRawCommandTypes)
		if err != nil {
			return nil, fmt.Errorf("reading commands: %w", err)
		}

		pathCommandTypes := splitCommandTypes(pathRawCommandTypes, count)

		var points []PathElement
		for i := byte(0); i < count; i++ {
			var line interface{}
			switch pathCommandTypes[i] {
			case pathCommandHLine:
				c, err := readFloatCoord(r)
				if err != nil {
					return nil, fmt.Errorf("reading hline coord: %w", err)
				}
				line = &HLine{c}
			case pathCommandVLine:
				c, err := readFloatCoord(r)
				if err != nil {
					return nil, fmt.Errorf("reading vline coord: %w", err)
				}
				line = &VLine{c}
			case pathCommandLine:
				p, err := readPoint(r)
				if err != nil {
					return nil, fmt.Errorf("reading point: %w", err)
				}
				line = &p
			case pathCommandCurve:
				c, err := readCurve(r)
				if err != nil {
					return nil, fmt.Errorf("reading curve: %w", err)
				}
				line = &c
			}
			points = append(points, line)
		}

		path.Elements = points
	default:
		var count uint8
		err := binary.Read(r, binary.LittleEndian, &count)
		if err != nil {
			return nil, fmt.Errorf("reading count for curves: %w", err)
		}
		var points []PathElement
		for i := byte(0); i < count; i++ {
			c, err := readCurve(r)
			if err != nil {
				return nil, fmt.Errorf("reading curve: %w", err)
			}
			points = append(points, &c)
		}
		path.Elements = points
	}

	return path, nil
}

func readCurve(r io.Reader) (Curve, error) {
	var c Curve
	p1, err := readPoint(r)
	if err != nil {
		return c, fmt.Errorf("reading first point: %w", err)
	}
	p2, err := readPoint(r)
	if err != nil {
		return c, fmt.Errorf("reading second point: %w", err)
	}
	p3, err := readPoint(r)
	if err != nil {
		return c, fmt.Errorf("reading third point: %w", err)
	}

	return Curve{p1, p2, p3}, nil
}
