package cmd

import "math/rand"

func RollDice(num, sides uint) int {
	var total int
	var i uint
	for i = 0; i < num; i++ {
		total += rand.Intn(int(sides)) + 1
	}
	return total
}
