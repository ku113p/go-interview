package postgres

import "go-interview/internal/user/domain"

func userToSQL(user *domain.User) *UserSQL {
	return &UserSQL{
		ID:         user.ID,
		CreatedAt:  user.CreatedAt,
		ExternalID: user.ExternalID,
	}
}

func (dto *UserSQL) ToDomain() *domain.User {
	return &domain.User{
		ID:         dto.ID,
		CreatedAt:  dto.CreatedAt,
		ExternalID: dto.ExternalID,
	}
}
