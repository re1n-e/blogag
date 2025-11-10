package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/re1n-e/blogag/internal/database"
)

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

func handleReset(s *state, _ command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to delete users table data: %v", err)
	}
	return err
}

func handleGetUsers(s *state, _ command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to retrive users: %v", err)
	}
	for _, user := range users {
		if user == s.cfg.Current_user_name {
			fmt.Printf("%s (current)\n", user)
		} else {
			fmt.Println(user)
		}
	}
	return nil
}

func handlerFetchFeed(s *state, cmd command) error {
	time_between_reqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("failed to parse arg to time: %v", err)
	}
	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			return err
		}
	}
}

func handleAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("usage blogag <title> <feed-url>")
	}
	user, err := s.db.GetUser(context.Background(), s.cfg.Current_user_name)
	if err != nil {
		return fmt.Errorf("failed to fetch given user: %v", err)
	}
	id, created_at, updated_at, title, url, user_id := uuid.New(), time.Now(), time.Now(), cmd.args[0], cmd.args[1], user.ID
	NewFeed := database.AddFeedParams{
		ID:        id,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
		Title:     title,
		Url:       url,
		UserID:    user_id,
	}
	feed, err := s.db.AddFeed(context.Background(), NewFeed)
	if err != nil {
		return fmt.Errorf("failed to create feed: %v", err)
	}
	createFeedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user_id,
		FeedID:    feed.ID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), createFeedFollow)
	if err != nil {
		return fmt.Errorf("failed to follow the given feed: %v", err)
	}
	fmt.Println(feed)
	return nil
}

func handleGetFeeds(s *state, _ command) error {
	feeds, err := s.db.GetAllFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to retrive feed from db: %v", err)
	}
	for _, feed := range feeds {
		userName, err := s.db.GetUserNameById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("failed to retrive user by id: %v", err)
		}
		fmt.Printf("Title: 		%s\n", feed.Title)
		fmt.Printf("Url:   		%s\n", feed.Url)
		fmt.Printf("Created By: %s\n", userName)
		fmt.Println()
	}
	return nil
}

func handleFeedFollow(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usgae: blogag <feed-url>")
	}
	feed_id, err := s.db.GetFeedIdByFeedUrl(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("failed to retrive feed url: %v", err)
	}
	user_id, err := s.db.GetUserIDbyName(context.Background(), s.cfg.Current_user_name)
	if err != nil {
		return fmt.Errorf("failed to retrive user id: %v", err)
	}
	createFeedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user_id,
		FeedID:    feed_id,
	}
	resp, err := s.db.CreateFeedFollow(context.Background(), createFeedFollow)
	if err != nil {
		return fmt.Errorf("failed to follow the given feed: %v", err)
	}
	fmt.Printf("Feed: %s\n", resp.FeedName)
	fmt.Printf("Followed by: %s\n\n", resp.UserName)
	return nil
}

func handleGetFeedFollowsForUser(s *state, _ command) error {
	userFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), s.cfg.Current_user_name)
	if err != nil {
		return fmt.Errorf("failed to retrive feed for given user: %v", err)
	}
	for _, name := range userFeeds {
		fmt.Println(name)
	}
	return nil
}

func handlerUnfollowFeed(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: blogag unfollow - <url>")
	}
	feedId, err := s.db.GetFeedIdByFeedUrl(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("failed to retrive feed id: %v", err)
	}
	userId, err := s.db.GetUserIDbyName(context.Background(), s.cfg.Current_user_name)
	if err != nil {
		return fmt.Errorf("failed to retrive user id: %v", err)
	}
	if err := s.db.UnfollowFeed(context.Background(), database.UnfollowFeedParams{UserID: userId, FeedID: feedId}); err != nil {
		return fmt.Errorf("failed to unfollow feed: %v", err)
	}
	return nil
}

func scrapeFeeds(s *state) error {
	feedUrl, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("failed to fetch next feed: %v", err)
	}
	feedId, err := s.db.GetFeedIdByFeedUrl(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("failed to get feed id: %v", err)
	}
	feed, err := fetchFeed(context.Background(), feedUrl)
	if err != nil {
		return err
	}
	fmt.Println("Feed fetched: ")
	fmt.Println(feed.Channel.Title)
	feedFetched := database.MarkFeedFetchedParams{
		ID:        feedId,
		UpdatedAt: time.Now(),
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
	if err := s.db.MarkFeedFetched(context.Background(), feedFetched); err != nil {
		return fmt.Errorf("failed to mark feed as fetched: %v", err)
	}
	if err := savePosts(s, feedId, feed); err != nil {
		return fmt.Errorf("failed to save post: %v", err)
	}
	return nil
}

func ParsePubDate(s string) (time.Time, error) {
	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC850,
	}

	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unknown date format: %s", s)
}

func savePosts(s *state, feedID uuid.UUID, rssFeed *RSSFeed) error {
	for _, item := range rssFeed.Channel.Item {

		// Parse pubDate
		pubTime, err := ParsePubDate(item.PubDate)
		if err != nil {
			return err
		}

		params := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: pubTime,
			FeedID:      feedID,
		}

		// Insert into DB
		err = s.db.CreatePost(context.Background(), params)
		if err != nil {
			if strings.Contains(err.Error(), "23505") {
				// Duplicate URL -> ignore silently
				continue
			}
			return err
		}
	}

	return nil
}

func handleBrowse(s *state, cmd command) error {
	var limit int32 = 2
	if len(cmd.args) == 1 {
		parsedLimit, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("failed to convert limit: %v", err)
		}
		limit = int32(parsedLimit)
	}
	uid, err := s.db.GetUserIDbyName(context.Background(), s.cfg.Current_user_name)
	if err != nil {
		return fmt.Errorf("failed to retirve user id: %v", err)
	}
	param := database.GetPostsForUserParams{
		UserID: uid,
		Limit:  limit,
	}
	posts, err := s.db.GetPostsForUser(context.Background(), param)
	if err != nil {
		return fmt.Errorf("failed to retrive posts: %v", err)
	}
	fmt.Println(posts)
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
