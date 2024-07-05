package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type User struct {
	Id    int
	Name  string
	Email string
}

func openDbConnection() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	cfg := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PWD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DB_ADDR"),
		DBName:               os.Getenv("DB_NAME"),
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		panic(pingErr)
	}
	fmt.Println("Connected to database!")
	return db
}

func getUsersToEmail(db *sql.DB) ([]User, error) {
	var users []User

	rows, err := db.Query("SELECT DISTINCT utilisateur_id FROM taches WHERE DATEDIFF(CURRENT_DATE(), date_echeance) > 0 && is_completed = 0")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func getEmailsOfUsers(db *sql.DB, users []User) ([]User, error) {
	var usersToEmail []User

	for _, user := range users {
		rows, err := db.Query("SELECT nom, email FROM users WHERE id = ?", user.Id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(&user.Name, &user.Email); err != nil {
				return nil, err
			}
			usersToEmail = append(usersToEmail, user)
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}
	return usersToEmail, nil
}
