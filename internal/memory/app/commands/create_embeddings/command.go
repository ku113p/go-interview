package createembeddings

import "github.com/google/uuid"

type Command struct {
	NodeID uuid.UUID
	Text   string
}

func NewCommand(nodeID uuid.UUID, text string) Command {
	return Command{
		NodeID: nodeID,
		Text:   text,
	}
}
