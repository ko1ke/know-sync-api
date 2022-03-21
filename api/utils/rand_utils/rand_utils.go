package rand_utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func MakeRandomUInt(maxInt int) uint {
	return uint(rand.Intn(maxInt))
}
