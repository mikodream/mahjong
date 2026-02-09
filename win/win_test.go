package win

import (
	"reflect"
	"testing"

	"github.com/feel-easy/mahjong/card"
)

// TestFindPairPos 测试找出对的位置是否正确
func TestFindPairPos(t *testing.T) {
	cards := []card.ID{1, 1, 1, 1, 5, 8, 8, 6, 7, 11, 11}
	findPos := FindPairPos(cards)
	// 注意：新逻辑下 1,1,1,1 会被识别为两对，分别在索引 0 和 2
	if !reflect.DeepEqual(findPos, []int{0, 2, 5, 9}) {
		t.Errorf("TestFindPairPos Error, findPos:%v", findPos)
	}
}

// 测试删除pair的判断
func TestRemovePair(t *testing.T) {
	cards := []card.ID{1, 1, 1, 1, 5, 6, 7, 8, 8, 11, 11}

	// 移除索引9的对子 (11, 11)
	if !reflect.DeepEqual(RemovePair(cards, 9), []card.ID{1, 1, 1, 1, 5, 6, 7, 8, 8}) {
		t.Error("removePair 验证失败 - case 1")
	}

	// 移除索引7的对子 (8, 8)
	if !reflect.DeepEqual(RemovePair(cards, 7), []card.ID{1, 1, 1, 1, 5, 6, 7, 11, 11}) {
		t.Error("removePair 验证失败 - case 2")
	}

	// 移除索引2的对子 (第2组 1, 1) -> 注意原数组是 [1, 1, 1, 1 ...]
	// RemovePair(2) 会移除 index 2 和 3
	if !reflect.DeepEqual(RemovePair(cards, 2), []card.ID{1, 1, 5, 6, 7, 8, 8, 11, 11}) {
		t.Error("removePair 验证失败 - case 3")
	}
}

// 测试是否全部是顺或者刻 (回溯算法测试)
func TestIsAllSequenceOrTriplet(t *testing.T) {
	var cards []card.ID

	// 1,1,1 (刻) + 2,3,4 (顺, 假设是同花色) + 3,4,5 (顺)
	// 注意：这里需要确保数字对应的花色是序数牌(万筒条)。
	// 假设测试用例里的数字都在万字牌范围内(1-9)

	// Case 1: 混合牌型
	cards = []card.ID{1, 1, 1, 2, 3, 4, 3, 4, 5}
	// 排序后: 1,1,1, 2,3,3,4,4,5
	// 拆解: [1,1,1] + [2,3,4] + [3,4,5] -> OK
	if !IsAllSequenceOrTriplet(cards) {
		t.Errorf("TestIsAllSequenceOrTriplet failed1:%v", cards)
	}

	// Case 2: 两顺
	cards = []card.ID{1, 2, 3, 2, 3, 4}
	if !IsAllSequenceOrTriplet(cards) {
		t.Errorf("TestIsAllSequenceOrTriplet failed2:%v", cards)
	}

	// Case 3: 复杂拆解 (旧贪心算法容易出错的 Case)
	// 2,3,4 + 3,3,3
	// 如果贪心算法先拿了 3,3,3，剩下 2,4 就挂了
	// 回溯算法应该先拿 2,3,4，剩下 3,3 (不是刻子) -> 回溯 -> 实际上这个例子需要配对

	// 让我们测一个 2,3,3,3,4
	// 这不是由纯粹的顺子和刻子组成的 (5张牌)，这通常是在有将牌移除后调用的
	// 我们测一个标准的：2,3,4 + 2,3,4
	cards = []card.ID{2, 2, 3, 3, 4, 4}
	if !IsAllSequenceOrTriplet(cards) {
		t.Errorf("TestIsAllSequenceOrTriplet failed3:%v", cards)
	}
}

func TestCanWin(t *testing.T) {
	var handCards []card.ID

	// 1. 七对 (Seven Pairs)
	handCards = []card.ID{
		card.MAHJONG_CRAK1, card.MAHJONG_CRAK1,
		card.MAHJONG_CRAK2, card.MAHJONG_CRAK2,
		card.MAHJONG_CRAK3, card.MAHJONG_CRAK3,
		card.MAHJONG_CRAK4, card.MAHJONG_CRAK4,
		card.MAHJONG_CRAK5, card.MAHJONG_CRAK5,
		card.MAHJONG_CRAK6, card.MAHJONG_CRAK6,
		card.MAHJONG_CRAK7, card.MAHJONG_CRAK7,
	}
	if !CanWin(handCards, nil) {
		t.Error("7对类型验证失败")
	}

	// 2. 破坏七对 (其中一张单着)
	handCards[13] = card.MAHJONG_CRAK9 // 最后一对变成 7, 9
	if CanWin(handCards, nil) {
		t.Error("7对错误用例验证失败")
	}

	// 3. 标准胡牌 (3n + 2) - 对对胡
	handCards = []card.ID{1, 1, 2, 2, 2, 3, 3, 3}
	if !CanWin(handCards, nil) {
		t.Error("标准对对胡验证失败")
	}

	// 4. 标准胡牌 - 顺子混搭
	// 11(将) + 123 + 234
	handCards = []card.ID{1, 1, 1, 2, 3, 2, 3, 4}
	if !CanWin(handCards, nil) {
		t.Error("标准顺子牌验证失败")
	}

	// 5. 单吊 (1张手牌，实际上 CanWin 接收的是 "手牌+摸到的牌"，所以一共2张)
	handCards = []card.ID{1, 1}
	if !CanWin(handCards, nil) {
		t.Error("单吊牌型验证失败")
	}
	handCards = []card.ID{1, 2}
	if CanWin(handCards, nil) {
		t.Error("错误单吊牌型验证失败")
	}

	// 6. [新增] 十三幺测试
	// 19万, 19筒, 19条, 东南西北, 中发白. 其中 1万 两张
	thirteenOrphans := []card.ID{
		card.MAHJONG_CRAK1, card.MAHJONG_CRAK1, card.MAHJONG_CRAK9,
		card.MAHJONG_DOT1, card.MAHJONG_DOT9,
		card.MAHJONG_BAM1, card.MAHJONG_BAM9,
		card.MAHJONG_EAST, card.MAHJONG_SOUTH, card.MAHJONG_WEST, card.MAHJONG_NORTH,
		card.MAHJONG_RED, card.MAHJONG_GREE, card.MAHJONG_WHITE,
	}
	if !CanWin(thirteenOrphans, nil) {
		t.Error("十三幺验证失败")
	}

	// 7. [新增] 错误的十三幺 (缺一张)
	thirteenOrphansFalse := make([]card.ID, len(thirteenOrphans))
	copy(thirteenOrphansFalse, thirteenOrphans)
	thirteenOrphansFalse[0] = card.MAHJONG_CRAK2 // 把1万改成2万
	if CanWin(thirteenOrphansFalse, nil) {
		t.Error("错误十三幺验证失败")
	}
}

func TestGetTingTiles(t *testing.T) {
	// 例子：1112345678999万
	// 也就是著名的“九莲宝灯”牌型，听所有同花色的牌(1-9万)
	handCards := []card.ID{
		card.MAHJONG_CRAK1, card.MAHJONG_CRAK1, card.MAHJONG_CRAK1,
		card.MAHJONG_CRAK2, card.MAHJONG_CRAK3, card.MAHJONG_CRAK4,
		card.MAHJONG_CRAK5, card.MAHJONG_CRAK6, card.MAHJONG_CRAK7,
		card.MAHJONG_CRAK8,
		card.MAHJONG_CRAK9, card.MAHJONG_CRAK9, card.MAHJONG_CRAK9,
	}

	tingList := GetTingTiles(handCards, nil)
	t.Logf("九莲宝灯听牌列表: %v", tingList)

	if len(tingList) != 9 {
		t.Errorf("九莲宝灯应该听9张牌，实际听了 %d 张", len(tingList))
	}
}
