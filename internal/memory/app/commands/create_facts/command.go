package createfacts

import (
	"github.com/google/uuid"
)

type Command struct {
	NodeID uuid.UUID
	Text   string
	Goal   string
}

func NewCommand(nodeID uuid.UUID, text string, goal string) Command {
	return Command{
		NodeID: nodeID,
		Text:   text,
		Goal:   goal,
	}
}
