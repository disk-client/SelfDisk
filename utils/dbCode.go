/*
 * @Author: your name
 * @Date: 2020-06-15 17:42:23
 * @LastEditTime: 2020-06-18 16:10:08
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /SelfDisk/utils/dbCode.go
 */

package utils

import (
	"database/sql"
)

// DisktDB PostgreSQL数据库连接
type DisktDB struct {
	Name     string
	Password string
	DBname   string
	Sslmode  string
}

// TheDB 数据库
var TheDB = DisktDB{Name: "postgres", Password: "postgres", DBname: "postgres", Sslmode: "disable"}

// SelectSQL 获取select语句的结果，且结果为0到多条
func (db DisktDB) SelectSQL(execsql string, args ...interface{}) *sql.Rows {
	var conn = new(sql.DB)
	var err error
	conn, err = sql.Open("postgres", "user="+db.Name+" dbname="+db.DBname+" password="+db.Password+" sslmode="+db.Sslmode)
	CheckErr(err)
	rows, err := conn.Query(execsql, args...)
	CheckErr(err)
	conn.Close()
	return rows
}

//GetOne 获取select语句的结果，且结果为1条
func (db DisktDB) GetOne(execsql string, params []interface{}, res []interface{}) error {
	var conn = new(sql.DB)
	var err error
	conn, err = sql.Open("postgres", "user="+db.Name+" dbname="+db.DBname+" password="+db.Password+" sslmode="+db.Sslmode)
	if err != nil {
		return err
	}
	err = conn.QueryRow(execsql, params...).Scan(res...)
	if err != nil {
		conn.Close()
		return err
	}
	err = conn.Close()
	return err
}

// UpdateSQL 进行数据库的update操作，不返回结果
func (db DisktDB) UpdateSQL(execsql string, params []interface{}) {
	var conn = new(sql.DB)
	var err error
	conn, err = sql.Open("postgres", "user="+db.Name+" dbname="+db.DBname+" password="+db.Password+" sslmode="+db.Sslmode)
	CheckErr(err)
	s, err := conn.Prepare(execsql)
	CheckErr(err)
	_, err = s.Exec(params...)
	CheckErr(err)
	conn.Close()
}

// DeleteSQL 进行数据库的delete操作，不返回结果
func (db DisktDB) DeleteSQL(execsql string, params []interface{}) {
	var conn = new(sql.DB)
	var err error
	conn, err = sql.Open("postgres", "user="+db.Name+" dbname="+db.DBname+" password="+db.Password+" sslmode="+db.Sslmode)
	CheckErr(err)
	s, err := conn.Prepare(execsql)
	CheckErr(err)
	_, err = s.Exec(params...)
	CheckErr(err)
	conn.Close()
}

// InsertSQL 进行数据库的insert操作，不返回结果
func (db DisktDB) InsertSQL(execsql string, params []interface{}) error {
	var conn = new(sql.DB)
	var err error
	conn, err = sql.Open("postgres", "user="+db.Name+" dbname="+db.DBname+" password="+db.Password+" sslmode="+db.Sslmode)
	defer conn.Close()
	if err != nil {
		return err
	}
	s, err := conn.Prepare(execsql)
	if err != nil {
		return err
	}
	_, err = s.Exec(params...)
	if err != nil {
		return err
	}
	return nil
}

// HandleSQLString 处理sql输出的可能为nil的字符串
func HandleSQLString(s sql.NullString) string {
	switch s.Valid {
	case true:
		return s.String
	default:
		return "-"
	}
}
