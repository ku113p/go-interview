package domain

type Content struct {
	value string
}

func NewContent(value string) Content {
	return Content{value}
}
