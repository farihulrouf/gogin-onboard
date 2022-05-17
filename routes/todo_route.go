
package routes

import (
    "restapi/controllers" //add this
    "github.com/gin-gonic/gin"
)

func TodoRoute(router *gin.Engine)  {
    router.POST("/api/todo", controllers.CreateTodo()) //add this
    router.GET("/api/todos", controllers.GetAllTodos()) //get all
    router.PUT("/api/todo/:todoId", controllers.EditATodo())
    router.GET("/api/todo/:todoId", controllers.GetTodo())


    router.POST("/api/job", controllers.CreateJob())
    router.GET("/api/jobs", controllers.GetAllJob())
    router.PUT("/api/job/:jobId", controllers.EditJob())
    router.GET("/api/job/:jobId", controllers.GetJob())


}

