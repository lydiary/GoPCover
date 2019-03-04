package main

import (
	_ "Coverage/models"
	"Coverage/utils"
	"Coverage/views"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	"log"
	"strings"
	"time"
)

var settings utils.Settings

func init() {
	settings = utils.GetSettings()
}

func main() {
	go crontabService()

	webService()
}

func webService() {
	gin.SetMode(gin.DebugMode)
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
		//v1.GET("/merge_or_diff_coverage/:test_id", views.MergeOrDiffCoverage)
	}
	log.Fatal(router.Run(":8081"))
}

func crontabService() {
	c := cron.New()

	spec := durationToCronSpec(settings.CrontabSettings.RedisPersistenceFrequency)
	c.AddFunc(spec, saveRedisDataToMysql)

	spec = durationToCronSpec(settings.CrontabSettings.RedisCheckCacheFrequency)
	c.AddFunc(spec, clearRedisCache)

	c.Start()
}

func durationToCronSpec(duration string) string {
	spec := ""
	if strings.HasSuffix(duration, "m") {
		spec = "0 0/" + duration[:len(duration)-1] + " * * * *"
	} else if strings.HasSuffix(duration, "h") {
		spec = "0 0 0/" + duration[:len(duration)-1] + " * * *"
	} else if strings.HasSuffix(duration, "s") {
		spec = "0/" + duration[:len(duration)-1] + " * * * * *"
	} else if strings.HasSuffix(duration, "d") {
		spec = "0 0 0 0/" + duration[:len(duration)-1] + " * *"
	}
	return spec
}

func saveRedisDataToMysql() {
	for _, key := range views.RedisClient.Keys("test-*") {
		var redisData views.RedisDataFormat
		serializedData, _ := views.RedisClient.GetValue(key)
		json.Unmarshal([]byte(serializedData), &redisData)

		if redisData.IsUpdated {
			fmt.Println("saving data to database: ", redisData)
			redisData.IsUpdated = false
			views.RedisClient.SetValue(key, redisData)
			views.ReportCoverageInfo(&redisData.CoverageInfo)
		}
	}
}

func clearRedisCache() {
	for _, key := range views.RedisClient.Keys("test-*") {
		var redisData views.RedisDataFormat
		serializedData, _ := views.RedisClient.GetValue(key)
		json.Unmarshal([]byte(serializedData), &redisData)

		//计算超时时间
		now := time.Now()
		expire := now.Sub(redisData.UpdateTime)
		mostDuration, _ := time.ParseDuration(settings.CrontabSettings.RedisClearCacheTimeout)

		if !redisData.IsUpdated && mostDuration < expire {
			fmt.Println("clearing unupdated cache: ", redisData)
			views.RedisClient.DeleteValue(key)
		}
	}
}
