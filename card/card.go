package card

import "github.com/feel-easy/mahjong/util"

func HaveGang(tiles []ID) (ID, bool) {
	cmap := NewCMap()
	cmap.SetTiles(tiles)
	for t, i := range cmap.GetTileMap() {
		if i == 4 {
			return t, true
		}
	}
	return 0, false
}

func HaveGangs(tiles []ID) []ID {
	cmap := NewCMap()
	cmap.SetTiles(tiles)
	gangs := make([]ID, 0)
	for t, i := range cmap.GetTileMap() {
		if i == 4 {
			gangs = append(gangs, t)
		}
	}
	return gangs
}

// 判断是不是可以吃
func CanChi(cards []ID, card ID) bool {
	if card == 0 {
		return false
	}
	return len(CanChiTiles(cards, card)) > 0
}

// 可以吃的牌
func CanChiTiles(cards []ID, card ID) [][]ID {
	if !card.IsSuit() {
		return [][]ID{}
	}
	ret := make([][]ID, 0)
	if IDInSlice(card+1, cards) && IDInSlice(card+2, cards) {
		ret = append(ret, []ID{card + 1, card + 2})
	}
	if IDInSlice(card-1, cards) && IDInSlice(card-2, cards) {
		ret = append(ret, []ID{card - 1, card - 2})
	}
	if IDInSlice(card-1, cards) && IDInSlice(card+1, cards) {
		ret = append(ret, []ID{card - 1, card + 1})
	}
	return ret
}

// 判断是不是可以碰
func CanPeng(cards []ID, card ID) bool {
	cmap := NewCMap()
	cmap.SetTiles(cards)
	return cmap.GetTileCnt(card) == 2
}

// 判断是不是可以杠
func CanGang(cards []ID, card ID) bool {
	cmap := NewCMap()
	cmap.SetTiles(cards)
	return cmap.GetTileCnt(card) == 4
}

// 判断是不是可以明杠
func CanMingGang(cards []ID, card ID) bool {
	cmap := NewCMap()
	cmap.SetTiles(cards)
	return cmap.GetTileCnt(card) == 3
}

// IsSuit 是否普通牌
// 普通牌是指万、筒、条
func IsSuit(card ID) bool {
	return card < MAHJONG_DOT_PLACEHOLDER
}

// GetSelfAndNeighborCards 获取自身或者相邻的一张牌, 结果需去重
// 不包括隔张
// 1条、1筒、1万只有自己和上一张
// 九条、九筒、九万只有自己和下一张
// 非万、筒、条 只有自己
func GetSelfAndNeighborCards(cards ...int) []int {
	result := []int{}
	for _, c := range cards {
		cardID := ID(c)
		result = append(result, c)
		// 非普通牌、只有自身
		if !cardID.IsSuit() {
			continue
		}
		if IDInSlice(cardID, LeftSideCards) {
			result = append(result, c+1)
		} else if IDInSlice(cardID, RightSideCards) {
			result = append(result, c-1)
		} else {
			result = append(result, c-1, c+1)
		}
	}
	return util.SliceUniqueInt(result)
}

// GetRelationTiles 获取有关联的牌
// 包括自己、相邻的、跳张
func GetRelationTiles(cards ...int) []int {
	result := []int{}
	for _, c := range cards {
		cardID := ID(c)
		result = append(result, c)
		// 非普通牌、只有自身
		if !cardID.IsSuit() {
			continue
		}

		if IDInSlice(cardID, LeftSideCards) {
			result = append(result, c+1, c+2)
		} else if IDInSlice(cardID, LeftSideNeighborCards) {
			result = append(result, c+1, c+2, c-1)
		} else {
			// 中间牌
			if IDInSlice(cardID, RightSideCards) {
				result = append(result, c-1, c-2)
			} else if IDInSlice(cardID, RightSideNeighborCards) {
				result = append(result, c-1, c-2, c+1)
			} else {
				result = append(result, c-1, c-2, c+1, c+2)
			}
		}
	}
	return util.SliceUniqueInt(result)
}

// IsHonor 是否风牌/字牌
func IsHonor(card ID) bool {
	return card >= MAHJONG_EAST
}

// IsTerminal 是否幺九牌
func IsTerminal(card ID) bool {
	if IsHonor(card) {
		return true
	}
	return IDInSlice(card, []ID{1, 9, 11, 19, 21, 29})
}

// IsYaoJiu IsTerminal 的别名
func IsYaoJiu(card ID) bool {
	return IsTerminal(card)
}

// IsGreen 是否绿一色牌
func IsGreen(card ID) bool {
	return card.IsBam() || card == MAHJONG_GREE
}

func IDInSlice(id ID, slice []ID) bool {
	for _, v := range slice {
		if v == id {
			return true
		}
	}
	return false
}

// Type methods

func (id ID) Int() int {
	return int(id)
}

func (id ID) IsSuit() bool {
	return id < MAHJONG_DOT_PLACEHOLDER
}

func (id ID) IsCrak() bool {
	return id >= MAHJONG_CRAK1 && id <= MAHJONG_CRAK9
}

func (id ID) IsBam() bool {
	return id >= MAHJONG_BAM1 && id <= MAHJONG_BAM9
}

func (id ID) IsDot() bool {
	return id >= MAHJONG_DOT1 && id <= MAHJONG_DOT9
}

func (id ID) IsHonor() bool {
	return id >= MAHJONG_EAST
}

func (id ID) Rank() int {
	return int(id) % 10
}
