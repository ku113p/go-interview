package domain

import common "go-interview/internal/common/domain"

type Transcript struct {
	common.UpdatableEntity
	Content common.Content
}

func NewTranscript(content common.Content) *Transcript {
	transcript := &Transcript{
		Content: content,
	}

	common.InitUpdatableEntity(&transcript.UpdatableEntity)

	return transcript
}
