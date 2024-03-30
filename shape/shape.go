package shape

import (
	"encoding/binary"
	"hvif/utils"
	"io"
)

type ShapeType uint8
type ShapeFlag uint8

const (
	ShapePathSource ShapeType = 0xa
)

const (
	ShapeFlagTransform ShapeFlag = 1 << (1 + iota)
	ShapeFlagHinting
	ShapeFlagLodScale
	ShapeFlagHasTransformers
	ShapeFlagTranslation
)

type Shape struct {
	Hinting    bool
	StyleID    uint8
	PathIDs    []uint8
	Transform  *utils.Transformable
	Translate  *utils.Translation
	LodScale   *utils.LodScale
	Transforms []utils.Transformer
}

func Read(r io.Reader) (Shape, error) {
	var s Shape
	var stype ShapeType
	binary.Read(r, binary.LittleEndian, &stype)
	if stype == ShapePathSource {
		var styleID uint8
		binary.Read(r, binary.LittleEndian, &styleID)
		s.StyleID = styleID

		var pathCount uint8
		binary.Read(r, binary.LittleEndian, &pathCount)
		for i := byte(0); i < pathCount; i++ {
			var pathID uint8
			binary.Read(r, binary.LittleEndian, &pathID)
			s.PathIDs = append(s.PathIDs, pathID)
		}

		var flags ShapeFlag
		binary.Read(r, binary.LittleEndian, &flags)
		if flags&ShapeFlagTransform != 0 {
			t := utils.ReadTransformable(r)
			s.Transform = &t
		}
		if flags&ShapeFlagTranslation != 0 {
			t := utils.ReadTranslation(r)
			s.Translate = &t
		}
		if flags&ShapeFlagLodScale != 0 {
			t := utils.ReadLodScale(r)
			s.LodScale = &t
		}
		if flags&ShapeFlagHasTransformers != 0 {
			var count uint8
			binary.Read(r, binary.LittleEndian, &count)
			for i := uint8(0); i < count; i++ {
				t, err := utils.ReadTransformer(r)
				if err != nil {
					return s, err
				}
				s.Transforms = append(s.Transforms, t)
			}
		}
		if flags&ShapeFlagHinting != 0 {
			s.Hinting = true
		}
	}
	return s, nil
}
