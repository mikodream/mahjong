package win

import (
	"sort"

	"github.com/mikodream/mahjong/card"
)

// CanWin 判断当前牌型是否是胡牌牌型
// 支持：标准胡牌(3n+2), 七对, 十三幺
func CanWin(handTiles, showTiles []card.ID) bool {
	// 复制并排序，以免修改原切片
	sortedTiles := make([]card.ID, len(handTiles))
	copy(sortedTiles, handTiles)
	sort.Slice(sortedTiles, func(i, j int) bool { return sortedTiles[i] < sortedTiles[j] })

	// 1. 判断十三幺 (Thirteen Orphans)
	// 十三幺必须是门清（没有碰/杠/吃，即 showTiles 为空）
	// 但考虑到有些游戏实现可能只传 handTiles，这里只检查手牌数量是否足够
	if len(sortedTiles) == 14 && IsThirteenOrphans(sortedTiles) {
		return true
	}

	// 2. 判断七对 (Seven Pairs)
	// 七对必须是14张牌
	if len(sortedTiles) == 14 {
		pairs := FindPairPos(sortedTiles)
		if len(pairs) == 7 {
			return true
		}
	}

	// 3. 标准胡牌逻辑 (x * (ABC/AAA) + DD)
	// 遍历每一个对子，尝试将其作为"将牌"(Eye)
	pairs := FindPairPos(sortedTiles)
	var lastPairTile card.ID = -1

	for _, idx := range pairs {
		pairTile := sortedTiles[idx]
		// 优化：如果手牌有 4 张 8万，会产生 3 个对子索引，我们只需要判断一次
		if pairTile == lastPairTile {
			continue
		}
		lastPairTile = pairTile

		// 移除这对将牌
		remainingTiles := RemovePair(sortedTiles, idx)

		// 使用回溯法检查剩下的牌是否全为顺子或刻子
		if IsAllSequenceOrTriplet(remainingTiles) {
			return true
		}
	}

	return false
}

// GetTingTiles 获取听牌列表
// 输入：handTiles (当前手里的牌，通常是 1, 4, 7, 10, 13 张)
// 输出：所有能胡的牌 ID 列表
func GetTingTiles(handTiles, showTiles []card.ID) []card.ID {
	tingList := make([]card.ID, 0)

	// 遍历麻将所有可能的 34 种牌
	for _, t := range card.AllTiles {
		tID := card.ID(t)
		// 模拟摸这一张牌：
		// 1. 复制当前手牌 (避免修改原切片)
		// 2. 加上这张牌 t
		tempHand := make([]card.ID, len(handTiles), len(handTiles)+1)
		copy(tempHand, handTiles)
		tempHand = append(tempHand, tID)

		// 检查加上这张牌后是否胡了
		if CanWin(tempHand, showTiles) {
			tingList = append(tingList, tID)
		}
	}

	return tingList
}

// IsThirteenOrphans 判断十三幺
// 牌型：19万, 19筒, 19条, 东南西北, 中发白 (共13种)，其中一种有2张(做将)，其余各1张
func IsThirteenOrphans(sortedTiles []card.ID) bool {
	if len(sortedTiles) != 14 {
		return false
	}

	// 十三幺所需要的 13 种牌的 ID
	orphans := map[card.ID]bool{
		card.MAHJONG_CRAK1: true, card.MAHJONG_CRAK9: true, // 1, 9 万
		card.MAHJONG_DOT1: true, card.MAHJONG_DOT9: true, // 1, 9 筒
		card.MAHJONG_BAM1: true, card.MAHJONG_BAM9: true, // 1, 9 条
		card.MAHJONG_EAST: true, card.MAHJONG_SOUTH: true, // 东, 南
		card.MAHJONG_WEST: true, card.MAHJONG_NORTH: true, // 西, 北
		card.MAHJONG_RED: true, card.MAHJONG_GREE: true, card.MAHJONG_WHITE: true, // 中, 发, 白
	}

	// 统计手牌中这13种牌的数量
	counts := make(map[card.ID]int)
	for _, t := range sortedTiles {
		if !orphans[t] {
			return false // 如果包含任何非么九字牌，直接失败
		}
		counts[t]++
	}

	// 必须涵盖所有 13 种牌
	if len(counts) != 13 {
		return false
	}

	// 检查是否有且仅有一种牌是 2 张 (作将)，其余都是 1 张
	hasEye := false
	for _, c := range counts {
		if c == 2 {
			if hasEye {
				return false // 只能有一个将
			}
			hasEye = true
		} else if c != 1 {
			return false // 数量不对
		}
	}

	return hasEye
}

// FindPairPos 找出所有对牌的位置 (输入必须已排序)
// 返回每对牌第一张的索引
func FindPairPos(sortedTiles []card.ID) []int {
	var pos = []int{}
	if len(sortedTiles) < 2 {
		return pos
	}
	length := len(sortedTiles) - 1
	for i := 0; i < length; i++ {
		if sortedTiles[i] == sortedTiles[i+1] {
			pos = append(pos, i)
			i++ // 跳过下一张，确保一对只算一次
		}
	}
	return pos
}

// RemovePair 从已排序的牌中，移除指定位置的一对，返回新切片
func RemovePair(sortedTiles []card.ID, pos int) []card.ID {
	remain := make([]card.ID, 0, len(sortedTiles)-2)
	remain = append(remain, sortedTiles[:pos]...)
	remain = append(remain, sortedTiles[pos+2:]...)
	return remain
}

// IsAllSequenceOrTriplet 回溯法判断是否全为顺子或刻子
func IsAllSequenceOrTriplet(tiles []card.ID) bool {
	// 递归终止条件：牌分完了，说明符合条件
	if len(tiles) == 0 {
		return true
	}

	// 1. 尝试作为刻子 (Triplet: AAA)
	if len(tiles) >= 3 && IsTriplet(tiles[0], tiles[1], tiles[2]) {
		// 移除刻子后递归
		if IsAllSequenceOrTriplet(tiles[3:]) {
			return true
		}
		// 如果这条路走不通，回溯，继续往下尝试顺子
	}

	// 2. 尝试作为顺子 (Sequence: ABC)
	// 只有万筒条(Suit)才能做顺子，字牌(Honor)不能
	first := tiles[0]
	if first.IsSuit() && len(tiles) >= 3 {
		// 在有序数组中寻找 first+1 和 first+2
		var secondIdx, thirdIdx = -1, -1

		// 寻找第二张 (first + 1)
		for i := 1; i < len(tiles); i++ {
			if tiles[i] == first+1 {
				secondIdx = i
				break
			} else if tiles[i] > first+1 {
				break // 排序过的，大了就不可能找到了
			}
		}

		// 只有找到了第二张，才找第三张 (first + 2)
		if secondIdx != -1 {
			for i := secondIdx + 1; i < len(tiles); i++ {
				if tiles[i] == first+2 {
					thirdIdx = i
					break
				} else if tiles[i] > first+2 {
					break
				}
			}
		}

		// 如果找到了完整的顺子
		if secondIdx != -1 && thirdIdx != -1 {
			// 构建移除这三张牌后的新切片
			remaining := make([]card.ID, 0, len(tiles)-3)
			remaining = append(remaining, tiles[1:secondIdx]...)          // 跳过第0张，取0到2之间的
			remaining = append(remaining, tiles[secondIdx+1:thirdIdx]...) // 取2到3之间的
			remaining = append(remaining, tiles[thirdIdx+1:]...)          // 取3之后的

			if IsAllSequenceOrTriplet(remaining) {
				return true
			}
		}
	}

	return false
}

// IsTriplet 是否刻子
func IsTriplet(a, b, c card.ID) bool {
	return a == b && b == c
}
