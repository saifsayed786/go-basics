package main

import (
	"fmt"
	_ "go-basics/docs"
	"go-basics/models"
	"go-basics/routes"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title GOLANG BASIC STRUCTURE
// @version 1.0
// @Description Golang basic API.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1/

func main() {
	godotenv.Load()          // Load env variables
	models.ConnectDataBase() // load db

	var router = make(chan *gin.Engine)
	go routes.GetRouter(router)

	r := <-router
	var port string = os.Getenv("SERVER_PORT")
	server_addr := fmt.Sprintf(":%s", port)
	r.Run(server_addr)

}
