package main

import (
    "fmt"
    "time"
    "errors"
    "context"
	"github.com/google/uuid"
    "github.com/sudonetizen/config"
    "github.com/sudonetizen/database"
)

type command struct {
    name string
    args []string 
}

type commands struct {
    mapCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
    function, ok := c.mapCommands[cmd.name]
    if !ok {return errors.New("command not found")}
    err := function(s, cmd)
    if err != nil {return err}
    return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
    c.mapCommands[name] = f 
}

//////////////////////////////////////////////////

func handlerGetUsers(s *state, cmd command) error {
    ctx := context.Background()
    if len(cmd.args) != 0 {return errors.New("provide only command")}
    
    users, err := s.db.GetUsers(ctx)    
    if err != nil {return err}

    for _, u := range users {
        if u == s.cfgP.Current_user_name {
            fmt.Printf("%v (current)\n", u)
        } else {fmt.Println(u)}
    }

    return nil
}


func handlerDeleteUsers(s *state, cmd command) error {
    ctx := context.Background() 
    if len(cmd.args) != 0 {return errors.New("provide only command")}

    err := s.db.DeleteUsers(ctx)
    if err != nil {
        fmt.Println("error with deleting of users")
        return err
    }
    fmt.Println("deleted users")
    return nil
}

func handlerLogin(s *state, cmd command) error {
    ctx := context.Background() 
    if len(cmd.args) != 1 {return errors.New("no username provided or more than one username")}
    _, err := s.db.GetUser(ctx, cmd.args[0])
    if err != nil {return errors.New("user doesnt exist")}

    config.SetUser(*(s.cfgP), cmd.args[0])
    fmt.Println("user has been set")
    return nil
}

func handlerRegister(s *state, cmd command) error {
    ctx := context.Background() 
    id := uuid.New()

    if len(cmd.args) != 1 {return errors.New("no username provided or more than one username")} 

    _, err := s.db.GetUser(ctx, cmd.args[0])
    if err == nil {return errors.New("user already exists")}

    _, err = s.db.CreateUser(ctx, database.CreateUserParams{id, time.Now().UTC(), time.Now().UTC(), cmd.args[0]})
    if err != nil {return errors.New("error with creation of user")}

    config.SetUser(*(s.cfgP), cmd.args[0])

    return nil 
}
