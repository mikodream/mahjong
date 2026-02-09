package game

import (
	"fmt"
	"sort"

	"github.com/feel-easy/mahjong/card"
	"github.com/feel-easy/mahjong/consts"
	"github.com/feel-easy/mahjong/tile"
)

// ShowCard 明牌
type ShowCard struct {
	opCode int       // 操作类型，对应吃、碰、杠对应的操作类型id
	target int       // 明牌对象，吃、碰、杠的牌是谁打出来的
	tiles  []card.ID // 关联的牌
	free   bool      // 是否付费，用于转弯杠，暂时用不上了
	show   bool
}

// NewShowCard 生成一个明牌
func NewShowCard(opCode, target int, tiles []card.ID, show bool, free bool) *ShowCard {
	sort.Slice(tiles, func(i, j int) bool { return tiles[i] < tiles[j] })
	return &ShowCard{opCode: opCode, target: target, tiles: tiles, show: show, free: free}
}

func (s *ShowCard) String() string {
	tileString := tile.ToTileString(s.tiles)
	if !s.show {
		tileString = "暗杠"
	}
	return fmt.Sprintf("[%v]%v", consts.OpCodeData[s.opCode], tileString)
}

func (s *ShowCard) StringOpen() string {
	tileString := tile.ToTileString(s.tiles)
	opName := consts.OpCodeData[s.opCode]
	if !s.show {
		opName = "暗杠"
	}
	return fmt.Sprintf("[%v]%v", opName, tileString)
}

// GetOpCode 获取明牌类型
func (s *ShowCard) GetOpCode() int {
	return s.opCode
}

// GetTiles 获取明牌类型
func (s *ShowCard) GetTiles() []card.ID {
	return s.tiles
}

// GetTile 返回明牌中的牌是什么
// 至于这个showCard是不是吃，需要外面的逻辑判断
func (s *ShowCard) GetTile() card.ID {
	return s.tiles[0]
}

// GetTarget 获取明牌对象
func (s *ShowCard) GetTarget() int {
	return s.target
}

// IsFree 是否免费
func (s *ShowCard) IsFree() bool {
	return s.free
}

// GetTilesLen 牌的数量
func (s *ShowCard) GetTilesLen() int {
	return len(s.tiles)
}

// ModifyPongToKong 将碰修改成杠
func (s *ShowCard) ModifyPongToKong(kongCode int, free bool) {
	s.opCode = kongCode
	s.free = free
	s.tiles = append(s.tiles, s.tiles[0])
	sort.Slice(s.tiles, func(i, j int) bool { return s.tiles[i] < s.tiles[j] })
}

// ModifyQiangKong 将kong设置为被抢的状态
func (s *ShowCard) ModifyQiangKong() {
	s.tiles = append([]card.ID{}, s.tiles[0:s.GetTilesLen()-1]...)
}

// IsPong 明牌是否是peng
func (s *ShowCard) IsPeng() bool {
	return s.opCode == consts.PENG
}

// IsPongTile 明牌是否是peng了这个牌
func (s *ShowCard) IsPengTile(tile card.ID) bool {
	return s.opCode == consts.PENG && s.tiles[0] == tile
}
