package main

import (
	"aoc2021/common/file"
	"aoc2021/day24/alu"
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func main() {

	numInput := make(chan int64, 1)
	inputRequest := make(chan bool, 1)

	alu := alu.NewALU(inputRequest, numInput)

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	cleanup := func(code int) {
		close(signals)
		close(inputRequest)
		close(numInput)
		os.Exit(code)
	}

	go func() {
		for {
			select {
			case <-signals:
				cleanup(1)
			case <-inputRequest:
				delim := '\n'
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("input> ")
				text, err := reader.ReadString(byte(delim))
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "%v", err)
					cleanup(-3)
				}

				tokens := strings.Split(text, string(delim))
				if len(tokens) > 0 {
					text = tokens[0]
				}

				if text == "exit" || text == "quit" {
					fmt.Println("exiting, goodbye!")
					cleanup(0)
				} else {
					val, err := strconv.ParseInt(text, 10, 64)
					if err != nil {
						_, _ = fmt.Fprintf(os.Stderr, "error parsing value: %s", err)
					}
					numInput <- val
				}
			}

		}
	}()

	lines, err := file.GetLines(os.Args[1])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Can not read file %s", os.Args[1])
		cleanup(-1)
	}

	for num, line := range lines {
		if len(line) > 0 {
			fmt.Println("processing instruction (", num, "): ", line)
			err = alu.Process(line)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%v", err)
				cleanup(-2)
			}
		}
	}

	fmt.Println(alu.String())
	cleanup(0)

}
