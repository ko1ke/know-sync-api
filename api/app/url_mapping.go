package app

import (
	"github.com/ko1ke/know-sync-api/controllers/procedures"
	"github.com/ko1ke/know-sync-api/controllers/users"
	"github.com/ko1ke/know-sync-api/middlewares"
)

func mapUrls() {
	router.POST("/sign_up", users.SignUp)
	router.POST("/sign_in", users.SignIn)
	router.POST("/auth_user", users.AuthUser)
	router.DELETE("/signout", users.SignOut)
	router.POST("/refresh", users.Refresh)
	router.DELETE("/users/:user_id", middlewares.TokenAuthMiddleware(), users.DeleteUser)
	router.GET("/procedures", middlewares.TokenAuthMiddleware(), procedures.GetProcedures)
	router.GET("/procedures/:procedure_id", middlewares.TokenAuthMiddleware(), procedures.GetProcedure)
	router.POST("/procedures", middlewares.TokenAuthMiddleware(), procedures.CreateProcedure)
	router.PUT("/procedures/:procedure_id", middlewares.TokenAuthMiddleware(), procedures.UpdateProcedure)
	router.PATCH("/procedures/:procedure_id", middlewares.TokenAuthMiddleware(), procedures.UpdateProcedure)
	router.DELETE("/procedures/:procedure_id", middlewares.TokenAuthMiddleware(), procedures.DeleteProcedure)
	router.GET("/public_procedures", procedures.GetPublicProcedures)
	router.GET("/public_procedures/:procedure_id", procedures.GetPublicProcedure)
}
