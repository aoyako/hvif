package transform

import (
	"encoding/binary"
	"io"
)

const MatrixSize int = 6

type Transformable struct {
	Matrix [MatrixSize]float32
}

func ReadTransformable(r io.Reader) Transformable {
	var t Transformable
	binary.Read(r, binary.LittleEndian, &t.Matrix)
	return t
}
