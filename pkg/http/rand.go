package http

import (
	"math/rand"
	"time"
)

// Range defines the range for the generated int value: [min, max]
type Range struct {
	min int
	max int
}

// Random defines the range
type Random struct {
	rand *rand.Rand
}

// NewRandom creates a random integer generator
func NewRandom() *Random {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return &Random{rand: r1}
}

func (random *Random) randn(r Range) int {
	return random.rand.Intn(r.max-r.min+1) + r.min
}

func (random *Random) randomSleep(min, max int) {
	range2 := Range{min: min, max: max}
	time.Sleep(time.Duration(random.randn(range2)) * time.Millisecond)
}
