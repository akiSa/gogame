package main

import (

)

func (m *Match) findChar (X, Y int) (target *Char) {
	return m.Map[X][Y].Unit
	// for _, team := range m.Teams {
	// 	for _, char := range team {
	// 		if char.X == X && char.Y == Y {
	// 			return &char
	// 		}
	// 	}
	// }
	// return nil
}
