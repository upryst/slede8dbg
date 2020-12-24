package vga

import (
	"image"
	"image/color"

	"gioui.org/op/paint"
)

const (
	Width  = 256
	Height = 256
	Colors = 256
)

type Framebuffer struct {
	image *image.Paletted

	imageOp      paint.ImageOp
	pendingVSync bool
}

func NewFramebuffer() *Framebuffer {
	fb := &Framebuffer{
		image: image.NewPaletted(image.Rect(0, 0, Width, Height),
			makeDefaultPalette()),
	}

	fb.imageOp = paint.NewImageOp(fb.image)

	return fb
}

// The name is misleading, but all we do here
// is making sure the texture caches are updated with
// the new fb.image
func (fb *Framebuffer) VSync() {
	fb.imageOp = paint.NewImageOp(fb.image)
	fb.pendingVSync = true
}

func (fb *Framebuffer) SetPaletteR(idx int, r uint8) {
	_, g, b, _ := fb.image.Palette[idx].RGBA()
	fb.image.Palette[idx] = color.RGBA{r, byte(g), byte(b), 0xff}
}

func (fb *Framebuffer) SetPaletteG(idx int, g uint8) {
	r, _, b, _ := fb.image.Palette[idx].RGBA()
	fb.image.Palette[idx] = color.RGBA{byte(r), g, byte(b), 0xff}
}

func (fb *Framebuffer) SetPaletteB(idx int, b uint8) {
	r, g, _, _ := fb.image.Palette[idx].RGBA()
	fb.image.Palette[idx] = color.RGBA{byte(r), byte(g), b, 0xff}
}

func (fb *Framebuffer) Read(addr uint16) byte {
	return fb.image.ColorIndexAt(int(addr&0xff), int(addr>>8))
}

func (fb *Framebuffer) Write(addr uint16, value byte) {
	fb.image.SetColorIndex(int(addr&0xff), int(addr>>8), value)
}

func makeDefaultPalette() color.Palette {
	p := make(color.Palette, Colors)

	for i := range p {
		if i < len(defaultVGAPalette) {
			p[i] = defaultVGAPalette[i]
		}
	}

	return p
}
