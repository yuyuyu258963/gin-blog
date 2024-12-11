package models

import (
	"fmt"
	"gin_example/pkg/setting"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

// 可以直接用 gorm.Model 来直接替代了，因为其中包含了这些常用的字段
type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
	DeletedOn  int `json:"deleted_on"`
}

func init() {
	var (
		err                                               error
		dbType, dbName, user, password, host, tablePrefix string
	)

	sec, err := setting.Cfg.GetSection("database")
	if err != nil {
		log.Fatal(2, "Fail to get section database: %v", err)
	}

	dbType = sec.Key("TYPE").MustString("mysql")
	dbName = sec.Key("NAME").MustString("NAME")
	user = sec.Key("USER").String()
	password = sec.Key("PASSWORD").String()
	host = sec.Key("HOST").String()
	tablePrefix = sec.Key("TABLE_PREFIX").String()

	// 打开与数据库的连接
	db, err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user,
		password,
		host,
		dbName,
	))

	if err != nil {
		log.Println(err)
	}
	log.Printf("success open database %s/%s", host, dbName)

	// 设置默认表名前缀
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix + defaultTableName
	}

	db.SingularTable(true)
	db.LogMode(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// 替换为自己的回调函数
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)
}

func CloseDB() {
	defer db.Close()
}

// 这么做的好处就是后面只要嵌入了Model，那么就可以触发对应的回调函数

// updateTimeStampForCreateCallback will set `CreatedOn` `ModifiedOn` when creating
// gorm.Scope 其实就类似之前实现的GeeORM的Schema
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		// 创建时间
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
			if createTimeField.IsBlank {
				createTimeField.Set(nowTime)
			}
		}

		// 修改时间
		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(nowTime)
			}
		}
	}
}

// updateTimeStampForUpdateCallback will set `ModifiedOn` when updating
func updateTimeStampForUpdateCallback(scop *gorm.Scope) {
	if _, ok := scop.FieldByName("ModifiedOn"); ok {
		scop.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

// 删除记录时的回调函数
func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")

		// 获取约定的删除字段，若存在则UPDATE软删除，若不存在则DELETE硬删除
		// scope.Search.Unscoped 是 GORM 中用于禁用全局作用域
		// （global scope）的一个设置。它主要有两个重要用途：
		if !scope.Search.Unscoped && hasDeletedOnField {
			// 若存在则软删除
			scope.Raw(fmt.Sprintf(
				"UPDATE %v Set %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(time.Now().Unix()),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			))
		} else {
			// 不存在则硬删除
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			))
		}
	}
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return "" + str
	}
	return ""
}
