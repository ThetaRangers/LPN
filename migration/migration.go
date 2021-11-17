package migration

import (
	"SDCC/utils"
	"time"
)

var costRead = 1
var costWrite = 2
var windowLength = 10 * time.Minute

var migrationThreashold = 1

const (
	ReadOperation  = 0
	WriteOperation = 1
	Master         = 0
	Replica        = 1
	External       = 2
)

// Does not require synchronization because it's called by a single thread
var master = make(map[string][]TimeOp)
var replica = make(map[string][]TimeOp)
var external = make(map[string][]TimeOp)

type TimeOp struct {
	time time.Time
	cost int
}

type KeyOp struct {
	Key  string
	Op   int
	Mode int
}

func Reset() {
	master = make(map[string][]TimeOp)
	replica = make(map[string][]TimeOp)
	external = make(map[string][]TimeOp)
}

func SetRead(key string, m map[string][]TimeOp) {
	updateAndEvaluate(key, costRead, m)
}

func SetWrite(key string, m map[string][]TimeOp) {
	updateAndEvaluate(key, costWrite, m)
}

func SetMigrated(key string) {
	master[key] = external[key]
	external[key] = make([]TimeOp, 0)
}

func SetExported(key string) {
	master[key] = make([]TimeOp, 0)
}

func updateAndEvaluate(key string, cost int, m map[string][]TimeOp) {
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

func GetCostMaster(key string, now time.Time) int {
	return GetCost(key, now, master)
}

func GetCostExternal(key string, now time.Time) int {
	return GetCost(key, now, external)
}

func GetCost(key string, now time.Time, m map[string][]TimeOp) int {
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

func EvaluateMigration() []string {
	migrationKeys := make([]string, 0)
	now := time.Now()

	for k, _ := range external {
		cost := GetCostExternal(k, now)

		// Find max cost
		if cost > utils.MigrationThreshold {
			migrationKeys = append(migrationKeys, k)
		}
	}

	return migrationKeys
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

func Evaluate() {

}

func ManagementThread(channel chan KeyOp, costR int, costW int, windowMinutes int) {
	var op KeyOp
	var m map[string][]TimeOp

	costRead = costR
	costWrite = costW
	windowLength = time.Minute * time.Duration(windowMinutes)

	for {
		op = <-channel

		switch op.Mode {
		case Master:
			m = master
			break
		case Replica:
			m = replica
			break
		case External:
			m = external
		}

		switch op.Op {
		case ReadOperation:
			SetRead(op.Key, m)
			break
		case WriteOperation:
			SetWrite(op.Key, m)
			break
		}
	}
}
