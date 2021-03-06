package gemrest

import (
	"encoding/json"
	"fmt"
	"regexp"
	"runtime"

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
	logger.Println(string(stack()))
	logger.Println("end here.")
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

var re_origin = regexp.MustCompile(`https?://`)

func (m *CrosService) Before(ctx *Context) bool {
	origin := string(ctx.Request.Header.Peek("Origin"))
	if origin != "" {
		ctx.Response.Header.Add("Access-Control-Allow-Origin", re_origin.ReplaceAllString(origin, ""))
		ctx.Response.Header.Add("Access-Control-Allow-Credentials", "true")
		ctx.Response.Header.Add("Access-Control-Allow-Headers", "x-auth-token,content-type")
	}
	return m.DefaultService.Before(ctx)
}

type DatabaseService struct {
	CrosService
	Db *xorm.Session
}

// database support
func (b *DatabaseService) Before(ctx *Context) bool {
	b.Db = Db.NewSession()
	// b.Db.Begin()
	return b.CrosService.Before(ctx)
}

func (b *DatabaseService) Finish(err interface{}) {
	// if err != nil {
	// 	b.Db.Commit()
	// } else {
	// 	b.Db.Rollback()
	// }
	b.Db.Close()
	b.CrosService.Finish(err)
}
