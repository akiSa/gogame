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
	if err = ws.ReadJSON(match); err != nil {
		panic("LOL WAT")
	}
	//Initialize the match, after that all messages will be in Message struct.
	
	match.Socket = ws
	match.ID = uuid.New()
	gloMatch = append(gloMatch, match)
	
	//msg := new (Message)
	var msg *Message
	var tickah *time.Ticker

	tickah = time.NewTicker(time.Second)
	for {
		<- tickah.C
		//After tick, do tick action.
		if turn, list := match.Tick(); turn {
			//Poll for new Actions, add to action list
			//send list
			//tickah.Stop() I'm dumb, channels are unbuffered.
			msg = new (Message)
			msg.Action = "turn"
			msg.Players = list
			match.Socket.WriteJSON(msg) //Send list of players that are to take an action.
			if err = ws.ReadJSON(msg); err != nil {
				panic(err)
			}
			//Read the action list, it will be in Message format, and will consist of
			// Action: "turn"
			// Actions: [ Actions ] ... maybe it should be map[int][]Action.. probably
			//Which will look like... { 51: []Action, 2: []Action }, unfortunately it'll go through it in numerical order, but i'll make it execute in non-ordered fashion, or something, shrug
			//msg.Actions = map[int][]Action
			for x, y := range msg.Actions {
				//x = char ID, y = []Action
				//Search for char with ID
				cl := match.Find(x) //link to char
				if cl != nil {
					cl.ACList.Actions = append( cl.ACList.Actions, y...)
				}else {
					//wat???
					panic("wat???")
				}
			}
		}
		match.Execute()
	}
}


// Read:
// 	if err = ws.ReadJSON(msg); err != nil {
// 		panic("LOL WAT")
// 	}
// MainLoop:
// 	for {
// 		turn = false
// 		//NOTE, DO A VALIDITY CHECK EVERYTIME YOU READ JSON
// 		if err = ws.ReadJSON(msg); err != nil {
// 			panic("LOL WAT")
// 		}
// 		//Do a parse of the message
// 		tickah = time.NewTicker(time.Second)
// 		select {
// 		case <- tickah.C:
// 			//First sweep, just add speed to CT's
// 			for x, y := range match.Teams {
// 				for x2, _ := range y {
// 					//Check the actions list, then for every tick, execute (note, add execute code)
// 					match.Teams[x][x2].CT += match.Teams[x][x2].Stats.Spd
// 					if !turn && y[x2].CT >= 100 && len(y[x2].ACList.Actions) == 0{
// 						//After adding, check the CT, if CT >= 100 then that player has a turn
// 						tickah.Stop(); turn = true;
// 					}
// 				}
// 			}
// 			if turn {
// 				//Someones turn, pause everything and whatnot
// 				//Send message to client saying whose turn it is
// 				//Messages to client will be so
// 				//{
// 				//  "Action": //string: turn, execute
// 				//  If action == turn:
// 				//  "Players": [ Char ID/int ]
// 				//  If action == execute:
// 				//  "Actions": [ { Action } ] // SX/SY will tell you which char it is
// 				//}
// 				msg := new (Message)
// 				var turnList []int
// 				for _, teams := range match.Teams {
// 					for _, char := range teams {
// 						if char.CT >= 100 {
// 							char.CT = 100
// 							turnList = append (turnList, char.ID)
// 						}
// 					}
// 				}
// 				msg.Players = turnList
// 				//Send message over match.Socket
// 				//Then wait for
// 				//
// 				//
// 				//
// 				continue MainLoop
// 			}

// 			match.Execute()
// 		}
// 		//execList = []*Char{}
// 		tickah = time.NewTicker(time.Second)
// 	}
// }
	// buf := new(bytes.Buffer)
	// buf.ReadFrom(req.Body)
	// //s := buf.String()
	// log.Println(buf.String())



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
