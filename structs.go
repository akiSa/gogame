package main;

import (

	//"code.google.com/p/go-uuid/uuid"
)
type Message struct {
	Action string //String: turn/execute
	Players []int `json:", omitempty"`
	//Actions []Action `json:", omitempty"`
	Actions map[int][]Action `json:", omitempty"`
}
type Conflict struct {
	Actions []*Action
	Type int //Move conflict (chars walking close/into each other), Clash, Attack conflict (A attacking B attacking C attacking A, etc)
}
//One thing I could do, represent the map as a [][]*Char, each XY is either a char or nil, if it's a char then dew damage or whatever, this is probably better because it speeds up the searching by a lot.
//Look at the map in respect to the map itself rather than in respect to the char.
//So basically, [][][]int, the third = [depth, type, unit]
//unit = a char ID, or 0... but this doesn't speed up searching if there IS a unit, it only speeds up searching if there is no unit (one check rather than checking the whole user list)
//So I'd need a Map structure
//[][][]...Tile
//Tile {
//Depth
//Type
//Unit *Char (pointer to the char on it)
//}
//That will speed up searching a lot, because you'd just do if map[x][y].Char != nil, blah


type Tile struct {
	Depth int
	Type int
	Unit *Char
}
type Char struct {
	ID int
	HP int
	CT int
	Stats statList
	X int // X Pos
	Y int // Y Pos
	ACList Actions//[]Action
	//Team int // 0 or 1
}
func (ch *Char) setXY (X int, Y int, m *Match) {
	m.Map[ch.X][ch.Y].Unit = nil
	ch.X = X; ch.Y = Y
	m.Map[X][Y].Unit = ch
}

type statList struct {
	Str int
	Vit int
	Int int
	Wis int
	Dex int
	Spd int
}
type Actions struct {
	Actions []Action
	TICKCT int
}
type Action struct {
	ID int //Char ID.. tagging purposes, makes life so easy -_-
	SX int //Starting X Pos
	SY int //Starting Y Pos
	EX int //Expected X Pos, or Target X Pos (for attacks)
	EY int //Expected Y Pos, or Target Y Pos (for attacks)
	Type []int //Two ints, first is Action Type: (Mobility, Attack, etc), second is action ID, ie what exactly is being done.
	CT int //CT cost, precalculated.
}
