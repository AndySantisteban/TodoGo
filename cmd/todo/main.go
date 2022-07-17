package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/AndySantisteban/todogo/cmd/service"
)

const (
	todoFile = ".todo.json"
)

func main() {

	add := flag.Bool("add", false, "add a task")
	clear := flag.Bool("clear", false, "clear all tasks")
	get := flag.Bool("get", false, "get all tasks")
	complete := flag.Int("complete", 0, "complete a task")
	remove := flag.Int("remove", 0, "remove a task")
	print := flag.Bool("print", false, "print all tasks")
	flag.Parse()

	todos := &service.Todos{}

	if err := todos.Load(todoFile); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	switch {
	case *add:
		task, err := getInput(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		todos.Add(task)
		err = todos.Get(todoFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

	case *remove > 0:
		err := todos.Remove(*remove)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		err = todos.Get(todoFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

	case *clear:
		err := todos.Clear()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = todos.Get(todoFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case *get:
		err := todos.Get(todoFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case *complete > 0:
		err := todos.Complete(*complete)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		err = todos.Get(todoFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	case *print:
		todos.Print()

	default:
		fmt.Println("no command")
		os.Exit(0)
	}
}

func getInput(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	scanner := bufio.NewScanner(r)

	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}
	text := scanner.Text()
	if len(text) == 0 {
		return "", errors.New("Ey you need to enter a task")
	}
	return text, nil

}
