package blend

import (
	"image/color"
)

// Threshold gives a color to use for values near Max.
type Threshold struct {
	Max   uint64
	Color color.Color
}

// PickColor picks a color in based on the given thresholds.
//
// Thresholds are expected to be in decreasing order.
//
// If value is greater than or equal to the greatest Max, the color
// from that threshold is returned.
//
// If value is less than the Max of the last threshold, the color from
// that threshold is returned.
//
// Otherwise, the color will be proportionally in between the color of
// the matching threshold and the next one, based on where value is
// between the two Max values.
func PickColor(thresholds []Threshold, value uint64) color.Color {
	for idx, cur := range thresholds {
		next := Threshold{Max: 0, Color: cur.Color}
		if idx < len(thresholds)-1 {
			next = thresholds[idx+1]
		}
		if value < next.Max {
			continue
		}
		p := 1 - (float32(value-next.Max) / float32(cur.Max-next.Max))
		return Gradient(cur.Color, next.Color, p)
	}

	// The slice is logically extended with an item with Max 0, no
	// uint64 input can be less than 0. We can only reach here if
	// there were no thresholds.
	return color.Black
}
