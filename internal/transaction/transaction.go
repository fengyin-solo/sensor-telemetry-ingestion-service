package transaction

import "errors"

type Backend interface {
	Commit() error
	Rollback() error
	AuditSuccess()
}

func Finish(backend Backend, operationErr error) error {
	if operationErr != nil {
		if rollbackErr := backend.Rollback(); rollbackErr != nil {
			return errors.Join(operationErr, rollbackErr)
		}
		return operationErr
	}
	if err := backend.Commit(); err != nil {
		return err
	}
	backend.AuditSuccess()
	return nil
}
