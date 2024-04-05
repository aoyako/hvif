package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

type styleType uint8

const (
	styleSolidColor styleType = 1 + iota
	styleGradient
	styleSolidColorNoAlpha
	styleSolidGray
	styleSolidGrayNoAlpha
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

type gradientFlag uint8

const (
	gradientFlagTransform gradientFlag = 1 << (1 + iota)
	gradientFlagNoAlpha
	gradientFlag16BitColors // Unused
	gradientFlagGrays
)

type Style any // Color | Gradient

type Color struct {
	Red   uint8
	Green uint8
	Blue  uint8
	Alpha uint8
}

type Gradient struct {
	Type          GradientType
	Transformable *Transformable
	Colors        []Color
	Offsets       []uint8
}

type solidColor struct {
	Red   uint8
	Green uint8
	Blue  uint8
	Alpha uint8
}

type solidColorNoAlpha struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

type solidGray struct {
	Gray  uint8
	Alpha uint8
}

type solidGrayNoAlpha struct {
	Gray uint8
}

func (sc solidColor) toColor() Color {
	//nolint
	return Color{
		Red:   sc.Red,
		Green: sc.Green,
		Blue:  sc.Blue,
		Alpha: sc.Alpha,
	}
}

func (scna solidColorNoAlpha) toColor() Color {
	return Color{
		Red:   scna.Red,
		Green: scna.Green,
		Blue:  scna.Blue,
		Alpha: 0xff,
	}
}

func (sg solidGray) toColor() Color {
	return Color{
		Red:   sg.Gray,
		Green: sg.Gray,
		Blue:  sg.Gray,
		Alpha: sg.Alpha,
	}
}

func (sgna solidGrayNoAlpha) toColor() Color {
	return Color{
		Red:   sgna.Gray,
		Green: sgna.Gray,
		Blue:  sgna.Gray,
		Alpha: 0xff,
	}
}

func readGradient(r io.Reader) (Gradient, error) {
	var gradient Gradient

	var gradientType GradientType
	err := binary.Read(r, binary.LittleEndian, &gradientType)
	if err != nil {
		return gradient, fmt.Errorf("reading type: %w", err)
	}
	gradient.Type = gradientType

	var gradientFlags gradientFlag
	err = binary.Read(r, binary.LittleEndian, &gradientFlags)
	if err != nil {
		return gradient, fmt.Errorf("reading flags: %w", err)
	}

	var ncolors uint8
	err = binary.Read(r, binary.LittleEndian, &ncolors)
	if err != nil {
		return gradient, fmt.Errorf("reading number of colors: %w", err)
	}

	if gradientFlags&gradientFlagTransform != 0 {
		t := ReadTransformable(r)
		gradient.Transformable = &t
	}

	for colorID := byte(0); colorID < ncolors; colorID++ {
		var offset uint8
		err := binary.Read(r, binary.LittleEndian, &offset)
		if err != nil {
			return gradient, fmt.Errorf("reading [%d] offset: %w", colorID, err)
		}

		var cType styleType
		switch {
		case gradientFlags&gradientFlagGrays != 0 && gradientFlags&gradientFlagNoAlpha != 0:
			cType = styleSolidGrayNoAlpha
		case gradientFlags&gradientFlagGrays != 0 && gradientFlags&gradientFlagNoAlpha == 0:
			cType = styleSolidGray
		case gradientFlags&gradientFlagGrays == 0 && gradientFlags&gradientFlagNoAlpha != 0:
			cType = styleSolidColorNoAlpha
		case gradientFlags&gradientFlagGrays == 0 && gradientFlags&gradientFlagNoAlpha == 0:
			cType = styleSolidColor
		default:
			return gradient, fmt.Errorf("figuring color [%d] type %d", colorID, gradientFlags)
		}

		color, err := readColor(r, cType)
		if err != nil {
			return gradient, fmt.Errorf("reading color [%d]: %w", colorID, err)
		}

		gradient.Colors = append(gradient.Colors, color)
		gradient.Offsets = append(gradient.Offsets, offset)
	}

	return gradient, nil
}

func readColor(r io.Reader, cType styleType) (Color, error) {
	switch cType {
	case styleSolidColor:
		var s solidColor
		err := binary.Read(r, binary.LittleEndian, &s)
		if err != nil {
			return Color{}, fmt.Errorf("reading solid color: %w", err)
		}

		return s.toColor(), nil
	case styleSolidColorNoAlpha:
		var sna solidColorNoAlpha
		err := binary.Read(r, binary.LittleEndian, &sna)
		if err != nil {
			return Color{}, fmt.Errorf("reading solid color without alpha: %w", err)
		}

		return sna.toColor(), nil
	case styleSolidGray:
		var sg solidGray
		err := binary.Read(r, binary.LittleEndian, &sg)
		if err != nil {
			return Color{}, fmt.Errorf("reading solid gray color: %w", err)
		}

		return sg.toColor(), nil
	case styleSolidGrayNoAlpha:
		var sgna solidGrayNoAlpha
		err := binary.Read(r, binary.LittleEndian, &sgna)
		if err != nil {
			return Color{}, fmt.Errorf("reading solid gray color without alpha: %w", err)
		}

		return sgna.toColor(), nil
	}

	return Color{}, fmt.Errorf("color %d not recognized", cType)
}

func readStyle(r io.Reader) (Style, error) {
	var styleType styleType
	err := binary.Read(r, binary.LittleEndian, &styleType)
	if err != nil {
		return nil, fmt.Errorf("reading style type: %w", err)
	}

	switch styleType {
	case styleSolidColor, styleSolidColorNoAlpha, styleSolidGray, styleSolidGrayNoAlpha:
		c, err := readColor(r, styleType)
		if err != nil {
			return nil, fmt.Errorf("reading color: %w", err)
		}

		return c, nil
	case styleGradient:
		gradient, err := readGradient(r)
		if err != nil {
			return nil, fmt.Errorf("reading gradient: %w", err)
		}

		return gradient, nil
	}

	return nil, fmt.Errorf("unknown style: %d", styleType)
}
