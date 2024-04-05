package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

type (
	shapeType uint8
	shapeFlag uint8
)

const (
	shapePathSource shapeType = 0xa
)

const (
	shapeFlagTransform shapeFlag = 1 << (1 + iota)
	shapeFlagHinting
	shapeFlagLodScale
	shapeFlagHasTransformers
	shapeFlagTranslation
)

type Shape struct {
	Hinting bool
	styleID uint8
	pathIDs []uint8
	// Transform  *utils.Transformable
	// Translate  *utils.Translation
	// LodScale   *utils.LodScale
	Transforms []Transformer
}

func readShape(r io.Reader) (*Shape, error) {
	s := &Shape{}
	var stype shapeType

	err := binary.Read(r, binary.LittleEndian, &stype)
	if err != nil {
		return nil, fmt.Errorf("reading type: %w", err)
	}

	if stype == shapePathSource {
		var styleID uint8
		err := binary.Read(r, binary.LittleEndian, &styleID)
		if err != nil {
			return nil, fmt.Errorf("reading style id: %w", err)
		}
		s.styleID = styleID

		var pathCount uint8
		err = binary.Read(r, binary.LittleEndian, &pathCount)
		if err != nil {
			return nil, fmt.Errorf("reading path count: %w", err)
		}
		for i := byte(0); i < pathCount; i++ {
			var pathID uint8
			err := binary.Read(r, binary.LittleEndian, &pathID)
			if err != nil {
				return nil, fmt.Errorf("reading path [%d] id: %w", i, err)
			}
			s.pathIDs = append(s.pathIDs, pathID)
		}

		var flags shapeFlag
		err = binary.Read(r, binary.LittleEndian, &flags)
		if err != nil {
			return nil, fmt.Errorf("reading flags: %w", err)
		}
		if flags&shapeFlagTransform != 0 {
			t, err := readAffine(r)
			if err != nil {
				return nil, fmt.Errorf("reading affine transformer: %w", err)
			}
			s.Transforms = append(s.Transforms, t)
		}
		if flags&shapeFlagTranslation != 0 {
			t, err := readTranslation(r)
			if err != nil {
				return nil, fmt.Errorf("read translation %w", err)
			}
			s.Transforms = append(s.Transforms, t)
		}
		if flags&shapeFlagLodScale != 0 {
			t, err := readLodScale(r)
			if err != nil {
				return nil, fmt.Errorf("read lod scale: %w", err)
			}
			s.Transforms = append(s.Transforms, t)
		}
		if flags&shapeFlagHasTransformers != 0 {
			var count uint8
			err := binary.Read(r, binary.LittleEndian, &count)
			if err != nil {
				return nil, fmt.Errorf("reading transformers count: %w", err)
			}
			for i := uint8(0); i < count; i++ {
				t, err := readTransformer(r)
				if err != nil {
					return nil, fmt.Errorf("reading transformer [%d]: %w", i, err)
				}
				s.Transforms = append(s.Transforms, t)
			}
		}
		if flags&shapeFlagHinting != 0 {
			s.Hinting = true
		}
	}

	return s, nil
}
