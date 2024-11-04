package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
    // "To_do_Task/database"


	_ "github.com/go-sql-driver/mysql"
)

type MySQLDatabase struct {
	DB *sql.DB
}

func (m *MySQLDatabase) InitDB(dsn string) error {
	var err error
	m.DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	if err := m.DB.Ping(); err != nil {
		return err
	}

	m.SetupDatabase()
	return nil
}

func (m *MySQLDatabase) SetupDatabase() {
	if _, err := m.DB.Exec("CREATE DATABASE IF NOT EXISTS To_Do_Task"); err != nil {
		log.Fatalf("Error creating database: %v", err)
	}

	if _, err := m.DB.Exec("USE abc"); err != nil {
		log.Fatalf("Error selecting database: %v", err)
	}

	createTableSQL := `
    CREATE TABLE IF NOT EXISTS tasks (
        id INT AUTO_INCREMENT PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        description TEXT NOT NULL,
        status ENUM('pending', 'completed') NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	if _, err := m.DB.Exec(createTableSQL); err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
}

func (m *MySQLDatabase) CreateTask(task Task) error {
	query := "INSERT INTO tasks (title, description, status) VALUES (?, ?, ?)"
	_, err := m.DB.Exec(query, task.Title, task.Description, task.Status)
	return err
}

func (m *MySQLDatabase) ShowTasks() ([]Task, error) {
	query := "SELECT id, title, description, status, created_at FROM tasks"
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ShowTasks: %v", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var createdAtBytes []byte

		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &createdAtBytes); err != nil {
			return nil, fmt.Errorf("ShowTasks (scanning): %v", err)
		}

		dateOnly := string(createdAtBytes)[:10]
		createdAt, err := time.Parse("2006-01-02", dateOnly)
		if err != nil {
			return nil, fmt.Errorf("ShowTasks (parsing created_at): %v", err)
		}
		task.CreatedAt = createdAt
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ShowTasks (rows error): %v", err)
	}

	return tasks, nil
}

func (m *MySQLDatabase) UpdateStatus(id int64) (error) {

	query := "UPDATE tasks SET status = 'Completed'  WHERE id = ?"

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve affected rows: %w", err)
	}
	if rowsAffected == 0 {
		log.Println("No task found with the specified ID, Maybe Its Completed Already.")
	} else {
		log.Printf("ID %d is updated", id)
	}

	return nil

}

func (m *MySQLDatabase) DeleteTask(id int64) (error) {
	
	query := "DELETE FROM tasks WHERE id = ?"
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("DeleteTask: Error executing query: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteTask: Error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		log.Printf("DeleteTask: No task found with ID %d", id)
	} else {
		log.Printf("DeleteTask: Task deleted successfully\n")
	}
    return nil
}

func (m *MySQLDatabase) Close() error {
	if m.DB != nil {
		err := m.DB.Close()
		if err != nil {
			return fmt.Errorf("failed to close the database: %w", err)
		}
	}
	return nil
}
