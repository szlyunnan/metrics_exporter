package model

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var querydata []QueryData

func NewMetricsDB(dbtype, dbpath string) *MetricsDB {
	return &MetricsDB{
		DbType: dbtype,
		DbPath: dbpath,
	}
}

func (mdb *MetricsDB) DBEngine() (*sqlx.DB, error) {
	db, err := sqlx.Open(mdb.DbType, mdb.DbPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (mdb *MetricsDB) DBQuery(sql string) ([]QueryData, error) {
	dbengine, err := mdb.DBEngine()
	if err != nil {
		return nil, err
	}
	defer dbengine.Close()

	rows, err := dbengine.Queryx(sql)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		qd := &QueryData{}
		if err := rows.StructScan(qd); err != nil {
			return nil, err
		}
		querydata = append(querydata, []QueryData{0: *qd}...)
	}
	return querydata, nil
}
