#gome操作说明
功能：
---
    1.编辑态同步增加和删除字段;
    2.编辑态数据类型的同步(删除旧字段和数据，新增新字段）；
##示例：
``````
type Cfg struct {
	Username    string `gorm:"username";"PRIMARY_KEY"`
	Pwd			string
	Address     string
	DBName		string
	Gender      string
}
type person struct {
	Name 	string
	Gender string
}

gorme.FreshDB(db,&Cfg{},&person{}) //db 是gorm连接成功的返回类型；
``````
##注意：
    同步操作只同步字段，不会保存历史数据，请谨慎操作！<br>
	test
	

