package constants

type ExpenseStatus string

const (
	ExpenseStatusPending   ExpenseStatus = "pending"
	ExpenseStatusApproved  ExpenseStatus = "approved"
	ExpenseStatusRejected  ExpenseStatus = "rejected"
	ExpenseStatusCompleted ExpenseStatus = "completed"
)

// does not require type conversion when used in domains
const (
	MinExpenseAmount  int64 = 10000
	MaxExpenseAmount  int64 = 50000000
	ApprovalThreshold int64 = 1000000
)
