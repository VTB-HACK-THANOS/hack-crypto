package models

// TrainersNotFoundError кастомная ошибка.
type AlreadyExistsError struct{}

// Error inheritdoc.
func (AlreadyExistsError) Error() string {
	return "already exists"
}
