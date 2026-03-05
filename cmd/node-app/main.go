package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/smirnovalles/yt/internal/node"
)

func main() {

	port := flag.Int("p", 0, "port to listen on")
	nodeAddress := flag.String("n", "", "address connection node")
	id := flag.String("id", "noname", "address connection node")
	message := flag.String("m", "", "send messege")

	flag.Parse()

	if *port <= 0 {
		fmt.Println("no port set")
		os.Exit(1)
	}

	ch := make(chan os.Signal, 2)

	go func() {
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	}()

	nn := node.New(*id, "localhost", *port)

	err := nn.Start()
	if err != nil {
		fmt.Printf("start server error: %s", err.Error())
	}

	fmt.Printf("start server...%s\n", nn)

	if *nodeAddress != "" {
		//nn.Connect(*nodeAddress)
	}

	if *message != "" && *nodeAddress != "" {
		nn.Send(*nodeAddress, []byte(*message))
	}

	<-ch

	fmt.Println("Terminated...")
}
