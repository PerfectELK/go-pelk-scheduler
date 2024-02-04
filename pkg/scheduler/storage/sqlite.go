package storage

import (
	"database/sql"
	"fmt"
	"github.com/PerfectELK/go-pelk-scheduler/pkg/scheduler"
	_ "github.com/mattn/go-sqlite3"
	"time"
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

	jds := &SqliteJobDataStorage{
		conn:      conn,
		jobsTable: jobsTable,
		db:        db,
	}

	err = jds.createTables()
	if err != nil {
		return nil, err
	}

	return jds, nil
}

func (s *SqliteJobDataStorage) createTables() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS jobs (
		    name varchar(255) primary key,
		    last_started integer default null
		);
	`)
	return err
}

func (s *SqliteJobDataStorage) ImportJobsData() []*scheduler.JobData {
	rows, err := s.db.Query("SELECT * FROM jobs")
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
		_, err := s.db.Exec(`
				INSERT INTO jobs (name, last_started) VALUES (?, ?) 
				ON CONFLICT (name) DO UPDATE SET last_started=? WHERE name=? 
			`,
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
