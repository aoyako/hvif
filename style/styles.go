package style

import (
	"encoding/binary"
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

type Style struct {
	isColor    bool
	isGradient bool
	color      Color
	gradient   Gradient
}

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

func (s Style) Gradient() (Gradient, bool) {
	return s.gradient, s.isGradient
}

func (s Style) Color() (Color, bool) {
	return s.color, s.isColor
}

func styleFromGradient(g Gradient) Style {
	return Style{isGradient: true, gradient: g}
}

func styleFromColor(c Color) Style {
	return Style{isColor: true, color: c}
}

func (sc SolidColor) ToColor() Color {
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

func Read(r io.Reader) Style {
	var styleType Type
	binary.Read(r, binary.LittleEndian, &styleType)
	switch styleType {
	case StyleSolidColor:
		var c SolidColor
		binary.Read(r, binary.LittleEndian, &c)
		return styleFromColor(c.ToColor())
	case StyleSolidColorNoAlpha:
		var c SolidColorNoAlpha
		binary.Read(r, binary.LittleEndian, &c)
		return styleFromColor(c.ToColor())
	case StyleSolidGray:
		var c SolidGray
		binary.Read(r, binary.LittleEndian, &c)
		return styleFromColor(c.ToColor())
	case StyleSolidGrayNoAlpha:
		var c SolidGrayNoAlpha
		binary.Read(r, binary.LittleEndian, &c)
		return styleFromColor(c.ToColor())
	case StyleGradient:
		var g Gradient
		var gradientType GradientType
		var gradientFlags GradientFlag
		var ncolors uint8
		binary.Read(r, binary.LittleEndian, &gradientType)
		binary.Read(r, binary.LittleEndian, &gradientFlags)
		binary.Read(r, binary.LittleEndian, &ncolors)

		g.Type = gradientType
		if gradientFlags&GradientFlagTransform != 0 {
			t := utils.ReadTransformable(r)
			g.Transformable = &t
		}

		for i := byte(0); i < ncolors; i++ {
			var color Color
			var offset uint8
			binary.Read(r, binary.LittleEndian, &offset)

			if gradientFlags&GradientFlagGrays != 0 {
				if gradientFlags&GradientFlagNoAlpha != 0 {
					var gc SolidGrayNoAlpha
					binary.Read(r, binary.LittleEndian, &gc)
					color = gc.ToColor()
				} else {
					var gc SolidGray
					binary.Read(r, binary.LittleEndian, &gc)
					color = gc.ToColor()
				}
			} else {
				if gradientFlags&GradientFlagNoAlpha != 0 {
					var gc SolidColorNoAlpha
					binary.Read(r, binary.LittleEndian, &gc)
					color = gc.ToColor()
				} else {
					var gc SolidColor
					binary.Read(r, binary.LittleEndian, &gc)
					color = gc.ToColor()
				}
			}

			g.Colors = append(g.Colors, GradientColor{
				StopOffset: offset,
				Color:      color,
			})
		}

		return styleFromGradient(g)
	}

	return Style{}
}
