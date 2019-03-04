package models

type File struct {
	FileId             int32  `gorm:"type:bigint(20) AUTO_INCREMENT;PRIMARY_KEY;NOT NULL;"`
	FileName           string `gorm:"type:varchar(255)" json:"file_name"`
	TestId             int32  `gorm:"type:int(11)" json:"test_id"`
	ModuleId           int32  `gorm:"type:bigint(20)" json:"module_id"`
	ModuleName         string `gorm:"type:varchar(255)" json:"module_name"`
	FilePath           string `gorm:"type:varchar(255)" json:"file_path"`
	RelPath            string `gorm:"type:varchar(255)" json:"rel_path"`
	SvnPath            string `gorm:"type:varchar(255)" json:"svn_path"`
	Revision           int16  `gorm:"type:int(11)" json:"revision"`
	SvnPathDiff        string `gorm:"type:varchar(255)" json:"svn_path_diff"`
	RevisionDiff       int16  `gorm:"type:int(11)" json:"revision_diff"`
	EffectiveLineCount int16  `gorm:"type:int(11)" json:"effective_line_count"`
	CoveredLineCount   int16  `gorm:"type:int(11)" json:"covered_line_count"`
}

func (File) TableName() string {
	return "t04_file"
}

type Files []File

func GetAllFiles() Files {
	var files Files
	dbInstance.Find(&files)
	return files
}

func GetSingleFileById(id int32) File {
	var file File
	dbInstance.Where("file_id = ?", id).First(&file)
	return file
}

func GetFilesByTestId(testId int32) Files {
	var files Files
	dbInstance.Find(&files, "test_id = ?", testId)
	return files
}

func AddFile(file *File) {
	dbInstance.Create(&file)
}

func DeleteFileByFileId(id int32) {
	var file File
	dbInstance.Where("file_id = ?", id).First(&file)
	dbInstance.Delete(&file)
}

func FindFile(condition map[string]interface{}) File {
	queryString := ""
	spliter := ""
	var args []interface{}
	for key, value := range condition {
		queryString += spliter + key + " = ? "
		args = append(args, value)
		spliter = "AND "
	}
	var file File
	dbInstance.Where(queryString, args...).First(&file)
	return file
}
