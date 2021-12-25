package database

import (
	"database/sql"
	"encoding/json"

	"github.com/gocraft/dbr/v2"
	_ "github.com/lib/pq"
)

var sessions map[string]*Session

type Session struct {
	session *dbr.Session
}

type InsertStatement struct {
	stmt *dbr.InsertStmt
}

type SelectStatement struct {
	stmt *dbr.SelectStmt
}

type UpdateStatement struct {
	stmt *dbr.UpdateStmt
}

type DeleteStatement struct {
	stmt *dbr.DeleteStmt
}

func NewSession(d *Dsn) *Session {
	dsn := d.Dsn()
	if sessions[dsn] != nil {
		return sessions[dsn]
	}
	db, err := dbr.Open("postgres", dsn, nil)
	if err != nil {
		panic(err)
	}
	s := db.NewSession(nil)
	ss := Session{session: s}
	if sessions == nil {
		sessions = map[string]*Session{}
	}
	sessions[dsn] = &ss
	return sessions[dsn]
}

func (s *Session) InsertInto(table string) *InsertStatement {
	stmt := s.session.InsertInto(table)
	return &InsertStatement{stmt: stmt}
}

func (s *InsertStatement) Columns(column ...string) *InsertStatement {
	stmt := s.stmt.Columns(column...)
	return &InsertStatement{stmt: stmt}
}

func (s *InsertStatement) Values(value ...interface{}) *InsertStatement {
	stmt := s.stmt.Values(value...)
	return &InsertStatement{stmt: stmt}
}

func (s *InsertStatement) Exec() (sql.Result, error) {
	return s.stmt.Exec()
}

func (s *Session) InsertRecord(record interface{}) (sql.Result, error) {
	table := TableName(record)
	j, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}
	var hash map[string]interface{}
	err = json.Unmarshal(j, &hash)
	if err != nil {
		return nil, err
	}
	column := []string{}
	value := []interface{}{}
	for k, v := range hash {
		column = append(column, k)
		value = append(value, v)
	}
	return s.InsertInto(table).Columns(column...).Values(value...).Exec()
}

func (s *Session) InsertRecordOrPanic(record interface{}) sql.Result {
	r, err := s.InsertRecord(record)
	if err != nil {
		panic(err)
	}
	return r
}

func (s *Session) Select(column ...string) *SelectStatement {
	st := s.session.Select(column...)
	return &SelectStatement{st}
}

func (s *SelectStatement) From(table interface{}) *SelectStatement {
	st := s.stmt.From(table)
	return &SelectStatement{st}
}

func (s *SelectStatement) Join(table interface{}, on interface{}) *SelectStatement {
	st := s.stmt.Join(table, on)
	return &SelectStatement{st}
}

func (s *SelectStatement) Where(query interface{}, value ...interface{}) *SelectStatement {
	st := s.stmt.Where(query, value...)
	return &SelectStatement{st}
}

func (s *SelectStatement) Load(value interface{}) (int, error) {
	return s.stmt.Load(value)
}

func (s *Session) Update(table string) *UpdateStatement {
	st := s.session.Update(table)
	return &UpdateStatement{st}
}

func (s *Session) DeleteFrom(table string) *DeleteStatement {
	st := s.session.DeleteFrom(table)
	return &DeleteStatement{st}
}

func (s *DeleteStatement) Exec() (sql.Result, error) {
	return s.stmt.Exec()
}

func (s *DeleteStatement) ExecOrPanic() sql.Result {
	r, err := s.Exec()
	if err != nil {
		panic(err)
	}
	return r
}
