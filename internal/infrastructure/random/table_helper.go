package random

import (
	"math/rand"
)

func CoinFlip(probability float64) bool {
	if probability > 1 {
		return true
	} else if probability < 0 {
		return false
	}

	var rollChance = int(probability * 100)
	var roll = rand.Intn(99) + 1

	if roll <= rollChance {
		return true
	} else {
		return false
	}
}
