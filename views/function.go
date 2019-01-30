package views

import (
	"Coverage/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetAllFunctions(c *gin.Context) {
	Functions := models.GetAllFunctions()
	c.JSON(200, gin.H{"data": Functions, "result": "Success"})
}

func GetSingleFunctionByFunctionId(c *gin.Context) {
	FunctionId, _ := strconv.Atoi(c.Param("function_id"))
	Function := models.GetSingleFunctionById(int32(FunctionId))
	if Function.FuncName == "" {
		c.JSON(200, gin.H{"data": "", "result": "Success"})
	} else {
		c.JSON(200, gin.H{"data": Function, "result": "Success"})
	}
}

func AddFunction(c *gin.Context) {
	var Function models.Function
	if err := c.BindJSON(&Function); err != nil {
		fmt.Printf(err.Error())
		return
	}
	models.AddFunction(&Function)
	c.JSON(200, gin.H{"result": "Success"})
}

func DeleteFunctionByFunctionId(c *gin.Context) {
	FunctionId, _ := strconv.Atoi(c.Param("function_id"))
	models.DeleteFunctionByFuncId(int32(FunctionId))
	c.JSON(200, gin.H{"result": "Success"})
}
