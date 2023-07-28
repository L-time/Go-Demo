package main

import "gorm.io/gorm"

func deletes(db *gorm.DB) {
	user := User{
		Age: "16",
	}

	db.Delete(&user)
	// DELETE from users where age = '16';

	db.Where("name = ?", "s").Delete(&user)
	// DELETE from users where name = 's' and age = '16';

	db.Delete(&user, 1)
	// DELETE from users where id = 1 and age = '16';

	db.Delete(&user, []int{1, 2, 3})
	// DELETE from users where id in (1,2,3) and age = '16';

	var users []User
	db.Unscoped().Where("age = '16'").Find(&users)
	// SELECT * FROM users WHERE age = '16';

	db.Unscoped().Delete(&user)
	// DELETE FROM users WHERE age = '16';

}
