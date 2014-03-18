package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"code.google.com/p/go-uuid/uuid"
	
	"net/http"
	"time"

	"log"

//	"bytes"
)


//Kinda like chrono trigger atb system, BUT with a tactical field
//Note: Movement would take up CT based on equips, weight etc, same for mobility skills
//Every tick, CT = CT + Spd (this will change), if CT >= 100, CT = 100 and char takes a turn
//If char queued turn CT total > his Spd cost, it becomes a timed queue, meaning it'll span over multiple turns, where turn count = X/Spd
//When a char gets hit, they gain CT based on damage taken
//Skills are gonna be basic, for example: Dash, Jump for mobility, Power fist (high damage, knockback), Lighter attacks give enemy less CT//Drain CT, but keep them in place/apply hitstun (enemy unit cannot act).
//Can only chain Mobility -> Light -> Power | Mobility -> Mobility
//Mobility to mobility has a large increase in CT cost, as it circumvents guard (another way to break guard)
//Light can chain to anything, Power uses more CT to chain to mobility, usually enders but chaining with mobility is essentially a roman.
//Also, hit gauge, the more hits you do the more --Fever(?) you gain, basically like akihiko gayness lol
//Can guard, power attacks break guard, light attacks do nothing, maybe autoguard? or guard is a skill which you can queue up, also counter mechanic.
//Burst mechanic, once per game, conditional
//Also, counter hits doing more damage and giving you CT (allowing for resets) <--- IDK, MAYBE, DUNNO HOW I'D IMPLEMENT THIS, maybe add a conditional counter mechanic, and if you read that then you're golden

//Equips:
//Magic weilders wear 5 point gloves, gloves determine which 5 elements you can weild, can only combine two (????), 'cause two gloves, but can put more energy from one hand to the next, and spell affinity etc.
//Can program spells

//Action Queue Create
func (ch *Char) ACQC () {
	ch.ACList.TICKCT = 0
	
	//Create move queue based on user input
	return
}
//When you perform an action, your speed is very important
//Your TICKCT = Speed - Action CT cost, if TICKCT is negative, let a tick pass and add Spd to it, if TICKCT is then positive, execute
//Action result, 0 = let a tick pass because not enough Speed to do this in this tick, 1 = execute
// func (ch *Char) Execute (ACR int) {
// 	//Blah blah execute
// 	if ch.ACList.Actions[0].CT + ch.ACList.TICKCT > ch.Stats.Spd {
// 		ch.ACList.Actions[0].CT = ch.Stats.Spd - ch.ACList.TICKCT // Difference between TICKCT and speed become the move CT (so you can 'charge up' moves, so to speak)
// 		ch.ACList.TICKCT = 0
// 		ACR = 0; return
// 	}else {
// 		ch.ACList.TICKCT += ch.ACList.Actions[0].CT
// 		//
// 		//Execute, whatever that means
// 		//
// 		ch.ACList.Actions = ch.ACList.Actions[1:] //Dequeue
// 		ACR = 1; return 
// 	}
// 	//Pop, then if len(ac.Actions) == 0, TICKCT = 0
	
// }
// func (ac *Actions) Pop (CT int) {
// 	CT = ac.Actions[0].CT
	
// }
var gloMatch []*Match
func remoteHandler(res http.ResponseWriter, req *http.Request) {
	var err error;

	ws, _ := websocket.Upgrade(res, req, nil, 1024, 1024)
	log.Printf("got websocket conn from %v\n", ws.RemoteAddr())

	match := new(Match)
	match.Socket = ws
	match.ID = uuid.New()
	gloMatch = append(gloMatch, match)
	
	
	//tickChan := time.NewTicker(time.Second).C
	var tickah *time.Ticker
	//var execList []*Char
	turn := false
MainLoop:
	for {
		//NOTE, DO A VALIDITY CHECK EVERYTIME YOU READ JSON
		if err = ws.ReadJSON(match); err != nil {
			panic("LOL WAT")
		}
		tickah = time.NewTicker(time.Second)
		select {
		case <- tickah.C:
			//First sweep, just add speed to CT's
			for x, y := range match.Teams {
				for x2, _ := range y {
					//Check the actions list, then for every tick, execute (note, add execute code)
					match.Teams[x][x2].CT += match.Teams[x][x2].Stats.Spd
					if !turn && y[x2].CT >= 100 && len(y[x2].ACList.Actions) == 0{
						//After adding, check the CT, if CT >= 100 then that player has a turn
						tickah.Stop(); turn = true;
					}
				}
			}
			if turn {
				//Someones turn, pause everything and whatnot
				//Send message to client saying whose turn it is
				//Messages to client will be so
				//{
				//  "Action": //string: turn, execute
				//  If action == turn:
				//  "Players": [ Char ID/int ]
				//  If action == execute:
				//  "Actions": [ { Action } ] // SX/SY will tell you which char it is
				//}
				msg := new (Message)
				var turnList []int
				for _, teams := range match.Teams {
					for _, char := range teams {
						if char.CT >= 100 {
							char.CT = 100
							turnList = append (turnList, char.ID)
						}
					}
				}
				msg.Players = turnList
				//Send message over match.Socket
				//
				//
				//
				//
				continue MainLoop
			}

			match.Execute()
		}
		//execList = []*Char{}
		tickah = time.NewTicker(time.Second)
	}
}
	// buf := new(bytes.Buffer)
	// buf.ReadFrom(req.Body)
	// //s := buf.String()
	// log.Println(buf.String())

//When executing, TICKCT += action.CT && if TICKCT > char speed, TICKCT = TICKCT - char speed and wait until next tick.
//confList = [ (Conflicts) ]
//Conflict = { [ Action ] }
func (m *Match) Execute() {
	var execList []*Char
	for _, y := range m.Teams {
		for _, char := range y {
			if len(char.ACList.Actions) != 0 {
				execList = append(execList, &char)
			}
		}		
	}
	if len(execList) > 0 {
		confList := confCheck(execList)
		//Find if any conflicts will be had first.
		if len(confList) > 0 {
			//If conflict, resolve all the conflicts in confList (resolution will be simple)
			//Will actually change the action list
		} 
		//Execute here
		msg := new (Message)
		msg.Action = "execute"
		acCH := make(chan Action)
		doneCH := make(chan bool)
		//var acCH chan Action
		//var doneCH chan bool
		for _, char := range execList {
			go func (ch *Char, acCH chan Action, doneCH chan bool) {
				//var x int
				for {
					ch.ACList.TICKCT += ch.ACList.Actions[0].CT
					if ch.ACList.TICKCT > ch.Stats.Spd {
						ch.ACList.TICKCT -= ch.Stats.Spd
						doneCH <- true
						return
					}
					//Actual execution now.
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
					acCH <- ch.ACList.Actions[0]
					//Dump ch.ACList.Actions[0]
					ch.ACList.Actions = ch.ACList.Actions[1:]
				}

			}(char, acCH, doneCH)
		}
		for x := 1; x < len(execList); x++ {
		chanLoop:
			for {
				select {
				case ac := <- acCH:
					msg.Actions = append(msg.Actions, ac)
				case <- doneCH:
					break chanLoop
				}
			}
		}
		//Send message
		m.Send(msg)
	}else {
		return
	}
}
func confCheck (eL []*Char) (confL []Conflict) {
	for _, char := range eL {
		//CHECK FOR CONFLICTS
		log.Println(char)
	}

	return
}
func main () {
	r := mux.NewRouter()
	r.HandleFunc("/ws", remoteHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./html/")))
	http.ListenAndServe(":8080", r)
	//http.ListenAndServe(":8080", http.FileServer(http.Dir("/home/akisa/mygo/gogame/html")))
}
