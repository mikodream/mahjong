package game

import (
	"github.com/mikodream/mahjong/card"
	"github.com/mikodream/mahjong/consts"
	"github.com/mikodream/mahjong/tile"
	"github.com/mikodream/mahjong/win"
)

type Game struct {
	players *PlayerIterator
	deck    *Deck
	pile    *Pile
}

func (g *Game) Players() *PlayerIterator {
	return g.players
}

func (g *Game) Deck() *Deck {
	return g.deck
}

func (g *Game) Pile() *Pile {
	return g.pile
}

func (g *Game) Next() *PlayerController {
	player := g.Players().Next()
	g.pile.SetCurrentPlayer(player)
	return player
}

func New(players []Player) *Game {
	return &Game{
		players: newPlayerIterator(players),
		deck:    NewDeck(),
		pile:    NewPile(),
	}
}

func (g *Game) GetPlayerTiles(id int) string {
	tiles := g.players.GetPlayerController(id).Hand()
	return tile.ToTileString(tiles)
}

func (g *Game) DealStartingTiles() {
	g.players.ForEach(func(player *PlayerController) {
		hand := g.deck.Draw(13)
		player.AddTiles(hand)
	})
}

func (g *Game) Current() *PlayerController {
	return g.players.Current()
}

func (g Game) ExtractState(player *PlayerController) State {
	playerSequence := make([]*PlayerController, 0)
	playerShowCards := make(map[string][]*ShowCard)
	specialPrivileges := make(map[int][]int)
	canWin := make([]*PlayerController, 0)
	originallyPlayer := g.pile.originallyPlayer
	topTile := g.pile.Top()
	g.players.ForEach(func(player *PlayerController) {
		playerSequence = append(playerSequence, player)
		playerShowCards[player.Name()] = player.GetShowCard()
		if _, ok := g.pile.SayNoPlayer()[player.ID()]; !ok &&
			topTile > 0 && g.pile.lastPlayer.ID() != player.ID() {
			handWithTop := make([]card.ID, len(player.Hand()))
			copy(handWithTop, player.Hand())
			handWithTop = append(handWithTop, topTile)
			if win.CanWin(handWithTop, player.GetShowCardTiles()) {
				canWin = append(canWin, player)
			}
			// Note: card package functions currently take int, will update to card.ID
			// Casting for now or updating card package later?
			// I will update card package to take card.ID.
			if card.CanMingGang(player.Hand(), topTile) {
				specialPrivileges[player.ID()] = append(specialPrivileges[player.ID()], consts.GANG)
			}
			if card.CanPeng(player.Hand(), topTile) {
				specialPrivileges[player.ID()] = append(specialPrivileges[player.ID()], consts.PENG)
			}
			if originallyPlayer.ID() == player.ID() &&
				card.CanChi(player.Hand(), topTile) {
				specialPrivileges[player.ID()] = append(specialPrivileges[player.ID()], consts.CHI)
			}
		}
	})
	return State{
		PlayerSequence:   playerSequence,
		PlayerShowCards:  playerShowCards,
		CurrentPlayer:    player, // Renamed ActivePlayer -> CurrentPlayer
		LastPlayedTile:   topTile,
		LastPlayer:       g.pile.LastPlayer(),
		OriginallyPlayer: originallyPlayer,
		// LastPlayedTileFrom: g.pile.LastPlayer().ID(), // Removed as not in State struct
		// AllPlayersID:       g.players.cycler.Elements(), // Removed as not in State struct
		PlayedTiles:       g.pile.Tiles(), // Added
		CurrentPlayerHand: player.Tiles(), // Added
		SpecialPrivileges: specialPrivileges,
		CanWin:            canWin,
	}
}
