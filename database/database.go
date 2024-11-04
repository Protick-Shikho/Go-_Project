package database

// import "time"

// type Task struct {
//     ID          int
//     Title       string
//     Description string
//     Status      string
//     CreatedAt   time.Time
// }

type Database interface {
    InitDB(dataSourceName string) error
    CreateTask(task Task) error
    ShowTasks() ([]Task, error)
    UpdateStatus(id int64) error
    DeleteTask(id int64) error
    Close() error
}
