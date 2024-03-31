package utils

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

func ReadFloatCoord(r io.Reader) (float32, error) {
	var val uint16
	var x uint8
	if err := binary.Read(r, binary.LittleEndian, &x); err != nil {
		return 0, fmt.Errorf("read first part: %w", err)
	}
	val = uint16(x)
	if val&0x80 != 0 {
		var xlow uint8
		binary.Read(r, binary.LittleEndian, &xlow)
		val = val & (0x80 - 1)
		val = (val << 8) | uint16(xlow)

		return float32(val)/102.0 - 128.0, nil
	}

	return float32(val) - 32.0, nil
}

func ReadFloat24(r io.Reader) float32 {
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

func ReadMatrix(r io.Reader, size int) []float32 {
	res := make([]float32, size)
	for i := 0; i < size; i++ {
		res[i] = ReadFloat24(r)
	}
	return res
}
