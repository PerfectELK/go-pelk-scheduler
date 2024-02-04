package storage

import (
	"database/sql"
	"fmt"
	"github.com/PerfectELK/go-pelk-scheduler/pkg/scheduler"
	"time"
)

type SqliteJobDataStorage struct {
	conn      string
	jobsTable string
	db        *sql.DB
}

func NewSqliteJobDataStorage(db *sql.DB, jobsTable string) (*SqliteJobDataStorage, error) {
	jds := &SqliteJobDataStorage{
		jobsTable: jobsTable,
		db:        db,
	}

	err := jds.createTables()
	if err != nil {
		return nil, err
	}

	return jds, nil
}

func (s *SqliteJobDataStorage) createTables() error {
	_, err := s.db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
		    name varchar(255) primary key,
		    last_started integer default null
		);
	`, s.jobsTable))
	return err
}

func (s *SqliteJobDataStorage) ImportJobsData() []*scheduler.JobData {
	rows, err := s.db.Query(fmt.Sprintf("SELECT * FROM %s", s.jobsTable))
	if err != nil {
		fmt.Println(fmt.Sprintf("sqlite.ImportJobsData error: %s", err))
		return nil
	}

	var arr []*scheduler.JobData
	for rows.Next() {
		var name string
		var lastStarted int64

		err := rows.Scan(&name, &lastStarted)
		if err != nil {
			return nil
		}

		arr = append(arr, &scheduler.JobData{
			Name:        name,
			LastStarted: time.Unix(lastStarted, 0),
		})
	}
	return arr
}

func (s *SqliteJobDataStorage) ExportJobsData(jobs []*scheduler.Job) error {
	for _, j := range jobs {
		_, err := s.db.Exec(fmt.Sprintf(`
				INSERT INTO %s (name, last_started) VALUES (?, ?) 
				ON CONFLICT (name) DO UPDATE SET last_started=? WHERE name=? 
			`, s.jobsTable),
			j.GetName(),
			j.GetLatStarted().Unix(),
			j.GetLatStarted().Unix(),
			j.GetName(),
		)
		if err != nil {
			return err
		}
	}
	return nil
}
