package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var dB *gorm.DB

func InitDB(dataSource string) *gorm.DB {
	var err error
	dB, err = gorm.Open("mysql", dataSource)
	if err != nil {
		log.Printf("Gorm New DB Error:%v", err)
		return nil
	}
	return dB
}
