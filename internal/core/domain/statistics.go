package domain

import "time"

type Statistics struct {
	TasksCreated               int
	TasksCompleted             int
	TasksCompletionRate        *float64
	TasksAverageCompletionTime *time.Duration
}

func NewStatistics(
	tasksCreated, tasksCompleted int,
	tasksCompletionRate *float64,
	tasksAverage *time.Duration,
) Statistics {
	return Statistics{
		TasksCreated:               tasksCreated,
		TasksCompleted:             tasksCompleted,
		TasksCompletionRate:        tasksCompletionRate,
		TasksAverageCompletionTime: tasksAverage,
	}
}

type StatisticsFilters struct {
	UserID *int
	From   *time.Time
	To     *time.Time
}
