package path

import (
	"encoding/binary"
	"io"
)

type Type uint8
type Path interface{}

const (
	PathHLine Type = iota
	PathVLine
	PathLine
	PathCurve
)

type HLine struct {
	X float32
}

type VLine struct {
	Y float32
}

type Line struct {
	Point Point
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

func Parse(r io.Reader) []Path {
	var path []Path
	var pathType Type
	binary.Read(r, binary.LittleEndian, &pathType)
	switch pathType {
	case PathHLine:
		var p HLine
		binary.Read(r, binary.LittleEndian, &p)
		path = append(path, p)
	}

	return path
}
