package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var x int

	fmt.Println("Create a Task (Enter 0)")
	fmt.Println("Show Task (Enter 1)")
	fmt.Println("Update Task Status (Enter 2)")
	fmt.Println("Delete Task (Enter 3)")
	fmt.Println("Exit (Enter 4)")
	fmt.Scan(&x)

	if x == 0 {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Title: ")
		title, _ := reader.ReadString('\n')
		title = strings.TrimSpace(title)

		fmt.Print("Description: ")
		description, _ := reader.ReadString('\n')
		description = strings.TrimSpace(description)

		newTask := Task{
			Title:       title,
			Description: description,
			Status:      "pending",
		}

		err := CreateTask(db, newTask)
		if err != nil {
			log.Fatal(err)
			return
		}

		fmt.Println("New task created")

	} else if x == 1 {
		fmt.Println("Current Tasks:")
		ShowTasks(db)

	} else if x == 2 {
		ShowTasks(db)
		fmt.Println("Which task do you want to Update? Enter the id")
		var idToUpdate int
		fmt.Scan(&idToUpdate)
		UpdateStatus(db, idToUpdate)

	} else if x == 3 {
		ShowTasks(db)
		fmt.Println("Which task do you want to Delete? Enter the id")
		var idToDelete int
		fmt.Scan(&idToDelete)
		DeleteTask(db, idToDelete)

	} else if x == 4 {
		fmt.Println("Exiting....")
	}
}
