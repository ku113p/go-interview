package domain

type Description struct {
	value string
}

func NewDescription(value string) Description {
	return Description{value}
}

func (d *Description) String() string {
	return d.value
}
