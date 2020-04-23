package sqlite

import (
	"database/sql"
	"errors"
	l "github.com/infinit-lab/yolanda/logutils"
	_ "github.com/mattn/go-sqlite3"
	"sync"
)

type Sqlite struct {
	db    *sql.DB
	mutex sync.Mutex
}

type Column struct {
	Name    string
	Type    string
	Default string
}

type Table struct {
	Name    string
	Columns []Column
}

func InitializeDatabase(path string) (*Sqlite, error) {
	l.TraceF("Initializing sqlite %s...", path)
	s := new(Sqlite)
	var err error
	s.db, err = sql.Open("sqlite3", path)
	if err != nil {
		l.Error("Failed to open sqlite. error: ", err)
		return nil, err
	}
	return s, nil
}

func (s *Sqlite) InitializeTable(t Table) error {
	rows, err := s.Query("PRAGMA table_info(`" + t.Name + "`)")
	if err != nil {
		return err
	}
	isExists := false
	var columns []Column
	for rows.Next() {
		isExists = true
		var column Column
		var cid int64
		var notnull int64
		var dflt sql.NullString
		var pk int64
		err := rows.Scan(&cid, &column.Name, &column.Type, &notnull, &dflt, &pk)
		if err != nil {
			l.Error("Failed to scan. error: ", err)
			return err
		}
		columns = append(columns, column)
	}
	_ = rows.Close()

	if !isExists {
		_, err := s.Exec("CREATE TABLE IF NOT EXISTS `" + t.Name + "` (" +
			"`id` INTEGER PRIMARY KEY AUTOINCREMENT" +
			")")
		if err != nil {
			return err
		}
	}

	for _, column := range t.Columns {
		isFind := false
		for _, c := range columns {
			if c.Name == column.Name {
				isFind = true
				break
			}
		}
		if !isFind {
			sqlString := "ALTER TABLE `" + t.Name + "` ADD COLUMN `" + column.Name + "` " + column.Type + " NOT NULL " +
				"DEFAULT '" + column.Default + "'"
			_, err := s.Exec(sqlString)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Sqlite) Prepare(query string) (stmt *sql.Stmt, err error) {
	if s.db == nil {
		l.Error("Database is nil")
		return nil, errors.New("数据库打开失败")
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	stmt, err = s.db.Prepare(query)
	if err != nil {
		l.Error("Failed to prepare sql. error: ", err)
		l.Error("Sql is ", query)
	}
	return
}

func (s *Sqlite) Exec(query string, args ...interface{}) (ret sql.Result, err error) {
	if s.db == nil {
		l.Error("Database is nil")
		return nil, errors.New("数据库打开失败")
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	ret, err = s.db.Exec(query, args...)
	if err != nil {
		l.Error("Failed to exec sql. error: ", err)
		l.Error("Sql is ", query, args)
	}
	return
}

func (s *Sqlite) Query(query string, args ...interface{}) (rows *sql.Rows, err error) {
	if s.db == nil {
		l.Error("Database is nil")
		return nil, errors.New("数据库打开失败")
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	rows, err = s.db.Query(query, args...)
	if err != nil {
		l.Error("Failed to query sql, error: ", err)
		l.Error("Sql is ", query, args)
	}
	return
}
