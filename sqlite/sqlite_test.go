package sqlite

import (
	"github.com/infinit-lab/yolanda/logutils"
	"testing"
)

var s *Sqlite

func TestInitializeDatabase(t *testing.T) {
	var err error
	s, err = InitializeDatabase("./test.db")
	if err != nil {
		t.Fatal("Failed to InitializeDatabase. error: ", err)
	}
}

func TestSqlite_InitializeTable(t *testing.T) {
	table := Table{
		Name: "test",
		Columns: []Column{
			{
				Name:    "name",
				Type:    "VARCHAR(32)",
				Default: "",
				Index:   true,
				Unique:  true,
			},
			{
				Name:    "age",
				Type:    "INTEGER",
				Default: "18",
				Index:   true,
			},
		},
	}

	err := s.InitializeTable(table)
	if err != nil {
		t.Fatal("Failed to InitializeTable. error: ", err)
	}
}

func TestSqlite_Prepare(t *testing.T) {
	stmt, err := s.Prepare("INSERT INTO `test` (`name`) VALUES (?)")
	if err != nil {
		t.Error("Failed to Prepare. error: ", err)
		return
	}
	_, err = stmt.Exec("test1")
	if err != nil {
		t.Error("Failed to Exec. error: ", err)
	}
	_ = stmt.Close()

	stmt, err = s.Prepare("SELECT `name`, `age` FROM `test` WHERE `name` = ? LIMIT 1")
	if err != nil {
		t.Error("Failed to Prepare. error: ", err)
	}
	defer func() {
		_ = stmt.Close()
	}()
	rows, err := stmt.Query("test1")
	if err != nil {
		t.Error("Failed to Query. error: ", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	if rows.Next() {
		var name string
		var age int64
		err := rows.Scan(&name, &age)
		if err != nil {
			t.Error("Failed to Scan. error: ", err)
		}
		if name != "test1" {
			t.Error("Name should be test1 not ", name)
		}
		if age != 18 {
			t.Error("Age should be 18 not ", age)
		}
		logutils.ErrorF("Name is %s. Age is %d", name, age)
	} else {
		t.Error("Result is empty.")
	}
}

func TestSqlite_Exec(t *testing.T) {
	_, err := s.Exec("INSERT INTO `test` (`name`, `age`) VALUES (?, ?)", "test2", 30)
	if err != nil {
		t.Error("Failed to Exec. error: ", err)
	}
}

func TestSqlite_Query(t *testing.T) {
	rows, err := s.Query("SELECT `name`, `age` FROM `test` WHERE `name` = ? LIMIT 1", "test2")
	if err != nil {
		t.Error("Failed to Query. error: ", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	if rows.Next() {
		var name string
		var age int64
		err := rows.Scan(&name, &age)
		if err != nil {
			t.Error("Failed to Scan. error: ", err)
		}
		if name != "test2" {
			t.Error("Name should be test2 not ", name)
		}
		if age != 30 {
			t.Error("Age should be 30 not ", age)
		}
		logutils.ErrorF("Name is %s. Age is %d", name, age)
	} else {
		t.Error("Result is empty.")
	}
}
