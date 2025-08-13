package main

import (
    "fmt"
    "errors"
    "github.com/sudonetizen/config"
)

type state struct {
   cfgP *config.Config  
}

type command struct {
    name string
    args []string 
}

func handlerLogin(s *state, cmd command) error {
    if len(cmd.args) == 0 {return errors.New("no username provided")}
    //cfg := (*(*s).pointer)
    config.SetUser(*(s.cfgP), cmd.args[0])
    fmt.Println("user has been set")
    return nil
}

type commands struct {
    mapCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
    function := c.mapCommands[cmd.name]
    err := function(s, cmd)
    if err != nil {return err}
    return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
    c.mapCommands[name] = f 
}
