package internal

type Error string

func (e Error) Error() string {
	return string(e)
}

func (e Error) Unwrap() error {
	return nil
}

const ErrNotFound = Error("not found")
