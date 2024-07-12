package commands

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"go.mongodb.org/mongo-driver/mongo"
	"gotui/internal/storage"
	"os"
	"time"
)

type DbConnection struct {
	Client         *mongo.Client
	UserRepository *storage.UserRepository
	Err            error
}

var preloadedUsers = []storage.User{
	{Email: "john@gmail.com"},
	{Email: "ana2@yahoo.com"},
}

func InitDatabase() tea.Msg {
	ctx := context.TODO()
	dbClient, err := storage.NewDbConnection(ctx, os.Getenv(""))
	if err != nil {
		return DbConnection{Err: err}
	}
	err = dbClient.Ping(ctx, nil)
	if err != nil {
		return DbConnection{Err: fmt.Errorf("failed pinging db: %v", err)}
	}

	userRepository, err := storage.NewUserRepository(dbClient.Database("gotui"))
	if err != nil {
		return DbConnection{Err: fmt.Errorf("failed creating user repository: %v", err)}
	}

	for _, user := range preloadedUsers {
		if err = userRepository.AddUserIfNotExists(ctx, user); err != nil {
			return err
		}
	}

	return DbConnection{
		Client:         dbClient,
		UserRepository: userRepository,
		Err:            nil,
	}
}

type GetUserByEmailMsg struct {
	User *storage.User
	Err  error
}

func GetUserByEmail(userRepo *storage.UserRepository, email string) tea.Cmd {
	return func() tea.Msg {
		user, err := userRepo.FindByEmail(context.TODO(), email)

		return GetUserByEmailMsg{
			User: user,
			Err:  err,
		}
	}
}

type GetTokenByUserEmail struct {
	Token string
	Err   error
}

func GetLatestTokenByUserEmail(userEmail string) tea.Cmd {
	return func() tea.Msg {
		// We simulate an I/O
		time.Sleep(time.Second * 2)

		timeTextBytes, err := time.Now().UTC().MarshalText()

		return GetTokenByUserEmail{Token: userEmail + string(timeTextBytes), Err: err}
	}
}
