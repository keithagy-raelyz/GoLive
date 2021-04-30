package db

import (
	"database/sql"
	"log"
)

type Database struct {
	b *sql.DB
}

// Initialize the DB schema
func (d *Database) InitializeDB(db *sql.DB) {
	d.b = db
	Query1 := `CREATE TABLE IF NOT EXISTS Users (
    UserID int NOT NULL AUTO_INCREMENT,
    Username VARCHAR(255) NOT NULL,
    Password VARCHAR(255) NOT NULL,
    Email varchar(255) NOT NULL,
    PRIMARY KEY (UserID)
	)`
	_, err := db.Exec(Query1)
	if err != nil {
		log.Fatal(err)
	}
	Query2 := `
	CREATE TABLE IF NOT EXISTS Merchants (
		MerchantID int NOT NULL AUTO_INCREMENT,
		Username VARCHAR(255) NOT NULL,
		Password VARCHAR(255) NOT NULL,
		Email varchar(255) NOT NULL,
		Description VARCHAR(255) NOT NULL,
		PRIMARY KEY (MerchantID)
	);`
	_, err = db.Exec(Query2)
	if err != nil {
		log.Fatal(err)
	}
	Query3 := `CREATE TABLE IF NOT EXISTS Products (
    ProductID int NOT NULL AUTO_INCREMENT,
    Product_Name VARCHAR(255) NOT NULL,
    Quantity int NOT NULL,
    Image varchar(255) NOT NULL,
    Price float not null,
    Description VARCHAR(255),
    MerchantID int NOT NULL,
    Foreign Key (MerchantID) REFERENCES Merchants (MerchantID),
    PRIMARY KEY (ProductID)
	);`
	_, err = db.Exec(Query3)
	if err != nil {
		log.Fatal(err)
	}
}
