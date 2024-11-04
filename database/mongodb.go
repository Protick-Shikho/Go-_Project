package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client     *mongo.Client
	collection *mongo.Collection
	dbName     string
	counterCol *mongo.Collection
}


func (m *MongoDB) getNextID() (int64, error) {
    filter := bson.M{"_id": "taskid"}
    update := bson.M{"$inc": bson.M{"seq": 1}}
    opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

    var result struct {
        Seq int64 `bson:"seq"`
    }
    err := m.counterCol.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&result)
    if err != nil {
        return 0, fmt.Errorf("failed to get next ID: %v", err)
    }
    return result.Seq, nil
}

// Function to insert a document with an auto-incremented ID
func (m *MongoDB) InsertTask(task string, status string) error {
    nextID, err := m.getNextID()
    if err != nil {
        return err
    }

    document := bson.M{"id": nextID, "task": task, "status": status}
    _, err = m.collection.InsertOne(context.TODO(), document)
    if err != nil {
        return fmt.Errorf("failed to insert task: %v", err)
    }
    return nil
}
// InitDB initializes the MongoDB connection using the provided DSN and database name.
func (m *MongoDB) InitDB(dsn string) error {
    var err error
    m.client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(dsn))
    if err != nil {
        return fmt.Errorf("failed to connect to MongoDB: %v", err)
    }

    // Check the connection
    if err = m.client.Ping(context.TODO(), nil); err != nil {
        return fmt.Errorf("failed to ping MongoDB: %v", err)
    }

    // Set the database name and initialize the collection
    m.dbName = "To_do_Task" // Set your database name here
    m.collection = m.client.Database(m.dbName).Collection("tasks")

    // Call SetupCounter to initialize the counter collection
    err = m.SetupCounter()
    if err != nil {
        return err
    }

    return nil
}

func (m *MongoDB) SetupCounter() error {
    // Set up the counter collection (e.g., named "counters")
    m.counterCol = m.client.Database(m.dbName).Collection("counters")

    // Initialize the counter for "tasks" collection if it doesn't exist
    filter := bson.M{"_id": "taskid"}
    update := bson.M{"$setOnInsert": bson.M{"seq": 0}}
    opts := options.Update().SetUpsert(true)

    _, err := m.counterCol.UpdateOne(context.TODO(), filter, update, opts)
    if err != nil {
        return fmt.Errorf("failed to setup counter: %v", err)
    }
    return nil
}

func (m *MongoDB) CreateTask(task Task) error {
    nextID, err := m.getNextID()
    if err != nil {
        return fmt.Errorf("failed to get next ID: %v", err)
    }

    task.ID = nextID // Set the auto-incremented ID
    task.CreatedAt = time.Now()

    _, err = m.collection.InsertOne(context.TODO(), task)
    if err != nil {
        return fmt.Errorf("failed to insert task: %v", err)
    }
    return nil
}

func (m *MongoDB) ShowTasks() ([]Task, error) {
	cursor, err := m.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("ShowTasks: %v", err)
	}
	defer cursor.Close(context.TODO())

	var tasks []Task
	for cursor.Next(context.TODO()) {
		var task Task
		if err := cursor.Decode(&task); err != nil {
			return nil, fmt.Errorf("ShowTasks (decoding): %v", err)
		}
		tasks = append(tasks, task)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("ShowTasks (cursor error): %v", err)
	}

	return tasks, nil
}

func (m *MongoDB) UpdateStatus(id int64) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": "Completed"}}

	result, err := m.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}

	if result.MatchedCount == 0 {
		log.Println("No task found with the specified ID, maybe it's completed already.")
	} else {
		log.Printf("ID %d is updated", id)
	}

	return nil
}

func (m *MongoDB) DeleteTask(id int64) error {
	filter := bson.M{"_id": id}
	result, err := m.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("DeleteTask: Error executing query: %v", err)
	}

	if result.DeletedCount == 0 {
		log.Printf("DeleteTask: No task found with ID %d", id)
	} else {
		log.Printf("DeleteTask: Task deleted successfully\n")
	}
	return nil
}

func (m *MongoDB) Close() error {
	return m.client.Disconnect(context.TODO())
}
