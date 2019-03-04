package views

import (
	"Coverage/models"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/pkg/errors"
	"math/rand"
	"path"
	"strconv"
	"strings"
	"time"
)

type UploadFileInfo struct {
	RevisionDiff string `json:"revision_diff"`
	SvnPathDiff  string `json:"svn_path_diff"`
	RelPath      string `json:"rel_path"`
	SvnPath      string `json:"svn_path"`
	Revision     string `json:"revision"`
	ModuleName   string `json:"module_name"`
}

type UploadFuncInfo struct {
	FuncName          string `json:"func_name"`
	FuncDesc          string `json:"func_desc"`
	EffectiveLineList string `json:"effective_line_list"`
	CoveredLineList   string `json:"covered_line_list"`
}

type UploadFileCoverageInfo struct {
	FileInfo       UploadFileInfo   `json:"file_info"`
	FuncDetailList []UploadFuncInfo `json:"func_detail"`
}

type UploadCoverage struct {
	CoverageDataList []UploadFileCoverageInfo `json:"coverage_data_list"`
	TestId           string                   `json:"test_id"`
	Product          string                   `json:"product"`
	Mode             string                   `json:"mode"`
}

type ConfigureJsonSubItem struct {
	Version string   `json:"version"`
	IpList  []string `json:"ip_list"`
}

type ConfigureJsonSubItems []ConfigureJsonSubItem

func (conf ConfigureJsonSubItems) Len() int {
	return len(conf)
}

func (conf ConfigureJsonSubItems) Less(i, j int) bool {
	return VersionCompare(conf[i].Version, conf[j].Version) < 0
}

func VersionCompare(version1, version2 string) int {
	version1Nos := strings.Split(version1, ".")
	version2Nos := strings.Split(version2, ".")

	if len(version1Nos) != len(version2Nos) {
		fmt.Println("版本格式错误，应该是xxx.xxx.xxx等以.隔开的格式")
		return 0
	}

	for index, version1No := range version1Nos {
		no1, err := strconv.Atoi(version1No)
		if err != nil {
			fmt.Println("版本格式错误，版本信息只能包含有效数字")
			return 0
		}

		no2, err := strconv.Atoi(version2Nos[index])
		if err != nil {
			fmt.Println("版本格式错误，版本信息只能包含有效数字")
			return 0
		}

		if no1 < no2 {
			return -1
		} else if no1 > no2 {
			return 1
		}
	}
	return 0
}

func (conf ConfigureJsonSubItems) Swap(i, j int) {
	conf[i], conf[j] = conf[j], conf[i]
}

type ConfigurationJson struct {
	TestData ConfigureJsonSubItems `json:"test_data"`
	DiffData ConfigureJsonSubItems `json:"diff_data"`
}

func CreateTestNo() string {
	t := time.Now()
	firstPart := fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
	secondPart := random(10000, 99999)
	return fmt.Sprintf("%s-%d-00", firstPart[2:], secondPart)
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func ReportModuleInfo(moduleName, version, processName, status string) (moduleId int32) {
	findModule := models.FindModule(map[string]interface{}{
		"module_name": moduleName, "version": version, "process_name": processName,
	})
	if findModule.ModuleId > 0 {
		return findModule.ModuleId
	}

	module := models.Module{
		ModuleName: moduleName, Version: version,
		ProcessName: processName, Status: status,
	}
	models.AddModule(&module)
	return module.ModuleId
}

func ReportFileInfo(uploadFileInfo *UploadFileInfo, testId, moduleId int32) (fileId int32) {
	revision, err := strconv.Atoi(uploadFileInfo.Revision)
	if err != nil {
		fmt.Println(err.Error())
		return 0;
	}
	revisionDiff, err := strconv.Atoi(uploadFileInfo.RevisionDiff)
	if err != nil {
		//revision diff不是有效的数字时，使用默认的0
		revisionDiff = 0
	}

	findFile := models.FindFile(map[string]interface{}{
		"module_name": uploadFileInfo.ModuleName, "file_name": GetFileNameFromPath(uploadFileInfo.RelPath),
		"test_id": testId, "module_id": moduleId, "file_path": uploadFileInfo.RelPath,
		"rel_path": uploadFileInfo.RelPath, "svn_path": uploadFileInfo.SvnPath,
		"revision": int16(revision), "svn_path_diff": uploadFileInfo.SvnPathDiff,
		"revision_diff": int16(revisionDiff),
	})
	if findFile.FileId > 0 {
		return findFile.FileId
	}

	file := models.File{
		ModuleName: uploadFileInfo.ModuleName, FileName: GetFileNameFromPath(uploadFileInfo.RelPath),
		TestId: testId, ModuleId: moduleId, FilePath: uploadFileInfo.RelPath,
		RelPath: uploadFileInfo.RelPath, SvnPath: uploadFileInfo.SvnPath,
		Revision: int16(revision), SvnPathDiff: uploadFileInfo.SvnPathDiff,
		RevisionDiff: int16(revisionDiff),
	}

	models.AddFile(&file)
	return file.FileId
}

func ReportFuncDetails(funcDetails []UploadFuncInfo, testId, moduleId, fileId int32, moduleName, relPath string) error {
	for _, funcDetail := range funcDetails {
		if err := ReportSingleFuncDetail(&funcDetail, testId, moduleId, fileId, moduleName, relPath); err != nil {
			return err
		}
	}
	return nil
}

func ReportSingleFuncDetail(funcDetail *UploadFuncInfo, testId, moduleId, fileId int32, moduleName, relPath string) error {
	findFunc := models.FindFunction(map[string]interface{}{
		"func_name": funcDetail.FuncName, "func_desc": funcDetail.FuncDesc, "rel_path": relPath,
		"test_id": testId, "module_id": moduleId, "file_id": fileId, "module_name": moduleName,
	})
	if findFunc.FuncId > 0 {
		findFunc.EffectiveLineList = funcDetail.EffectiveLineList
		findFunc.CoveredLineList = funcDetail.CoveredLineList
		models.UpdateFunction(&findFunc)
		return nil
	}

	function := models.Function{
		FuncName: funcDetail.FuncName, FuncDesc: funcDetail.FuncDesc, RelPath: relPath,
		EffectiveLineList: funcDetail.EffectiveLineList, CoveredLineList: funcDetail.CoveredLineList,
		TestId: testId, ModuleId: moduleId, FileId: fileId, ModuleName: moduleName,
	}
	models.AddFunction(&function)
	if function.FileId <= 0 {
		logs.Error("insert single function into database failed")
		return errors.New("insert single function into database failed")
	}
	return nil
}

func GetFileNameFromPath(filePath string) string {
	return path.Base(filePath)
}

func ReportCoverageInfo(coverageInfo *UploadCoverage) error {
	testId, err := strconv.Atoi(coverageInfo.TestId)
	if err != nil {
		logs.Error("failed to parse test id")
		return errors.New("failed to parse test id")
	}
	test := models.GetSingleTestById(int32(testId))
	for _, fileCoverageInfo := range coverageInfo.CoverageDataList {
		fileInfo := fileCoverageInfo.FileInfo
		funcDetails := fileCoverageInfo.FuncDetailList
		moduleId := ReportModuleInfo(fileInfo.ModuleName, test.Version, "", "")
		if moduleId <= 0 {
			logs.Error("failed to report module info.")
			return errors.New("failed to report module info.")
		}
		fileId := ReportFileInfo(&fileInfo, int32(testId), moduleId)
		if fileId <= 0 {
			logs.Error("failed to report file info.")
			return errors.New("failed to report file info.")
		}
		if ReportFuncDetails(funcDetails, int32(testId), moduleId, fileId, fileInfo.ModuleName, fileInfo.RelPath) != nil {
			logs.Error("failed to report function info.")
			return errors.New("failed to report function info.")
		}
	}
	return nil
}

func IntArrayToString(a []int, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

type SvnDiffItem struct {
	AnchorLineList  [][3]int `json:"anchor_line_list"`
	RelPath         string   `json:"rel_path"`
	ChangedLineList []int    `json:"changed_line_list"`
}

type SvnDiffInfo struct {
	DiffInfo []SvnDiffItem ``
}

func GetSvnDiffInfo(newVersionSvnPath, oldVersionSvnPath string) *SvnDiffInfo {
	// TODO
	var svnDiffInfo SvnDiffInfo
	return &svnDiffInfo
}

func GetRevisionFromVersion(version string) string {
	tests := models.GetTestsByVersion(version)
	if len(tests) <= 0 {
		return ""
	}

	fileInfos := models.GetFilesByTestId(tests[0].TestId)
	if len(fileInfos) <= 0 {
		return ""
	}

	return strconv.Itoa(int(fileInfos[0].Revision))
}

func GetRevisionFromTestId(testId int32) string {
	fileInfos := models.GetFilesByTestId(testId)
	if len(fileInfos) <= 0 {
		return ""
	}

	return strconv.Itoa(int(fileInfos[0].Revision))
}

func FindSvnDiffItemFromRelPath(svnDiffInfo *SvnDiffInfo, relPath string) *SvnDiffItem {
	for _, diffItem := range svnDiffInfo.DiffInfo {
		if diffItem.RelPath == relPath {
			return &diffItem
		}
	}
	return nil
}