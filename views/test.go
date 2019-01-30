package views

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"Coverage/models"
	"strconv"
)

func GetAllTests(c *gin.Context) {
	tests := models.GetAllTests()
	c.JSON(200, gin.H{"data": tests, "result": "Success"})
}

func GetSingleTestByTestId(c *gin.Context) {
	testId, _ := strconv.Atoi(c.Param("test_id"))
	test := models.GetSingleTestById(int32(testId))
	if test.TestNo == "" {
		c.JSON(200, gin.H{"data": "", "result": "Success"})
	} else {
		c.JSON(200, gin.H{"data": test, "result": "Success"})
	}
}

func AddTest(c *gin.Context) {
	var test models.Test
	if err := c.BindJSON(&test); err != nil {
		fmt.Printf(err.Error())
		return
	}
	models.AddTest(&test)
	c.JSON(200, gin.H{"result": "Success"})
}

func DeleteTestByTestId(c *gin.Context) {
	testId, _ := strconv.Atoi(c.Param("test_id"))
	models.DeleteTestById(int32(testId))
	c.JSON(200, gin.H{"result": "Success"})
}
