package gemrest

import (
	"log"
	"strconv"

	"reflect"
	"regexp"
)

type TableInterface interface {
	TableName() string
}

type ModelService struct {
	DatabaseService
	Table TableInterface
}

var (
	orderReg = regexp.MustCompile(`-([^,]+)`)
)

func defaultWFunc(ctx *Context) string {
	return string(ctx.QueryArgs().Peek("where"))
}
func defaultOFunc(ctx *Context) string {
	return orderReg.ReplaceAllString(string(ctx.QueryArgs().Peek("order")), `${1} desc`)
}

func NewModelService(t TableInterface) *ModelService {
	return &ModelService{Table: t}
}

func (m *ModelService) Get(wFunc func(*Context) string) (interface{}, string) {
	if m.Table == nil {
		return make([]interface{}, 0), "need Table"
	}
	if wFunc == nil {
		wFunc = defaultWFunc
	}
	one := reflect.New(reflect.TypeOf(m.Table).Elem()).Interface()
	m.Db.Where(wFunc(m.Ctx))
	m.Db.Get(one)
	return one, ""
}
func (m *ModelService) GetById(id int) (interface{}, string) {
	if m.Table == nil {
		return make([]interface{}, 0), "need Table"
	}
	one := reflect.New(reflect.TypeOf(m.Table).Elem()).Interface()
	m.Db.Id(id).Get(one)
	return one, ""
}
func (m *ModelService) Find(wFunc, oFunc func(*Context) string) ([]interface{}, string) {
	if m.Table == nil {
		return make([]interface{}, 0), "need Table"
	}
	if wFunc == nil {
		wFunc = defaultWFunc
	}
	if oFunc == nil {
		oFunc = defaultOFunc
	}
	query := m.Ctx.QueryArgs()
	page, _ := strconv.Atoi(string(query.Peek("page")))
	size, _ := strconv.Atoi(string(query.Peek("size")))
	if size == 0 {
		size = 10
	}
	data := make([]interface{}, size)
	m.Db.Where(wFunc(m.Ctx))
	m.Db.OrderBy(oFunc(m.Ctx))
	m.Db.Limit(size, size*page)
	n := 0
	log.Println(size)
	m.Db.Iterate(m.Table, func(i int, item interface{}) error {
		data[i] = item
		n++
		return nil
	})
	return data[:n], ""
}
