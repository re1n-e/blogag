package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/re1n-e/blogag/internal/config"
	"github.com/re1n-e/blogag/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}
	if user, err := s.db.GetUser(context.Background(), cmd.args[0]); err == nil {
		if err := s.cfg.SetUser(user.Name); err != nil {
			return err
		}
		fmt.Println("user has been logged in")
		return nil
	}
	return fmt.Errorf("user not found")
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}
	id, created_at, updated_at, name := uuid.New(), time.Now(), time.Now(), cmd.args[0]
	newUser := database.CreateUserParams{
		ID:        id,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
		Name:      name,
	}
	user, err := s.db.CreateUser(context.Background(), newUser)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	println("User has been created")
	fmt.Println(user)
	if err := s.cfg.SetUser(cmd.args[0]); err != nil {
		return err
	}
	fmt.Println("user has been set")
	return nil
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if fn, ok := c.cmds[cmd.name]; ok {
		return fn(s, cmd)
	}
	return fmt.Errorf("command with cmd name: %v is not a registered cmd", cmd.name)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("not enought arguments")
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
	if err := cmds.run(&st, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := config.Read(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(config.Current_user_name)
	fmt.Println(config.Db_url)
}
