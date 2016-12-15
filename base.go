package gemrest

import (
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

func (d *DefaultService) Before(ctx *gem.Context) bool {
	d.Ctx = ctx
	return true
}
func (d *DefaultService) After(ctx *gem.Context, data interface{}, msg string) {
	ctx.JSON(200, ApiMsg{data, msg})
}
func (d *DefaultService) Finish(ctx *gem.Context, err interface{}) {
	if err != nil {
		printStack()
		ctx.JSON(500, ApiMsg{err, fmt.Sprintf("Internal Server Error:%v", err)})
	}
}
func stack() []byte {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	return buf[:n]
}
func printStack() {
	log.Println(string(stack()))
	log.Println("end here.")
}
