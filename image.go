package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Image struct {
	styles []Style
	pathes []*Path
	shapes []*Shape
}

func ReadImage(r io.Reader) (*Image, error) {
	img := &Image{}

	magic := make([]uint8, 4)
	_, err := io.ReadFull(r, magic)
	if err != nil {
		return nil, fmt.Errorf("reading magic: %w", err)
	}

	if string(magic) != "ncif" {
		return nil, fmt.Errorf("magic should be ncif, found: %s", magic)
	}

	var styleCount uint8
	err = binary.Read(r, binary.LittleEndian, &styleCount)
	if err != nil {
		return nil, fmt.Errorf("reading styles count: %w", err)
	}

	for i := uint8(0); i < styleCount; i++ {
		s, err := readStyle(r)
		if err != nil {
			return nil, fmt.Errorf("reading style [%d]: %w", i, err)
		}
		img.styles = append(img.styles, s)
	}

	var pathCount uint8
	err = binary.Read(r, binary.LittleEndian, &pathCount)
	if err != nil {
		return nil, fmt.Errorf("reading pathes count: %w", err)
	}

	for i := uint8(0); i < pathCount; i++ {
		p, err := readPath(r)
		if err != nil {
			return nil, fmt.Errorf("reading path [%d]: %w", i, err)
		}
		img.pathes = append(img.pathes, p)
	}

	var shapeCount uint8
	err = binary.Read(r, binary.LittleEndian, &shapeCount)
	if err != nil {
		return nil, fmt.Errorf("reading shapes count: %w", err)
	}

	for i := uint8(0); i < shapeCount; i++ {
		s, err := readShape(r)
		if err != nil {
			return nil, fmt.Errorf("reading shape [%d]: %w", i, err)
		}
		img.shapes = append(img.shapes, s)
	}

	return img, nil
}

func (i *Image) GetStyles() []Style {
	return i.styles
}

func (i *Image) GetPathes() []*Path {
	return i.pathes
}

func (i *Image) GetShapes() []*Shape {
	return i.shapes
}

func (i *Image) GetShapeStyle(s *Shape) Style {
	return i.styles[s.styleID]
}

func (i *Image) GetShapePathes(s *Shape) []*Path {
	var res []*Path
	for pid := range s.pathIDs {
		res = append(res, i.pathes[pid])
	}

	return res
}
