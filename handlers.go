package main

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/Isudin/gator/feed"
	"github.com/Isudin/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("%v command expects username argument", cmd.Name)
	}

	username := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Println("Username " + username + " has been saved to configuration")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("%v command expects username argument", cmd.Name)
	}

	username := cmd.Args[0]
	id := uuid.New()
	params := database.CreateUserParams{
		ID:        id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	user, _ := s.db.GetUser(context.Background(), username)
	if user.Name != "" {
		return fmt.Errorf("user '%v' already exists", user.Name)
	}

	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Printf("User %v has been created\n", username)
	fmt.Printf("UUID: %v\n", user.ID)
	fmt.Printf("CreatedAt: %v\n", user.CreatedAt)
	fmt.Printf("UpdatedAt: %v\n", user.UpdatedAt)

	return nil
}

func handlerReset(s *state, _ command) error {
	err := s.db.DeleteUsers(context.Background())
	return err
}

func handlerUsers(s *state, _ command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	if len(users) == 0 {
		fmt.Println("No users found")
	}

	for _, user := range users {
		if s.cfg.CurrentUser == user {
			user += " (current)"
		}
		fmt.Printf("* %v\n", user)
	}

	return nil
}

func handlerAggregate(s *state, cmd command) error {
	// if len(cmd.Args) < 1 {
	// 	return fmt.Errorf("%v command expects rss link argument", cmd.Name)
	// }

	// uri := cmd.Args[0]
	uri := "https://www.wagslane.dev/index.xml"
	_, err := url.ParseRequestURI(uri)
	if err != nil {
		return err
	}

	feed, err := feed.FetchFeed(context.Background(), uri)
	if err != nil {
		return err
	}

	fmt.Printf("%v - %v\n", feed.Channel.Title, feed.Channel.Description)
	for _, item := range feed.Channel.Item {
		fmt.Printf("%v\n%v\n\n", item.Title, item.Description)
	}
	//fmt.Println(feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("%v command expects name and url arguments", cmd.Name)
	}

	feedName := cmd.Args[0]
	uri := cmd.Args[1]
	_, err := url.ParseRequestURI(uri)
	if err != nil {
		return err
	}

	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUser)
	if err != nil {
		return err
	}
	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       uri,
		UserID:    currentUser.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Println("Feed has been created:")
	fmt.Printf("Name: %v\n", feed.Name)
	fmt.Printf("UUID: %v\n", feed.ID)
	fmt.Printf("CreatedAt: %v\n", feed.CreatedAt)
	fmt.Printf("UpdatedAt: %v\n", feed.UpdatedAt)
	fmt.Printf("Url: %v\n", feed.Url)
	fmt.Printf("UserID: %v\n", feed.UserID)

	cmd.Args = []string{feed.Url}
	err = handlerFollow(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func handlerListFeeds(s *state, _ command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("Name: %v\n", feed.Name)
		fmt.Printf("Url: %v\n", feed.Url)
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		fmt.Print("User: ")
		if err != nil {
			fmt.Println("<user not found>")
		} else {
			fmt.Println(user.Name)
		}
		fmt.Println("---")
	}

	return nil
}

func handlerFollow(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("%v command expects url argument", cmd.Name)
	}

	uri := cmd.Args[0]
	_, err := url.ParseRequestURI(uri)
	if err != nil {
		return err
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), uri)
	if err != nil {
		return err
	}

	user, err := s.db.GetUserByID(context.Background(), feed.UserID)
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("Feed %v of user %v followed\n", feedFollow.FeedName, feedFollow.UserName)
	return nil
}
