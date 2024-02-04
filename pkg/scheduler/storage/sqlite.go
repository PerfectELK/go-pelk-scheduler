package storage

import (
	"database/sql"
	"github.com/PerfectELK/go-pelk-scheduler/pkg/scheduler"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteJobDataStorage struct {
	conn      string
	jobsTable string
	db        *sql.DB
}

func NewSqliteJobDataStorage(conn string, jobsTable string) (*SqliteJobDataStorage, error) {
	db, err := sql.Open("sqlite3", conn)
	if err != nil {
		return nil, err
	}

	return &SqliteJobDataStorage{
		conn:      conn,
		jobsTable: jobsTable,
		db:        db,
	}, nil
}

func (s *SqliteJobDataStorage) ImportJobsData() []*scheduler.JobData {
	return nil
}

func (s *SqliteJobDataStorage) ExportJobsData(jobs []*scheduler.Job) error {
	for _, j := range jobs {
		_, err := s.db.Exec("INSERT INTO jobs (name, last_started) VALUES (?, ?)", j.GetName(), j.GetLatStarted().Unix())
		if err != nil {
			return err
		}
	}
	return nil
}
