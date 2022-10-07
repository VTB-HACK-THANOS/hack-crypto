package models

// TrainersNotFoundError кастомная ошибка.
type NotFoundError struct{}

// Error inheritdoc.
func (NotFoundError) Error() string {
	return "not found"
}
