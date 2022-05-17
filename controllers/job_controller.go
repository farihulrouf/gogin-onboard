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


var jobCollection *mongo.Collection = configs.GetCollection(configs.DB, "jobs")
//var validate = validator.New()

func GetAllJob() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        var jobs []models.Job
        defer cancel()

        results, err := jobCollection.Find(ctx, bson.M{})

        if err != nil {
            c.JSON(http.StatusInternalServerError, responses.DataResponse{Status: http.StatusInternalServerError,
             Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        //reading from the db in an optimal way
        defer results.Close(ctx)
        for results.Next(ctx) {
            var singleJob models.Job
            if err = results.Decode(&singleJob); err != nil {
                c.JSON(http.StatusInternalServerError, responses.DataResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            }
          
            jobs = append(jobs, singleJob)
        }

        c.JSON(http.StatusOK,
            responses.DataResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": jobs}},
        )
    }
}

func EditJob() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        jobId := c.Param("jobId")
        var job models.Job
        defer cancel()
        objId, _ := primitive.ObjectIDFromHex(jobId)

        if err := c.BindJSON(&job); err != nil {
            c.JSON(http.StatusBadRequest, responses.DataResponse{Status: http.StatusBadRequest, 
                Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        //use the validator library to validate required fields
        if validationErr := validate.Struct(&job); validationErr != nil {
            c.JSON(http.StatusBadRequest, 
                responses.DataResponse{Status: http.StatusBadRequest, 
                Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
            return
        }

        update := bson.M{"title": job.Title, "desc": job.Desc, "depart": job.Depart, "no": job.No}
        result, err := jobCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
        if err != nil {
            c.JSON(http.StatusInternalServerError, 
                responses.DataResponse{Status: http.StatusInternalServerError, 
                Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        if err != nil {
            c.JSON(http.StatusInternalServerError, 
                responses.DataResponse{Status: http.StatusInternalServerError, 
                Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        //get updated user details
        var updatedJob models.Job
        if result.MatchedCount == 1 {
            err := jobCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedJob)
            if err != nil {
                c.JSON(http.StatusInternalServerError, 
                    responses.DataResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
                return
            }
        }

        c.JSON(http.StatusOK, responses.DataResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedJob}})
    

    }
}

func CreateJob() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        var job models.Job
        defer cancel()
    

    //validate the request body
    if err := c.BindJSON(&job); err != nil {
        c.JSON(http.StatusBadRequest, responses.DataResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
        return
    }

    //use the validator library to validate required fields
    if validationErr := validate.Struct(&job); validationErr != nil {
        c.JSON(http.StatusBadRequest, responses.DataResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
        return
    }

    newJob := models.Job{
        Id:       primitive.NewObjectID(),
        Title:     job.Title,
        Desc:      job.Desc,
        Depart:    job.Depart,
        No:        job.No,       
    }

    result, err := jobCollection.InsertOne(ctx, newJob)

    if err != nil {
        c.JSON(http.StatusInternalServerError, responses.DataResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
        return
     }

    c.JSON(http.StatusCreated, responses.DataResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
  
  }

}