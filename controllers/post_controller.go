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


var postCollection *mongo.Collection = configs.GetCollection(configs.DB, "posts")
func CreatePost() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel:= context.WithTimeout(context.Background(), 10*time.Second)
        var post models.Post
        defer cancel()

        //Validate the request body
        if err := c.BindJSON(&post); err != nil {
            c.JSON(http.StatusBadRequest, responses.DataResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        //use validaro libaray to validate required fields
         if validationErr := validate.Struct(&post); validationErr != nil {
            c.JSON(http.StatusBadRequest, responses.DataResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
            return
        }

        //Create Post
        newPost := models.Post{
            Id: primitive.NewObjectID(),
            Title:  post.Title,
            Content: post.Content,
        }

        result, err := postCollection.InsertOne(ctx, newPost)
        if err != nil {
            c.JSON(http.StatusInternalServerError, responses.DataResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        c.JSON(http.StatusCreated, responses.DataResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
    }
}

func GetAllPost() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        var posts []models.Post
        defer cancel()

        results, err := postCollection.Find(ctx, bson.M{})

        if err != nil {
            c.JSON(http.StatusInternalServerError, responses.DataResponse{Status: http.StatusInternalServerError,
             Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        //reading from the db in an optimal way
        defer results.Close(ctx)
        for results.Next(ctx) {
            var singlePost models.Post
            if err = results.Decode(&singlePost); err != nil {
                c.JSON(http.StatusInternalServerError, responses.DataResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            }
          
            posts = append(posts, singlePost)
        }

        c.JSON(http.StatusOK,
            responses.DataResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": posts}},
        )
    }
}