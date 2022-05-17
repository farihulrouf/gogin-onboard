package models
import "go.mongodb.org/mongo-driver/bson/primitive"

type Job struct {
    Id       primitive.ObjectID   `json:"id,omitempty"`
    Title      string             `json:"title,omitempty" validate:"required"`
    Desc       string             `json:"desc,omitempty" validate:"required"`
    Depart     string             `json:"depart,omitempty" validate:"required"`
    No         int                `json:"no,omitempty" validate:"required"`
}