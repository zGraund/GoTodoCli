package models

import (
	"database/sql"
	"time"

	database "github.com/zGraund/TodoCli/internal/db"
)

type Todo struct {
	ID uint

	Name        string
	CreatedAt   time.Time `gorm:"autoCreateTime:false"`
	Completed   bool
	CompletedAt sql.NullTime
}

var db = database.Get()

// Get all Todos for a given day
func GetByDay(t time.Time, offset int) ([]Todo, error) {
	todos := []Todo{}
	dayOffset := time.Hour * 24 * time.Duration(offset)
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Add(dayOffset)
	end := start.Add(24*time.Hour - time.Second)
	result := db.Where("created_at BETWEEN ? AND ?", start, end).Find(&todos)
	return todos, result.Error
}

func Create(n string, offset int) (Todo, error) {
	dayOffset := time.Hour * 24 * time.Duration(offset)
	date := time.Now().Add(dayOffset)
	todo := Todo{
		Name:      n,
		Completed: false,
		CreatedAt: date,
	}
	result := db.Create(&todo)
	return todo, result.Error
}

func (t Todo) FilterValue() string { return "[ ] " + t.Name }

func (t Todo) Title() string {
	if t.Completed {
		return "[✓] " + t.Name
	}
	return "[ ] " + t.Name
}

func (t Todo) Description() string {
	if t.Completed {
		return "Completed at " + t.CompletedAt.Time.Format("Monday 02/01/2006 15:04")
	}
	return "In progress"
}

func (t *Todo) Delete() error {
	result := db.Delete(&t)
	return result.Error
}

// Toggle status between completed and to do
func (t *Todo) SetStatus() error {
	if !t.Completed {
		t.CompletedAt.Time = time.Now()
		t.CompletedAt.Valid = true
	} else {
		t.CompletedAt.Valid = false
	}
	t.Completed = !t.Completed

	result := db.Save(&t)
	return result.Error
}

func (t *Todo) SetName(n string) error {
	t.Name = n
	result := db.Save(&t)
	return result.Error
}

// Return the status as a string
func (t Todo) Status() string {
	s := t.Name + " set as "
	if t.Completed {
		s += "completed!"
	} else {
		s += "to-do!"
	}
	return s
}
