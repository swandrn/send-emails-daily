package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(24 * time.Hour)
	// ticker := time.NewTicker(5 * time.Second)
	db := openDbConnection()

	for {
		select {
		case <-ticker.C:
			users, err := getUsersToEmail(db)
			if err != nil {
				panic(err)
			}

			users, err = getEmailsOfUsers(db, users)
			if err != nil {
				panic(err)
			}

			for _, user := range users {
				wg.Add(1)
				go sendEmail(user.Email, user.Name)
				fmt.Printf("Sent email to %s\n", user.Email)
			}
		}
	}
}
