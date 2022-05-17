package main
import (
        "github.com/gin-gonic/gin"
        "restapi/configs"
        "restapi/routes"
)

func main() {
        router := gin.Default()
        /*
        router.GET("/", func(c *gin.Context) {
                c.JSON(200, gin.H {
                        "data" : "Its Oe From Gin-gonix & mondoDB",
                })
        })
        */
        //connect to mongodb
        configs.ConnectDB()


        routes.TodoRoute(router) //add this

        router.Run(":7000")
}


