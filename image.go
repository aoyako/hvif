package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"slices"
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

	for i := range styleCount {
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

	for i := range pathCount {
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

	for i := range shapeCount {
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
	if s.styleID == nil {
		return nil
	}

	return i.styles[*s.styleID]
}

func (i *Image) GetShapePathes(s *Shape) []*Path {
	res := make([]*Path, 0, len(s.pathIDs))
	for pid := range s.pathIDs {
		res = append(res, i.pathes[pid])
	}

	return res
}

func (i *Image) AddStyle(s Style) {
	i.styles = append(i.styles, s)
}

func (i *Image) AddPath(p *Path) {
	i.pathes = append(i.pathes, p)
}

func (i *Image) AddShape(sp *Shape) {
	i.shapes = append(i.shapes, sp)
}

func (i *Image) RemoveStyle(s Style) {
	styleID := slices.Index(i.styles, s)
	if styleID == -1 {
		return
	}
	i.styles = slices.Delete(i.styles, styleID, styleID+1)

	for _, sp := range i.shapes {
		// Shift style by one
		if sp.styleID != nil && *sp.styleID > uint8(styleID) {
			*sp.styleID--
		}
		if sp.styleID != nil && *sp.styleID == uint8(styleID) {
			sp.styleID = nil
		}
	}
}

func (i *Image) RemovePath(p *Path) {
	pathID := slices.Index(i.pathes, p)
	if pathID == -1 {
		return
	}
	i.pathes = slices.Delete(i.pathes, pathID, pathID+1)

	for _, sp := range i.shapes {
		// TODO: Optimize to 0 allocs?
		newPathIDs := make([]uint8, 0, len(sp.pathIDs))
		for _, pid := range sp.pathIDs {
			if pid != uint8(pathID) {
				if pid > uint8(pathID) {
					pid--
				}
				newPathIDs = append(newPathIDs, pid)
			}
		}
		sp.pathIDs = newPathIDs
	}
}

func (i *Image) RemoveShape(sp *Shape) {
	shapeID := slices.Index(i.shapes, sp)
	if shapeID == -1 {
		return
	}
	i.shapes = slices.Delete(i.shapes, shapeID, shapeID+1)
}
