package domain

import (
	common "go-interview/internal/common/domain"
)

type Vector struct {
	common.Entity
	value []float64
}

func NewVector(value []float64) *Vector {
	vector := &Vector{
		value: value,
	}
	common.InitEntity(&vector.Entity)
	return vector
}
