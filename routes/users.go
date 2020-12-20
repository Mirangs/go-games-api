package routes

import (
	"github.com/Mirangs/bm-go-test-task/controllers"
)

func (r Routes) Users() {
	users := r.Router.Group("/api/users")
	users.GET("/", controllers.UsersController{Client: r.Client}.GetUsers)
	users.GET("/:id", controllers.UsersController{Client: r.Client}.GetUserById)
	users.POST("/", controllers.UsersController{Client: r.Client}.CreateUser)
	users.PUT("/:id", controllers.UsersController{Client: r.Client}.UpdateUser)
	users.DELETE("/:id", controllers.UsersController{Client: r.Client}.DeleteUser)

	r.Router.GET("/api/users-rating", controllers.UsersController{Client: r.Client}.GetUsersRating)
	r.Router.GET("/api/games-statistics/:id", controllers.UsersController{Client: r.Client}.GetGamesStatistics)
}
