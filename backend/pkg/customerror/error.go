package customerror

import "fmt"

type CustomError struct {
	Module   string
	Endpoint string
	Message  string
}

func (customError CustomError) Error() string {
	return fmt.Sprintf("ERROR|%s|%s:%s", customError.Endpoint, customError.Module, customError.Message)
}

func (customError *CustomError) AppendModule(module string) {
	customError.Module = module + "." + customError.Module
}

func NewError(module, endpoint, message string) error {
	return CustomError{
		Module:   module,
		Endpoint: endpoint,
		Message:  message,
	}
}
