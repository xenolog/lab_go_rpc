package main

import (
	// "errors"
	// "fmt"
	"github.com/xenolog/lab_go_rpc/simple_task"
	kcp "github.com/xtaci/kcp-go"
	cli "gopkg.in/urfave/cli.v1"
	"gopkg.in/xenolog/go-tiny-logger.v1"
	"net"
	"net/rpc"
	"os"
	// "strings"
	// "time"
)

const (
	Version = "0.0.1"
)

var (
	Log *logger.Logger
	App *cli.App
	err error
)

func init() {
	// Setup logger
	Log = logger.New()

	// Configure CLI flags and commands
	App = cli.NewApp()
	App.Name = "RPC calls testing"
	App.Version = Version
	App.EnableBashCompletion = true
	// App.Usage = "Specify entry point of tree and got subtree for simple displaying"
	App.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable debug mode. Show more output",
		},
		cli.StringFlag{
			Name:  "url, u",
			Value: "udp://127.0.0.1:4001",
			Usage: "Specify URL for connection or listen",
		},
	}
	App.Commands = []cli.Command{{
		Name:   "server",
		Usage:  "run server",
		Action: runServer,
	}, {
		Name:   "client",
		Usage:  "connect to server, ask to run simple job",
		Action: runClient,
	}}
	App.Before = func(c *cli.Context) error {
		if c.GlobalBool("debug") {
			Log.SetMinimalFacility(logger.LOG_D)
		} else {
			Log.SetMinimalFacility(logger.LOG_I)
		}
		Log.Debug("EtcdTree started.")
		return nil
	}
	App.CommandNotFound = func(c *cli.Context, cmd string) {
		Log.Printf("Wrong command '%s'", cmd)
		os.Exit(1)
	}
}

func main() {
	App.Run(os.Args)
}

func runServer(c *cli.Context) error {
	task := new(simple_task.Tasks)
	rpcserver := rpc.NewServer()
	rpcserver.Register(task)
	//l, err := kcp.ListenWithOptions(":1234", block BlockCrypt, dataShards, parityShards int) (*Listener, error)
	l, err := kcp.Listen(":1234")
	if err != nil {
		Log.Error("listen error:", err)
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			Log.Error("connection error:", err)
		}
		Log.Info("RPC call from: %s will be served", conn.RemoteAddr())
		go rpcserver.ServeConn(conn)
	}
}

func runClient(c *cli.Context) {
	serverAddr := []string{c.GlobalString("url")}
	var (
		conn net.Conn
		err  error
	)
	// create connection
	if conn, err = kcp.Dial(serverAddr[0]); err != nil {
		Log.Fail("dialing:", err)
	}
	// create RPC client
	srv := rpc.NewClient(conn)
	args := &simple_task.Args{7, 8}
	var reply int
	// RPC call
	if err := srv.Call("Tasks.Task1", args, &reply); err != nil {
		Log.Fail("RPC error:", err)
	}
	Log.Info("RPC call: %d*%d=%d", args.A, args.B, reply)
	// RPC call
	if err := srv.Call("Tasks.Task2", args, &reply); err != nil {
		Log.Fail("RPC error:", err)
	}
	Log.Info("RPC call: %d+%d=%d", args.A, args.B, reply)
}
