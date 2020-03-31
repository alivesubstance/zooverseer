package zk

import (
	"github.com/alivesubstance/zooverseer/core"
	"log"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	c := make(chan Node)
	connInfo := core.JsonConnInfo{
		Host: "localhost",
		Port: 2181,
	}
	go Get("/", &connInfo, c)
	node := <-c
	log.Println(node.name)
	for _, child := range node.children {
		log.Println(child.name)
	}

}

func TestNil(t *testing.T) {
	var c []string
	if len(c) == 0 {
		log.Print("empty")
	} else {
		log.Print("not empty")
	}
}

func TestSleep(t *testing.T) {
	log.Print("start")

	for i := 0; i < 10; i++ {
		log.Print(i)
		time.Sleep(300 * time.Millisecond)

		if i == 2 {
			go sleep()
		}
	}

	log.Print("end")
}

func sleep() {
	log.Print("before sleep")
	time.Sleep(1 * time.Second)
	log.Print("after sleep")
}
