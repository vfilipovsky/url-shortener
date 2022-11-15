package randomizer

//go:generate mockgen -source=./randomizer.go -destination=./mock/mock_randomizer.go

import (
	"math/rand"
	"time"
)

var Chars = []rune("01234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var Numbers = []rune("0123456789")

type Randomizer interface {
	Random(length int, col []rune) string
}

type randomizer struct{}

func New() Randomizer {
	return &randomizer{}
}

func (r *randomizer) Random(length int, col []rune) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, length)

	for i := range b {
		b[i] = col[rand.Intn(len(col))]
	}

	return string(b)
}
