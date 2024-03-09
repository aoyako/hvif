package style

type Type uint8
type GradientType uint8

const (
	StyleSolidColor Type = 1 + iota
	StyleGradient
	StyleSolidColorNoAlpha
	StyleSolidGray
	StyleSolidGrayNoAlpha
)

const (
	GradientLinear GradientType = iota
	GradientCircular
	GradientDiamond
	GradientConic
	GradientXY
	GradientSqrtXY
)

type SolidColor struct {
	Red   uint8
	Green uint8
	Blue  uint8
	Alpha uint8
}

type Gradient struct {
	Type   GradientType
	Colors []GradientColor
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

type GradientColor struct {
	StopOffset uint8
	Alpha      uint8
	Red        uint8
	Green      uint8
	Blue       uint8
}
