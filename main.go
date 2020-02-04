package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"log"
	"metrics_exporter/config"
	"metrics_exporter/model"
	"metrics_exporter/utils"
	"net/http"
	"os"
)

var (
	conf   config.Configure
	listen string
	host   string
	port   int
	path   string
	help   bool
)

func main() {

	dbfullpath := fmt.Sprintf("%s/%s", conf.Metrics.DbPath, conf.Metrics.DbName)

	mdb := model.NewMetricsDB("sqlite3", dbfullpath)
	data, err := mdb.DBQuery("select si_pkg, si_code, si_msg, si_build from `aax_sign-1`")
	if err != nil {
		panic(err)
	}

	worker := utils.NewMetricsCollector(data)
	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(worker)

	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.HandlerFor(
		//prometheus.DefaultGatherer,
		reg,
		promhttp.HandlerOpts{},
	)).Methods("GET")

	log.Println(fmt.Sprintf("Server listen on: %s", listen))
	log.Fatal(http.ListenAndServe(listen, router))

}

func init() {
	flag.BoolVar(&help, "h", false, "help")
	flag.StringVar(&host, "host", "127.0.0.1", "The hostname or IP on which the REST server will listen")
	flag.IntVar(&port, "port", 8090, "The port on which the REST server will listen")
	flag.StringVar(&path, "config", "./conf.d", "config path directory, defaults is current directory ./conf.d")

	flag.Parse()
	listen = fmt.Sprintf("%s:%d", host, port)
	if help {
		flag.Usage()
		os.Exit(0)
	}

	vper := viper.New()
	vper.SetConfigName("metrics")
	vper.AddConfigPath(path)
	vper.AddConfigPath(".")
	vper.AddConfigPath("/tmp")
	vper.SetConfigType("yaml")

	err := vper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	if err := vper.Unmarshal(&conf); err != nil {
		panic(fmt.Errorf("Unmarshal error file: %s \n", err))
	}
}
