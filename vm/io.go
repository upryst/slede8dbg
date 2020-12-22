package vm

type IOHandler interface {
	ReadPort(addr uint16) byte
	WritePort(addr uint16, value byte)
}

type dummyIOHandler struct{}

func (h dummyIOHandler) ReadPort(addr uint16) byte {
	return 0
}

func (h dummyIOHandler) WritePort(addr uint16, value byte) {
}
