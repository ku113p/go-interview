package transport

import (
	"context"
	"time"
)

type Envelope struct {
	TaskID   string
	Ctx      context.Context
	Module   string
	Action   string
	Payload  any
	Metadata map[string]string
}

func NewEnvelope(
	ctx context.Context,
	taskID, module, action string,
	payload any,
	metadata map[string]string,
) Envelope {
	if ctx == nil {
		ctx = context.Background()
	}

	return Envelope{
		TaskID:   taskID,
		Ctx:      ctx,
		Module:   module,
		Action:   action,
		Payload:  payload,
		Metadata: cloneMetadata(metadata),
	}
}

type Result struct {
	TaskID     string
	Module     string
	Action     string
	Value      any
	Err        error
	FinishedAt time.Time
	Duration   time.Duration
}

func NewResult(
	taskID, module, action string,
	value any,
	err error,
	startedAt time.Time,
) Result {
	finished := time.Now().UTC()

	return Result{
		TaskID:     taskID,
		Module:     module,
		Action:     action,
		Value:      value,
		Err:        err,
		FinishedAt: finished,
		Duration:   finished.Sub(startedAt),
	}
}

type HandlerFunc func(ctx context.Context, env Envelope) (any, error)

func cloneMetadata(src map[string]string) map[string]string {
	if len(src) == 0 {
		return nil
	}

	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}

	return dst
}
