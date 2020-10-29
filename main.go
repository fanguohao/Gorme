package  main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"gorm_extern/gorme"
	"time"
)
type SchemaHistory struct {
	Id          int64
	Version     string
	Description string
	Types       string
	Script      string
	CheckSum    string
	InstalledOn time.Time
	Success     bool
}
type Cfg struct {
	Username    string `gorm:"-";"PRIMARY_KEY"`
	Pwd			string
	Address     int32
	DBName		int32
	Gender     	string  `gorm:"gen"`
}

type person struct {
	Name 	string
	Gender  string
}

func main()  {
	db,err := gorm.Open("mysql","root:Fan003174@/elmcms?charset=utf8")
	defer db.Close()
	if err!=nil{
		fmt.Println("数据库连接失败！")
		return
	}
	//db.AutoMigrate(&Cfg{},&person{})
	gorme.FreshDB(db,&Cfg{})
	//db.CreateTable(&Cfg{})
}



