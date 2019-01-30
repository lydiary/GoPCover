package views

import (
	"Coverage/models"
	"fmt"
	"github.com/gin-gonic/gin"
)

func GetTestInfo(c *gin.Context) {
	var test models.Test
	if err := c.BindJSON(&test); err != nil {
		c.JSON(200, gin.H{"result": "failed", "reason": "parameter error"})
		return
	}

	if test.IP == "" || test.Version == "" {
		c.JSON(200, gin.H{"result": "failed", "reason": "version or ip is empty"})
		return
	}

	findTest := models.FindTest(map[string]interface{}{"version": test.Version, "ip": test.IP})
	if findTest.TestNo == "" {
		test.TestNo = CreateTestNo()
		models.AddTest(&test)
		findTest = test
	}
	c.JSON(200, gin.H{"result": "success", "data": map[string]interface{}{"test_id": findTest.TestId, "test_no": findTest.TestNo}})
}

func ParseCoverageDataList(c *gin.Context) {
	var coverageDataList UploadCoverage
	if err := c.BindJSON(&coverageDataList); err != nil {
		fmt.Println(err.Error())
		c.JSON(200, gin.H{"result": "failed", "reason": "upload data format error."})
		return
	}
	//fmt.Println(coverageDataList)
	if err := ReportCoverageInfo(&coverageDataList); err != nil {
		c.JSON(200, gin.H{"result": "failed", "reason": err.Error()})
		return
	}
	//go ReportCoverageInfo(&coverageDataList)
	c.JSON(200, gin.H{"result": "success"})
}

func MergeOrDiffCoverage(c *gin.Context) {

}