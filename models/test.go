package models

type Test struct {
	TestId            int32  `gorm:"type:bigint(20) auto_increment;primary_key;not null"`
	TestNo            string `gorm:"type:varchar(25)" json:"test_no"`
	Author            string `gorm:"type:varchar(255)" json:"author"`
	Version           string `gorm:"type:varchar(255)" json:"version"`
	VersionDiff       string `gorm:"type:varchar(255)" json:"version_diff"`
	ContrastVersion   string `gorm:"type:varchar(255)" json:"contrast_version"`
	TestType          string `gorm:"type:enum('all', 'diff');DEFAULT:'all'" json:"test_type"`
	CoverageType      string `gorm:"enum('diff','all','ip','merge','other');DEFAULT:'ip'" json:"coverage_type"`
	ProjectType       string `gorm:"enum('debug','release');DEFAULT:'release'" json:"project_type"`
	VersionType       string `gorm:"enum('version','upgrade');DEFAULT:'version'" json:"version_type"`
	Status            string `gorm:"enum('new','preprocess','running','finished','exception');DEFAULT:'new'" json:"status"`
	Description       string `gorm:"type:varchar(255)" json:"description"`
	Remark            string `gorm:"type:varchar(1023)" json:"remark"`
	Reason            string `gorm:"type:varchar(255)" json:"reason"`
	CoveredLineNum    int16  `gorm:"type:int(11)" json:"covered_line_num"`
	EffectiveLineNum  int16  `gorm:"type:int(11)" json:"effective_line_num"`
	ConfigurationJson string `gorm:"type:longtext" json:"configuration_json"`
	IP                string `gorm:"type:varchar(32)" json:"ip"`
	NeedUpdate        bool   `gorm:"type:tinyint(1)" json:"need_update"`
}

func (Test) TableName() string {
	return "t01_test"
}

type Tests []Test

func GetAllTests() Tests {
	var tests Tests
	dbInstance.Find(&tests)
	return tests
}

func GetSingleTestById(id int32) Test {
	var test Test
	dbInstance.Where("test_id = ?", id).First(&test)
	return test
}

func AddTest(test *Test) {
	dbInstance.Create(&test)
}

func DeleteTestById(id int32) {
	var test Test
	dbInstance.Where("test_id = ?", id).First(&test)
	dbInstance.Delete(&test)
}

func FindTest(condition map[string]interface{}) Test {
	queryString := ""
	spliter := ""
	var args []interface{}
	for key, value := range condition {
		queryString += spliter + key + " = ? "
		args = append(args, value)
		spliter = "AND "
	}
	var test Test
	dbInstance.Where(queryString, args...).First(&test)
	return test
}

func GetTestsByVersion(version string) Tests {
	var tests Tests
	dbInstance.Find(&tests, "version = ?", version)
	return tests
}

func GetTestByVersionIp(version, ip string) Test {
	var test Test
	dbInstance.Find(&test, "version = ? and ip = ?", version, ip)
	return test
}

func UpdateTest(test *Test) {
	dbInstance.Save(test)
}