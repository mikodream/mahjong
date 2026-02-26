package game

import "github.com/mikodream/mahjong/card"

type Player interface {
	PlayerID() int
	NickName() string
	Play(tiles []card.ID, gameState State) (card.ID, error)
	Take(tiles []card.ID, gameState State) (int, []card.ID, error)
}
