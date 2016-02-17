package quobar

// Resolution maps physical measurements to pixels. This is useful for
// drawing widgets at a suitable size regardless to display pixel
// density.
type Resolution float32

// NewResolution computes the resolution of a screen with the given
// dimensions.
func NewResolution(pixels uint16, mm uint32) Resolution {
	return Resolution(float32(pixels) / float32(mm))
}

// Pixels converts millimeter measurements into pixels.
func (p Resolution) Pixels(mm float32) int {
	return int(float32(p)*mm + 0.5)
}
