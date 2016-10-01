package simple_task

import (
	"math/rand"
	"time"
)

type Args struct {
	A, B int
}

type TaskResult struct {
	Result   int
	Duration time.Duration
}

type Tasks bool

func init() {
	rand.Seed(time.Now().Unix())
}

func sleep() {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
}

func (t *Tasks) Task1(args *Args, reply *TaskResult) error {
	startTime := time.Now()
	rv := args.A * args.B
	sleep()
	*reply = TaskResult{Result: rv, Duration: time.Since(startTime)}
	return nil
}
func (t *Tasks) Task2(args *Args, reply *TaskResult) error {
	startTime := time.Now()
	rv := args.A + args.B
	sleep()
	*reply = TaskResult{Result: rv, Duration: time.Since(startTime)}
	return nil
}
