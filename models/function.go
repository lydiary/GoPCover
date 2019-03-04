package models

type Function struct {
	FuncId              int32  `gorm:"type:bigint(20) AUTO_INCREMENT;PRIMARY_KEY;NOT NULL;"`
	FuncName            string `gorm:"type:text" json:"func_name"`
	FuncDesc            string `gorm:"type:varchar(255)" json:"func_desc"`
	TestId              int32  `gorm:"type:int(11)" json:"test_id"`
	ModuleId            int32  `gorm:"type:bigint(20)" json:"module_id"`
	ModuleName          string `gorm:"type:varchar(255)" json:"module_name"`
	FileId              int32  `gorm:func"type:bigint(20)" json:"file_id"`
	RelPath             string `gorm:"type:varchar(255)" json:"rel_path"`
	StartLine           int32  `gorm:"type:int(11)" json:"start_line"`
	EndLine             int32  `gorm:"type:int(11)" json:"end_line"`
	EffectiveLineList   string `gorm:"type:text" json:"effective_line_list"`
	CoveredLineList     string `gorm:"type:text" json:"covered_line_list"`
	DiffLineList        string `gorm:"type:text" json:"diff_line_list"`
	CoveredDiffLineList string `gorm:"type:text" json:"covered_diff_line_list"`
	EffectiveLineNum    int32  `gorm:"type:int(11)" json:"effective_line_num"`
	CoveredLineNum      int32  `gorm:"type:int(11)" json:"covered_line_num"`
	DiffLineNum         int32  `gorm:"type:int(11)" json:"diff_line_num"`
	CoveredDiffLineNum  int32  `gorm:"type:int(11)" json:"covered_diff_line_num"`
}

func (Function) TableName() string {
	return "t05_func"
}

type Functions []Function

func GetAllFunctions() Functions {
	var functions Functions
	dbInstance.Find(&functions)
	return functions
}

func GetFunctionsByFileId(fileId int32) Functions {
	var functions Functions
	dbInstance.Find(&functions, "file_id = ?", fileId)
	return functions
}

func GetSingleFunctionById(id int32) Function {
	var function Function
	dbInstance.Where("func_id = ?", id).First(&function)
	return function
}

func AddFunction(function *Function) {
	dbInstance.Create(&function)
}

func DeleteFunctionByFuncId(id int32) {
	var function Function
	dbInstance.Where("file_id = ?", id).First(&function)
	dbInstance.Delete(&function)
}

func FindFunction(condition map[string]interface{}) Function {
	queryString := ""
	spliter := ""
	var args []interface{}
	for key, value := range condition {
		queryString += spliter + key + " = ? "
		args = append(args, value)
		spliter = "AND "
	}
	var function Function
	dbInstance.Where(queryString, args...).First(&function)
	return function
}

func UpdateFunction(function *Function) {
	dbInstance.Save(function)
}
