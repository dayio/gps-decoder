package source

type SDRSource struct {
	deviceContext interface{}
}

func (s *SDRSource) Read(buffer []int8) error {
	return nil
}

func (s *SDRSource) Close() error {
	return nil
}
