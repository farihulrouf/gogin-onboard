package main
import (
        "github.com/gin-gonic/gin"
        "restapi/configs"
        "restapi/routes"
)

func main() {
        router := gin.Default()
        //connect to mongodb
        configs.ConnectDB()


        routes.TodoRoute(router) //add this

        router.Run(":7000")
}


