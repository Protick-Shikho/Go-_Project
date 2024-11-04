package cmd

import (
	"To_do_Task/database"
	"bufio"
	"fmt"
	"log"
	"strconv"
	"time"

	"os"
	"strings"

	"github.com/spf13/cobra"
)

var db database.Database

func SetDatabase(database database.Database) {
	db = database
}

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new task",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Title: ")
		title, _ := reader.ReadString('\n')
		title = strings.TrimSpace(title)

		fmt.Print("Description: ")
		description, _ := reader.ReadString('\n')
		description = strings.TrimSpace(description)

		newTask := database.Task{
			Title:       title,
			Description: description,
			Status:      "pending",
            CreatedAt: time.Now(),
		}

		err := db.CreateTask(newTask)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("New task created")
	},
}

var ShowCmd = &cobra.Command{
    Use:   "show",
    Short: "Show all tasks",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Current Tasks:")
        tasks, err := db.ShowTasks()
        if err != nil {
            log.Fatalf("Error showing tasks: %v", err)
        }
        for _, task := range tasks {
            fmt.Printf("ID: %d, Title: %s, Description: %s, Status: %s, Created at: %s\n",
                task.ID, task.Title, task.Description, task.Status, task.CreatedAt.Format("2006-01-02"))
        }
    },
}


var UpdateCmd = &cobra.Command{
    Use:   "update [id]",
    Short: "Update a task's status",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {

        idToUpdate, err := strconv.Atoi(args[0]) 
        if err != nil {
            log.Fatalf("Invalid task ID: %v", err)
        }
        
        db.UpdateStatus(int64(idToUpdate))

    },
}
var DeleteCmd = &cobra.Command{
    Use:   "delete [id]",
    Short: "Delete a task",
    Run: func(cmd *cobra.Command, args []string) {
        idToDelete, err := strconv.Atoi(args[0]) 
        if err != nil {
            log.Fatalf("Invalid task ID: %v", err)
        }
        
        db.DeleteTask(int64(idToDelete)) //cannot use idToDelete (variable of type int) as int64 value in argument to db.DeleteTask
    },
}
