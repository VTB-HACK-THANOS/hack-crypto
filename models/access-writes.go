package models

type AccessWrites int

const (
	NotAuthorizedUser AccessWrites = 0
	RegularUser       AccessWrites = 1
	PrivilegedUser    AccessWrites = 2
)
