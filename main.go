package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var stdin = bufio.NewReader(os.Stdin)

func less(a, b string) bool {
	for {
		fmt.Printf("\na) %s\nb) %s\nWhich is higher priority? [a, b] > ", a, b)
		line, err := stdin.ReadString('\n')
		if err != nil {
			panic(err)
		}

		choice := strings.ToLower(strings.TrimSpace(line))
		// these return the reverse of what "less" should mean, because we want
		// to store highest->lowest priority
		if choice == "a" {
			return true
		} else if choice == "b" {
			return false
		}
		fmt.Println(`Invalid choice, must be "a" or "b", try again`)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s filename\n", os.Args[0])
		return
	}

	var in []string
	{
		file, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}

		fileBR := bufio.NewReader(file)
		for {
			line, err := fileBR.ReadString('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}

			if el := strings.TrimSpace(line); el != "" {
				in = append(in, el)
			}
		}
		file.Close()
	}

	if len(in) == 0 {
		panic("no input")
	}

	tree := new(Tree)
	for _, el := range in {
		tree.Insert(el, less)
	}

	fmt.Println("You're done sorting! Hit enter for the output (highest priority to lowest)")
	if _, err := stdin.ReadString('\n'); err != nil {
		panic(err)
	}

	tree.Traverse(tree.Root, func(n *Node) {
		fmt.Printf("%s\n", n.Value)
	})
}
