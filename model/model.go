package model

type MetricsDB struct {
	DbType string
	DbPath string
}

type QueryData struct {
	SiPkg   string `db:"si_pkg"`
	SiMsg   string `db:"si_msg"`
	SiCode  int64  `db:"si_code"`
	SiBuild int64  `db:"si_build"`
}
