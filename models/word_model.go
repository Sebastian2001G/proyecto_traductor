package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Word struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	WordText    string
	Translation string
}
