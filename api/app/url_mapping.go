package app

import (
	"know-sync-api/controllers/procedures"
	"know-sync-api/controllers/users"
	"know-sync-api/middlewares"
)

func mapUrls() {
	router.POST("/sign_up", users.SignUp)
	router.POST("/sign_in", users.SignIn)
	router.POST("/auth_user", users.AuthUser)
	router.DELETE("/signout", users.SignOut)
	router.POST("/refresh", users.Refresh)
	router.DELETE("/users/:user_id", middlewares.TokenAuthMiddleware(), users.DeleteUser)
	router.GET("/procedures/:procedure_id", procedures.GetProcedure)
	router.GET("/procedures", procedures.GetProcedures)
	router.POST("/procedures", procedures.CreateProcedure)
	router.PUT("/procedures/:procedure_id", procedures.UpdateProcedure)
	router.PATCH("/procedures/:procedure_id", procedures.UpdateProcedure)
	router.DELETE("/procedures/:procedure_id", procedures.DeleteProcedure)
}
