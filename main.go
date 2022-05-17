package main
import (
        "github.com/gin-gonic/gin"
        "restapi/configs"
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

        router.Run(":7000")
}


/*
package main

import (
    "gin-mongo-api/configs"
    "gin-mongo-api/routes" //add this
    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()
    
    //run database
    configs.ConnectDB()

    //routes
    routes.UserRoute(router) //add this

    router.Run("localhost:6000")
}
*/
