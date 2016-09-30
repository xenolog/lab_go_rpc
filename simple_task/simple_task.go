package simple_task

import (
	"math/rand"
	"time"
)

type Args struct {
	A, B int
}

type Tasks int

func sleep() {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
}

func (t *Tasks) Task1(args *Args, reply *int) error {
	*reply = args.A * args.B
	sleep()
	return nil
}
func (t *Tasks) Task2(args *Args, reply *int) error {
	*reply = args.A + args.B
	sleep()
	return nil
}
