package path

import (
	"encoding/binary"
	"fmt"
	"hvif/utils"
	"io"
	"math"
)

type Type uint8

type Path struct {
	isClosed bool

	Elements []PathElement
}

type PathFlag uint8
type PathCommandType uint8

const pathCommandSizeBits = 2
const byteSizeBits = 8

const (
	PathFlagClosed PathFlag = 1 << (1 + iota)
	PathFlagUsesCommands
	PathFlagNoCurves
)

const (
	PathCommandHLine PathCommandType = iota
	PathCommandVLine
	PathCommandLine
	PathCommandCurve
)

type PathElement interface{}

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

func readPoint(r io.Reader) Point {
	var p Point
	p.X = utils.ReadFloatCoord(r)
	p.Y = utils.ReadFloatCoord(r)
	return p
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

func Read(r io.Reader) Path {
	var path Path
	var flag PathFlag
	binary.Read(r, binary.LittleEndian, &flag)
	path.isClosed = flag&PathFlagClosed != 0

	if flag&PathFlagNoCurves != 0 {
		var count uint8
		binary.Read(r, binary.LittleEndian, &count)

		var points []PathElement
		for i := byte(0); i < count; i++ {
			points = append(points, readPoint(r))
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
				line = HLine{utils.ReadFloatCoord(r)}
			case PathCommandVLine:
				line = VLine{utils.ReadFloatCoord(r)}
			case PathCommandLine:
				line = Line{readPoint(r)}
			case PathCommandCurve:
				line = Curve{readPoint(r), readPoint(r), readPoint(r)}
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
			points = append(points, Curve{readPoint(r), readPoint(r), readPoint(r)})
		}
		path.Elements = points
	}

	return path
}
