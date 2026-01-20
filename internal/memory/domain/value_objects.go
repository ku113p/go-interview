package domain

type Vector []float64

func NewVector(value []float64) Vector {
	return value
}

type Content string

func NewContent(value string) Content {
	return Content(value)
}

type Info string

func NewInfo(value string) Info {
	return Info(value)
}
