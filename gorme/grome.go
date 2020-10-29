/**
 * @Author fanguohao
 * @Description // 数据库字段编辑态同步脚本；
 * @Date 11:35 2020/9/17
 **/
package gorme

import (
	"fmt"
	//_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"reflect"
	"strings"
)

type Result struct {
	Field   string `gorm:"column:Field"`
	Type    string `gorm:"column:Type"`
	Null    string `gorm:"column:Null"`
	Key     string `gorm:"column:Key"`
	Default string `gorm:"column:Default"`
	Extra   string `gorm:"column:Extra"`
}

var (
	FieldAndType  map[string]string  //源码中的字段参数和类型
	FieldAndTypeInDB  map[string]string //数据库中的字段参数和类型
)

type ModelStruct struct {
	ModelType     reflect.Type
}

//func GetModelStruct() *ModelStruct {
//	return gorm.Scope.getModelStruct(scope, make([]*StructField, 0))
//}




func  freshTables(db *gorm.DB,scope *gorm.Scope) *gorm.Scope {
	var tags []string
	for _, field := range scope.GetModelStruct().StructFields {
		if field.IsNormal {
			sqlTag := scope.Dialect().DataTypeOf(field)
			tags = append(tags, scope.Quote(field.DBName)+" "+sqlTag)
		}
	}

	FieldAndTypeInDB = make(map[string]string)

	var result []Result
	db.Raw(fmt.Sprintf("DESCRIBE %s",scope.QuotedTableName())).Scan(&result)

	for i:=0; i<len(result);i++{
		FieldAndTypeInDB[result[i].Field] = result[i].Type
	}

	FieldAndType = make(map[string]string)

	n:= len(tags) //源码中字段数
	for i:=0;i<n;i++{
		arr := strings.Fields(tags[i])
		filed :=arr[0][1:len(arr[0])-1]
		FieldAndType[filed]= arr[1]
	}

	// 字段增加同步；
	for key,value := range FieldAndType{
		if _,ok := FieldAndTypeInDB[key];ok{
			typeCheck(scope,key,value)
		}else {
			println(key+"在数据库中不存在"+"开始写入数据库！")
			scope.Raw(fmt.Sprintf("alter table %s add (%s %s)",scope.QuotedTableName(),key,value)).Exec()
		}
	}

	// 字段删除同步；
	for key,_ := range FieldAndTypeInDB{
		if value,ok := FieldAndType[key];ok{
			typeCheck(scope,key,value)
		}else{
			println(key+"在数据库中冗余,"+"开始删除字段！")
			scope.Raw(fmt.Sprintf("alter table %s drop column %s",scope.QuotedTableName(),key)).Exec()
			//if ok:=strings.HasSuffix(key,"del"); !ok{
			//	println(key+"在数据库中冗余,"+"开始删除字段！")
			//	// 如果[key]_del 已经存在，对[key]_del真删除；
			//	if _,ok:= fieldAndTypeInDB[key+"_del"];ok{
			//		scope.Raw(fmt.Sprintf("alter table %s drop column %s",scope.QuotedTableName(),key)).Exec()
			//	}
			//	scope.Raw(fmt.Sprintf("alter table %s change %s %s %s",scope.QuotedTableName(),key,key+"_del",fieldAndTypeInDB[key])).Exec()
			//}
		}
	}

	////类型检测
	//for key,value := range FieldAndType{
	//	if _,ok := FieldAndTypeInDB[key];ok{
	//		ok:=strings.Compare(FieldAndTypeInDB[key],FieldAndType[key])==0
	//		if !ok{
	//			println(key+"在数据库的类型为："+FieldAndTypeInDB[key]+";在源码中的类型为："+FieldAndType[key])
	//			println("开始删除旧类型字段，新加新类型字段！")
	//			//scope.Raw(fmt.Sprintf("alter table %s change %s %s %s",scope.QuotedTableName(),key,key+"_del",fieldAndTypeInDB[key])).Exec()
	//			scope.Raw(fmt.Sprintf("alter table %s drop column %s",scope.QuotedTableName(),key)).Exec()
	//			scope.Raw(fmt.Sprintf("alter table %s add (%s %s)",scope.QuotedTableName(),key,value)).Exec()
	//		}
	//	}
	//}

	return scope
}

/*
  @param  key: 源码字段名；
  @param  value: 源码字段类型；
 */

func  typeCheck(scope *gorm.Scope, key string, value string) {
	ok:=strings.Compare(strings.ToLower(FieldAndTypeInDB[key]),strings.ToLower(FieldAndType[key]))==0
	if !ok{
		if value == "boolean"&& (FieldAndTypeInDB[key] == "tinyint(1)") {

		}else if strings.Contains(value,"int") && (strings.Contains(FieldAndTypeInDB[key],"int")){

		}else{
			println(key+"在数据库的类型为："+FieldAndTypeInDB[key]+";在源码中的类型为："+FieldAndType[key]+"开始类型同步操作！")
			//scope.Raw(fmt.Sprintf("alter table %s change %s %s %s",scope.QuotedTableName(),key,key+"_del",fieldAndTypeInDB[key])).Exec()
			scope.Raw(fmt.Sprintf("alter table %s drop column %s",scope.QuotedTableName(),key)).Exec()
			scope.Raw(fmt.Sprintf("alter table %s add (%s %s)",scope.QuotedTableName(),key,value)).Exec()
			FieldAndType[key]=value     // 类型更新
			FieldAndTypeInDB[key]=value
		}

	}
}


func FreshDB(DB *gorm.DB,models ...interface{})  {
	println("开始源码与数据库字段同步操作！")
	db := DB.Unscoped()
	for _,model := range models{
		DbScope:= db.NewScope(model)
		freshTables(DB,DbScope)
	}
	println("字段同步完成！")
}

