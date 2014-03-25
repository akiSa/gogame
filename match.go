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
//When executing, TICKCT += action.CT && if TICKCT > char speed, TICKCT = TICKCT - char speed and wait until next tick.
//confList = [ (Conflicts) ]
//Conflict = { [ Action ] }
func (m *Match) Execute () {
	var exList []*Char
	for x, _ := range m.Teams {
		for y, _ := range m.Teams[x] {
			if len(m.Teams[x][y].ACList.Actions) != 0 {
				exList = append(exList, &m.Teams[x][y])
			}
		}
	}
	if len(exList) > 0 {
		//If there are more than 0 characters ready to execute..
		if conflicts := confCheck(exList); len(conflicts) > 0 {
			//If there are any conflicts..
			//Resolve..
		}
		msg := new (Message)
		msg.Action = "execute"
		acCH := make(chan Action)
		doneCH := make(chan bool)

		for x, _ := range exList {
			//For each char to be executed..
			go func (ch *Char, acCH chan Action, doneCH chan bool) {
				//Make a thread to handle it..
				for {
					if len(ch.ACList.Actions) == 0 {
						doneCH <- true //Action list expended, DONE.
						return
					}
					ch.ACList.TICKCT += ch.ACList.Actions[0].CT
					if ch.ACList.TICKCT > ch.Stats.Spd {
						ch.ACList.TICKCT -= ch.Stats.Spd
						doneCH <- true //Done executing for this tick
						return
					}
					switch ch.ACList.Actions[0].Type[0] {
					case 0: //Mobility, raw X/Y change.
						//Check movement validity (Well, frontend will handle that, to ensure that movement is valid before the user is able to enter it.)
						ch.setXY(ch.ACList.Actions[0].EX, ch.ACList.Actions[0].EY, m)
					case 1: //Attack
						target := m.findChar(ch.ACList.Actions[0].EX, ch.ACList.Actions[0].EY)
						if target == nil {
							//Then the skill is not being targetted, perform it anyway and check the skills' attack radius to see if any would be hit... idk why I even have this check >>
							//Need to get the attack list
							//Attack = Name, ID, Radius, Span (horiz only, vert only (in respect to unit), area (like a bomb)), Damage modifier, Main stat (str/int etc), anything else I can think of.. magic will be different ofc
						} else {
							//Do stuff
							target.HP -= 1 //placeholder
						}
					}
					//EXECUTION DONE
					acCH <- ch.ACList.Actions[0]

					ch.ACList.Actions = ch.ACList.Actions[1:]
				}
			}(exList[x], acCH, doneCH)
		}
		//for each char id in exlist, have a for loop that collects the actions..?.. nope, it comes over one single channel.. tag the action itself? shouldn't need to
		for x := 1; x < len(exList); x++ {
			for {
				select {
				case ac := <- acCH:
					msg.Actions[ac.ID] = append(msg.Actions[ac.ID], ac)
				case <- doneCH:
					break
				}
			}
		}
		
	}
}
