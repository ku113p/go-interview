package domain

type Title struct {
	value string
}

func NewTitle(value string) Title {
	return Title{value}
}

func (t *Title) String() string {
	return t.value
}
