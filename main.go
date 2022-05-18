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


        routes.RouteGo(router) 

        router.Run(":8080")
}


