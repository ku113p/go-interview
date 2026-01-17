package domain

type Goal struct {
	value string
}

func NewGoal(value string) Goal {
	return Goal{value}
}
