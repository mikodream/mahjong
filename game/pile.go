package game

import "github.com/mikodream/mahjong/card"

type Pile struct {
	tiles            []card.ID
	lastPlayer       *PlayerController
	originallyPlayer *PlayerController
	currentPlayer    *PlayerController
	sayNoPlayer      map[int]*PlayerController
}

func (p *Pile) SetCurrentPlayer(player *PlayerController) {
	p.currentPlayer = player
}

func (p *Pile) CurrentPlayer() *PlayerController {
	return p.currentPlayer
}

func (p *Pile) AddSayNoPlayer(player *PlayerController) {
	if p.sayNoPlayer == nil {
		p.sayNoPlayer = make(map[int]*PlayerController)
	}
	p.sayNoPlayer[player.ID()] = player
}

func (p *Pile) SayNoPlayer() map[int]*PlayerController {
	return p.sayNoPlayer
}

func NewPile() *Pile {
	return &Pile{tiles: make([]card.ID, 0, 144)}
}

func (p *Pile) SetOriginallyPlayer(player *PlayerController) {
	p.originallyPlayer = player
	p.sayNoPlayer = make(map[int]*PlayerController)
}

func (p *Pile) OriginallyPlayer() *PlayerController {
	return p.originallyPlayer
}

func (p *Pile) SetLastPlayer(player *PlayerController) {
	p.lastPlayer = player
}

func (p *Pile) LastPlayer() *PlayerController {
	return p.lastPlayer
}

func (p *Pile) Add(tile card.ID) {
	p.tiles = append(p.tiles, tile)
}

func (p *Pile) Tiles() []card.ID {
	tiles := make([]card.ID, len(p.tiles))
	copy(tiles, p.tiles)
	return tiles
}

func (p *Pile) ReplaceTop(tile card.ID) {
	p.tiles[len(p.tiles)-1] = tile
}

func (p *Pile) Top() card.ID {
	pileSize := len(p.tiles)
	if pileSize == 0 {
		return 0
	}
	return p.tiles[pileSize-1]
}

func (d *Pile) BottomDrawOne() card.ID {
	tile := d.tiles[len(d.tiles)-1]
	d.tiles = d.tiles[0 : len(d.tiles)-1]
	return tile
}
