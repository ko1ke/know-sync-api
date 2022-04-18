package procedures

import (
	"fmt"
	"know-sync-api/controllers/users"

	"know-sync-api/domain/procedures"
	"know-sync-api/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetProcedure(c *gin.Context) {
	procedureID, idErr := getProcedureID(c.Param("procedure_id"))
	if idErr != nil {
		logrus.Error(idErr)
		c.JSON(http.StatusBadRequest, gin.H{"error": idErr.Error()})
		return
	}

	procedure, getErr := services.GetProcedureItem(uint(procedureID))

	isOwn, ownErr := isOwnProcedure(c.Request, &procedure.Procedure)

	if !isOwn {
		logrus.Error(ownErr)
		c.JSON(http.StatusForbidden, gin.H{"error": "not your procedure"})
		return
	}

	if getErr != nil {
		logrus.Error(getErr)
		c.JSON(http.StatusNotFound, gin.H{"error": getErr.Error()})
		return
	}

	c.JSON(http.StatusOK, procedure)
}

func GetPublicProcedure(c *gin.Context) {
	procedureID, idErr := getProcedureID(c.Param("procedure_id"))
	if idErr != nil {
		logrus.Error(idErr)
		c.JSON(http.StatusBadRequest, gin.H{"error": idErr.Error()})
		return
	}

	procedure, getErr := services.GetProcedureItem(uint(procedureID))

	if !procedure.Publish {
		logrus.Errorf("procedure %v is not public", procedure.ID)
		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("procedure %v is not public", procedure.ID)})
		return
	}

	if getErr != nil {
		logrus.Error(getErr)
		c.JSON(http.StatusNotFound, gin.H{"error": getErr.Error()})
		return
	}

	c.JSON(http.StatusOK, procedure)
}

func GetProcedures(c *gin.Context) {
	page, pageErr := getPage(c.Query("page"))
	if pageErr != nil {
		logrus.Error(pageErr)
		c.JSON(http.StatusBadRequest, gin.H{"error": pageErr.Error()})
		return
	}

	var limit int = 10
	var offset int = limit * (page - 1)

	user, err := users.GetUserFromToken(c.Request)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	procedures, getErr := services.GetProcedures(limit, offset, c.Query("keyword"), user.ID)
	if getErr != nil {
		logrus.Error(getErr)
		c.JSON(http.StatusNotFound, gin.H{"error": getErr.Error()})
		return
	}

	pagination, _ := services.GetPagination(page, limit, len(*procedures))

	c.JSON(http.StatusOK, gin.H{"procedures": procedures, "pagination": pagination})
}

func GetPublicProcedures(c *gin.Context) {
	page, pageErr := getPage(c.Query("page"))
	if pageErr != nil {
		logrus.Error(pageErr)
		c.JSON(http.StatusBadRequest, gin.H{"error": pageErr.Error()})
		return
	}

	var limit int = 10
	var offset int = limit * (page - 1)

	procedures, getErr := services.GetPublicProcedures(limit, offset, c.Query("keyword"))
	if getErr != nil {
		logrus.Error(getErr)
		c.JSON(http.StatusNotFound, gin.H{"error": getErr.Error()})
		return
	}

	pagination, _ := services.GetPagination(page, limit, len(*procedures))

	c.JSON(http.StatusOK, gin.H{"procedures": procedures, "pagination": pagination})
}

func CreateProcedure(c *gin.Context) {
	var procedure procedures.Procedure
	if err := c.ShouldBindJSON(&procedure); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	user, err := users.GetUserFromToken(c.Request)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	procedure.UserID = user.ID
	newProcedure, saveErr := services.CreateProcedure(procedure)
	if saveErr != nil {
		logrus.Error(saveErr)
		c.JSON(http.StatusBadRequest, gin.H{"error": saveErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, newProcedure)
}

func UpdateProcedure(c *gin.Context) {
	procedureID, idErr := getProcedureID(c.Param("procedure_id"))
	if idErr != nil {
		logrus.Error(idErr)
		c.JSON(http.StatusNotFound, gin.H{"error": idErr.Error()})
		return
	}

	currentProcedure, getErr := services.GetProcedure(uint(procedureID))

	if getErr != nil {
		logrus.Error(getErr)
		c.JSON(http.StatusNotFound, gin.H{"error": getErr.Error()})
		return
	}

	isOwn, ownErr := isOwnProcedure(c.Request, currentProcedure)
	if !isOwn {
		logrus.Error(ownErr)
		c.JSON(http.StatusForbidden, gin.H{"error": "not your procedure"})
		return
	}

	var procedure procedures.Procedure
	if err := c.ShouldBindJSON(&procedure); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	procedure.ID = uint(procedureID)

	isPartial := c.Request.Method == http.MethodPatch

	newProcedure, err := services.UpdateProcedure(isPartial, procedure)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newProcedure)
}

func getProcedureID(procedureIDParam string) (uint, error) {
	procedureID, procedureErr := strconv.ParseUint(procedureIDParam, 10, 64)
	if procedureErr != nil {
		return 0, procedureErr
	}
	return uint(procedureID), nil
}

func getPage(pageParam string) (int, error) {
	if pageParam == "" {
		return 1, nil
	}
	page, pageErr := strconv.Atoi(pageParam)
	if pageErr != nil {
		return 0, pageErr
	}
	return page, nil
}

func DeleteProcedure(c *gin.Context) {
	procedureID, idErr := getProcedureID(c.Param("procedure_id"))
	if idErr != nil {
		logrus.Error(idErr)
		c.JSON(http.StatusNotFound, gin.H{"error": idErr.Error()})
		return
	}

	procedure, getErr := services.GetProcedure(uint(procedureID))
	if getErr != nil {
		logrus.Error(getErr)
		c.JSON(http.StatusNotFound, gin.H{"error": getErr.Error()})
		return
	}

	isOwn, ownErr := isOwnProcedure(c.Request, procedure)
	if !isOwn {
		logrus.Error(ownErr)
		c.JSON(http.StatusForbidden, gin.H{"error": "not your procedure"})
		return
	}

	if err := services.DeleteProcedure(procedureID); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, procedure)
}

func isOwnProcedure(r *http.Request, procedure *procedures.Procedure) (bool, error) {
	user, err := users.GetUserFromToken(r)
	if err != nil {
		return false, err
	}
	if procedure.UserID != user.ID {
		return false, nil
	}
	return true, nil
}
