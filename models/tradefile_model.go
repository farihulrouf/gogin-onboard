package models
import ( 
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Tradefile struct {
    Id       primitive.ObjectID      `json:"id,omitempty"`
    Expire       time.Time           `json:"expire,omitempty" validate:"required"`
    Tprice       float64             `json:"tprice,omitempty" validate:"required"`
    Rprice       float64             `json:"rprice,omitempty" validate:"required"`
    Contractcode string              `json:"contractcode,omitempty" validate:"required"`
    CreatedAt    time.Time           `json:”created_at,omitempty” bson:”created_at”`
    UpdatedAt    time.Time           `json:”updated_at,omitempty” bson:”updated_at”`
    
}