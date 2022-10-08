package models

type AccessWrites int

const (
	NotAuthorizedUser AccessWrites = 0
	RegularUser       AccessWrites = 1
	ManagerUser       AccessWrites = 2
	AdminUser         AccessWrites = 3
)
