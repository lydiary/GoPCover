package main

import (
	_ "Coverage/models"
	"Coverage/views"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	//gin.SetMode(gin.DebugMode)
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/tests", views.GetAllTests)
		v1.GET("/test/:test_id", views.GetSingleTestByTestId)
		v1.DELETE("/test/:test_id", views.DeleteTestByTestId)
		v1.POST("/test", views.AddTest)

		v1.GET("/modules", views.GetAllModules)
		v1.GET("/module/:module_id", views.GetSingleModuleByModuleId)
		v1.DELETE("/module/:module_id", views.DeleteModuleByModuleId)
		v1.POST("/module", views.AddModule)

		v1.GET("/files", views.GetAllFiles)
		v1.GET("/file/:file_id", views.GetSingleFileByFileId)
		v1.DELETE("/file/:file_id", views.DeleteFileByFileId)
		v1.POST("/file", views.AddFile)

		v1.GET("/functions", views.GetAllFunctions)
		v1.GET("/function/:function_id", views.GetSingleFunctionByFunctionId)
		v1.DELETE("/function/:function_id", views.DeleteFunctionByFunctionId)
		v1.POST("/function", views.AddFunction)

		v1.POST("/get_test_info", views.GetTestInfo)
		v1.POST("/parse_coverage_data_list", views.ParseCoverageDataList)
		v1.GET("/merge_or_diff_coverage/:test_id", views.MergeOrDiffCoverage)
	}
	log.Fatal(router.Run(":8081"))
}
