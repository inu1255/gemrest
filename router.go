package gemrest

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/go-gem/gem"
)

var (
	Router = gem.NewRouter()
)

type Context struct {
	*gem.Context
}

type ApiService interface {
	Before(ctx *Context) bool
	Finish(err interface{})
	After(data interface{}, msg string)
}

// call the api
func makeHandlerFunc(m reflect.Method, call []convertFunc) gem.HandlerFunc {
	return func(gctx *gem.Context) {
		ctx := &Context{Context: gctx}
		n := len(call)
		params := make([]reflect.Value, n)
		params[0] = call[0](ctx, "")
		service := params[0].Interface().(ApiService)
		defer func() {
			err := recover()
			if err != nil {
				logger.Println("recover", err)
			}
			service.Finish(err)
		}()
		for i := 1; i < n; i++ {
			params[i] = call[i](ctx, strconv.Itoa(i))
		}
		if service.Before(ctx) {
			logger.Println(params[0].Type(), m.Name, params[1:])
			out := m.Func.Call(params)
			data := out[0].Interface()
			msg := out[1].Interface()
			if msg == nil {
				service.After(data, "")
			} else {
				service.After(data, msg.(error).Error())
			}
		}
	}
}

// bind a router for service's method
// method satisfy func(in ...interface{}) (interface{},error) will be export
// if there is slice/struct/ptr in "params in" ,export a POST router
func Bind(prefix string, service ApiService) int {
	t := reflect.TypeOf(service)
	numMethod := t.NumMethod()
	// instCall := newInstCall(t.Elem())
	instCall := copyInstCall(service)
	n := 0
	for i := 0; i < numMethod; i++ {
		m := t.Method(i)
		flag, path, call := convertMethodParams(prefix, m)
		if flag == -1 {
			continue
		}
		n++
		call[0] = instCall
		if flag == 1 {
			logger.Println("\x1b[34mpost\x1b[0m", path)
			Router.POST(path, makeHandlerFunc(m, call))
		} else {
			logger.Println("\x1b[32mget \x1b[0m", path)
			Router.GET(path, makeHandlerFunc(m, call))
		}
	}
	return n
}

func Reg(service ApiService) int {
	prefix := reflect.TypeOf(service).String()
	if prefix[0] == '*' {
		prefix = prefix[1:]
	}
	prefix = strings.ToLower(prefix)
	n := len(prefix)
	if n > 7 && prefix[n-7:] == "service" {
		n -= 7
		prefix = prefix[:n]
	}
	for i := 0; i < n; i++ {
		if prefix[i] == '.' {
			prefix = strings.Join([]string{prefix[:i], "/", prefix[i+1:]}, "")
			break
		}
	}
	return Bind("/"+prefix, service)
}

func Start(host string) {
	srv := gem.New(host, Router.Handler())
	logger.Fatal(srv.ListenAndServe())
}
