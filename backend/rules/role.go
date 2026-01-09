package rules

import (
	"backend/constants"
	"backend/models"
	"errors"
)

var (
	ErrForbidden = errors.New("forbidden action")
)

func IsManager(user *models.User) bool {
	return user.Role == constants.UserRoleManager
}

func IsUser(user *models.User) bool {
	return user.Role == constants.UserRoleUser
}

func CanApproveRole(user *models.User) error {
	if !IsManager(user) {
		return ErrForbidden
	}
	return nil
}
