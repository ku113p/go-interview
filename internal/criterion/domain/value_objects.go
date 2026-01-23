package domain

type Description string

func NewDescription(value string) Description {
	return Description(value)
}

func (d Description) String() string {
	return string(d)
}
