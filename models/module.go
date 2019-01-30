package models

type Module struct {
	ModuleId    int32  `gorm:"type:bigint(20) AUTO_INCREMENT;PRIMARY_KEY;NOT NULL"`
	ModuleName  string `gorm:"type:varchar(255)" json:"module_name"`
	ProcessName string `gorm:"type:varchar(255)" json:"process_name"`
	TestId      int32  `gorm:"type:int(16)" json:"test_id"`
	Version     string `gorm:"type:varchar(255)" json:"version"`
	Status      string `gorm:"type:enum('new','running','finished','exception');DEFAULT:'new'" json:"status"`
}

type Modules []Module

func (Module) TableName() string {
	return "t02_module"
}

func GetAllModules() Modules {
	var modules Modules
	dbInstance.Find(&modules)
	return modules
}

func GetSingleModuleById(id int32) Module {
	var module Module
	dbInstance.Where("module_id = ?", id).First(&module)
	return module
}

func AddModule(module *Module) {
	dbInstance.Create(&module)
}

func DeleteModuleByModuleId(id int32) {
	var module Module
	dbInstance.Where("module_id = ?", id).First(&module)
	dbInstance.Delete(&module)
}

func FindModule(condition map[string]interface{}) Module {
	queryString := ""
	spliter := ""
	var args []interface{}
	for key, value := range condition {
		queryString += spliter + key + " = ? "
		args = append(args, value)
		spliter = "AND "
	}
	var module Module
	dbInstance.Where(queryString, args...).First(&module)
	return module
}
