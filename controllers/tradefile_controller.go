package controllers
import (
    "context"
    "restapi/configs"
    "restapi/models"
    "restapi/responses"
    "net/http"
    "time"
    "go.mongodb.org/mongo-driver/bson"
    "github.com/gin-gonic/gin"
    //"github.com/go-playground/validator/v10"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

var tradefileCollection *mongo.Collection = configs.GetCollection(configs.DB, "tradefiles")


func GetAllTradeFile() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        var tradefiles []models.Tradefile
        defer cancel()

        results, err := tradefileCollection.Find(ctx, bson.M{})

        if err != nil {
            c.JSON(http.StatusInternalServerError, responses.DataResponse{Status: http.StatusInternalServerError,
             Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        //reading from the db in an optimal way
        defer results.Close(ctx)
        for results.Next(ctx) {
            var singleTradefile models.Tradefile
            if err = results.Decode(&singleTradefile); err != nil {
                c.JSON(http.StatusInternalServerError, responses.DataResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            }
          
            tradefiles = append(tradefiles, singleTradefile)
        }

        c.JSON(http.StatusOK,
            responses.DataResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": tradefiles}},
        )
    }
}


func CreateFile() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        var tradefile models.Tradefile
        defer cancel()

        //validate the request body
        if err := c.BindJSON(&tradefile); err != nil {
            c.JSON(http.StatusBadRequest, responses.DataResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        //use the validator library to validate required fields
        if validationErr := validate.Struct(&tradefile); validationErr != nil {
            c.JSON(http.StatusBadRequest, responses.DataResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
            return
        }

        newTradedfile := models.Tradefile{
            Id:       	primitive.NewObjectID(),
            Expire:     	tradefile.Expire,
            Tprice: 		tradefile.Tprice,
            Rprice:			tradefile.Rprice,
            Contractcode:	tradefile.Contractcode,
        }
      
        result, err := tradefileCollection.InsertOne(ctx, newTradedfile)
        if err != nil {
            c.JSON(http.StatusInternalServerError, responses.DataResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        c.JSON(http.StatusCreated, responses.DataResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
    }
}