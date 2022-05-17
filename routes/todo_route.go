
package routes

import (
    "restapi/controllers" //add this
    "github.com/gin-gonic/gin"
)

func TodoRoute(router *gin.Engine)  {
    router.POST("/api/todo", controllers.CreateTodo()) //add this
    router.GET("/api/todos", controllers.GetAllTodos())
}

