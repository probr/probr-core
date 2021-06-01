package run

import "fmt"

// ServicePackError retains an error object and the name of the pack that generated it
type ServicePackError struct {
	ServicePack string
	Err         error
}

// ServicePackErrors holds a list of errors and an Error() method
// so it adheres to the standard Error interface
type ServicePackErrors struct {
	Errors []ServicePackError
}

func (e *ServicePackErrors) Error() string {
	return fmt.Sprintf("Service Pack Errors: %v", e.Errors)
}
