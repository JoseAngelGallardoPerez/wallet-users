package commands

import (
	"fmt"
	"net/url"
)

type command struct {
	name        string
	usage       string
	description string
}

var helpCommand = command{
	name:        "help",
	usage:       "help",
	description: "Display available commands",
}

func (command) Handle(url.Values) {}

func (c command) Name() string {
	return c.name
}

func (c command) Description() string {
	return fmt.Sprintf("%s:\nUsage:\t%s\nDescription:\t%s\n\n", c.name, c.usage, c.description)
}
