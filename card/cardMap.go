package card

import (
	"sort"
	"sync"
)

// CMap 牌面=>数量
type CMap struct {
	Mux   *sync.RWMutex
	tiles map[ID]int
}

// NewCMap 初始化一个TileMap
func NewCMap() *CMap {
	return &CMap{
		Mux:   &sync.RWMutex{},
		tiles: make(map[ID]int),
	}
}

// SetTiles 初始化手牌
func (cm *CMap) SetTiles(tiles []ID) {
	cm.Mux.Lock()
	defer cm.Mux.Unlock()
	for _, tile := range tiles {
		cm.tiles[tile]++
	}
}

// GetTileMap 读取所有牌的列表
// 这里不会主动加锁，在外面用的话，如果用于range，需要手动加锁
func (cm *CMap) GetTileMap() map[ID]int {
	return cm.tiles
}

// AddTile 添加手牌
func (cm *CMap) AddTile(tile ID, cnt int) {
	cm.Mux.Lock()
	defer cm.Mux.Unlock()
	cm.tiles[tile] += cnt
}

// DelTile 删除手牌
func (cm *CMap) DelTile(tile ID, cnt int) bool {
	cm.Mux.Lock()
	defer cm.Mux.Unlock()
	if cm.tiles[tile] > cnt {
		cm.tiles[tile] -= cnt
	} else if cm.tiles[tile] == cnt {
		delete(cm.tiles, tile)
	} else {
		return false
	}
	return true
}

// ToSlice 转成slice
func (cm *CMap) ToSlice() []ID {
	cm.Mux.RLock()
	defer cm.Mux.RUnlock()
	tiles := []ID{}
	for tile, cnt := range cm.tiles {
		for i := 0; i < cnt; i++ {
			tiles = append(tiles, tile)
		}
	}
	sort.Slice(tiles, func(i, j int) bool { return tiles[i] < tiles[j] })
	return tiles
}

// ToSortedSlice 转成slice并排序
func (cm *CMap) ToSortedSlice() []ID {
	tiles := cm.ToSlice()
	sort.Slice(tiles, func(i, j int) bool { return tiles[i] < tiles[j] })
	return tiles
}

// GetUnique 获取独立的牌
func (cm *CMap) GetUnique() []ID {
	cm.Mux.RLock()
	defer cm.Mux.RUnlock()
	tiles := []ID{}
	for tile := range cm.tiles {
		tiles = append(tiles, tile)
	}
	return tiles
}

// GetTileCnt 获取某张牌的数量
func (cm *CMap) GetTileCnt(tile ID) int {
	cm.Mux.RLock()
	defer cm.Mux.RUnlock()
	return cm.tiles[tile]
}
