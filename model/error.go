package model

// ErrNotFound is returned when a requested entity is not found.
type ErrNotFound struct{}

func (e *ErrNotFound) Error() string {
	return "resource not found"
}

func (e *ErrNotFound) Is(target error) bool {
	_, ok := target.(*ErrNotFound)
	return ok
}
