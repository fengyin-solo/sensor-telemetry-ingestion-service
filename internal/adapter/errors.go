package adapter

import "fmt"

type RemoteError struct {
	Code       string
	RetryAfter int
}

func (e *RemoteError) Error() string { return fmt.Sprintf("remote response %s", e.Code) }

type RejectedError struct{ Cause error }

func (e *RejectedError) Error() string { return "reading rejected: " + e.Cause.Error() }
func (e *RejectedError) Unwrap() error { return e.Cause }

type TemporaryError struct {
	Cause      error
	RetryAfter int
}

func (e *TemporaryError) Error() string { return "temporary delivery failure: " + e.Cause.Error() }
func (e *TemporaryError) Unwrap() error { return e.Cause }

func Translate(err error) error {
	remote, ok := err.(*RemoteError)
	if !ok {
		return err
	}
	return fmt.Errorf("downstream delivery failed: %v", remote)
}
