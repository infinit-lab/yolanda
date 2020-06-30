package database

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/infinit-lab/yolanda/config"
	l "github.com/infinit-lab/yolanda/logutils"
	"reflect"
	"strings"
	"time"
)

var pool *sql.DB

func init() {
	l.Trace("Initializing database...")

	go func() {
		for {
			url := config.GetString("mysql.url")
			l.TraceF("Url is %s", url)
			if url == "" {
				l.Error("Failed to get mysql url")
				time.Sleep(500 * time.Millisecond)
				continue
			}
			var err error
			pool, err = sql.Open("mysql", url)
			if err != nil {
				pool = nil
				l.ErrorF("Failed to open %s. error: %s", url, err.Error())
				time.Sleep(5 * time.Second)
				continue
			}
			l.TraceF("Success to open %s", url)
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
		var tempArgs []interface{}
		for _, arg := range args {
			if reflect.ValueOf(arg).Kind() == reflect.Ptr {
				tempArgs = append(tempArgs, reflect.ValueOf(arg).Elem().Interface())
			} else {
				tempArgs = append(tempArgs, arg)
			}
		}
		l.Error("Sql is ", query, tempArgs)
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
		var tempArgs []interface{}
		for _, arg := range args {
			if reflect.ValueOf(arg).Kind() == reflect.Ptr {
				tempArgs = append(tempArgs, reflect.ValueOf(arg).Elem().Interface())
			} else {
				tempArgs = append(tempArgs, arg)
			}
		}
		l.Error("Sql is ", query, tempArgs)
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

type tagAndField struct {
	tags []string
	fields []interface{}
}

func reflectStruct(v reflect.Value, omit string, tagAndField *tagAndField) {
	if v.Kind() != reflect.Struct {
		return
	}
	t := v.Type()
	n := t.NumField()
	for i := 0; i < n; i++ {
		f := v.Field(i)
		if f.Kind() == reflect.Struct {
			reflectStruct(f, omit, tagAndField)
			continue
		}
		field := t.Field(i)
		tag := field.Tag.Get("db")
		if tag == "" {
			continue
		}
		patterns := strings.Split(tag, ",")
		t := patterns[0]
		isSkip := false
		for i, pattern := range patterns {
			if i == 0 {
				continue
			}
			if strings.TrimSpace(pattern) == omit {
				isSkip = true
				break
			}
		}
		if isSkip {
			continue
		}
		if reflect.ValueOf(f.Interface()).Kind() == reflect.Ptr {
			continue
		}
		if !f.CanAddr() {
			continue
		}
		tagAndField.tags = append(tagAndField.tags, t)
		tagAndField.fields = append(tagAndField.fields, f.Addr().Interface())
	}
}

func reflectValue(value interface{}, omit string) (tags []string, fields []interface{}, err error) {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Ptr {
		return nil, nil, errors.New("value must be ptr")
	}
	if value == nil {
		return nil, nil, errors.New("value should not be nil")
	}
	e := v.Type().Elem()
	if e.Kind() != reflect.Struct {
		return nil, nil, errors.New("value must be struct")
	}
	tagAndFields := new(tagAndField)
	reflectStruct(v.Elem(), omit, tagAndFields)
	tags = tagAndFields.tags
	fields = tagAndFields.fields
	if len(tags) == 0 {
		return nil, nil, errors.New("no tag")
	}
	err = nil
	return
}

func getSqlString(keys map[string]string, value interface{}, tableName string) (string, []interface{}, error) {
	tags, fields, err := reflectValue(value, "omitget")
	if err != nil {
		l.Error("Failed to reflectValue. error: ", err)
		return "", nil, err
	}
	sqlString := "SELECT "
	for i, tag := range tags {
		sqlString += "IFNULL(`" + tag + "`, '')"
		if i != len(tags) - 1 {
			sqlString += ", "
		}
	}
	sqlString += " FROM " + tableName
	if len(keys) > 0 {
		sqlString += " WHERE "
		isFirst := true
		for key, v := range keys {
			if !isFirst {
				sqlString += " AND "
			}
			isFirst = false
			sqlString += "`" + key + "` = '" + v +"'"
		}
	}
	return sqlString, fields, nil
}

func SingleTableGet(keys map[string]string, value interface{}, tableName string) error {
	sqlString, fields, err := getSqlString(keys, value, tableName)
	if err != nil {
		return err
	}
	sqlString += " LIMIT 1"
	rows, err := Query(sqlString)
	if err != nil {
		return err
	}
	defer func() {
		_ = rows.Close()
	}()

	if !rows.Next() {
		return errors.New("not found")
	}

	err = rows.Scan(fields...)
	if err != nil {
		l.Error("Failed to Scan. error: ", err)
		return err
	}
	return nil
}

func SingleTableGetList(keys map[string]string, value interface{}, tableName string) ([]interface{}, error) {
	sqlString, _, err := getSqlString(keys, value, tableName)
	if err != nil {
		return nil, err
	}
	sqlString += " ORDER BY `id` DESC"
	rows, err := Query(sqlString)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var values []interface{}
	for rows.Next() {
		newValue := reflect.New(reflect.TypeOf(value).Elem()).Interface()
		_, fields, err := reflectValue(newValue, "omitget")
		if err != nil {
			return nil, err
		}
		err = rows.Scan(fields...)
		if err != nil {
			return nil, err
		}
		values = append(values, newValue)
	}
	return values, nil
}

func SingleTableCreate(value interface{}, tableName string) error {
	tags, fields, err := reflectValue(value, "omitcreate")
	if err != nil {
		l.Error("Failed to reflectValue. error: ", err)
		return err
	}
	sqlString := "INSERT INTO " + tableName + " ("
	for i, tag := range tags {
		sqlString += "`" + tag + "`"
		if i != len(tags) - 1 {
			sqlString += ", "
		}
	}
	sqlString += ") VALUES ("
	for i, _ := range tags {
		sqlString += "?"
		if i != len(tags) - 1 {
			sqlString += ", "
		}
	}
	sqlString += ")"
	_, err = Exec(sqlString, fields...)
	return err
}

func SingleTableUpdate(keys map[string]string, value interface{}, tableName string) error {
	tags, fields, err := reflectValue(value, "omitupdate")
	if err != nil {
		l.Error("Failed to reflectValue. error: ", err)
		return err
	}
	sqlString := "UPDATE " + tableName + " SET "
	for i, tag := range tags {
		sqlString += "`"
		sqlString += tag
		sqlString += "` = ?"
		if i != len(tags) - 1 {
			sqlString += ", "
		}
	}
	if len(keys) > 0 {
		sqlString += " WHERE "
		isFirst := true
		for key, v := range keys {
			if !isFirst {
				sqlString += " AND "
			}
			isFirst = false
			sqlString += "`" + key + "` = ?"
			fields = append(fields, v)
		}
	}
	_, err = Exec(sqlString, fields...)
	return err
}

func SingleTableDelete(keys map[string]string, tableName string) error {
	if len(keys) == 0 {
		return errors.New("keys cannot be empty")
	}
	var values []interface{}
	sqlString := "DELETE FROM " + tableName + " WHERE "
	isFirst := true
	for key, v := range keys {
		if !isFirst {
			sqlString += " AND "
		}
		isFirst = false
		sqlString += "`" + key + "` = ?"
		values = append(values, v)
	}
	_, err := Exec(sqlString, values...)
	return err
}
