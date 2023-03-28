package database

import (
	"context"
	"fmt"

	"github.com/dg/acordia/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo *MongoRepo) InsertUser(ctx context.Context, user *models.InsertUser) (profile *models.Profile, err error) {
	collection := repo.client.Database("Acordia").Collection("users")
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		fmt.Println("1")
		fmt.Println(err.Error())
		return nil, err
	}
	oid := result.InsertedID.(primitive.ObjectID)
	profile, err = repo.GetUserById(ctx, oid.Hex())
	if err != nil {
		fmt.Println("2")
		fmt.Println(err.Error())
		return nil, err
	}
	return profile, nil
}
func (repo *MongoRepo) GetUserById(ctx context.Context, id string) (*models.Profile, error) {
	collection := repo.client.Database("Acordia").Collection("users")
	var user models.User
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	// Find one and populate company
	err = collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	// Populate profile
	var profile = models.Profile{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Image:     user.Image,
		DesertRef: user.DesertRef,
	}
	return &profile, nil
}
func (repo *MongoRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	collection := repo.client.Database("Acordia").Collection("users")
	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (repo *MongoRepo) ListUsers(ctx context.Context, companyId string) ([]models.Profile, error) {
	collection := repo.client.Database("Acordia").Collection("users")
	oid, err := primitive.ObjectIDFromHex(companyId)
	cursor, err := collection.Find(ctx, bson.M{"company": oid})
	if err != nil {
		return nil, err
	}
	var users []models.User
	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, err
	}
	var profiles []models.Profile
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		// Populate profile
		var profile = models.Profile{
			Id:        user.Id,
			Name:      user.Name,
			Email:     user.Email,
			Image:     user.Image,
			DesertRef: user.DesertRef,
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}
func (repo *MongoRepo) UpdateUser(ctx context.Context, data models.UpdateUser) (*models.Profile, error) {
	collection := repo.client.Database("Acordia").Collection("users")
	oid, err := primitive.ObjectIDFromHex(data.Id)
	if err != nil {
		return nil, err
	}
	update := bson.M{
		"$set": bson.M{},
	}
	iterableData := map[string]interface{}{
		"name":      data.Name,
		"email":     data.Email,
		"image":     data.Image,
		"desertref": data.DesertRef,
	}
	for key, value := range iterableData {
		if value != "" {
			update["$set"].(bson.M)[key] = value
		}
	}
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": oid}, update).Err()
	if err != nil {
		return nil, err
	}
	profile, err := repo.GetUserById(ctx, data.Id)
	if err != nil {
		return nil, err
	}
	return profile, nil
}
func (repo *MongoRepo) DeleteUser(ctx context.Context, id string) error {
	collection := repo.client.Database("Acordia").Collection("users")
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
