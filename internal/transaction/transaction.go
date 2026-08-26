package transaction

type Backend interface {
	Commit() error
	Rollback() error
	AuditSuccess()
}

func Finish(backend Backend, operationErr error) (resultErr error) {
	if operationErr != nil {
		// Failure path: roll back and surface the original error. Do not
		// commit or audit as success — the deferred commit/audit below is
		// intentionally only attached to the success path.
		_ = backend.Rollback()
		return operationErr
	}
	defer func() {
		if commitErr := backend.Commit(); commitErr != nil {
			resultErr = commitErr
			return
		}
		backend.AuditSuccess()
	}()
	return nil
}
