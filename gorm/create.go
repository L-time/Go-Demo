package main

import (
	"fmt"
	"gorm.io/gorm"
)

func create(db *gorm.DB) {
	user := User{
		Name:     "Nick",
		Age:      "19",
		NickName: "AAA",
	}

	users := []*User{
		&User{
			Name:     "A",
			Age:      "15",
			NickName: "a",
		},
		&User{
			Name:     "B",
			Age:      "16",
			NickName: "b",
		},
	}

	result := db.Create(&user)
	//如果想要判断创建结果是否成功，只需要调用result.Error即可
	if result.Error != nil {
		fmt.Println("创建失败")
	}
	//返回记录的ID
	fmt.Println("Id = ", user.ID)
	//返回插入记录的条数
	fmt.Println("Rows = ", result.RowsAffected)

	//批量创建
	result = db.Create(&users)
	if result.Error != nil {
		fmt.Println("创建失败")
	}

	//返回插入记录的条数
	fmt.Println("Rows = ", result.RowsAffected)
}
