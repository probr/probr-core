package core

import "fmt"

// ServicePackError ...
type ServicePackError struct {
	ServicePack string
	Err         error
}

// ServicePackErrors ...
type ServicePackErrors struct {
	SPErrs []ServicePackError
}

func (e *ServicePackErrors) Error() string {
	return fmt.Sprintf("Service Pack Errors: %v", e.SPErrs)
}
