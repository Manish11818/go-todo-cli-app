package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // todo, in-progress, done
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

const dataFile = "tasks.json"

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Usage: task-cli <command> [arguments]")
		return
	}

	command := args[1]
	tasks := loadTasks()

	switch command {
	case "add":
		if len(args) < 3 {
			fmt.Println("Usage: task-cli add \"task description\"")
			return
		}
		desc := args[2]
		id := getNextID(tasks)
		now := time.Now()
		task := Task{ID: id, Description: desc, Status: "todo", CreatedAt: now, UpdatedAt: now}
		tasks = append(tasks, task)
		saveTasks(tasks)
		fmt.Printf("Task added successfully (ID: %d)\n", id)

	case "update":
		if len(args) < 4 {
			fmt.Println("Usage: task-cli update <id> \"new description\"")
			return
		}
		id, _ := strconv.Atoi(args[2])
		found := false
		for i := range tasks {
			if tasks[i].ID == id {
				tasks[i].Description = args[3]
				tasks[i].UpdatedAt = time.Now()
				found = true
				break
			}
		}
		if found {
			saveTasks(tasks)
			fmt.Println("Task updated.")
		} else {
			fmt.Println("Task not found.")
		}

	case "delete":
		if len(args) < 3 {
			fmt.Println("Usage: task-cli delete <id>")
			return
		}
		id, _ := strconv.Atoi(args[2])
		newTasks := []Task{}
		found := false
		for _, t := range tasks {
			if t.ID != id {
				newTasks = append(newTasks, t)
			} else {
				found = true
			}
		}
		if found {
			saveTasks(newTasks)
			fmt.Println("Task deleted.")
		} else {
			fmt.Println("Task not found.")
		}

	case "mark-in-progress":
		changeStatus(tasks, args, "in-progress")

	case "mark-done":
		changeStatus(tasks, args, "done")

	case "list":
		if len(args) == 2 {
			printTasks(tasks)
		} else {
			filter := args[2]
			filtered := []Task{}
			for _, t := range tasks {
				if t.Status == filter {
					filtered = append(filtered, t)
				}
			}
			printTasks(filtered)
		}

	default:
		fmt.Println("Unknown command.")
	}
}

// 🔁 Helper: Change task status
func changeStatus(tasks []Task, args []string, newStatus string) {
	if len(args) < 3 {
		fmt.Printf("Usage: task-cli mark-%s <id>\n", newStatus)
		return
	}
	id, _ := strconv.Atoi(args[2])
	found := false
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Status = newStatus
			tasks[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}
	if found {
		saveTasks(tasks)
		fmt.Println("Task updated.")
	} else {
		fmt.Println("Task not found.")
	}
}

// 🔄 Load tasks from JSON file
func loadTasks() []Task {
	var tasks []Task
	file, err := os.ReadFile(dataFile)
	if err != nil {
		// If file doesn't exist, create it
		_ = os.WriteFile(dataFile, []byte("[]"), 0644)
		return tasks
	}
	json.Unmarshal(file, &tasks)
	return tasks
}

// 💾 Save tasks to file
func saveTasks(tasks []Task) {
	data, _ := json.MarshalIndent(tasks, "", "  ")
	_ = os.WriteFile(dataFile, data, 0644)
}

// ➕ Get next ID
func getNextID(tasks []Task) int {
	max := 0
	for _, t := range tasks {
		if t.ID > max {
			max = t.ID
		}
	}
	return max + 1
}

// 📋 Print tasks
func printTasks(tasks []Task) {
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}
	for _, t := range tasks {
		fmt.Printf("[%d] %s | %s | Created: %s | Updated: %s\n",
			t.ID, t.Description, strings.ToUpper(t.Status),
			t.CreatedAt.Format("2006-01-02 15:04"),
			t.UpdatedAt.Format("2006-01-02 15:04"),
		)
	}
}
