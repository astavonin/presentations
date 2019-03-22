package main

import (
	"fmt"
	"math/rand"
	"os/exec"
	"time"
)

type response interface {}

type requestCommand int

const(
	genNum requestCommand = iota
	runCmd
)

type request struct {
	cmd requestCommand
	data interface {}
	out chan response
}

func doExec(req request) {
	time.Sleep(1*time.Second)

	cmd := req.data.(string)
	out, err := exec.Command(cmd).Output()
	if err != nil {
		req.out <- err.Error()
	} else {
		req.out <- string(out)
	}
}

func processRequest(req request)  {
	switch req.cmd {
	case genNum:
		req.out <- rand.Intn(100)
	case runCmd:
		go doExec(req)
	}
}

func eventLoop(done chan struct{}, inCh chan request) {
	for {
		select {
		case req := <-inCh:
			processRequest(req)
		case <-done:
			break
		}
	}
}

func main() {
	done := make(chan struct{})
	defer close(done)
	inCh := make(chan request, 10)
	defer close(inCh)

	go eventLoop(done, inCh)

	resp1 := make(chan response, 1)

	inCh <- request{
		runCmd,
		"ls",
		resp1,
	}

	resp2 := make(chan response, 1)
	inCh <- request{
		genNum,
		"resp1",
		resp2,
	}

	fmt.Println(<-resp2)
	fmt.Println(<-resp1)

	done <- struct{}{}
}