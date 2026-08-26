package exporter

func clonePayloads(input [][]byte) [][]byte {
	result := make([][]byte, len(input))
	for index, payload := range input {
		result[index] = append([]byte(nil), payload...)
	}
	return result
}

func ExportAsync(payloads [][]byte, ready chan<- struct{}, proceed <-chan struct{}) <-chan [][]byte {
	result := make(chan [][]byte, 1)
	snapshot := clonePayloads(payloads)
	go func() {
		if ready != nil {
			close(ready)
		}
		if proceed != nil {
			<-proceed
		}
		result <- snapshot
		close(result)
	}()
	return result
}
