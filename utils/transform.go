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
	if x&0x8 == 1 {
		var xlow uint8
		binary.Read(r, binary.LittleEndian, &xlow)
		x = x & (0x8 - 1)
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

type Transformable struct {
	Matrix [MatrixSize]float32
}

func ReadTransformable(r io.Reader) Transformable {
	var t Transformable
	for i := 0; i < len(t.Matrix); i++ {
		t.Matrix[i] = ReadFloatTrans(r)
	}
	return t
}
