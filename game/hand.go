package game

import "github.com/mikodream/mahjong/card"

type Hand struct {
	tiles []card.ID
}

func NewHand() *Hand {
	return &Hand{tiles: make([]card.ID, 0, 17)}
}

func (h *Hand) AddTiles(tiles []card.ID) {
	h.tiles = append(h.tiles, tiles...)
}

func (h *Hand) Tiles() []card.ID {
	tiles := make([]card.ID, len(h.tiles))
	copy(tiles, h.tiles)
	return tiles
}

func (h *Hand) Empty() bool {
	return len(h.tiles) == 0
}

func (h *Hand) RemoveTile(tile card.ID) {
	for index, tileInHand := range h.tiles {
		if tileInHand == tile {
			h.tiles[index] = h.tiles[len(h.tiles)-1]
			h.tiles = h.tiles[:len(h.tiles)-1]
			return
		}
	}
}

func (h *Hand) RemoveTiles(tiles []card.ID) {
	for _, t := range tiles {
		h.RemoveTile(t)
	}
}

func (h *Hand) Size() int {
	return len(h.tiles)
}
