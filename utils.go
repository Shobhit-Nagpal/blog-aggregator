package main

import (
	"fmt"

	"github.com/Shobhit-Nagpal/blog-aggregator/internal/db"
)

func printUsers(users []db.User, currentUser string) {
	for _, user := range users {
		if currentUser == user.Name {
			fmt.Printf("* %s (current)\n", user.Name)
      continue
		}

		fmt.Printf("* %s\n", user.Name)
	}
}
