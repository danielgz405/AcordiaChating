package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Channel struct {
	Id                  primitive.ObjectID `bson:"_id" json:"_id"`
	Users               []Profile          `bson:"users" json:"users"`
	Color               string             `bson:"color" json:"color"`
	Background          string             `bson:"background" json:"background"`
	DesertRefBackground string             `bson:"desert_ref_background" json:"desert_ref_background"`
	Image               string             `bson:"image" json:"image"`
	DesertRefImage      string             `bson:"desert_ref_image" json:"desert_ref_image"`
	CreateDate          string             `bson:"create_date" json:"create_date"`
	Description         string             `bson:"description" json:"description"`
	Name                string             `bson:"name" json:"name"`
	Messages            []ChannelMessage   `bson:"messages" json:"messages"`
}

type ChannelMessage struct {
	User        Profile `bson:"user" json:"user"`
	Date        string  `bson:"date" json:"date"`
	Description string  `bson:"description" json:"description"`
	Image       string  `bson:"image" json:"image"`
	DesertRef   string  `bson:"desert_ref" json:"desert_ref"`
}

type InsertChannel struct {
	Users               []Profile        `bson:"users" json:"users"`
	Color               string           `bson:"color" json:"color"`
	Background          string           `bson:"background" json:"background"`
	DesertRefBackground string           `bson:"desert_ref_background" json:"desert_ref_background"`
	Image               string           `bson:"image" json:"image"`
	DesertRefImage      string           `bson:"desert_ref_image" json:"desert_ref_image"`
	CreateDate          string           `bson:"create_date" json:"create_date"`
	Description         string           `bson:"description" json:"description"`
	Name                string           `bson:"name" json:"name"`
	Messages            []ChannelMessage `bson:"messages" json:"messages"`
}

type UpdateChannel struct {
	Name                string           `bson:"name" json:"name"`
	Description         string           `bson:"description" json:"description"`
	Color               string           `bson:"color" json:"color"`
	Background          string           `bson:"background" json:"background"`
	DesertRefBackground string           `bson:"desert_ref_background" json:"desert_ref_background"`
	Image               string           `bson:"image" json:"image"`
	DesertRefImage      string           `bson:"desert_ref_image" json:"desert_ref_image"`
}