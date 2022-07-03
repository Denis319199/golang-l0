package main

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"os"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("command arguments: <url> <cluster_id> <dir>")
		return
	}

	fmt.Println("url:", os.Args[1])
	fmt.Println("cluster-id:", os.Args[2])
	fmt.Println("dir:", os.Args[3])

	var path string
	runeDirPath := []rune(os.Args[3])
	if runeDirPath[len(runeDirPath)-1] == os.PathSeparator {
		path = os.Args[3]
	} else {
		path = os.Args[3] + string(os.PathSeparator)
	}

	nats := setupNatsStreaming(os.Args[1], os.Args[2])
	defer func(nats stan.Conn) {
		_ = nats.Close()
	}(nats)

	items, err := os.ReadDir(os.Args[3])
	checkError(err)
	for _, item := range items {
		if !item.IsDir() {
			fileContent, err := os.ReadFile(path + item.Name())
			checkError(err)

			err = nats.Publish("myevent", fileContent)
			checkError(err)
		}
	}

	fmt.Println("All files has been transmitted successfully")
}
