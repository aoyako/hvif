package main

type Image struct {
	styles []Style
}

func (i *Image) GetStyles() []Style {
	return i.styles
}

func (i *Image) AddStyle() {}
