package main

import (
    "os"
    "fmt"
    "github.com/sudonetizen/config"
)


func main() {
    cfg := config.Read()
    stt := state{cfgP: &cfg}

    cmm := commands{mapCommands: make(map[string]func(*state, command) error)} 
    cmm.register("login", handlerLogin)

    if len(os.Args) < 2 {
        fmt.Println("no arguments were provided")
        os.Exit(1)  
    }

    cmnd := command{name: os.Args[1], args: os.Args[2:]} 
    err := cmm.run(&stt, cmnd)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    //fmt.Println(os.Args[1], os.Args[2:])
    //fmt.Println(cmm.mapCommands)
    //fmt.Println(*stt.pointer)
}
