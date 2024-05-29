package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Word struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"` //EL ID SE ASIGNA AUTOMATICAMENTE
	WordText    string
	Translation string
}
