package actions

import (
	"backend/constants"
	"backend/models"
)

type ExpenseAuditLogInput struct {
	ExpenseID  int64
	ActorID    *int64
	FromStatus constants.ExpenseStatus
	ToStatus   constants.ExpenseStatus
	Reason     string
}

func ExpenseAuditLog(input ExpenseAuditLogInput) (audit *models.ExpenseAuditLog, err error) {
	fromStatus := input.FromStatus

	var actorID *int64
	if input.ActorID != nil {
		actorID = input.ActorID
	}

	audit = &models.ExpenseAuditLog{
		ExpenseID:  input.ExpenseID,
		ActorID:    actorID,
		FromStatus: fromStatus,
		ToStatus:   input.ToStatus,
		Reason:     input.Reason,
	}

	return audit, nil
}
