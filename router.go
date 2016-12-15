package gemrest

import (
	"log"
	"reflect"
	"strconv"

	"github.com/go-gem/gem"
)

var (
	Router = gem.NewRouter()
)

type ApiService interface {
	Before(ctx *gem.Context) bool
	Finish(ctx *gem.Context, err interface{})
	After(ctx *gem.Context, data interface{}, msg string)
}

func makeHandlerFunc(apply reflect.Value, call []convertFunc) gem.HandlerFunc {
	return func(ctx *gem.Context) {
		n := len(call)
		params := make([]reflect.Value, n)
		params[0] = call[0](ctx, "")
		service := params[0].Interface().(ApiService)
		defer func() {
			err := recover()
			log.Println("recover", err)
			service.Finish(ctx, err)
		}()
		for i := 1; i < n; i++ {
			params[i] = call[i](ctx, strconv.Itoa(i))
		}
		if service.Before(ctx) {
			log.Println(params)
			out := apply.Call(params)
			data := out[0].Interface()
			msg := out[1].String()
			service.After(ctx, data, msg)
		}
	}
}

func Bind(prefix string, service ApiService) {
	t := reflect.TypeOf(service)
	numMethod := t.NumMethod()
	instCall := newInstCall(t.Elem())
	for i := 0; i < numMethod; i++ {
		m := t.Method(i)
		flag, path, call := convertMethodParams(prefix, m)
		if flag == -1 {
			continue
		}
		call[0] = instCall
		if flag == 1 {
			log.Println("post", path)
			Router.POST(path, makeHandlerFunc(m.Func, call))
		} else {
			log.Println("get", path)
			Router.GET(path, makeHandlerFunc(m.Func, call))
		}
	}
}

func Start(host string) {
	srv := gem.New(host, Router.Handler())
	log.Fatal(srv.ListenAndServe())
}
