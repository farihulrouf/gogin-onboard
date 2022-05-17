package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Todo struct {
    Id       primitive.ObjectID `json:"id,omitempty"`
    Title     string             `json:"title,omitempty" validate:"required"`
    Desc      string             `json:"desc,omitempty" validate:"required"`
}