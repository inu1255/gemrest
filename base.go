package gemrest

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"

	"reflect"
	"github.com/go-gem/gem"
	"github.com/go-xorm/xorm"
)

func stack() []byte {
	buf := make([]byte, 10240)
	n := runtime.Stack(buf, false)
	if n > 738 {
		copy(buf[23:], buf[738:n])
		return buf[:n-715]
	}
	return buf[:n]
}
func printStack() {
	log.Println(string(stack()))
	log.Println("end here.")
}

// ApiMsg ...
type ApiMsg struct {
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

// DefaultService return json
type DefaultService struct {
	Ctx *Context
}

func (d *DefaultService) json(data interface{}, msg string) {
	body, err := json.MarshalIndent(&ApiMsg{data, msg}, "", "    ")
	if err != nil {
		d.json(err, fmt.Sprintf("Internal Server Error: %v", err))
		return
	}
	d.Ctx.RequestCtx.Response.Header.SetStatusCode(200)
	d.Ctx.RequestCtx.Response.Header.SetContentType(gem.HeaderContentTypeJSON)
	d.Ctx.RequestCtx.Response.SetBody(body)
}
func (d *DefaultService) Before(ctx *Context) bool {
	d.Ctx = ctx
	return true
}
func (d *DefaultService) After(data interface{}, msg string) {
	d.json(data, msg)
}
func (d *DefaultService) Finish(err interface{}) {
	if err != nil {
		printStack()
		d.json(err, fmt.Sprintf("Internal Server Error:%v", err))
	}
}

// cross origin
type CrosService struct {
	DefaultService
}

func (m *CrosService) Before(ctx *Context) bool {
	ctx.Response.Header.Add("Access-Control-Allow-Origin", string(ctx.Request.Host()))
	ctx.Response.Header.Add("Access-Control-Allow-Credentials", "true")
	ctx.Response.Header.Add("Access-Control-Allow-Headers", "x-auth-token,content-type")
	return m.DefaultService.Before(ctx)
}

type DatabaseService struct {
	DefaultService
	Db *xorm.Session
}

// database support
func (b *DatabaseService) Before(ctx *Context) bool {
	b.Db = Db.NewSession()
	// b.Db.Begin()
	return b.DefaultService.Before(ctx)
}

func (b *DatabaseService) Finish(err interface{}) {
	// if err != nil {
	// 	b.Db.Commit()
	// } else {
	// 	b.Db.Rollback()
	// }
	b.Db.Close()
	b.DefaultService.Finish(err)
}

type ModelService struct {
	DatabaseService
	Table interface{}
}

func defaultWFunc(ctx *Context) string {
	return string(ctx.QueryArgs().Peek("where"))
}
func defaultOFunc(ctx *Context) string {
	return string(ctx.QueryArgs().Peek("order"))
}
func (m *ModelService) Get(wFunc func(*Context) string) (interface{}, string){
	if m.Table == nil {
		return make([]interface{}, 0), "need Table"
	}
	if wFunc == nil {
		wFunc = defaultWFunc
	}
	one := reflect.New(reflect.TypeOf(m.Table).Elem()).Interface()
	m.Db.Where(wFunc(m.Ctx))
	m.Db.Get(one)
	return one,""
}
func (m *ModelService) Find(wFunc,oFunc func(*Context) string) ([]interface{}, string) {
	if m.Table == nil {
		return make([]interface{}, 0), "need Table"
	}
	if wFunc == nil {
		wFunc = defaultWFunc
	}
	if oFunc==nil {
		oFunc = defaultOFunc
	}
	size := 10
	data := make([]interface{}, size)
	m.Db.Where(wFunc(m.Ctx))
	m.Db.OrderBy(oFunc(m.Ctx))
	m.Db.Limit(size)
	n := 0
	m.Db.Iterate(m.Table, func(i int, item interface{}) error {
		data[i] = item
		n++
		return nil
	})
	return data[:n], ""
}