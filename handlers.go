package main

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

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
	if len(cmd.Args) < 1 {
		return fmt.Errorf("%v command expects time argument", cmd.Name)
	}

	timeBetweenRequests := cmd.Args[0]
	duration, err := time.ParseDuration(timeBetweenRequests)
	if err != nil {
		return fmt.Errorf("%v; error parsing time argument", err)
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)

	ticker := time.NewTicker(duration)
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("%v command expects name and url arguments", cmd.Name)
	}

	feedName := cmd.Args[0]
	uri := cmd.Args[1]
	_, err := url.ParseRequestURI(uri)
	if err != nil {
		return err
	}

	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       uri,
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Println("Feed has been created:")
	fmt.Printf("Name: %v\n", feed.Name)
	fmt.Printf("Url: %v\n", feed.Url)

	cmd.Args = []string{feed.Url}
	err = handlerFollow(s, cmd, user)
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

func handlerFollow(s *state, cmd command, user database.User) error {
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

	feedUser, err := s.db.GetUserByID(context.Background(), feed.UserID)
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

	fmt.Printf("Feed '%v' of user '%v' followed\n", feedFollow.FeedName, feedUser.Name)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return err
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds followed")
		return nil
	}

	fmt.Println("Currently followed feeds:")
	for _, feed := range feeds {
		fmt.Printf("- %v\n", feed.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("%v command expects url argument", cmd.Name)
	}

	uri := cmd.Args[0]
	_, err := url.ParseRequestURI(uri)
	if err != nil {
		return err
	}

	params := database.DeleteFeedFollowParams{UserID: user.ID, Url: uri}
	return s.db.DeleteFeedFollow(context.Background(), params)
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.Args) >= 1 {
		parsedNumber, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("parameter should be an integer")
		}

		limit = parsedNumber
	}

	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}
	posts, err := s.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf("%v: \n%v\n\n", post.Title, post.Description)
	}

	return nil
}
