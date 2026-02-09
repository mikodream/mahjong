package game

import (
	"fmt"
	"sort"
	"strings"

	"github.com/feel-easy/mahjong/card"
	"github.com/feel-easy/mahjong/tile"
	"github.com/feel-easy/mahjong/ting"
)

type State struct {
	LastPlayer        *PlayerController
	OriginallyPlayer  *PlayerController
	CurrentPlayer     *PlayerController
	LastPlayedTile    card.ID
	PlayedTiles       []card.ID
	CurrentPlayerHand []card.ID
	PlayerSequence    []*PlayerController
	PlayerShowCards   map[string][]*ShowCard
	SpecialPrivileges map[int][]int
	CanWin            []*PlayerController
	// Adding fields used in game.go ExtractState if needed, but better to fix game.go.
	// game.go uses: ActivePlayer (matches CurrentPlayer?), LastPlayedTileFrom, AllPlayersID.
	// I prefer adding them here if server relies on them.
	// But checking server usage is hard.
	// I will stick to what state.go has, and fix game.go.
}

func (s State) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("playedTiles:%s", tile.ToTileString(s.PlayedTiles)))
	var playerStatuses []string

	for _, player := range s.PlayerSequence {
		playerStatus := fmt.Sprintf("%s:", player.Name())
		if canTing, _ := ting.CanTing(player.Hand(), player.GetShowCardTiles()); canTing {
			playerStatus += "(听)"
		}
		if showCards, ok := s.PlayerShowCards[player.Name()]; ok && len(showCards) > 0 {
			for _, showCard := range showCards {
				playerStatus += fmt.Sprintf("%s ", showCard.String())
			}
		}
		playerStatuses = append(playerStatuses, playerStatus)
	}
	drew := s.CurrentPlayerHand[len(s.CurrentPlayerHand)-1]

	// Separate melds from hand
	standingHand := make([]card.ID, len(s.CurrentPlayerHand))
	copy(standingHand, s.CurrentPlayerHand)
	var myShowCards []*ShowCard
	if scs, ok := s.PlayerShowCards[s.CurrentPlayer.Name()]; ok {
		myShowCards = scs
		for _, sc := range scs {
			for _, t := range sc.tiles {
				for i, v := range standingHand {
					if v == t {
						standingHand = append(standingHand[:i], standingHand[i+1:]...)
						break
					}
				}
			}
		}
	}

	sort.Slice(standingHand, func(i, j int) bool { return standingHand[i] < standingHand[j] })

	// Check Ting
	canTing := false
	var tingTiles []card.ID
	tingStatus := ""

	// If standingHand has 3n+2 tiles (e.g. 14), we need to discard one to Ting.
	if len(standingHand)%3 == 2 {
		// Identify which discards lead to Ting
		type DiscardTing struct {
			Discard card.ID
			Ting    []card.ID
		}
		var suggestions []DiscardTing

		// Use specific unique tiles to avoid redundant checks
		checked := make(map[card.ID]bool)
		for _, discard := range standingHand {
			if checked[discard] {
				continue
			}
			checked[discard] = true

			// Create a temp hand with this tile removed
			tempHand := make([]card.ID, len(standingHand)-1)
			k := 0
			removed := false
			for _, v := range standingHand {
				if !removed && v == discard {
					removed = true
					continue
				}
				tempHand[k] = v
				k++
			}

			if ok, tings := ting.CanTing(tempHand, GetShowCardTiles(myShowCards)); ok {
				suggestions = append(suggestions, DiscardTing{Discard: discard, Ting: tings})
			}
		}

		if len(suggestions) > 0 {
			tingStatus = "(听)"
			var parts []string
			for _, s := range suggestions {
				var tStrs []string
				for _, t := range s.Ting {
					tStrs = append(tStrs, tile.Tile(t).String())
				}
				parts = append(parts, fmt.Sprintf("打 %s 听 %s", tile.Tile(s.Discard), strings.Join(tStrs, " ")))
			}
			tingStatus += " " + strings.Join(parts, "; ")
		}

	} else {
		// Normal check (e.g. 13 tiles)
		canTing, tingTiles = ting.CanTing(standingHand, GetShowCardTiles(myShowCards))
		if canTing {
			tingStatus = "(听)"
			tingStr := []string{}
			for _, t := range tingTiles {
				tingStr = append(tingStr, tile.Tile(t).String())
			}
			tingStatus += fmt.Sprintf(" %s", strings.Join(tingStr, " "))
		}
	}

	if s.LastPlayer != nil {
		lines = append(lines, fmt.Sprintf("ShowCards:\n%s ", strings.Join(playerStatuses, "\n")))
		lines = append(lines, fmt.Sprintf("%s played: %s", s.LastPlayer.Name(), tile.Tile(s.LastPlayedTile).String()))
	}
	if len(standingHand)%3 == 2 {
		// Calculate drew from standingHand end, because we sorted standingHand, so drew might be mixed.
		// Actually s.CurrentPlayerHand last element is the drew one.
		// But standingHand is sorted.
		// Visual display: Just show "Your drew" if we have 14 tiles (modulo 3 == 2).
		// The `drew` variable is from raw `CurrentPlayerHand`, so it is correct as the last added tile.
		lines = append(lines, fmt.Sprintf("Your drew: %s ", tile.Tile(drew)))
	}

	// Format Hand + Melds
	handStr := tile.ToTileString(standingHand)
	if len(myShowCards) > 0 {
		meldStrs := []string{}
		for _, sc := range myShowCards {
			meldStrs = append(meldStrs, sc.StringOpen())
		}
		handStr += " | " + strings.Join(meldStrs, " ")
	}

	lines = append(lines, fmt.Sprintf("Your hand: %s", handStr))
	if tingStatus != "" {
		lines = append(lines, tingStatus)
	}
	return strings.Join(lines, "\n") + "\n"
}

func GetShowCardTiles(scs []*ShowCard) []card.ID {
	ret := make([]card.ID, 0, len(scs)*4)
	for _, t := range scs {
		ret = append(ret, t.tiles...)
	}
	return ret
}
