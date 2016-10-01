package main

import (
	"github.com/xenolog/lab_go_rpc/simple_task"
	kcp "github.com/xtaci/kcp-go"
	smux "github.com/xtaci/smux"
	cli "gopkg.in/urfave/cli.v1"
	"gopkg.in/xenolog/go-tiny-logger.v1"
	"net"
	"net/rpc"
	"os"
	"time"
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
		Log.Debug("Started.")
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

func pingServer(conn net.Conn) error {
	var (
		sess     *smux.Session
		err      error
		upStream *smux.Stream
		dnStream *smux.Stream
	)
	defer func() {
		conn.Close()
	}()
	// setup session (one per client)
	if sess, err = smux.Server(conn, nil); err != nil {
		Log.Error("SMUX server bind error:", err)
		return err
	}
	// waiting for ping
	if upStream, err = sess.AcceptStream(); err != nil {
		Log.Error("SMUX accept error:", err)
		return err
	}
	buf := make([]byte, 4)
	upStream.Read(buf)
	Log.Info("received '%s'", buf)
	upStream.Write([]byte("pong"))
	upStream.Read(buf)
	Log.Info("received '%s'", buf)
	if dnStream, err = sess.OpenStream(); err != nil {
		Log.Error("OpenStream error:", err)
		return err
	}
	dnStream.Write([]byte("Ping"))
	dnStream.Read(buf)
	Log.Info("received '%s'", buf)

	return nil
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
		Log.Info("Incoming connection from: %s will be served", conn.RemoteAddr())
		go pingServer(conn)
	}
	return nil
}

func runClient(c *cli.Context) error {
	serverAddr := []string{c.GlobalString("url")}
	var (
		conn     net.Conn
		err      error
		sess     *smux.Session
		upStream *smux.Stream
		dnStream *smux.Stream
	)
	// create connection
	if conn, err = kcp.Dial(serverAddr[0]); err != nil {
		Log.Fail("dialing:", err)
		return err
	}
	defer func() {
		conn.Close()
	}()
	// start multiplexer
	if sess, err = smux.Client(conn, nil); err != nil {
		Log.Fail("Multiplexer:", err)
		return err
	}
	// start session
	if upStream, err = sess.OpenStream(); err != nil {
		Log.Error("OpenStream error:", err)
		return err
	}
	// dialog
	buf := make([]byte, 4)
	upStream.Write([]byte("ping"))
	upStream.Read(buf)
	Log.Info("received '%s'", buf)
	upStream.Write([]byte("wait"))
	// waiting reversed connection
	if dnStream, err = sess.AcceptStream(); err != nil {
		Log.Error("SMUX accept error:", err)
		return err
	}
	dnStream.Read(buf)
	Log.Info("received '%s'", buf)
	dnStream.Write([]byte("Pong"))
	// End of communication
	upStream.Close()
	time.Sleep(time.Second)
	dnStream.Close()
	return nil
}
