package postgres

import (
	"go-interview/internal/message/domain"
)

func rawDataToSQL(raw *domain.RawData) *RawDataSQL {
	return &RawDataSQL{
		ID:        raw.ID,
		CreatedAt: raw.CreatedAt,
		S3Path:    raw.S3Path,
		MediaType: string(raw.MediaType),
	}
}

func (dto *RawDataSQL) ToDomain() *domain.RawData {
	return &domain.RawData{
		ID:        dto.ID,
		CreatedAt: dto.CreatedAt,
		S3Path:    dto.S3Path,
		MediaType: domain.MediaType(dto.MediaType),
	}
}

func transcriptToSQL(t *domain.Transcript) *TranscriptSQL {
	return &TranscriptSQL{
		ID:        t.ID,
		CreatedAt: t.CreatedAt,
		NodeID:    t.NodeID,
		RawDataID: t.RawDataID,
		Content:   string(t.Content),
	}
}

func (dto *TranscriptSQL) ToDomain() *domain.Transcript {
	return &domain.Transcript{
		ID:        dto.ID,
		CreatedAt: dto.CreatedAt,
		NodeID:    dto.NodeID,
		RawDataID: dto.RawDataID,
		Content:   domain.Content(dto.Content),
	}
}
