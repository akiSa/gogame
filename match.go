package main

import (
	"github.com/gorilla/websocket"
)

type Match struct {
	Teams [2][1]Char //Two teams of 1 chars each, shrug
	//Map [][][]int //[X][y][depth, type (marsh etc) ]
	Map [][]Tile
	//Timer, or rather, a ticker, which will pause everytime a player gets a turn.
	ID string //uuid
	Socket *websocket.Conn `json:'-'`
}
func (m *Match) Send (msg *Message) (err error){
	err = m.Socket.WriteJSON(msg)
	return
}
func (m *Match) Find (ID int) *Char {
	for x, team := range m.Teams {
		for y, _ := range team {
			if m.Teams[x][y].ID == ID {
				return &m.Teams[x][y]
			}
		}
	}

	return nil
}
func (m *Match) Tick () (turn bool, list []int) {
	turn = false; //No funny business >.>
	for x, team := range m.Teams {
		for y, _ := range team {
			m.Teams[x][y].CT += m.Teams[x][y].Stats.Spd
			if team[y].CT >= 100 {
				//Every action has a cost and this will execute right after the ticks/retrieve action list anyway.
				turn = true
				list = append(list, team[y].ID)
			}
		}
	}
	return;
}
