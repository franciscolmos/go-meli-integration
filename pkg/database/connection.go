package database
import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB() gorm.DB{
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root:@tcp(127.0.0.1:3306)/meli_app?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return *db
}
