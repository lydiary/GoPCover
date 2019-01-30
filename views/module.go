package views

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"Coverage/models"
	"strconv"
)

func GetAllModules(c *gin.Context) {
	modules := models.GetAllModules()
	c.JSON(200, gin.H{"data": modules, "result": "Success"})
}

func GetSingleModuleByModuleId(c *gin.Context) {
	moduleId, _ := strconv.Atoi(c.Param("module_id"))
	module := models.GetSingleModuleById(int32(moduleId))
	if module.ModuleName == "" {
		c.JSON(200, gin.H{"data": "", "result": "Success"})
	} else {
		c.JSON(200, gin.H{"data": module, "result": "Success"})
	}
}

func AddModule(c *gin.Context) {
	var module models.Module
	if err := c.BindJSON(&module); err != nil {
		fmt.Printf(err.Error())
		return
	}
	models.AddModule(&module)
	c.JSON(200, gin.H{"result": "Success"})
}

func DeleteModuleByModuleId(c *gin.Context) {
	moduleId, _ := strconv.Atoi(c.Param("module_id"))
	models.DeleteModuleByModuleId(int32(moduleId))
	c.JSON(200, gin.H{"result": "Success"})
}
