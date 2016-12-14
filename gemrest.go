package gemrest

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var (
	configFile = flag.String("conf", "config.json", "General configuration file")
	conf       = config{Database: database{Driver: "mysql", DataSource: "root@/gemrest"}, Dev: true}
	Db         *xorm.Engine
)

func init() {
	flag.Parse()
	if _, err := os.Stat(*configFile); err == nil {
		file, err := os.Open(*configFile)
		if err != nil {
			log.Fatalln("open configuration file error ", err)
		}
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&conf)
		if err != nil {
			log.Fatalln("decode configuration file error", err)
		}
	} else {
		body, err := json.MarshalIndent(conf, "", "    ")
		if err == nil {
			log.Println("no configuration,writing")
			ioutil.WriteFile(*configFile, body, 0644)
		}
	}
	var err error
	Db, _ = xorm.NewEngine(conf.Database.Driver, conf.Database.DataSource)
	if err != nil {
		log.Fatalln("open database false", err)
	}
	if conf.Dev {
		Db.ShowSQL(true)
	}

}
