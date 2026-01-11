package models

import (
	"backend/constants"
	"time"
)

type ExpenseAuditLog struct {
	ID         int64                   `json:"id"`
	ExpenseID  int64                   `json:"expense_id"`
	ActorID    *int64                  `json:"actor_id"`
	FromStatus constants.ExpenseStatus `json:"from_status"`
	ToStatus   constants.ExpenseStatus `json:"to_status"`
	Reason     string                  `json:"reason"`
	CreatedAt  time.Time               `json:"created_at"`
}
