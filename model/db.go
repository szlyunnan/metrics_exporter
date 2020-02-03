package model

import (
	"database/sql"
)

var querydata []QueryData

func New(dbtype, dbpath string) *MetricsDB {
	return &MetricsDB{
		DbType: dbtype,
		DbPath: dbpath,
	}
}

func (mdb *MetricsDB) DBEngine() (*sql.DB, error) {
	db, err := sql.Open(mdb.DbType, mdb.DbPath)
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

	rows, err := dbengine.Query(sql)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		qd := &QueryData{}
		if err := rows.Scan(qd.SiPkg, qd.SiCode, qd.SiMsq); err != nil {
			return nil, err
		}
		querydata = append(querydata, []QueryData{0: *qd}...)

	}
	return querydata, nil
}
