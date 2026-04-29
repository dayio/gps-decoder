package source

type IQSource interface {
	Read(buffer []int8) error
	Close() error
}
