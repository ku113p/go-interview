package domain

import (
	common "go-interview/internal/common/domain"
)

type MediaType string

const (
	Video MediaType = "video"
	Audio MediaType = "audio"
	Text  MediaType = "text"
)

type RawData struct {
	common.Entity
	S3Path    string
	MediaType MediaType
}

func NewRawData(s3Path string, mediaType MediaType) *RawData {
	rawData := &RawData{
		S3Path:    s3Path,
		MediaType: mediaType,
	}

	common.InitEntity(&rawData.Entity)

	return rawData
}
