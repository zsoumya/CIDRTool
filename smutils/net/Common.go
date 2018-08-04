package net

import (
	"errors"
	"fmt"
	"math"
)

func pow2(n byte) uint32 {
	pow := uint32(1)
	for i := byte(1); i <= n; i++ {
		pow = pow * 2
	}

	return pow
}

func blockSize(diff uint32) byte {
	exp := math.Log2(float64(diff))
	return 32 - byte(exp)
}

func splitCountToBlockSize(count uint) (byte, error) {
	f := math.Log2(float64(count))
	if f != float64(int(f)) {
		return 0, errors.New(fmt.Sprintf("invalid split count (has to be power of 2): %d", count))
	}

	return byte(f), nil
}
