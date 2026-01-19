package domain

type Title string

func NewTitle(s string) Title {
	return Title(s)
}

func (t Title) String() string {
	return string(t)
}

type Goal string

func NewGoal(s string) Goal {
	return Goal(s)
}

func (g Goal) String() string {
	return string(g)
}

type Description string

func NewDescription(s string) Description {
	return Description(s)
}

func (d Description) String() string {
	return string(d)
}
