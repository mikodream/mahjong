package game

type PlayerIterator struct {
	players map[int]*PlayerController
	cycler  *Cycler
}

func (i *PlayerIterator) GetPlayerController(id int) *PlayerController {
	return i.players[id]
}

func newPlayerIterator(players []Player) *PlayerIterator {
	var playerIDs []int
	playerMap := make(map[int]*PlayerController, len(players))
	for _, player := range players {
		playerID := player.PlayerID()
		playerIDs = append(playerIDs, playerID)
		playerMap[playerID] = NewPlayerController(player)
	}
	return &PlayerIterator{
		players: playerMap,
		cycler:  NewCycler(playerIDs),
	}
}

func (i *PlayerIterator) Current() *PlayerController {
	return i.players[i.cycler.Current()]
}

func (i *PlayerIterator) ForEach(function func(player *PlayerController)) {
	for range i.players {
		function(i.Current())
		i.Next()
	}
}

func (i *PlayerIterator) Next() *PlayerController {
	return i.players[i.cycler.Next()]
}
