# blog_agg
`Gator` RSS feed aggregator CLI in Go

- install Go from here -> https://go.dev/doc/install
- install Postgresql -> https://learn.microsoft.com/en-us/windows/wsl/tutorials/wsl-database#install-postgresql
- install Gator
```go install github.com/sudonetizen/blog_agg@latest```
this will install blog_agg binary at ~/go/bin/ path


- create Gator config file at `~/.gatorconfig.json`
- set up cofig file with follow structure:
```json
{
    "db_url": "postgres://username:password@host:port/database"
}
```
default port is 5432 and database name 

- after these settings blog_agg should be ready to run, just type blog_agg and press enter it should print that no arguments provided

### commands
```
# register user
blog_agg register name_of_user

# login already registered user 
blog_agg login name_of_user

# get list of all registered users and current shows logged user
blog_agg users 

# delete all registered users
blog_agg reset 

# adds a feed to Feeds table and creates a follow for this feed by current logged user at Feed_follows table
blog_agg addfeed "name_of_feed" "url_of_feed"

# follows given feed by its url for logged user if it already exists at Feeds table 
blog_agg follow "url_of_feed"

# shows followed feeds for logged user 
blog_agg following

# shows all feeds from Feeds table
blog_agg feeds

# unfollows a feed by given url from logged user at Feed_follows table (unfollow doesnt affect Feeds table)
blog_agg unfollow "url_of_feed"

# loops feeds from Feeds table based on last fetched time non-stop with given time interval in seconds and saves posts from feed to Posts table 
blog_agg agg x_number

# shows certain number of posts from Posts table
blog_agg browse x_number_of_posts
```
