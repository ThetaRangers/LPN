package migration

import (
	"log"
	"time"
)

var costRead = 1
var costWrite = 2
var windowLength = 10 * time.Minute

const (
	ReadOperation  = 0
	WriteOperation = 1
)

// Does not require synchronization because it's called by a single thread
var m = make(map[string][]TimeOp)

type TimeOp struct {
	time time.Time
	cost int
}

type KeyOp struct {
	Key string
	Op  int
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
	for _, x := range slice {
		cost += x.cost
	}

	return cost
}

func deleteExpired(input []TimeOp, now time.Time) []TimeOp {
	tmp := input
	var p int

	windowStart := now.Add(-windowLength)
	//windowStart := now

	index := 0
	for {
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

func ManagementThread(channel chan KeyOp, costR int, costW int, windowMinutes int) {
	var op KeyOp

	costRead = costR
	costWrite = costW
	windowLength = time.Minute * time.Duration(windowMinutes)

	for {
		op = <-channel

		switch op.Op {
		case ReadOperation:
			SetRead(op.Key)
			break
		case WriteOperation:
			SetWrite(op.Key)
			break
		}

		for k, _ := range m {
			log.Println("Cost for: ", k, GetCost(k, time.Now()))
		}

	}
}
