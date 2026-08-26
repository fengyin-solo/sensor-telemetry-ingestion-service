package transaction

type Backend interface {
	Commit() error
	Rollback() error
	AuditSuccess()
}

func Finish(backend Backend, operationErr error) (resultErr error) {
	defer func() {
		if commitErr := backend.Commit(); commitErr != nil {
			resultErr = commitErr
		}
		backend.AuditSuccess()
	}()
	if operationErr != nil {
		_ = backend.Rollback()
		return operationErr
	}
	return nil
}
