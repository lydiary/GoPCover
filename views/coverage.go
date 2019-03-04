package views

import (
	"Coverage/models"
	"Coverage/redisclient"
	"Coverage/utils"
	"encoding/json"
	"fmt"
	"github.com/deckarep/golang-set"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"sort"
	"strconv"
	"strings"
	"time"
)

var RedisClient redisclient.Client

func init() {
	settings := utils.GetSettings()
	addr := fmt.Sprintf("%s:%d",
		settings.RedisSettings.Host, settings.RedisSettings.Port)
	RedisClient.NewClient(&redis.Options{
		Addr:       addr,
		PoolSize:   100,
		MaxRetries: 2,
		Password:   "",
		DB:         0,
	})
}

type RedisDataFormat struct {
	UpdateTime   time.Time      `json:"update_time"`
	IsUpdated    bool           `json:"is_updated"`
	CoverageInfo UploadCoverage `json:"coverage_info"`
}

func SaveCoverageInfoToRedis(redisData RedisDataFormat) error {
	testId, _ := strconv.Atoi(redisData.CoverageInfo.TestId)
	key := GetRedisKey(testId)
	return RedisClient.SetValue(key, redisData)
	//fmt.Printf("saving coverage info to redis, result: %v, err: %v\n", result, err)
}

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

	// 检测test_id是否有效
	testId, _ := strconv.Atoi(coverageDataList.TestId)
	test := models.GetSingleTestById(int32(testId))
	if test.TestNo == "" {
		c.JSON(200, gin.H{"data": "Invalid test id", "result": "Failed"})
		return
	}

	//fmt.Println(coverageDataList)
	//优化: 将数据保存到redis中，定时将redis中的数据写入到数据库
	//if err := ReportCoverageInfo(&coverageDataList); err != nil {
	//	c.JSON(200, gin.H{"result": "failed", "reason": err.Error()})
	//	return
	//}

	redisData := RedisDataFormat{UpdateTime: time.Now(), IsUpdated: true, CoverageInfo: coverageDataList}
	if err := SaveCoverageInfoToRedis(redisData); err != nil {
		c.JSON(200, gin.H{"Data": "Save data to redis failed.", "result": "Failed"})
		return
	}
	c.JSON(200, gin.H{"result": "Success"})
}

func MergeOrDiffCoverage(testId int32) (*RedisDataFormat, error) {
	// 根据传进来的testId从test表中找到配置信息，如果找到了就处理
	// 但是如果传进来的testId并非是通过页面配置产生的话，暂不处理
	test := models.GetSingleTestById(testId)
	if test.TestNo == "" || test.ConfigurationJson == "" {
		// 不存在该testId或者不是通过页面配置产生的
		return nil, errors.New("Test ID错误")
	}

	configuration := GetConfigurationJsonFromTest(test)
	fmt.Println("Configuration: ", configuration)
	if len(configuration.TestData) <= 0 {
		return nil, errors.New("配置文件格式错误，测试版本数据为空")
	}

	// 下面计算合并或者差异的覆盖率信息
	var calculatedCovInfo *RedisDataFormat

	// 如果没有比较的版本信息，是merge操作，否则是diff操作
	if len(configuration.DiffData) <= 0 {
		test.TestType = "merge"
		calculatedCovInfo = MergeCoverageInfo(configuration)
	} else {
		test.TestType = "diff"
		calculatedCovInfo = DiffCoverageInfo(configuration)
	}
	models.UpdateTest(&test)
	return calculatedCovInfo, nil
}

func DiffCoverageInfo(configuration *ConfigurationJson) *RedisDataFormat {
	sort.Sort(configuration.DiffData)
	sort.Sort(configuration.TestData)

	// 先合并测试版本的覆盖率数据，然后和对比版本中最后一个版本比较
	mergedCoverageInfo := MergeCoverageInfo(configuration)
	diffResult := diffCoverageInfoInternal(mergedCoverageInfo, configuration.DiffData[len(configuration.DiffData)-1])
	return diffResult
}

func diffCoverageInfoInternal(data *RedisDataFormat, diffVersion ConfigureJsonSubItem) *RedisDataFormat {
	// TODO
	return nil
}

func MergeCoverageInfo(configuration *ConfigurationJson) *RedisDataFormat {
	sort.Sort(configuration.TestData)

	var mergedCoverageInfo *RedisDataFormat
	for index, versionInfo := range configuration.TestData {
		versionCoverage := MergeCoverageInfoForSameVersion(&versionInfo)
		if index == 0 {
			mergedCoverageInfo = versionCoverage
			continue
		}

		changedVersionCover := ChangeCoverInfoFromOlderVersionToNewer(mergedCoverageInfo, configuration.TestData[index-1].Version)
		*mergedCoverageInfo = MergeSameVersionCoverageInfoForTwoIps(*versionCoverage, *changedVersionCover)
	}
	return mergedCoverageInfo
}

func ChangeCoverInfoFromOlderVersionToNewer(oldVersionCoverInfo *RedisDataFormat, newerVersion string) *RedisDataFormat {
	// 数据是空的，什么都不做
	if len(oldVersionCoverInfo.CoverageInfo.CoverageDataList) <= 0 {
		return oldVersionCoverInfo
	}

	newTests := models.GetTestsByVersion(newerVersion)
	if len(newTests) <= 0 {
		return oldVersionCoverInfo
	}
	newTestId := newTests[0].TestId

	svnPath := oldVersionCoverInfo.CoverageInfo.CoverageDataList[0].FileInfo.SvnPath
	relPath := oldVersionCoverInfo.CoverageInfo.CoverageDataList[0].FileInfo.RelPath
	svnPathBase := svnPath[0:len(svnPath)-len(relPath)] + "service@"
	oldVersionSvnPath := svnPathBase + oldVersionCoverInfo.CoverageInfo.CoverageDataList[0].FileInfo.Revision
	newerVersionRevision := GetRevisionFromTestId(newTestId)
	if newerVersionRevision == "" {
		return oldVersionCoverInfo
	}
	newerVersionSvnPath := svnPathBase + newerVersionRevision

	// 读取新版本的覆盖率数据
	newVersionCoverInfo := ReadCoverageFromStorage(newTestId)
	if newVersionCoverInfo == nil {
		return oldVersionCoverInfo
	}

	svnDiffInfo := GetSvnDiffInfo(newerVersionSvnPath, oldVersionSvnPath)
	//svnDiffInfo = &SvnDiffInfo{DiffInfo: []SvnDiffItem{SvnDiffItem{
	//	AnchorLineList: [][3]int{{1, 1, 0}, {2, 2, 0}, {3, 3, -1}, {4, 3, 1}, {4, 4, 0}, {5, 5, -1}, {6, 5, 0}, {7, 6, 0}, {8, 7, 0}, {9, 8, 1}, {9, 9, 0}, {10, 10, 0}},
	//	RelPath:"yes",
	//}}}
	// 根据diff的结果转换行
	oldVersionCoverInfo = changeCoverInfoLinesFromSvnDiffInfo(oldVersionCoverInfo, newVersionCoverInfo, svnDiffInfo)
	return oldVersionCoverInfo
}

func changeCoverInfoLinesFromSvnDiffInfo(oldVersionCoverInfo,
	newVersionCoverInfo *RedisDataFormat, svnDiffInfo *SvnDiffInfo) *RedisDataFormat {

	// 先将新版本的File Info填充到老版本中
	for index, fileCoverageInfo := range oldVersionCoverInfo.CoverageInfo.CoverageDataList {
		findNewFileCovInfo := findFileInfoByRelPath(newVersionCoverInfo, fileCoverageInfo.FileInfo.RelPath)
		oldVersionCoverInfo.CoverageInfo.CoverageDataList[index].FileInfo = findNewFileCovInfo.FileInfo
	}

	for _, svnDiffItem := range svnDiffInfo.DiffInfo {
		findOldFileCovInfo := findFileInfoByRelPath(oldVersionCoverInfo, svnDiffItem.RelPath)
		findNewFileCovInfo := findFileInfoByRelPath(newVersionCoverInfo, svnDiffItem.RelPath)
		if findOldFileCovInfo == nil || findNewFileCovInfo == nil {
			continue
		}

		for index, funcDetail := range findNewFileCovInfo.FuncDetailList {
			findOldFuncCovInfo := FindFuncInfoFromFileCoverageInfo(
				findOldFileCovInfo, funcDetail.FuncName, funcDetail.FuncDesc)

			// 老版本中不存在这个函数
			if findOldFileCovInfo == nil {
				continue
			}

			findOldFuncCovInfo.EffectiveLineList = funcDetail.EffectiveLineList
			findOldFileCovInfo.FuncDetailList[index].CoveredLineList =
				changeLineFromOldToNewByDiffInfo(findOldFuncCovInfo.CoveredLineList, &svnDiffItem)
		}
	}
	return oldVersionCoverInfo
}

func changeLineFromOldToNewByDiffInfo(lines string, svnDiffItem *SvnDiffItem) string {
	lineNos := strings.Split(lines, " ")
	maxLine, _ := strconv.Atoi(lineNos[len(lineNos)-1])
	lineMap := changeAnchorLineListToLineMap(svnDiffItem.AnchorLineList, maxLine)

	var lineList IntSlice
	for _, line := range lineNos {
		lineInt, _ := strconv.Atoi(line)
		newLine := lineMap[lineInt]
		if newLine == 0 { // anchor_line_list中该行已被删除
			continue
		}
		lineList = append(lineList, newLine)
	}
	sort.Sort(lineList)
	return IntArrayToString(lineList, " ")
}

func changeAnchorLineListToLineMap(anchorLineList [][3]int, maxLine int) map[int]int {
	recordMap := make(map[int]int)

	for _, anchorLine := range anchorLineList {
		if anchorLine[2] == 0 {
			recordMap[anchorLine[0]] = anchorLine[1]
		} else if anchorLine[2] == -1 {
			recordMap[anchorLine[0]] = 0
		}
	}

	step := 0
	for i := 1; i <= maxLine; i++ {
		if val, ok := recordMap[i]; !ok {
			recordMap[i] = i + step
		} else {
			step = val - i
		}
	}

	return recordMap
}

func findFileInfoByRelPath(coverInfo *RedisDataFormat, relPath string) *UploadFileCoverageInfo {
	for _, item := range coverInfo.CoverageInfo.CoverageDataList {
		if item.FileInfo.RelPath == relPath {
			return &item
		}
	}
	return nil
}

func MergeCoverageInfoForSameVersion(versionInfo *ConfigureJsonSubItem) *RedisDataFormat {
	// 如果ip list为空，合并该版本对应的所有的ip版本
	if len(versionInfo.IpList) <= 0 {
		tests := models.GetTestsByVersion(versionInfo.Version)
		var ipList []string
		for _, test := range tests {
			ipList = append(ipList, test.IP)
		}
		versionInfo.IpList = ipList
	}
	return MergeCoverageInfoForListIpsOfVersion(versionInfo)
}

func MergeCoverageInfoForAllIpsOfVersion(version string) *RedisDataFormat {
	var mergedCoverageInfo RedisDataFormat

	tests := models.GetTestsByVersion(version)
	for _, test := range tests {
		coverInfo := ReadCoverageFromStorage(test.TestId)
		mergedCoverageInfo = MergeSameVersionCoverageInfoForTwoIps(mergedCoverageInfo, *coverInfo)
	}
	return &mergedCoverageInfo
}

func ReadCoverageFromStorage(testId int32) *RedisDataFormat {
	coverInfo := ReadCoverageInfoFromRedis(testId)
	if coverInfo == nil {
		coverInfo = ReadCoverageInfoFromDatabase(testId)
		if SaveCoverageInfoToRedis(*coverInfo) != nil {
			fmt.Println("Save data to redis failed.")
		}
	}
	return coverInfo
}

func MergeSameVersionCoverageInfoForTwoIps(coverInfoOne RedisDataFormat, coverInfoAnother RedisDataFormat) RedisDataFormat {
	if len(coverInfoOne.CoverageInfo.CoverageDataList) <= 0 {
		return coverInfoAnother
	}

	if len(coverInfoAnother.CoverageInfo.CoverageDataList) <= 0 {
		return coverInfoOne
	}

	// 两个列表合并时，合并one和another中都存在的数据，以one为基础
	for outerIndex, coverOne := range coverInfoOne.CoverageInfo.CoverageDataList {
		coverOneFileInfo := coverOne.FileInfo
		coverAnotherFileCoverInfo := FindFileCoverInfoFromFileInfo(&coverInfoAnother, &coverOneFileInfo)
		if coverAnotherFileCoverInfo == nil {
			continue
		}

		for innerIndex, funcInfo := range coverOne.FuncDetailList {
			anotherFuncInfo := FindFuncInfoFromFileCoverageInfo(coverAnotherFileCoverInfo,
				funcInfo.FuncName, funcInfo.FuncDesc)
			if anotherFuncInfo == nil { // 没有找到该函数信息
				continue
			}

			mergedEffectiveLineList := MergeLineList(funcInfo.EffectiveLineList, anotherFuncInfo.EffectiveLineList)
			mergedCoveredLineList := MergeLineList(funcInfo.CoveredLineList, anotherFuncInfo.CoveredLineList)

			coverInfoOne.CoverageInfo.CoverageDataList[outerIndex].FuncDetailList[innerIndex].EffectiveLineList = mergedEffectiveLineList
			coverInfoOne.CoverageInfo.CoverageDataList[outerIndex].FuncDetailList[innerIndex].CoveredLineList = mergedCoveredLineList
		}
	}

	// 两个列表合并时，将another中独有的信息保存到one中
	for outerIndex, coverAnother:= range coverInfoAnother.CoverageInfo.CoverageDataList {
		coverAnotherFileInfo := coverAnother.FileInfo
		coverOneFileCoverInfo := FindFileCoverInfoFromFileInfo(&coverInfoOne, &coverAnotherFileInfo)
		if coverOneFileCoverInfo == nil {
			coverInfoOne.CoverageInfo.CoverageDataList = append(coverInfoOne.CoverageInfo.CoverageDataList, coverAnother)
			continue
		}

		for _, funcInfo := range coverAnother.FuncDetailList {
			oneFuncInfo := FindFuncInfoFromFileCoverageInfo(coverOneFileCoverInfo,
				funcInfo.FuncName, funcInfo.FuncDesc)
			if oneFuncInfo == nil { // 没有找到该函数信息,添加到结果函数集中
				coverInfoOne.CoverageInfo.CoverageDataList[outerIndex].FuncDetailList =
					append(coverInfoOne.CoverageInfo.CoverageDataList[outerIndex].FuncDetailList, funcInfo)
				continue
			}
		}
	}
	return coverInfoOne
}

func MergeLineList(lineList1, lineList2 string) string {
	line1 := strings.Split(lineList1, " ")
	line2 := strings.Split(lineList2, " ")

	mergedLineSet := mapset.NewSet()

	for _, line := range line1 {
		lineNo, _ := strconv.Atoi(line)
		mergedLineSet.Add(lineNo)
	}

	for _, line := range line2 {
		lineNo, _ := strconv.Atoi(line)
		mergedLineSet.Add(lineNo)
	}

	var mergedLineList IntSlice
	for _, val := range mergedLineSet.ToSlice() {
		mergedLineList = append(mergedLineList, val.(int))
	}
	sort.Sort(mergedLineList)
	return IntArrayToString(mergedLineList, " ")
}

type IntSlice []int

func (v IntSlice) Len() int           { return len(v) }
func (v IntSlice) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v IntSlice) Less(i, j int) bool { return v[i] < v[j] }

func FindFileCoverInfoFromFileInfo(coverInfo *RedisDataFormat, fileInfo *UploadFileInfo) *UploadFileCoverageInfo {
	for _, cover := range coverInfo.CoverageInfo.CoverageDataList {
		if cover.FileInfo == *fileInfo {
			return &cover
		}
	}
	return nil
}

func FindFuncInfoFromFileCoverageInfo(fileCoverageInfo *UploadFileCoverageInfo, funcName, funcDesc string) *UploadFuncInfo {
	for _, funcInfo := range fileCoverageInfo.FuncDetailList {
		if funcInfo.FuncName == funcName && funcInfo.FuncDesc == funcDesc {
			return &funcInfo
		}
	}
	return nil
}

func MergeCoverageInfoForListIpsOfVersion(versionInfo *ConfigureJsonSubItem) *RedisDataFormat {
	var mergedCoverageInfo RedisDataFormat
	isFirstTime := true

	for _, ip := range versionInfo.IpList {
		test := models.GetTestByVersionIp(versionInfo.Version, ip)

		coverInfo := ReadCoverageInfoFromRedis(test.TestId)
		if coverInfo == nil {
			coverInfo = ReadCoverageInfoFromDatabase(test.TestId)
			if SaveCoverageInfoToRedis(*coverInfo) != nil {
				fmt.Println("Save data to redis failed.")
			}
		}

		if isFirstTime {
			isFirstTime = false
			mergedCoverageInfo = *coverInfo
		} else {
			mergedCoverageInfo = MergeSameVersionCoverageInfoForTwoIps(mergedCoverageInfo, *coverInfo)
		}
	}
	return &mergedCoverageInfo
}

func GetConfigurationJsonFromTest(test models.Test) *ConfigurationJson {
	var configuration ConfigurationJson
	_ = json.Unmarshal([]byte(test.ConfigurationJson), &configuration)
	return &configuration
}

func ReadCoverageInfoFromRedis(testId int32) *RedisDataFormat {
	key := GetRedisKey(int(testId))
	serializedValue, err := RedisClient.GetValue(key)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	var redisData RedisDataFormat
	err = json.Unmarshal([]byte(serializedValue), &redisData)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return &redisData
}

func ReadCoverageInfoFromDatabase(testId int32) *RedisDataFormat {
	test := models.GetSingleTestById(testId)
	if test.TestNo == "" { // 数据库中不存在该testId对应的数据
		return nil
	}

	var redisData RedisDataFormat
	redisData.CoverageInfo.TestId = strconv.Itoa(int(testId))
	redisData.CoverageInfo.CoverageDataList = ReadCoverageByTestId(testId)
	return &redisData
}

func GetRedisKey(testId int) string {
	return "test-" + strconv.Itoa(testId)
}

func ReadCoverageByTestId(testId int32) []UploadFileCoverageInfo {
	var fileCoverages []UploadFileCoverageInfo

	for _, file := range models.GetFilesByTestId(testId) {
		var covInfo UploadFileCoverageInfo
		covInfo.FileInfo = GetFileInfoOfFile(&file)
		covInfo.FuncDetailList = GetFuncDetailOfFile(&file)
		fileCoverages = append(fileCoverages, covInfo)
	}
	return fileCoverages
}

func GetFileInfoOfFile(file *models.File) UploadFileInfo {
	return UploadFileInfo{
		RevisionDiff: strconv.Itoa(int(file.RevisionDiff)), SvnPathDiff: file.SvnPathDiff,
		RelPath: file.RelPath, SvnPath: file.SvnPath,
		Revision: strconv.Itoa(int(file.Revision)), ModuleName: file.ModuleName,
	}
}

func GetFuncDetailOfFile(file *models.File) []UploadFuncInfo {
	var funcDetails []UploadFuncInfo
	for _, function := range models.GetFunctionsByFileId(file.FileId) {
		funcInfo := UploadFuncInfo{
			FuncName: function.FuncName, FuncDesc: function.FuncDesc,
			EffectiveLineList: function.EffectiveLineList,
			CoveredLineList:   function.CoveredLineList,
		}
		funcDetails = append(funcDetails, funcInfo)
	}
	return funcDetails
}
