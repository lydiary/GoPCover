package views

import (
	"Coverage/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetAllFiles(c *gin.Context) {
	files := models.GetAllFiles()
	c.JSON(200, gin.H{"data": files, "result": "Success"})
}

func GetSingleFileByFileId(c *gin.Context) {
	FileId, _ := strconv.Atoi(c.Param("file_id"))
	File := models.GetSingleFileById(int32(FileId))
	if File.FileName == "" {
		c.JSON(200, gin.H{"data": "", "result": "Success"})
	} else {
		c.JSON(200, gin.H{"data": File, "result": "Success"})
	}
}

func AddFile(c *gin.Context) {
	var File models.File
	if err := c.BindJSON(&File); err != nil {
		fmt.Printf(err.Error())
		return
	}
	models.AddFile(&File)
	c.JSON(200, gin.H{"result": "Success"})
}

func DeleteFileByFileId(c *gin.Context) {
	FileId, _ := strconv.Atoi(c.Param("file_id"))
	models.DeleteFileByFileId(int32(FileId))
	c.JSON(200, gin.H{"result": "Success"})
}
