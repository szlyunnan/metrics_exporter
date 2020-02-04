package config

import (
	"fmt"
	"time"
)

type Config struct {
	DbPath string `mapstructure:"db_path"`
	DbName string `mapstructure:"db_name"`
	TbName string `mapstructure:"tb_name"`
	SiName string `mapstructure:"si_name"`
}

type Configure struct {
	Metrics Config `mapstructure:"metices"`
}

func New() *Configure {
	// 本地开发环境配置
	// tb 的名称是按照tb-prefix-day 构成的
	// tb 名称栗子: metrics_tb-28	·	···
	timeObj := time.Now()
	return &Configure{
		Metrics: Config{
			DbPath: "/Users/sam/data/storage",
			DbName: "metrics_prometheus.db",
			TbName: fmt.Sprintf("%s-%d", "metrics_tb", timeObj.Day()),
			SiName: "metrics-si",
		},
	}
}
