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

func Read(r io.Reader) {
	var stype ShapeType
	binary.Read(r, binary.LittleEndian, &stype)
	if stype == ShapePathSource {
		var styleID uint8
		binary.Read(r, binary.LittleEndian, &styleID)

		var pathCount uint8
		binary.Read(r, binary.LittleEndian, &pathCount)
		for i := byte(0); i < pathCount; i++ {
			var pathID uint8
			binary.Read(r, binary.LittleEndian, &pathID)
		}

		var flags ShapeFlag
		binary.Read(r, binary.LittleEndian, &flags)
		if flags&ShapeFlagTransform != 0 {
			utils.ReadTransformable(r)
		}
		if flags&ShapeFlagTranslation != 0 {
			utils.ReadTranslation(r)
		}
		if flags&ShapeFlagLodScale != 0 {
			utils.ReadLodScale(r)
		}
		if flags&ShapeFlagHasTransformers != 0 {
			var count uint8
			binary.Read(r, binary.LittleEndian, &count)
			for i := uint8(0); i < count; i++ {
				utils.ReadTransformer(r)
			}
		}
		if flags&ShapeFlagHinting != 0 {
		}
	}
}
