package repository

import (
	"context"

	"github.com/dg/acordia/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateChannel(ctx context.Context, data models.InsertChannel) (*models.Channel, error) {
	return implementation.CreateChannel(ctx, data)
}

func UpdateChannel(ctx context.Context, id string, data models.UpdateChannel) (*models.Channel, error) {
	return implementation.UpdateChannel(ctx, id, data)
}

func DeleteChannel(ctx context.Context, id string) error {
	return implementation.DeleteChannel(ctx, id)
}

func AddUserToChannel(ctx context.Context, userId string, channelId string) (*models.Channel, error) {
	return implementation.AddUserToChannel(ctx, userId, channelId)
}

func RemoveUser(ctx context.Context, channelId string, userId string) (*models.Channel, error) {
	return implementation.RemoveUser(ctx, channelId, userId)
}

func AddMessagesToChannel(ctx context.Context, data *models.ChannelMessage, channelId string) (*models.Channel, error) {
	return implementation.AddMessagesToChannel(ctx, data, channelId)
}

func ListOfChannels(ctx context.Context, usOid primitive.ObjectID) ([]models.Channel, error) {
	return implementation.ListOfChannels(ctx, usOid)
}
