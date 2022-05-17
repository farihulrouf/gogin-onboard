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
    "github.com/go-playground/validator/v10"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

var todoCollection *mongo.Collection = configs.GetCollection(configs.DB, "todos")
var validate = validator.New()

func GetTodo() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        todoId := c.Param("todoId")
        var todo models.Todo
        defer cancel()
      
        objId, _ := primitive.ObjectIDFromHex(todoId)
      
        err := todoCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&todo)
        if err != nil {
            c.JSON(http.StatusInternalServerError, responses.DataResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }
      
        c.JSON(http.StatusOK, responses.DataResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": todo}})
    }
}

func GetAllTodos() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        var todos []models.Todo
        defer cancel()

        results, err := todoCollection.Find(ctx, bson.M{})

        if err != nil {
            c.JSON(http.StatusInternalServerError, responses.DataResponse{Status: http.StatusInternalServerError,
             Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        //reading from the db in an optimal way
        defer results.Close(ctx)
        for results.Next(ctx) {
            var singleTodo models.Todo
            if err = results.Decode(&singleTodo); err != nil {
                c.JSON(http.StatusInternalServerError, responses.DataResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            }
          
            todos = append(todos, singleTodo)
        }

        c.JSON(http.StatusOK,
            responses.DataResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": todos}},
        )
    }
}

func EditATodo() gin.HandlerFunc {
    return func(c *gin.Context) {
       ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        todoId := c.Param("todoId")
        var todo models.Todo
        defer cancel()
        objId, _ := primitive.ObjectIDFromHex(todoId)

        //validate the request body
        if err := c.BindJSON(&todo); err != nil {
            c.JSON(http.StatusBadRequest, responses.DataResponse{Status: http.StatusBadRequest, 
                Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        //use the validator library to validate required fields
        if validationErr := validate.Struct(&todo); validationErr != nil {
            c.JSON(http.StatusBadRequest, 
                responses.DataResponse{Status: http.StatusBadRequest, 
                Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
            return
        }

        update := bson.M{"title": todo.Title, "desc": todo.Desc}
        result, err := todoCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
        if err != nil {
            c.JSON(http.StatusInternalServerError, 
                responses.DataResponse{Status: http.StatusInternalServerError, 
                Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        //get updated user details
        var updatedTodo models.Todo
        if result.MatchedCount == 1 {
            err := todoCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedTodo)
            if err != nil {
                c.JSON(http.StatusInternalServerError, 
                    responses.DataResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
                return
            }
        }

        c.JSON(http.StatusOK, responses.DataResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedTodo}})
        
    }
}

func CreateTodo() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        var todo models.Todo
        defer cancel()

        //validate the request body
        if err := c.BindJSON(&todo); err != nil {
            c.JSON(http.StatusBadRequest, responses.DataResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        //use the validator library to validate required fields
        if validationErr := validate.Struct(&todo); validationErr != nil {
            c.JSON(http.StatusBadRequest, responses.DataResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
            return
        }

        newTodo := models.Todo{
            Id:       primitive.NewObjectID(),
            Title:     todo.Title,
            Desc: todo.Desc,
        }
      
        result, err := todoCollection.InsertOne(ctx, newTodo)
        if err != nil {
            c.JSON(http.StatusInternalServerError, responses.DataResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        c.JSON(http.StatusCreated, responses.DataResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
    }
}