package capture

import "telemetry/internal/transaction"

type Resource interface {
	Read() error
	Close() error
}

type Opener interface {
	Open(name string) (Resource, error)
}

func processOne(opener Opener, name string) (err error) {
	resource, err := opener.Open(name)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := resource.Close(); err == nil {
			err = closeErr
		}
	}()
	return resource.Read()
}

func Process(opener Opener, backend transaction.Backend, names []string) error {
	var operationErr error
	for _, name := range names {
		// processOne opens, reads, and closes each resource within its own
		// scope, so the handle is released before the next item is opened.
		// A loop-body defer would accumulate until this function returns and
		// exhaust the handle quota before reaching later items.
		if err := processOne(opener, name); err != nil {
			operationErr = err
			break
		}
	}
	return transaction.Finish(backend, operationErr)
}
