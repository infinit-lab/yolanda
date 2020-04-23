package database

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/infinit-lab/yolanda/config"
	l "github.com/infinit-lab/yolanda/logutils"
	"time"
)

var pool *sql.DB

func init() {
	l.Trace("Initializing database...")
	url := config.GetString("mysql.url")
	l.TraceF("Url is %s", url)
	if url == "" {
		l.Error("Failed to get mysql url")
		return
	}

	maxOpenConns := config.GetInt("mysql.maxOpenConns")
	l.TraceF("Max open connections is %d", maxOpenConns)
	if maxOpenConns == 0 {
		maxOpenConns = 20
		l.TraceF("Max open connections reset to %d", maxOpenConns)
	}

	maxIdleConns := config.GetInt("mysql.maxIdleConns")
	l.TraceF("Max idle connections is %d", maxIdleConns)
	if maxIdleConns == 0 {
		maxIdleConns = 5
		l.TraceF("Max idle connections reset to %d", maxIdleConns)
	}

	maxLifetime := config.GetInt("mysql.maxLifetime")
	l.TraceF("Max lifetime is %d", maxLifetime)
	if maxLifetime == 0 {
		maxLifetime = 120
		l.TraceF("Max lifetime reset to %d", maxLifetime)
	}

	go func() {
		for {
			var err error
			pool, err = sql.Open("mysql", url)
			if err != nil {
				pool = nil
				l.ErrorF("Failed to open %s. error: %s", url, err.Error())
				time.Sleep(5 * time.Second)
				continue
			}
			l.TraceF("Success to open %s", url)
			pool.SetMaxOpenConns(maxOpenConns)
			pool.SetMaxIdleConns(maxIdleConns)
			pool.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Second)

			err = pool.Ping()
			if err == nil {
				for {
					err := pool.Ping()
					if err != nil {
						l.Error("Mysql disconnected!!!. error: ", err)
						_ = pool.Close()
						pool = nil
						break
					}
					time.Sleep(2 * time.Second)
				}
			} else {
				l.Error("Failed to connect mysql, error: ", err)
			}

			time.Sleep(5 * time.Second)
		}
	}()
}

func Exec(query string, args ...interface{}) (sql.Result, error) {
	if pool == nil {
		l.Error("Mysql pool pointer is nil")
		return nil, errors.New("Mysql未连接")
	}
	ret, err := pool.Exec(query, args...)
	if err != nil {
		l.Error("Exec sql error: ", err)
		l.Error("Sql is ", query, args)
	}
	return ret, err
}

func Query(query string, args ...interface{}) (*sql.Rows, error) {
	if pool == nil {
		l.Error("Mysql pool pointer is nil")
		return nil, errors.New("Mysql未连接")
	}
	rows, err := pool.Query(query, args...)
	if err != nil {
		l.Error("Query sql error: ", err)
		l.Error("Sql is ", query, args)
	}
	return rows, err
}

func Begin() (*sql.Tx, error) {
	if pool == nil {
		l.Error("Mysql pool pointer is null")
		return nil, errors.New("Mysql未连接")
	}
	return pool.Begin()
}
