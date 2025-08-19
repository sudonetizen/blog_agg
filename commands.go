package main

import (
    "io"
    "fmt"
    "html"
    "time"
    "errors"
    "context"
    "net/http"
    "encoding/xml"
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

type RSSFeed struct {
    Channel struct {
        Title       string    `xml:"title"`
        Link        string    `xml:"link"`
        Description string    `xml:"description"`
        Item        []RSSItem `xml:"item"`
    } `xml:"channel"` 
}

type RSSItem struct {
    Title       string `xml:"title"`
    Link        string `xml:"link"`
    Description string `xml:"description"`
    PubDate     string `xml:"pubDate"`
}

func fetchFeed(s *state, cmd command) error {
    if len(cmd.args) != 1 {return errors.New("no url provided or more than one url")}

    req, err := http.NewRequestWithContext(context.Background(), "GET", cmd.args[0], nil)
    if err != nil {return err} 

    req.Header.Set("User-Agent", "gator")
    
    client := http.Client{}
    res, err := client.Do(req)
    if err != nil {return err}
    defer res.Body.Close()

    data, err := io.ReadAll(res.Body)
    if err != nil {return err} 
    
    rss := RSSFeed{}
    err = xml.Unmarshal(data, &rss)
    if err != nil {return fmt.Errorf("error with unmarshal: %v", err)}

    fmt.Printf("Title: %v\n", html.UnescapeString(rss.Channel.Title))
    fmt.Printf("Link: %v\n", html.UnescapeString(rss.Channel.Link))
    fmt.Printf("Description: %v\n\n", html.UnescapeString(rss.Channel.Description))

    for _, i := range rss.Channel.Item {
        fmt.Printf("Title: %v\n", html.UnescapeString(i.Title))
        fmt.Printf("Published Date: %v\n", html.UnescapeString(i.PubDate))
        fmt.Printf("Link: %v\n", html.UnescapeString(i.Link))
        fmt.Printf("Description: %v\n\n", html.UnescapeString(i.Description))
    }

    return nil
}

func handlerUserFeeds(s *state, cmd command) error {
    if len(cmd.args) != 0 {return errors.New("only command")}
    
    feeds, err := s.db.GetFeedFollowsForUser(context.Background(), s.cfgP.Current_user_name)
    if err != nil {return fmt.Errorf("error with feeds: %v\n", err)}
    
    for _, f := range feeds {
        fmt.Printf("Feed: << %v >> followed by user %v\n",  f.Name_2.String, f.Name.String)
    }

    return nil
}


func handlerFeedFollow(s *state, cmd command) error {
    if len(cmd.args) != 1 {return errors.New("no or more url provided")}
    
    userID, err := s.db.GetUserId(context.Background(), s.cfgP.Current_user_name)
    if err != nil {return fmt.Errorf("error with getting user id: %v\n", err)}
    feedID, err := s.db.GetFeedId(context.Background(), cmd.args[0])
    if err != nil {return fmt.Errorf("error with getting feed id: %v\n", err)}

    feedFollow, err := s.db.CreateFeedFollow(
        context.Background(),
        database.CreateFeedFollowParams{
            uuid.New(), 
            time.Now().UTC(), 
            time.Now().UTC(), 
            userID,
            feedID, 
        },
    )
    if err != nil {return fmt.Errorf("error with feed follow: %v\n", err)}
    
    fmt.Printf("ID: %v\nCreated: %v\nUpdated: %v\nUser ID: %v\nFeed ID: %v\nUser Name: %v\nFeed Name: %v\n",
               feedFollow.ID,
               feedFollow.CreatedAt,
               feedFollow.UpdatedAt,
               feedFollow.UserID,
               feedFollow.FeedID, 
               feedFollow.Name.String,
               feedFollow.Name_2.String,
    )

    return nil
}

func handlerFeeds(s *state, cmd command) error {
    if len(cmd.args) != 0 {return errors.New("only command")}

    feeds, err := s.db.GetFeeds(context.Background())
    if err != nil {return err}    

    for _, i := range feeds {
        fmt.Printf("Name: %v\nUrl: %v\nAdded by User: %v\n\n", i.Name, i.Url, i.Name_2.String) 
    } 
    return nil 
}

func handlerAddFeed(s *state, cmd command) error {
    ctx := context.Background()
    if len(cmd.args) != 2 {return errors.New("less or more args provided")}
    userID, err := s.db.GetUserId(ctx, s.cfgP.Current_user_name)
    if err != nil {return fmt.Errorf("error with getting user id: %v\n", err)}

    _, err = s.db.CreateFeed(
        ctx, 
        database.CreateFeedParams{
            uuid.New(), 
            time.Now().UTC(), 
            time.Now().UTC(), 
            cmd.args[0], 
            cmd.args[1],  
            userID,
        },
    )
    if err != nil {return err} 

    cmd.args = []string{cmd.args[1]}
    err = handlerFeedFollow(s, cmd)
    if err != nil {return err} 
    
    return nil
}


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
