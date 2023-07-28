package main

import "gorm.io/gorm"

func update(db *gorm.DB) {
	user := User{}

	db.First(&user)

	//拿到记录后我们直接更改记录即可
	user.Name = "Luna"
	db.Save(&user)

	//有一个特性，如果你传入的结构体内没有包含主键的话，那么此时Save会调用Create方法

	userWithoutId := User{
		Name: "123",
	}
	//这里便是Create方法，相当于SQL的INSERT
	db.Save(&userWithoutId)

	userWithId := User{Model: gorm.Model{ID: 1}, Name: "s"}
	//这里便是Save方法，相当于SQL的UPDATE
	db.Save(&userWithId)
}
