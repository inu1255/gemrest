package gemrest

import (
	"log"

	"github.com/go-gem/gem"
)

var (
	Router = gem.NewRouter()
)

type Api interface{
    Bind(perfix string,*gem.Router)
}

type ApiSlice Api[]

type Service interface{
    GetApiSlice() ApiSlice
}

func Bind(prefix string,service Service) {
    for _,item := range service.GetApiSlice(){
        item.Bind(perfix,Router)
    }
}

func Start(host string) {
	srv := gem.New(host, Router.Handler())
	log.Fatal(srv.ListenAndServe())
}
