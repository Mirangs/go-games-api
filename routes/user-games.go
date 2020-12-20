package routes

import (
	"github.com/Mirangs/bm-go-test-task/controllers"
)

func (r Routes) UserGames() {
	userGames := r.Router.Group("/api/user-games")
	userGames.GET("/", controllers.UserGamesController{Client: r.Client}.GetUserGames)
	userGames.GET("/:id", controllers.UserGamesController{Client: r.Client}.GetUserGameById)
	userGames.POST("/", controllers.UserGamesController{Client: r.Client}.CreateUserGame)
	userGames.PUT("/:id", controllers.UserGamesController{Client: r.Client}.UpdateUserGame)
	userGames.DELETE("/:id", controllers.UserGamesController{Client: r.Client}.DeleteUserGame)
}
