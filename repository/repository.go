package repository

import (
	"context"

	"github.com/dg/acordia/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	//Users
	InsertUser(ctx context.Context, user *models.InsertUser) (*models.Profile, error)
	GetUserById(ctx context.Context, id string) (*models.Profile, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	ListUsers(ctx context.Context) ([]models.Profile, error)
	UpdateUser(ctx context.Context, data models.UpdateUser) (*models.Profile, error)
	DeleteUser(ctx context.Context, id string) error

	//channels
	CreateChannel(ctx context.Context, data models.InsertChannel) (*models.Channel, error)
	UpdateChannel(ctx context.Context, id string, data models.UpdateChannel) (*models.Channel, error)
	DeleteChannel(ctx context.Context, id string) error
	AddUserToChannel(ctx context.Context, userId string, channelId string) (*models.Channel, error)
	RemoveUser(ctx context.Context, channelId string, userId string) (*models.Channel, error)
	AddMessagesToChannel(ctx context.Context, data *models.ChannelMessage, channelId string) (*models.Channel, error)
	ListOfChannels(ctx context.Context, usOid primitive.ObjectID) ([]models.Channel, error)

	//Close the connection
	Close() error
}

var implementation Repository

// Repo
func SetRepository(repository Repository) {
	implementation = repository
}

// Close the connection
func Close() error {
	return implementation.Close()
}
