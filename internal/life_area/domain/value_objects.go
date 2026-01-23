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
