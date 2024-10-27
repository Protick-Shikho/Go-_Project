package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func CreateTask(db *sql.DB, task Task) error {
	query := "INSERT INTO tasks (title, description, status, created_at) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, task.Title, task.Description, task.Status, time.Now())

	if err != nil {
		return fmt.Errorf("CreateTask: %v", err)
	}

	return nil
}

func ShowTasks(db *sql.DB) {
	query := "SELECT id, title, description, status, created_at FROM tasks"
	rows, err := db.Query(query)

	if err != nil {
		log.Fatalf("ShowTasks: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		var createdAtBytes []byte

		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &createdAtBytes); err != nil {
			log.Fatalf("ShowTasks (scanning): %v", err)
		}

		dateOnly := string(createdAtBytes)[:10]

		createdAt, err := time.Parse("2006-01-02", dateOnly)
		if err != nil {
			log.Fatalf("ShowTasks (parsing created_at): %v", err)
		}
		task.CreatedAt = createdAt		

		fmt.Printf("ID: %d, Title: %s, Description: %s, Status: %s, Created at: %s\n", task.ID, task.Title, task.Description, task.Status, task.CreatedAt.Format("2006-01-02"))
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("ShowTasks (rows error): %v", err)
	}
}

func DeleteTask(db *sql.DB, id int) {
	query := "DELETE FROM tasks WHERE id = ?"
	result, err := db.Exec(query, id)
	if err != nil {
		fmt.Printf("DeleteTask: Error executing query: %v\n", err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("DeleteTask: Error getting rows affected: %v\n", err)
		return
	}

	if rowsAffected == 0 {
		fmt.Printf("DeleteTask: No task found with ID %d\n", id)
	} else {
		fmt.Printf("DeleteTask: Task deleted successfully\n")
	}
}

func UpdateStatus(db *sql.DB, id int) error {
	query := "UPDATE tasks SET status = 'Completed' WHERE id = ?"

	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve affected rows: %w", err)
	}
	if rowsAffected == 0 {
		log.Println("No task found with the specified ID, maybe it's completed already.")
	} else {
		log.Print("Updated to status Completed")
	}

	return nil
}
