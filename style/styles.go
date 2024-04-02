package style

import (
	"encoding/binary"
	"fmt"
	"io"

	"hvif/utils"
)

type Type uint8

const ColorChannelMaxValue = 0xff

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
	//nolint
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
		Alpha: ColorChannelMaxValue,
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
		Alpha: ColorChannelMaxValue,
	}
}

func readGradient(r io.Reader) (Gradient, error) {
	var gradient Gradient
	var gradientType GradientType
	var gradientFlags GradientFlag
	var ncolors uint8
	err := binary.Read(r, binary.LittleEndian, &gradientType)
	if err != nil {
		return gradient, fmt.Errorf("reading type: %w", err)
	}
	err = binary.Read(r, binary.LittleEndian, &gradientFlags)
	if err != nil {
		return gradient, fmt.Errorf("reading flags: %w", err)
	}
	err = binary.Read(r, binary.LittleEndian, &ncolors)
	if err != nil {
		return gradient, fmt.Errorf("reading number of colors: %w", err)
	}

	gradient.Type = gradientType

	if gradientFlags&GradientFlagTransform != 0 {
		t := utils.ReadTransformable(r)
		gradient.Transformable = &t
	}

	for colorID := byte(0); colorID < ncolors; colorID++ {
		var color Color
		var offset uint8
		err := binary.Read(r, binary.LittleEndian, &offset)
		if err != nil {
			return gradient, fmt.Errorf("reading color [%d] offset: %w", colorID, err)
		}

		if gradientFlags&GradientFlagGrays != 0 {
			if gradientFlags&GradientFlagNoAlpha != 0 {
				var gc SolidGrayNoAlpha
				err := binary.Read(r, binary.LittleEndian, &gc)
				if err != nil {
					return gradient, fmt.Errorf("reading color [%d] solid gray without alpha: %w", colorID, err)
				}
				color = gc.ToColor()
			} else {
				var gc SolidGray
				err := binary.Read(r, binary.LittleEndian, &gc)
				if err != nil {
					return gradient, fmt.Errorf("reading color [%d] solid gray: %w", colorID, err)
				}
				color = gc.ToColor()
			}
		} else {
			if gradientFlags&GradientFlagNoAlpha != 0 {
				var gc SolidColorNoAlpha
				err := binary.Read(r, binary.LittleEndian, &gc)
				if err != nil {
					return gradient, fmt.Errorf("reading color [%d] solid without alpha: %w", colorID, err)
				}
				color = gc.ToColor()
			} else {
				var gc SolidColor
				err := binary.Read(r, binary.LittleEndian, &gc)
				if err != nil {
					return gradient, fmt.Errorf("reading color [%d] solid: %w", colorID, err)
				}
				color = gc.ToColor()
			}
		}

		gradient.Colors = append(gradient.Colors, GradientColor{
			StopOffset: offset,
			Color:      color,
		})
	}

	return gradient, nil
}

func Read(r io.Reader) (Style, error) {
	var styleType Type
	err := binary.Read(r, binary.LittleEndian, &styleType)
	if err != nil {
		return nil, fmt.Errorf("reading style type: %w", err)
	}

	switch styleType {
	case StyleSolidColor:
		var s SolidColor
		err := binary.Read(r, binary.LittleEndian, &s)
		if err != nil {
			return nil, fmt.Errorf("reading solid color: %w", err)
		}

		return s.ToColor(), nil
	case StyleSolidColorNoAlpha:
		var sna SolidColorNoAlpha
		err := binary.Read(r, binary.LittleEndian, &sna)
		if err != nil {
			return nil, fmt.Errorf("reading solid color without alpha: %w", err)
		}

		return sna.ToColor(), nil
	case StyleSolidGray:
		var sg SolidGray
		err := binary.Read(r, binary.LittleEndian, &sg)
		if err != nil {
			return nil, fmt.Errorf("reading solid gray color: %w", err)
		}

		return sg.ToColor(), nil
	case StyleSolidGrayNoAlpha:
		var sgna SolidGrayNoAlpha
		err := binary.Read(r, binary.LittleEndian, &sgna)
		if err != nil {
			return nil, fmt.Errorf("reading solid gray color without alpha: %w", err)
		}

		return sgna.ToColor(), nil
	case StyleGradient:
		gradient, err := readGradient(r)
		if err != nil {
			return nil, fmt.Errorf("reading gradient: %w", err)
		}

		return gradient, nil
	}

	return nil, fmt.Errorf("unknown style: %d", styleType)
}
