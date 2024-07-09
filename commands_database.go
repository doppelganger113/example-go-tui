package main

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"go.mongodb.org/mongo-driver/mongo"
	"gotui/internal/storage"
	"os"
	"time"
)

type dbConnection struct {
	client         *mongo.Client
	userRepository *storage.UserRepository
	err            error
}

var preloadedUsers = []storage.User{
	{Email: "john@gmail.com"},
	{Email: "ana2@yahoo.com"},
}

func initDatabase() tea.Msg {
	ctx := context.TODO()
	dbClient, err := storage.NewDbConnection(ctx, os.Getenv(""))
	if err != nil {
		return dbConnection{err: err}
	}
	err = dbClient.Ping(ctx, nil)
	if err != nil {
		return dbConnection{err: fmt.Errorf("failed pinging db: %v", err)}
	}

	userRepository, err := storage.NewUserRepository(dbClient.Database("gotui"))
	if err != nil {
		return dbConnection{err: fmt.Errorf("failed creating user repository: %v", err)}
	}

	for _, user := range preloadedUsers {
		if err = userRepository.AddUserIfNotExists(ctx, user); err != nil {
			return err
		}
	}

	return dbConnection{
		client:         dbClient,
		userRepository: userRepository,
		err:            nil,
	}
}

type getUserByEmailMsg struct {
	user *storage.User
	err  error
}

func getUserByEmail(userRepo *storage.UserRepository, email string) tea.Cmd {
	return func() tea.Msg {
		user, err := userRepo.FindByEmail(context.TODO(), email)

		return getUserByEmailMsg{
			user: user,
			err:  err,
		}
	}
}

type getTokenByUserEmail struct {
	token string
	err   error
}

func getLatestTokenByUserEmail(userEmail string) tea.Cmd {
	return func() tea.Msg {
		// We simulate an I/O
		time.Sleep(time.Second * 2)

		timeTextBytes, err := time.Now().UTC().MarshalText()

		return getTokenByUserEmail{token: userEmail + string(timeTextBytes), err: err}
	}
}
