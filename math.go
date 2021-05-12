package sparky

import (
	"math"
)

// MaxFloat32 returns the larger of the passed values.
func MaxFloat32(nums ...float32) float32 {
	max := float32(0)
	for _, n := range nums {
		if n > max {
			max = n
		}
	}
	return max
}

// MinFloat32 returns the smaller of the passed values.
func MinFloat32(nums ...float32) float32 {
	min := float32(math.MaxFloat32)
	for _, n := range nums {
		if n < min {
			min = n
		}
	}
	return min
}
