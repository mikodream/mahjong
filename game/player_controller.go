package game

import (
	"github.com/feel-easy/mahjong/card"
	"github.com/feel-easy/mahjong/consts"
	"github.com/feel-easy/mahjong/event"
)

type playerController struct {
	player    Player
	hand      *Hand
	showCards []*ShowCard
}

func newPlayerController(player Player) *playerController {
	return &playerController{
		player:    player,
		hand:      NewHand(),
		showCards: make([]*ShowCard, 0, 5),
	}
}

func (c *playerController) DarkGang(tile card.ID) {
	c.showCards = append(c.showCards, NewShowCard(consts.GANG, 0, []card.ID{tile, tile, tile, tile}, false, false))
}

func (c *playerController) operation(op, target int, tiles []card.ID) {
	c.showCards = append(c.showCards, NewShowCard(op, target, tiles, true, false))
}

func (c *playerController) GetShowCard() []*ShowCard {
	return c.showCards
}

func (c *playerController) FindShowCard(id card.ID) *ShowCard {
	for _, sc := range c.showCards {
		for _, tile := range sc.tiles {
			if tile == id {
				return sc
			}
		}
	}
	return nil
}

func (c *playerController) GetShowCardTiles() []card.ID {
	ret := make([]card.ID, 0, len(c.showCards)*4)
	for _, t := range c.showCards {
		ret = append(ret, t.tiles...)
	}
	return ret
}

func (c *playerController) AddTiles(tiles []card.ID) {
	c.hand.AddTiles(tiles)
}

func (c *playerController) TryTopDecking(deck *Deck) {
	extraCard := deck.DrawOne()
	c.AddTiles([]card.ID{extraCard})
	event.PlayTile.Emit(event.PlayTilePayload{
		PlayerName: c.player.NickName(),
		Tile:       extraCard,
	})
}

func (c *playerController) TryBottomDecking(deck *Deck) {
	extraCard := deck.BottomDrawOne()
	c.AddTiles([]card.ID{extraCard})
	event.PlayTile.Emit(event.PlayTilePayload{
		PlayerName: c.player.NickName(),
		Tile:       extraCard,
	})
}

func (c *playerController) Hand() []card.ID {
	tiles := c.Tiles()
	return sliceDel(tiles, c.GetShowCardTiles()...)
}

func (c *playerController) Tiles() []card.ID {
	return c.hand.Tiles()
}
func (c *playerController) LastTile() card.ID {
	return c.hand.Tiles()[len(c.hand.Tiles())-1]
}

func (c *playerController) Name() string {
	return c.player.NickName()
}

func (c *playerController) ID() int {
	return c.player.PlayerID()
}
func (c *playerController) Player() *Player {
	return &c.player
}

func (c *playerController) Take(gameState State, deck *Deck, pile *Pile) (int, bool, error) {
	tiles := make([]card.ID, 0, len(c.Hand())+1)
	tiles = append(tiles, c.Hand()...)
	tiles = append(tiles, pile.Top())
	op, tiles, err := c.player.Take(tiles, gameState)
	if err != nil {
		return op, false, err
	}
	if len(tiles) == 0 {
		switch op {
		case consts.CHI:
			c.TryTopDecking(deck)
		case consts.PENG:
			if gameState.OriginallyPlayer.ID() == c.ID() {
				c.TryTopDecking(deck)
			}
		case consts.GANG:
			if gameState.OriginallyPlayer.ID() == c.ID() {
				c.TryTopDecking(deck)
			}
		}
		pile.AddSayNoPlayer(c)
		return op, false, nil
	}
	c.AddTiles([]card.ID{pile.BottomDrawOne()})
	c.operation(op, int(pile.LastPlayer().ID()), tiles)
	return op, true, nil
}

func (c *playerController) Play(gameState State) (card.ID, error) {
	selectedTile, err := c.player.Play(c.Hand(), gameState)
	if err != nil {
		return 0, err
	}
	c.hand.RemoveTile(selectedTile)
	return selectedTile, nil
}

func (c *playerController) RemoveTile(tile card.ID) {
	c.hand.RemoveTile(tile)
}

func (c *playerController) RemoveTiles(tiles []card.ID) {
	c.hand.RemoveTiles(tiles)
}

func sliceDel(slice []card.ID, elems ...card.ID) []card.ID {
	for _, e := range elems {
		for i, v := range slice {
			if v == e {
				slice = append(slice[:i], slice[i+1:]...)
				break
			}
		}
	}
	return slice
}
