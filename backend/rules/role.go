package rules

import (
	"backend/constants"
	"backend/models"
	"errors"
)

var (
	ErrForbidden = errors.New("forbidden action")
)

func CanApproveRole(user *models.User) error {
	if user.Role != constants.UserRoleManager {
		return ErrForbidden
	}
	return nil
}
