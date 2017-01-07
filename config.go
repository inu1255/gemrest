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

type database struct {
	Driver     string
	DataSource string
}

type config struct {
	Database database
	Dev      bool
	DocBind  string // doc api url
}

var (
	configFile = flag.String("conf", "config.json", "General configuration file")
	conf       = config{Database: database{Driver: "mysql", DataSource: "root@/gemrest"}, Dev: true}
	Db         *xorm.Engine
	logger     = log.New(os.Stdout, "", log.Ltime|log.Llongfile)
)

func SetLogger(l *log.Logger) {
	logger = l
}

func init() {
	flag.Parse()
	if _, err := os.Stat(*configFile); err == nil {
		file, err := os.Open(*configFile)
		if err != nil {
			logger.Fatalln("open configuration file error ", err)
		}
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&conf)
		if err != nil {
			logger.Fatalln("decode configuration file error", err)
		}
	} else {
		body, err := json.MarshalIndent(conf, "", "    ")
		if err == nil {
			logger.Println("no configuration,writing")
			ioutil.WriteFile(*configFile, body, 0644)
		}
	}
	var err error
	Db, err = xorm.NewEngine(conf.Database.Driver, conf.Database.DataSource)
	if err != nil {
		logger.Fatalln("open database false", err)
	}
	cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
	Db.SetDefaultCacher(cacher)
	if conf.Dev {
		Db.ShowSQL(true)
	}

}
