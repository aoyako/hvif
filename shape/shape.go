package shape

type Shape struct {
	StyleIndex  uint8
	PathIndeces []uint8
	Modifiers   []ShapeModifier
}

type ShapeModifier struct {
	TransformMatrix Matrix
}

type Matrix struct {
}
