package main

import (
    "os"
    "fmt"
    "log"
    "database/sql"
    _ "github.com/lib/pq"   
    "github.com/sudonetizen/config"
    "github.com/sudonetizen/database"
)

type state struct {
   cfgP *config.Config  
   db *database.Queries
}

func main() {
    cfg := config.Read()
    stt := state{cfgP: &cfg}

    db, err := sql.Open("postgres", stt.cfgP.Db_url)
    if err != nil {
      log.Fatalf("error: %v", err)
    }
    dbQueries := database.New(db)

    stt.db = dbQueries

    cmm := commands{mapCommands: make(map[string]func(*state, command) error)} 
    cmm.register("login", handlerLogin)
    cmm.register("register", handlerRegister)
    cmm.register("reset", handlerDeleteUsers)
    cmm.register("users", handlerGetUsers)
    cmm.register("agg", fetchFeed)
    cmm.register("addfeed", handlerAddFeed)
    cmm.register("feeds", handlerFeeds)
    cmm.register("follow", handlerFeedFollow)
    cmm.register("following", handlerUserFeeds)

    if len(os.Args) < 2 {
        fmt.Println("no arguments were provided")
        os.Exit(1)  
    }
    
    cmnd := command{name: os.Args[1], args: os.Args[2:]} 
    err = cmm.run(&stt, cmnd)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

}
