package user

import "github.com/Bipolar-Penguin/bff-website/pkg/domain"

type Repository interface {
	Find(userID string) (domain.User, error)
	Save(user domain.User) (domain.User, error)
}
