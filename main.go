package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/re1n-e/blogag/internal/config"
	"github.com/re1n-e/blogag/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}
	config := config.NewConfig()
	if err := config.Read(); err != nil {
		fmt.Println(err)
	}
	st := state{}
	st.cfg = &config
	cmd := command{
		name: args[1],
		args: args[2:],
	}
	db, err := sql.Open("postgres", config.Db_url)
	if err != nil {
		fmt.Printf("Err connecting to db: %v\n", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	st.db = dbQueries
	cmds := commands{
		cmds: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handleReset)
	cmds.register("users", handleGetUsers)
	cmds.register("agg", handlerFetchFeed)
	cmds.register("addfeed", handleAddFeed)
	cmds.register("feeds", handleGetFeeds)
	cmds.register("follow", handleFeedFollow)
	cmds.register("following", handleGetFeedFollowsForUser)
	cmds.register("unfollow", handlerUnfollowFeed)
	if err := cmds.run(&st, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := config.Read(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
