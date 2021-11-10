package utils

const (
	JoinOp     = 0
	LeaveOp    = 1
	ContainsOp = 2
	LenOp      = 3
)

type localOperation struct {
	operation rune
	ip        string
	ch        chan interface{}
}

type ClusterRoutine struct {
	cluster ClusterSet
	ch      chan localOperation
}

func (c ClusterRoutine) run() {
	var op localOperation
	for {
		op = <-c.ch
		switch op.operation {
		case JoinOp:
			c.cluster.Add(op.ip)
		case LeaveOp:
			c.cluster.Remove(op.ip)
		case ContainsOp:
			op.ch <- c.cluster.Contains(op.ip)
		case LenOp:
			op.ch <- c.cluster.Len()
		}
	}
}

func (c ClusterRoutine) Join(ip string) {
	c.ch <- localOperation{operation: JoinOp, ip: ip}
}

func (c ClusterRoutine) Leave(ip string) {
	c.ch <- localOperation{operation: LeaveOp, ip: ip}
}

func (c ClusterRoutine) Contains(ip string) bool {
	channel := make(chan interface{})
	c.ch <- localOperation{operation: ContainsOp, ip: ip, ch: channel}
	tmp := <-channel
	return tmp.(bool)
}

func (c ClusterRoutine) Len() int {
	channel := make(chan interface{})
	c.ch <- localOperation{operation: LenOp, ch: channel}
	tmp := <-channel
	return tmp.(int)
}

func NewClusterRoutine() ClusterRoutine {
	routine := ClusterRoutine{cluster: NewClusterSet(), ch: make(chan localOperation)}
	go routine.run()
	return routine
}
