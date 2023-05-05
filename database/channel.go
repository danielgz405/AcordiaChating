package database

import (
	"context"

	"github.com/dg/acordia/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo *MongoRepo) CreateChannel(ctx context.Context, data models.InsertChannel) (*models.Channel, error) {
	collection := repo.client.Database("Acordia").Collection("channels")
	result, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	oid := result.InsertedID.(primitive.ObjectID)
	channel, err := repo.GetChannelById(ctx, oid.Hex())
	if err != nil {
		return nil, err
	}
	return channel, nil
}

func (repo *MongoRepo) GetChannelById(ctx context.Context, id string) (*models.Channel, error) {
	collection := repo.client.Database("Acordia").Collection("channels")
	var channel models.Channel
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&channel)
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

func (repo *MongoRepo) UpdateChannel(ctx context.Context, id string, data models.UpdateChannel) (*models.Channel, error) {
	collection := repo.client.Database("Acordia").Collection("channels")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	update := bson.M{
		"$set": bson.M{},
	}
	iterableData := map[string]interface{}{
		"name":                  data.Name,
		"description":           data.Description,
		"color":                 data.Color,
		"background":            data.Background,
		"desert_ref_background": data.DesertRefBackground,
		"image":                 data.Image,
		"desert_ref_image":      data.DesertRefImage,
	}
	for key, value := range iterableData {
		if value != nil && value != "" {
			update["$set"].(bson.M)[key] = value
		}
	}
	_, err = collection.UpdateOne(ctx, bson.M{"_id": oid}, update)
	if err != nil {
		return nil, err
	}
	UpdateChannel, err := repo.GetChannelById(ctx, id)
	if err != nil {
		return nil, err
	}
	return UpdateChannel, nil
}

func (repo *MongoRepo) DeleteChannel(ctx context.Context, id string) error {
	collection := repo.client.Database("Acordia").Collection("channels")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	return nil
}

func (repo *MongoRepo) AddUserToChannel(ctx context.Context, userId string, channelId string) (*models.Channel, error) {
	collection := repo.client.Database("Acordia").Collection("channels")
	oid, err := primitive.ObjectIDFromHex(channelId)
	if err != nil {
		return nil, err
	}
	porfile, err := repo.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}
	_, err = collection.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$addToSet": bson.M{"users": porfile}})
	if err != nil {
		return nil, err
	}
	updateUser, err := repo.GetChannelById(ctx, channelId)
	if err != nil {
		return nil, err
	}
	return updateUser, nil
}

func (repo *MongoRepo) RemoveUser(ctx context.Context, channelId string, userId string) (*models.Channel, error) {
	collection := repo.client.Database("Acordia").Collection("channels")
	oid, err := primitive.ObjectIDFromHex(channelId)
	if err != nil {
		return nil, err
	}
	usOid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}
	_, err = collection.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$pull": bson.M{"users": bson.M{"_id": usOid}}})
	if err != nil {
		return nil, err
	}
	updateUser, err := repo.GetChannelById(ctx, channelId)
	if err != nil {
		return nil, err
	}
	return updateUser, nil
}

func (repo *MongoRepo) AddMessagesToChannel(ctx context.Context, data *models.ChannelMessage, channelId string) (*models.Channel, error) {
	collection := repo.client.Database("Acordia").Collection("channels")
	oid, err := primitive.ObjectIDFromHex(channelId)
	if err != nil {
		return nil, err
	}
	_, err = collection.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$push": bson.M{"messages": data}})
	if err != nil {
		return nil, err
	}
	updateUser, err := repo.GetChannelById(ctx, channelId)
	if err != nil {
		return nil, err
	}
	return updateUser, nil
}

func (repo *MongoRepo) ListOfChannels(ctx context.Context, usOid primitive.ObjectID) ([]models.Channel, error) {
	collection := repo.client.Database("Acordia").Collection("channels")
	users := []primitive.ObjectID{usOid}
	cursor, err := collection.Find(ctx, bson.M{"users._id": bson.M{"$in": users}})
	if err != nil {
		return nil, err
	}
	channels := []models.Channel{}
	if err = cursor.All(ctx, &channels); err != nil {
		return nil, err
	}
	return channels, nil
}
