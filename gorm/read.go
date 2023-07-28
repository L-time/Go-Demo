package main

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func read(db *gorm.DB) {
	user := User{}

	// 获取第一条记录（主键升序）
	db.First(&user)
	// SELECT * FROM users ORDER BY id LIMIT 1;
	fmt.Println(user)

	// 获取一条记录，没有指定排序字段
	db.Take(&user)
	// SELECT * FROM users LIMIT 1;
	fmt.Println(user)

	// 获取最后一条记录（主键降序）
	db.Last(&user)
	// SELECT * FROM users ORDER BY id DESC LIMIT 1;
	fmt.Println(user)

	result := db.First(&user)
	fmt.Println(result.RowsAffected) // 返回找到的记录数
	if result.Error != nil {         // returns error or nil
		fmt.Println(result.Error)
	}

	// 检查 ErrRecordNotFound 错误
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("找不到记录")
	}

	var users []User

	// Get first matched record
	db.Where("name = ?", "Nick").First(&user)
	// SELECT * FROM users WHERE name = 'Nick' ORDER BY id LIMIT 1;
	fmt.Println(user)

	// Get all matched records
	db.Where("name <> ?", "A").Find(&users)
	// SELECT * FROM users WHERE name <> 'A';
	fmt.Println(users)

	// IN
	db.Where("name IN ?", []string{"A", "B"}).Find(&users)
	// SELECT * FROM users WHERE name IN ('A','B');
	fmt.Println(users)

	// LIKE
	db.Where("name LIKE ?", "%Ni%").Find(&users)
	// SELECT * FROM users WHERE name LIKE '%Ni%';
	fmt.Println(users)

	// AND
	db.Where("name = ? AND age >= ?", "Nick", "10").Find(&users)
	// SELECT * FROM users WHERE name = 'Nick' AND age >= 10;
	fmt.Println(users)

	// BETWEEN
	db.Where("age BETWEEN ? AND ?", "5", "15").Find(&users)
	// SELECT * FROM users WHERE age BETWEEN '5' AND '15';
	fmt.Println(users)

}
