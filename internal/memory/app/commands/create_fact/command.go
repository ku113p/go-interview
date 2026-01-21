package createfact

import (
	"time"

	"github.com/google/uuid"
)

type Command struct {
	NodeID   uuid.UUID
	Info     string
	DateTime *time.Time
}

func NewCommand(nodeID uuid.UUID, info string, dateTime *time.Time) Command {
	return Command{
		NodeID:   nodeID,
		Info:     info,
		DateTime: dateTime,
	}
}
