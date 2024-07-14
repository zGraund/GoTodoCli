package models

import (
	"fmt"
	"time"
)

type Stats struct {
	totalCompleted int
	bestDay        time.Time
	dayWithMost    time.Time
}

func NewStats() Stats {
	s := Stats{}
	rows, _ := db.Model(&Todo{}).Where("completed = ?", true).Group("date(created_at)").Rows()
	defer rows.Close()
	result := Todo{}
	for rows.Next() {
		db.ScanRows(rows, &result)
		fmt.Println(result)
	}
	return s
}

func (s *Stats) updateCompleted() {
	var result int64
	db.Where(&Todo{Completed: true}).Count(&result)
	s.totalCompleted = int(result)
}

func (s *Stats) updateBestDay() {
	var result int64
	db.Where(&Todo{Completed: true}).Count(&result)
	s.totalCompleted = int(result)
}
