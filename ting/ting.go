package ting

import (
	"github.com/feel-easy/mahjong/card"
	"github.com/feel-easy/mahjong/win"
)

// CanTing 判断牌型是否可以听牌
// 返回是否可听、听什么
func CanTing(handCards, showCards []card.ID) (bool, []card.ID) {
	var canTing = false
	tingCards := make([]card.ID, 0)

	// 循环将可能听的牌，带入到手牌，再用胡牌算法检测是否可胡
	for _, t := range GetMaybeTing(handCards, showCards) {
		tempHand := make([]card.ID, len(handCards), len(handCards)+1)
		copy(tempHand, handCards)
		tempHand = append(tempHand, t)
		if win.CanWin(tempHand, showCards) {
			canTing = true
			tingCards = append(tingCards, t)
		}
	}
	return canTing, tingCards
}

// GetMaybeTing 获取哪些牌是可能听的
// 东南西北、花色等，只有自身
// 边张只有自身或者上下张的某一张
// 其他的是自身和上下张
// 如果有明牌，且明牌是3张的话，则明牌也可能是胡的
func GetMaybeTing(handCards, showCards []card.ID) []card.ID {
	// 转换 ID 到 int 传给 card 包 (兼容旧方法)
	// 理想状况下应该更新 card 包方法签名，但暂保持兼容
	handInts := make([]int, len(handCards))
	for i, v := range handCards {
		handInts[i] = int(v)
	}

	maybeInts := card.GetSelfAndNeighborCards(handInts...)
	maybeCards := make([]card.ID, len(maybeInts))
	for i, v := range maybeInts {
		maybeCards[i] = card.ID(v)
	}

	if len(showCards) == 3 &&
		showCards[0] == showCards[1] && showCards[1] == showCards[2] {

		// 检查 showCards[0] 是否已在 maybeCards 中
		exists := false
		for _, v := range maybeCards {
			if v == showCards[0] {
				exists = true
				break
			}
		}

		if !exists {
			maybeCards = append(maybeCards, showCards[0])
		}
	}
	return maybeCards
}

// GetTingMap 获取可听的列表
// key: 打什么
// value: 听哪些
func GetTingMap(handCards, showCards []card.ID) map[card.ID][]card.ID {
	tingMap := make(map[card.ID][]card.ID)

	// 去重手牌
	uniqueHands := make(map[card.ID]bool)
	for _, c := range handCards {
		uniqueHands[c] = true
	}

	for playCard := range uniqueHands {
		// 删除一张牌
		tempHand := sliceDel(handCards, playCard)

		if ting, tingCards := CanTing(tempHand, showCards); ting {
			tingMap[playCard] = tingCards
		}
	}
	return tingMap
}

func sliceDel(slice []card.ID, value card.ID) []card.ID {
	if slice == nil {
		return slice
	}
	for i, j := range slice {
		if j == value {
			newSlice := make([]card.ID, len(slice)-1)
			copy(newSlice, slice[:i])
			copy(newSlice[i:], slice[i+1:])
			return newSlice
		}
	}
	return slice
}
