package gemrest

import (
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var (
	Db     *xorm.Engine
	logger = log.New(os.Stdout, "[\x1b[30m gemrest \x1b[0m] ", log.Ltime|log.Lshortfile)
)

func SetLogger(l *log.Logger) {
	logger = l
	logger.SetPrefix("[\x1b[30m gemrest \x1b[0m] ")
}

func SetDb(driver, datasource string) {
	var err error
	Db, err = xorm.NewEngine(driver, datasource)
	if err != nil {
		logger.Fatalln("open database false", err)
	}
	cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
	Db.SetDefaultCacher(cacher)
}
