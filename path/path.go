package path

import (
	"encoding/binary"
	"fmt"
	"hvif/utils"
	"io"
	"math"
)

const pathCommandSizeBits = 2
const byteSizeBits = 8

type PathFlag uint8

const (
	PathFlagClosed PathFlag = 1 << (1 + iota)
	PathFlagUsesCommands
	PathFlagNoCurves
)

type PathCommandType uint8

const (
	PathCommandHLine PathCommandType = iota
	PathCommandVLine
	PathCommandLine
	PathCommandCurve
)

type PathElement any

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

type Line struct {
	Point Point
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
	x, err := utils.ReadFloatCoord(r)
	if err != nil {
		return p, fmt.Errorf("read x coord: %w", err)
	}
	y, err := utils.ReadFloatCoord(r)
	if err != nil {
		return p, fmt.Errorf("read x coord: %w", err)
	}
	p.X = x
	p.Y = y
	return p, nil
}

func splitCommandTypes(rawTypes []uint8, count uint8) []PathCommandType {
	pct := make([]PathCommandType, 0, count)
	const pctsPerByte = (byteSizeBits / pathCommandSizeBits)
	for i := uint8(0); i < count; i++ {
		segment := i / pctsPerByte
		shift := i % pctsPerByte

		commandType := (rawTypes[segment] >> (shift * pathCommandSizeBits)) & 0x3
		pct = append(pct, PathCommandType(commandType))
	}

	return pct
}

func Read(r io.Reader) (Path, error) {
	var path Path
	var flag PathFlag
	binary.Read(r, binary.LittleEndian, &flag)
	path.isClosed = flag&PathFlagClosed != 0

	if flag&PathFlagNoCurves != 0 {
		var count uint8
		binary.Read(r, binary.LittleEndian, &count)

		var points []PathElement
		for i := byte(0); i < count; i++ {
			p, err := readPoint(r)
			if err != nil {
				return path, fmt.Errorf("read point: %w", err)
			}
			points = append(points, p)
		}
		path.Elements = points

	} else if flag&PathFlagUsesCommands != 0 {
		var count uint8
		binary.Read(r, binary.LittleEndian, &count)

		// Each command is 2 bits, aligned in a byte
		bytesForCommandTypes := uint8(math.Ceil(pathCommandSizeBits * float64(count) / byteSizeBits))
		pathRawCommandTypes := make([]uint8, bytesForCommandTypes)
		binary.Read(r, binary.LittleEndian, &pathRawCommandTypes)
		pathCommandTypes := splitCommandTypes(pathRawCommandTypes, count)

		var points []PathElement
		for i := byte(0); i < count; i++ {
			var line interface{}
			switch pathCommandTypes[i] {
			case PathCommandHLine:
				c, err := utils.ReadFloatCoord(r)
				if err != nil {
					return path, fmt.Errorf("read hline coord: %w", err)
				}
				line = HLine{c}
			case PathCommandVLine:
				c, err := utils.ReadFloatCoord(r)
				if err != nil {
					return path, fmt.Errorf("read vline coord: %w", err)
				}
				line = VLine{c}
			case PathCommandLine:
				p, err := readPoint(r)
				if err != nil {
					return path, fmt.Errorf("read point: %w", err)
				}
				line = Line{p}
			case PathCommandCurve:
				p1, err := readPoint(r)
				if err != nil {
					return path, fmt.Errorf("read first point of curve: %w", err)
				}
				p2, err := readPoint(r)
				if err != nil {
					return path, fmt.Errorf("read second point of curve: %w", err)
				}
				p3, err := readPoint(r)
				if err != nil {
					return path, fmt.Errorf("read third point of curve: %w", err)
				}
				line = Curve{p1, p2, p3}
			default:
				fmt.Println(pathCommandTypes[i])
			}
			points = append(points, line)
		}

		path.Elements = points

	} else {
		var count uint8
		err := binary.Read(r, binary.LittleEndian, &count)
		if err != nil {
			panic(err)
		}
		var points []PathElement
		for i := byte(0); i < count; i++ {
			p1, err := readPoint(r)
			if err != nil {
				return path, fmt.Errorf("read first point of curve: %w", err)
			}
			p2, err := readPoint(r)
			if err != nil {
				return path, fmt.Errorf("read second point of curve: %w", err)
			}
			p3, err := readPoint(r)
			if err != nil {
				return path, fmt.Errorf("read third point of curve: %w", err)
			}
			points = append(points, Curve{p1, p2, p3})
		}
		path.Elements = points
	}

	return path, nil
}
