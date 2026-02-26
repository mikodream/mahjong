package ting

import (
	"reflect"
	"sort"
	"testing"

	"github.com/mikodream/mahjong/card"
)

// 测试可能听的牌
func TestMaybeTing(t *testing.T) {
	handCards := []card.ID{1}
	showcards := []card.ID{}
	maybeCards := []card.ID{}

	maybeCards = GetMaybeTing(handCards, showcards)
	sort.Slice(maybeCards, func(i, j int) bool { return maybeCards[i] < maybeCards[j] })
	if !reflect.DeepEqual(maybeCards, []card.ID{1, 2}) {
		t.Error("验证可能听的牌失败1")
	}
	handCards = []card.ID{1, 2}
	maybeCards = GetMaybeTing(handCards, showcards)
	sort.Slice(maybeCards, func(i, j int) bool { return maybeCards[i] < maybeCards[j] })
	if !reflect.DeepEqual(maybeCards, []card.ID{1, 2, 3}) {
		t.Error("验证可能听的牌失败2")
	}

	handCards = []card.ID{6, 9}
	maybeCards = GetMaybeTing(handCards, showcards)
	sort.Slice(maybeCards, func(i, j int) bool { return maybeCards[i] < maybeCards[j] })
	if !reflect.DeepEqual(maybeCards, []card.ID{5, 6, 7, 8, 9}) {
		t.Error("验证可能听的牌失败3")
	}

	handCards = []card.ID{6, 9}
	showcards = []card.ID{3, 3, 3}
	maybeCards = GetMaybeTing(handCards, showcards)
	sort.Slice(maybeCards, func(i, j int) bool { return maybeCards[i] < maybeCards[j] })
	if !reflect.DeepEqual(maybeCards, []card.ID{3, 5, 6, 7, 8, 9}) {
		t.Error("验证可能听的牌失败4")
	}

	// 11 = BAM1. 11,12,13,14,15,16. Neighbors: 11 has 11,12. 16 has 15,16,17.
	// So 11..17. 3 is added from showcards.
	handCards = []card.ID{11, 11, 12, 12, 13, 13, 14, 14, 15, 15, 16}
	showcards = []card.ID{3, 3, 3}
	maybeCards = GetMaybeTing(handCards, showcards)
	sort.Slice(maybeCards, func(i, j int) bool { return maybeCards[i] < maybeCards[j] })
	if !reflect.DeepEqual(maybeCards, []card.ID{3, 11, 12, 13, 14, 15, 16, 17}) {
		t.Error("验证可能听的牌失败4")
	}

	handCards = []card.ID{11, 11, 12, 12, 13, 13, 14, 14, 15, 15, 16}
	showcards = []card.ID{8, 8, 8, 8}
	maybeCards = GetMaybeTing(handCards, showcards)
	sort.Slice(maybeCards, func(i, j int) bool { return maybeCards[i] < maybeCards[j] })
	if !reflect.DeepEqual(maybeCards, []card.ID{11, 12, 13, 14, 15, 16, 17}) {
		t.Error("验证可能听的牌失败5")
	}
}

// 测试能不能听牌
func TestCanTing(t *testing.T) {
	handCards := []card.ID{1}
	showcards := []card.ID{}
	isTing := false
	tingCards := []card.ID{}

	isTing, tingCards = CanTing(handCards, showcards)
	if !isTing || !reflect.DeepEqual(tingCards, []card.ID{1}) {
		t.Error("验证叫牌失败1")
	}

	handCards = []card.ID{1, 2, 3, 4}
	showcards = []card.ID{}
	isTing, tingCards = CanTing(handCards, showcards)
	sort.Slice(tingCards, func(i, j int) bool { return tingCards[i] < tingCards[j] })
	if !isTing || !reflect.DeepEqual(tingCards, []card.ID{1, 4}) {
		t.Error("验证叫牌失败2")
	}

	handCards = []card.ID{1, 1, 2, 2, 3, 3, 4, 4, 5, 5}
	showcards = []card.ID{8, 8, 8}
	isTing, tingCards = CanTing(handCards, showcards)
	sort.Slice(tingCards, func(i, j int) bool { return tingCards[i] < tingCards[j] })

	// Expect {1, 2, 4, 5} as 8 is not a winner.
	if !isTing || !reflect.DeepEqual(tingCards, []card.ID{1, 2, 4, 5}) {
		t.Errorf("验证叫牌失败3, got %v", tingCards)
	}

	handCards = []card.ID{1, 1, 1, 2, 3, 4, 5, 6, 7, 8, 9, 9, 9}
	showcards = []card.ID{}
	isTing, tingCards = CanTing(handCards, showcards)
	sort.Slice(tingCards, func(i, j int) bool { return tingCards[i] < tingCards[j] })
	if !isTing || !reflect.DeepEqual(tingCards, []card.ID{1, 2, 3, 4, 5, 6, 7, 8, 9}) {
		t.Error("验证叫牌失败4")
	}

	handCards = []card.ID{1, 1, 3, 3, 5, 5, 7, 7, 9, 9, 11, 11, 18}
	showcards = []card.ID{}
	isTing, tingCards = CanTing(handCards, showcards)
	sort.Slice(tingCards, func(i, j int) bool { return tingCards[i] < tingCards[j] })
	if !isTing || !reflect.DeepEqual(tingCards, []card.ID{18}) {
		t.Error("验证叫牌失败5")
	}
}
