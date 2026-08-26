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
	for _, name := range names {
		if err := processOne(opener, name); err != nil {
			return transaction.Finish(backend, err)
		}
	}
	return transaction.Finish(backend, nil)
}
