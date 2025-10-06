package model

// Сообщение валидации сервиса
type Validation struct {
    Property string
    Message string
}

// Результат сервиса, данные или валидация
type ServiceDataResult[TD any] struct {
	Result   bool
	Message  string
	Validations* []Validation
	Code *string
	Data *TD
}

// Результат сервиса, список или валидация
type ServiceListResult[TD any] struct {
	Result   bool
	Message  string
	Validations* []Validation
	Code *string
	Count int
	Items *[]TD
}