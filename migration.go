package main

import (
	"fmt"
	"time"
)

var costRead = 1
var costWrite = 2
var windowLenght = 10 * time.Minute

// Does not require synchronization because it's called by a single thread
var m = make(map[string][]TimeOp)

type TimeOp struct {
	time time.Time
	cost int
}

func SetRead(key string) {
	updateAndEvaluate(key, costRead)
}

func SetWrite(key string) {
	updateAndEvaluate(key, costWrite)
}

func updateAndEvaluate(key string, cost int) {
	slice := m[key]
	now := time.Now()
	current := TimeOp{time: now, cost: cost}

	if len(slice) == 0 {
		sl := make([]TimeOp, 0)
		m[key] = append(sl, current)

		return
	} else {
		m[key] = append(deleteExpired(slice, now), current)
	}

}

func GetCost(key string, now time.Time) int {
	slice := m[key]
	var cost int
	if len(slice) == 0 {
		return 0
	}

	slice = deleteExpired(slice, now)
	fmt.Println(slice)
	for _, x := range slice {
		cost += x.cost
	}

	return cost
}

func deleteExpired(input []TimeOp, now time.Time) []TimeOp {
	tmp := input
	var p int

	windowStart := now.Add(-windowLenght)
	//windowStart := now

	index := 0
	for true {
		if len(tmp) == 0 {
			break
		}

		if len(tmp) == 1 {
			a := tmp[0]
			if a.time.Before(windowStart) || a.time.Equal(windowStart) {
				index++
				break
			}
		}

		p = len(tmp) / 2
		a := tmp[p]

		if a.time.Equal(windowStart) {
			index += p + 1
			break
		}

		if a.time.Before(windowStart) {
			index += p
			tmp = tmp[p:]
		} else {
			tmp = tmp[:p]
		}

	}

	return input[index:]
}

func main() {
	SetRead("a")
	SetWrite("a")

	fmt.Println(GetCost("a", time.Now()))

}
