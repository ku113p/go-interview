package domain

type Content string

func NewContent(value string) Content {
	return Content(value)
}
