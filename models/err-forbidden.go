package models

// TrainersNotFoundError кастомная ошибка.
type ForbiddenError struct{}

// Error inheritdoc.
func (*ForbiddenError) Error() string {
	return "forbidden"
}
