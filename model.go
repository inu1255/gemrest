package gemrest

import (
	"errors"
	"strconv"

	"github.com/inu1255/gohelper"

	"reflect"
	"regexp"
)

type TableInterface interface {
	GetDetail() interface{}
	GetSearch() interface{}
}

type ModelService struct {
	DatabaseService
	Table TableInterface
}

var (
	orderReg = regexp.MustCompile(`-([^,]+)`)
)

func defaultWhatFunc(ctx *Context) string {
	return string(ctx.QueryArgs().Peek("what"))
}
func defaultWhereFunc(ctx *Context) string {
	return string(ctx.QueryArgs().Peek("where"))
}
func defaultOrderFunc(ctx *Context) string {
	return orderReg.ReplaceAllString(string(ctx.QueryArgs().Peek("order")), `${1} desc`)
}

func NewModelService(t TableInterface) *ModelService {
	return &ModelService{Table: t}
}

func (m *ModelService) SetTable(t TableInterface) {
	m.Table = t
}

func (m *ModelService) Get(whatFunc, whereFunc func(*Context) string) (interface{}, error) {
	if m.Table == nil {
		return make([]interface{}, 0), errors.New("need Table")
	}
	if whatFunc == nil {
		whatFunc = defaultWhatFunc
	}
	if whereFunc == nil {
		whereFunc = defaultWhereFunc
	}
	one := reflect.New(reflect.TypeOf(m.Table).Elem()).Interface()
	m.Db.Cols(whatFunc(m.Ctx))
	m.Db.Where(whereFunc(m.Ctx))
	_, err := m.Db.Get(one)
	return one, err
}
func (m *ModelService) GetById(id string) (interface{}, error) {
	if m.Table == nil {
		return make([]interface{}, 0), errors.New("need Table")
	}
	one := reflect.New(reflect.TypeOf(m.Table).Elem()).Interface()
	m.Db.Id(id).Get(one)
	return one.(TableInterface).GetDetail(), nil
}

func (this *ModelService) DelById(id string) (interface{}, error) {
	if this.Table == nil {
		return make([]interface{}, 0), errors.New("need Table")
	}
	if gohelper.IsZero(id) {
		return nil, errors.New("id错误")
	}
	_, err := this.Db.Id(id).Delete(this.Table)
	return nil, err
}

func (m *ModelService) Find(whatFunc, whereFunc, orderFunc func(*Context) string) ([]interface{}, error) {
	if m.Table == nil {
		return make([]interface{}, 0), errors.New("need Table")
	}
	if whatFunc == nil {
		whatFunc = defaultWhatFunc
	}
	if whereFunc == nil {
		whereFunc = defaultWhereFunc
	}
	if orderFunc == nil {
		orderFunc = defaultOrderFunc
	}
	query := m.Ctx.QueryArgs()
	page, _ := strconv.Atoi(string(query.Peek("page")))
	size, _ := strconv.Atoi(string(query.Peek("size")))
	if size == 0 {
		size = 10
	}
	data := make([]interface{}, size)
	m.Db.Cols(whatFunc(m.Ctx))
	m.Db.Where(whereFunc(m.Ctx))
	m.Db.OrderBy(orderFunc(m.Ctx))
	m.Db.Limit(size, size*page)
	n := 0
	err := m.Db.Iterate(m.Table, func(i int, item interface{}) error {
		data[i] = item.(TableInterface).GetSearch()
		n++
		return nil
	})
	return data[:n], err
}
