package db

import (
	"database/sql"
	"e-library/models"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB(syncDB bool) {
	// dbDriver := nacos.DBDriver
	// dbUser := nacos.DBUser
	// dbPass := nacos.DBPassword
	// dbHost := nacos.DBHost
	// dbPort := nacos.DBPort
	// dbName := nacos.DBName

	dbDriver := "mysql"
	dbUser := "admin"
	dbPass := "123456"
	dbHost := "localhost"
	dbPort := "3306"
	dbName := "cloud"

	// Construct data source name (DSN)
	dsn := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8"

	// Register MySQL database driver
	orm.RegisterDriver("mysql", orm.DRMySQL)

	// Register default database
	orm.RegisterDataBase("default", dbDriver, dsn)

	// Open a database connection
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		logs.Error("[InitDB] Open DB fail")
		return
	}

	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		logs.Error("[InitDB] Ping DB fail")
		return
	}
	aliasName := "default"

	err = orm.SetDataBaseTZ(aliasName, time.Local)
	if err != nil {
		logs.Error(err)
	}

	if syncDB {
		err = orm.RunSyncdb(aliasName, false, true)
		if err != nil {
			logs.Error(err)
		}
	}

	// Check and insert initial data
	insertInitialBooks()

	if err != nil {
		logs.Info("[InitDB] Init DB Success")
	}

	// orm.Debug = true
}

func GetDB() *sql.DB {
	return db
}

// insertInitialBooks checks if the table is empty and inserts default books.
func insertInitialBooks() {
	o := orm.NewOrm()

	// Check if the RentingBook table is empty
	var count int
	err := o.Raw("SELECT COUNT(*) FROM book WHERE Deleted = 0").QueryRow(&count)
	if err != nil {
		logs.Error("[insertInitialBooks] Failed to count books: %v", err)
		return
	}

	if count == 0 {
		logs.Info("[insertInitialBooks] No books found, inserting default records")

		books := []models.RentingBook{
			{Title: "Hello", AvailableCopies: 5, MaxCopies: 5, CommStruct: models.CommStruct{CreateTime: uint64(time.Now().Unix())}},
			{Title: "Clean Code", AvailableCopies: 3, MaxCopies: 3, CommStruct: models.CommStruct{CreateTime: uint64(time.Now().Unix())}},
			{Title: "Book 1", AvailableCopies: 10, MaxCopies: 10, CommStruct: models.CommStruct{CreateTime: uint64(time.Now().Unix())}},
			{Title: "Book 2", AvailableCopies: 15, MaxCopies: 15, CommStruct: models.CommStruct{CreateTime: uint64(time.Now().Unix())}},
			{Title: "Book 3", AvailableCopies: 20, MaxCopies: 20, CommStruct: models.CommStruct{CreateTime: uint64(time.Now().Unix())}},
		}

		for _, book := range books {
			_, err := o.Insert(&book)
			if err != nil {
				logs.Error("[insertInitialBooks] Failed to insert book %s: %v", book.Title, err)
			}
		}

		logs.Info("[insertInitialBooks] Default books inserted successfully")
	} else {
		logs.Info("[insertInitialBooks] Books already exist, skipping insert")
	}
}
