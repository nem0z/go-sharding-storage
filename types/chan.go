package types

type WrappedChan struct {
	Data []byte
	Chan chan []byte
}
