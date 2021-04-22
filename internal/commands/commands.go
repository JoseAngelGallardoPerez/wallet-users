package commands

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
)

type Handler interface {
	Handle(args url.Values)
	Name() string
	Description() string
}

var commands = map[string]Handler{
	"help": helpCommand,
}

func process(cmd string) {
	u, err := url.Parse(cmd)
	if err != nil {
		log.Fatal("unable to parse command: " + cmd + "; command must be valid url format e.g. cmd?param=value")
	}
	if command, exist := commands[u.Path]; exist {
		if command.Name() == "help" {
			for _, v := range commands {
				fmt.Print(v.Description())
			}
			return
		}
		q, _ := url.ParseQuery(u.RawQuery)
		command.Handle(q)
		return
	}
	log.Fatal("command not found: " + u.Path)
}

func AddCommand(cmd Handler) {
	_, ok := commands[cmd.Name()]
	if !ok {
		commands[cmd.Name()] = cmd
	}
}

func Run() {
	cmd := flag.String("cmd", "", "-cmd \"commandName?param1=value1&param2=value2...paramN=valueN\"")
	flag.Parse()
	if *cmd != "" {
		process(*cmd)
		os.Exit(0)
	}
}
