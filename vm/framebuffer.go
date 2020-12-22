package vm

type FramebufferHandler interface {
	Read(addr uint16) byte
	Write(addr uint16, value byte)
	VSync()
}

type dummyFrameBuffer struct{}

func (fb dummyFrameBuffer) Read(addr uint16) byte {
	return 0
}

func (fb dummyFrameBuffer) Write(addr uint16, value byte) {
}

func (fb dummyFrameBuffer) VSync() {
}
