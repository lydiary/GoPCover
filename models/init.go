package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var dbInstance *gorm.DB

func init() {
	var err error
	dbInstance, err = gorm.Open("mysql", "root:1qaz@WSX@tcp(localhost:3306)/d_coverage?charset=utf8&parseTime=true")
	if err != nil {
		fmt.Println(err.Error())
		panic("can not open database.")
	}

	dbInstance.AutoMigrate(&Test{})
	dbInstance.AutoMigrate(&Module{})
	dbInstance.AutoMigrate(&File{})
	dbInstance.AutoMigrate(&Function{})
}
