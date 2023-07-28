package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	user     = "root"
	password = "123456"
	addr     = "127.0.0.1:3306"
	dbs      = "test"
)

type User struct {
	gorm.Model
	Name     string
	Age      string
	NickName string
}

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, addr, dbs)
	//db 便是我们的数据库对象
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("连接失败")
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		fmt.Println("自动迁移失败")
	}

	create(db)
	read(db)
	update(db)
	deletes(db)

}
