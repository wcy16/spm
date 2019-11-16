// This file contain functions for operating a database
package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"spm/config"
	"time"
)

// database connection
var conn *gorm.DB

// all of the tables registered
// we use this map to perform reflection or iterate the db tables
var table = make(map[string]interface{})

// interface to seed the database
type seeder interface {
	seed(db *gorm.DB) bool
}

// add foreign key constrain
type addForeignKeyer interface {
	addForeignKey(db *gorm.DB)
}

// delete additional tables created by many2many relations
type droper interface {
	drop(db *gorm.DB)
}

// basic database model that ignore time information when perform json marshalling
type IgnoreTimeModel struct {
	ID        uint       `gorm:"primary_key"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-";sql:"index"`
}

// path for seed files
var seedPath = "model/seed/"

// connect to the database
func Connect() {
	cfg := config.Read("mysql")
	user := cfg["user"]
	password := cfg["password"]
	address := cfg["address"]
	dbname := cfg["dbname"]

	var err error
	conn, err = gorm.Open("mysql",
		fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			user, password, address, dbname))

	if err != nil {
		panic(err)
		//Conn.Close()
	}
	//defer Conn.Close()

	log.Println("database connected")
}

// close the database
func Close() error {
	return conn.Close()
}

////////////////////////////////////////////
// General DB operations
///////////////////////////////////////////

// migrate tables without delete the database
func Migrate() {
	tx := conn.Begin()

	for k, v := range table {
		log.Printf("migrating %s.\n", k)
		tx.AutoMigrate(v)
	}

	// add foreign key constrain
	for k, v := range table {
		if s, ok := v.(addForeignKeyer); ok == true {
			s.addForeignKey(tx)
			log.Printf("add foreign key constrain %s.\n", k)
		}
	}

	tx.Commit()
	log.Println("migrate done.")
}

// drop and create all of the tables, reset the db
func Reset() {
	tx := conn.Begin()
	// disable foreign key check
	tx.Exec("SET FOREIGN_KEY_CHECKS = 0;")
	for k, v := range table {
		log.Printf("drop %s.\n", k)
		if tx.HasTable(v) {
			tx.DropTable(v)
		}
	}

	for k, v := range table {
		if s, ok := v.(droper); ok == true {
			s.drop(tx)
			log.Printf("drop tables related to %s.\n", k)
		}
	}
	// enable foreign key check
	tx.Exec("SET FOREIGN_KEY_CHECKS = 1;")
	tx.Commit()

	Migrate()

	log.Println("reset done.")
}

// reset and seed all of the tables
func Seed() {
	Reset()
	tx := conn.Begin()
	for k, v := range table {
		if s, ok := v.(seeder); ok == true {
			if s.seed(tx) {
				log.Printf("seeding %s.\n", k)
			} else {
				tx.Rollback()
				log.Fatal("error seeding " + k)
			}
		}
	}
	tx.Commit()
	log.Println("seed done.")
}
