package config

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

func NewSQLite(configuration Config) *sql.DB {

	// Setup DB
	driver := configuration.Get("DB_CONNECTION")
	databaseName := configuration.Get("DB_DATABASE")
	dsn := fmt.Sprintf("./%v.db", databaseName)
	db, err := sql.Open(driver, dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// users := `DROP TABLE users;`
	// user_details := `DROP TABLE user_details;`
	// users_email := `DROP TABLE email_verifications;`

	// users_email := `CREATE TABLE email_verifications(
	// 	"ID" INT,
	// 	"Email" TEXT,
	// 	"Signature" TEXT,
	// 	"Expired" INT);`

	// users := `CREATE TABLE users (
	//     id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	//     "Name" TEXT,
	//     "Username" TEXT,
	//     "Email" TEXT,
	//     "Password" INT,
	// 	"Role"  INT,s
	//     "Phone" TEXT,
	//     "Gender" INT,
	//     "DisabilityType" INT,
	// 	"Address" TEXT,
	//     "Birthdate" INT,
	//     "Image" TEXT,
	// 	"email_verification" INT,
	//     "Description" TEXT,
	// 	"created_at" INT,
	// 	"updated_at" INT);`

	// user_details := `CREATE TABLE user_details (
	// 		"user_id" INT,
	// 		"phone" TEXT,
	// 		"gender" INT,
	// 		"type_of_disability" INT,
	// 		"address" TEXT,
	// 		"birthdate" Birthdate,
	// 		"image" TEXT,
	// 		"description" TEXT);`

	// query1, err := db.Prepare(user_details)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// query1.Exec()
	// fmt.Println("Table created successfully!")

	// query2, err := db.Prepare(users_email)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// query2.Exec()
	// fmt.Println("Table created successfully!")

	// query, err := db.Prepare(users)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// query.Exec()
	// fmt.Println("Table created successfully!")

	// Limit connection with db pooling
	setMaxIdleConns, err := strconv.Atoi(configuration.Get("SQLITE_POOL_MIN"))
	if err != nil {
		panic(err)
	}
	setMaxOpenConns, err := strconv.Atoi(configuration.Get("SQLITE_POOL_MAX"))
	if err != nil {
		panic(err)
	}
	setConnMaxIdleTime, err := strconv.Atoi(configuration.Get("SQLITE_MAX_IDLE_TIME_SECOND"))
	if err != nil {
		panic(err)
	}
	setConnMaxLifetime, err := strconv.Atoi(configuration.Get("SQLITE_MAX_LIFE_TIME_SECOND"))
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(setMaxIdleConns)                                    // minimal connection
	db.SetMaxOpenConns(setMaxOpenConns)                                    // maximal connection
	db.SetConnMaxLifetime(time.Duration(setConnMaxIdleTime) * time.Minute) // unused connections will be deleted
	db.SetConnMaxIdleTime(time.Duration(setConnMaxLifetime) * time.Minute) // connection that can be used
	return db
}
