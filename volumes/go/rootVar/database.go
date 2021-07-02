package rootVar

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DbConn *gorm.DB

type DBSetting struct {
	DB_ACCOUNT  string `json:"DB_ACCOUNT"`
	DB_PASSWORD string `json:"DB_PASSWORD"`
	DB_HOST     string `json:"DB_HOST"`
	DB_NAME     string `json:"DB_NAME"`
	DB_PORT     string `json:"DB_PORT"`
}

func InitDatabaseOrm() *sql.DB {
	var err error
	// 打开文件
	file, err := ioutil.ReadFile("./config/database.json")

	if err != nil {
		fmt.Println("can't load db config file ", err)
	}

	db := DBSetting{}
	_ = json.Unmarshal([]byte(file), &db)

	//"root:34erDFcv@tcp(mariaDb:3306)/stocks?charset=utf8mb4&parseTime=True&loc=Asia%2fTaipei"
	dsn := db.DB_ACCOUNT + ":" + db.DB_PASSWORD + "@tcp(" + db.DB_HOST + ":" + db.DB_PORT + ")/" + db.DB_NAME + "?charset=utf8mb4&parseTime=True&loc=Asia%2fTaipei"
	//fmt.Println(dsn)
	DbConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Database Connect Failed (DB)", err)
	}

	sqlConnect, err := DbConn.DB()
	if err != nil {
		fmt.Println("Database Connect Failed (sqlConnect)", err)
	}

	return sqlConnect
}
