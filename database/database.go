package database

import (
	"fmt"
	"log"

	"github.com/RTradeLtd/Temporal/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var UploadObj *models.Upload
var UserObj *models.User
var PaymentObj *models.Payment
var IpnsObj *models.IPNS

type DatabaseManager struct {
	DB     *gorm.DB
	Upload *models.UploadManager
}

func Initialize(dbPass, dbURL, dbUser string) (*DatabaseManager, error) {
	dbm := DatabaseManager{}
	db, err := OpenDBConnection(dbPass, dbURL, dbUser)
	if err != nil {
		return nil, err
	}
	dbm.DB = db
	dbm.RunMigrations()
	return &dbm, nil
}

func (dbm *DatabaseManager) RunMigrations() {
	dbm.DB.AutoMigrate(UploadObj)
	dbm.DB.AutoMigrate(UserObj)
	dbm.DB.AutoMigrate(PaymentObj)
	// gorm will default table to name of ip_ns
	// so we will override with ipns
	check := dbm.DB.Table("ipns").CreateTable(IpnsObj)
	// get string representation of the error
	errString := fmt.Sprintf("%s", check.Error)
	// if the error is not nil, and the message isn't that table already exists, we can continue
	if check.Error != nil && errString != `pq: relation "ipns" already exists` {
		fmt.Printf("error encountered %T\n", check.Error)
		log.Fatal(check.Error)
	}
	dbm.DB.AutoMigrate(IpnsObj)
	//dbm.DB.Model(userObj).Related(uploadObj.Users)
}

// OpenDBConnection is used to create a database connection
func OpenDBConnection(dbPass, dbURL, dbUser string) (*gorm.DB, error) {
	if dbUser == "" {
		dbUser = "postgres"
	}
	// look into whether or not we wil disable sslmode
	dbConnURL := fmt.Sprintf("host=%s port=5432 user=%s dbname=temporal password=%s", dbURL, dbUser, dbPass)
	db, err := gorm.Open("postgres", dbConnURL)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func OpenTestDBConnection(dbPass string) (*gorm.DB, error) {
	dbConnURL := fmt.Sprintf("host=127.0.0.1 port=5432 user=postgres dbname=temporal password=%s", dbPass)
	db, err := gorm.Open("postgres", dbConnURL)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// CloseDBConnection is used to close a db
func CloseDBConnection(db *gorm.DB) {
	db.Close()
}
