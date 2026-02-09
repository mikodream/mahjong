package win

import (
	"testing"

	"github.com/feel-easy/mahjong/card"
)

func BenchmarkWin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		handCards := []card.ID{6, 7, 9, 9, 12, 12, 13, 14, 15, 15, 17, 26, 27, 28}
		CanWin(handCards, nil)
	}
}
