package style

import (
	"encoding/binary"
	"fmt"
	"hvif/utils"
	"io"
)

type Type uint8

const (
	StyleSolidColor Type = 1 + iota
	StyleGradient
	StyleSolidColorNoAlpha
	StyleSolidGray
	StyleSolidGrayNoAlpha
)

type GradientType uint8

const (
	GradientLinear GradientType = iota
	GradientCircular
	GradientDiamond
	GradientConic
	GradientXY
	GradientSqrtXY
)

type GradientFlag uint8

const (
	GradientFlagTransform GradientFlag = 1 << (1 + iota)
	GradientFlagNoAlpha
	GradientFlag16BitColors // Unused
	GradientFlagGrays
)

type Style any

type Color struct {
	Red   uint8
	Green uint8
	Blue  uint8
	Alpha uint8
}

type GradientColor struct {
	StopOffset uint8
	Color
}

type Gradient struct {
	Type          GradientType
	Transformable *utils.Transformable
	Colors        []GradientColor
}

type SolidColor struct {
	Red   uint8
	Green uint8
	Blue  uint8
	Alpha uint8
}

type SolidColorNoAlpha struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

type SolidGray struct {
	Gray  uint8
	Alpha uint8
}

type SolidGrayNoAlpha struct {
	Gray uint8
}

func (sc SolidColor) ToColor() Color {
	// nolint
	return Color{
		Red:   sc.Red,
		Green: sc.Green,
		Blue:  sc.Blue,
		Alpha: sc.Alpha,
	}
}

func (scna SolidColorNoAlpha) ToColor() Color {
	return Color{
		Red:   scna.Red,
		Green: scna.Green,
		Blue:  scna.Blue,
		Alpha: 0xff,
	}
}

func (sg SolidGray) ToColor() Color {
	return Color{
		Red:   sg.Gray,
		Green: sg.Gray,
		Blue:  sg.Gray,
		Alpha: sg.Alpha,
	}
}

func (sgna SolidGrayNoAlpha) ToColor() Color {
	return Color{
		Red:   sgna.Gray,
		Green: sgna.Gray,
		Blue:  sgna.Gray,
		Alpha: 0xff,
	}
}

func Read(r io.Reader) (Style, error) {
	var styleType Type
	err := binary.Read(r, binary.LittleEndian, &styleType)
	if err != nil {
		return nil, err
	}
	switch styleType {
	case StyleSolidColor:
		var c SolidColor
		err := binary.Read(r, binary.LittleEndian, &c)
		if err != nil {
			return nil, fmt.Errorf("reading solid color: %w", err)
		}
		return c.ToColor(), nil
	case StyleSolidColorNoAlpha:
		var c SolidColorNoAlpha
		err := binary.Read(r, binary.LittleEndian, &c)
		if err != nil {
			return nil, fmt.Errorf("reading solid color without alpha: %w", err)
		}
		return c.ToColor(), nil
	case StyleSolidGray:
		var c SolidGray
		err := binary.Read(r, binary.LittleEndian, &c)
		if err != nil {
			return nil, fmt.Errorf("reading solid gray color: %w", err)
		}
		return c.ToColor(), nil
	case StyleSolidGrayNoAlpha:
		var c SolidGrayNoAlpha
		err := binary.Read(r, binary.LittleEndian, &c)
		if err != nil {
			return nil, fmt.Errorf("reading solid gray color without alpha: %w", err)
		}
		return c.ToColor(), nil
	case StyleGradient:
		var g Gradient
		var gradientType GradientType
		var gradientFlags GradientFlag
		var ncolors uint8
		err := binary.Read(r, binary.LittleEndian, &gradientType)
		if err != nil {
			return nil, fmt.Errorf("reading gradient type: %w", err)
		}
		err = binary.Read(r, binary.LittleEndian, &gradientFlags)
		if err != nil {
			return nil, fmt.Errorf("reading gradient flags: %w", err)
		}
		err = binary.Read(r, binary.LittleEndian, &ncolors)
		if err != nil {
			return nil, fmt.Errorf("reading gradient number of colors: %w", err)
		}

		g.Type = gradientType
		if gradientFlags&GradientFlagTransform != 0 {
			t := utils.ReadTransformable(r)
			g.Transformable = &t
		}

		for i := byte(0); i < ncolors; i++ {
			var color Color
			var offset uint8
			err := binary.Read(r, binary.LittleEndian, &offset)
			if err != nil {
				return nil, fmt.Errorf("reading gradient [%d] color offset: %w", i, err)
			}

			if gradientFlags&GradientFlagGrays != 0 {
				if gradientFlags&GradientFlagNoAlpha != 0 {
					var gc SolidGrayNoAlpha
					err := binary.Read(r, binary.LittleEndian, &gc)
					if err != nil {
						return nil, fmt.Errorf("reading gradient [%d] solid gray color without alpha: %w", i, err)
					}
					color = gc.ToColor()
				} else {
					var gc SolidGray
					err := binary.Read(r, binary.LittleEndian, &gc)
					if err != nil {
						return nil, fmt.Errorf("reading gradient [%d] solid gray color: %w", i, err)
					}
					color = gc.ToColor()
				}
			} else {
				if gradientFlags&GradientFlagNoAlpha != 0 {
					var gc SolidColorNoAlpha
					err := binary.Read(r, binary.LittleEndian, &gc)
					if err != nil {
						return nil, fmt.Errorf("reading gradient [%d] solid color without alpha: %w", i, err)
					}
					color = gc.ToColor()
				} else {
					var gc SolidColor
					err := binary.Read(r, binary.LittleEndian, &gc)
					if err != nil {
						return nil, fmt.Errorf("reading gradient [%d] solid color: %w", i, err)
					}
					color = gc.ToColor()
				}
			}

			g.Colors = append(g.Colors, GradientColor{
				StopOffset: offset,
				Color:      color,
			})
		}

		return g, nil
	}

	return nil, fmt.Errorf("unknown style: %d", styleType)
}
