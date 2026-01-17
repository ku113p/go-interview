package domain

type Goal struct {
	value string
}

func NewGoal(value string) Goal {
	return Goal{value}
}

func (g *Goal) String() string {
	return g.value
}
