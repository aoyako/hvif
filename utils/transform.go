package utils

import (
	"encoding/binary"
	"io"
	"math"
)

func ReadFloatCoord(r io.Reader) float32 {
	var val uint16
	var x uint8
	binary.Read(r, binary.LittleEndian, &x)
	val = uint16(x)
	if val&0x80 != 0 {
		var xlow uint8
		binary.Read(r, binary.LittleEndian, &xlow)
		val = val & (0x80 - 1)
		val = (val << 8) | uint16(xlow)

		return float32(val)/102.0 - 128.0
	}

	return float32(val) - 32.0
}

func ReadFloatTrans(r io.Reader) float32 {
	var b1 uint8
	var b2 uint8
	var b3 uint8
	binary.Read(r, binary.LittleEndian, &b1)
	binary.Read(r, binary.LittleEndian, &b2)
	binary.Read(r, binary.LittleEndian, &b3)

	value := uint32(uint32(b1)<<16 | uint32(b2)<<8 | uint32(b3))
	if value == 0 {
		return 0.0
	}

	sign := ((value & 0b100000000000000000000000) >> 23)
	expo := ((value & 0b011111100000000000000000) >> 17) - 32
	mant := ((value & 0b000000011111111111111111) >> 6)

	bits := (sign << 31) | ((expo + 127) << 23) | mant
	return math.Float32frombits(bits)
}

const MatrixSize = 6

type TransformerType uint8

const (
	TransformerAffine TransformerType = 20 + iota
	TransformerContour
	TransformerPerspective
	TransformerStroke
)

type Transformable struct {
	Matrix [MatrixSize]float32
}

type Translation struct {
	X float32
	Y float32
}

type LodScale struct {
	MinS float32
	MaxS float32
}

func ReadTransformable(r io.Reader) Transformable {
	var t Transformable
	for i := 0; i < len(t.Matrix); i++ {
		t.Matrix[i] = ReadFloatTrans(r)
	}
	return t
}

func ReadTransformer(r io.Reader) Transformable {
	var ttype TransformerType
	binary.Read(r, binary.LittleEndian, &ttype)
}

func ReadTranslation(r io.Reader) Translation {
	var t Translation
	t.X = ReadFloatCoord(r)
	t.Y = ReadFloatCoord(r)
	return t
}

func ReadLodScale(r io.Reader) LodScale {
	var ls LodScale
	var scale uint8
	binary.Read(r, binary.LittleEndian, &scale)
	ls.MinS = float32(scale) / 63.75
	binary.Read(r, binary.LittleEndian, &scale)
	ls.MaxS = float32(scale) / 63.75

	return ls
}
