package gemrest

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"

	"github.com/go-gem/gem"
)

type ApiMsg struct {
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type DefaultService struct {
	Ctx *gem.Context
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
func (d *DefaultService) Before(ctx *gem.Context) bool {
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
func stack() []byte {
	buf := make([]byte, 10240)
	n := runtime.Stack(buf, false)
	return buf[:n]
}
func printStack() {
	log.Println(string(stack()))
	log.Println("end here.")
}
