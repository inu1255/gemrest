package gemrest

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type convertFunc func(*Context, string) reflect.Value

// only new
func newInstCall(t reflect.Type) convertFunc {
	return func(*Context, string) reflect.Value {
		return reflect.New(t)
	}
}

// new and copy
func copyInstCall(src ApiService) convertFunc {
	return func(*Context, string) reflect.Value {
		return CloneValue(reflect.ValueOf(src))
	}
}
func CloneValue(src reflect.Value) reflect.Value {
	var dst reflect.Value
	if src.Kind() == reflect.Ptr || src.Kind() == reflect.Interface {
		src = src.Elem()
	}
	dst = reflect.New(src.Type())
	cloneValue(src, dst.Elem())
	return dst
}
func cloneValue(src, dst reflect.Value) {
	switch src.Kind() {
	case reflect.Struct:
		n := src.NumField()
		for i := 0; i < n; i++ {
			f := src.Field(i)
			cloneValue(f, dst.Field(i))
		}
	default:
		if dst.CanSet() {
			dst.Set(src)
		}
	}
}

func newString(ctx *Context, index string) reflect.Value {
	return reflect.ValueOf(ctx.Param(index))
}
func newInt(ctx *Context, index string) reflect.Value {
	if r, e := strconv.Atoi(ctx.Param(index)); e == nil {
		return reflect.ValueOf(r)
	}
	return reflect.ValueOf(0)
}
func newInt64(ctx *Context, index string) reflect.Value {
	if r, e := strconv.ParseInt(ctx.Param(index), 10, 64); e == nil {
		return reflect.ValueOf(r)
	}
	return reflect.ValueOf(0)
}
func newFloat32(ctx *Context, index string) reflect.Value {
	if r, e := strconv.ParseFloat(ctx.Param(index), 32); e == nil {
		return reflect.ValueOf(float32(r))
	}
	return reflect.ValueOf(float32(0))
}
func newFloat64(ctx *Context, index string) reflect.Value {
	if r, e := strconv.ParseFloat(ctx.Param(index), 64); e == nil {
		return reflect.ValueOf(r)
	}
	return reflect.ValueOf(0.0)
}
func newJsonCall(t reflect.Type) convertFunc {
	return func(ctx *Context, index string) reflect.Value {
		v := reflect.New(t)
		if err := json.Unmarshal(ctx.Request.Body(), v.Interface()); err == nil {
			return v.Elem()
		} else {
			panic(err)
		}
	}
}
func newNilCall(t reflect.Type) convertFunc {
	return func(ctx *Context, index string) reflect.Value {
		return reflect.Zero(t)
	}
}

var (
	name2route = regexp.MustCompile(`([a-z]|^)[A-Z]`)
)

func nameToRoute(from string) string {
	if len(from) == 2 {
		return from[:1] + "-" + strings.ToLower(from[1:])
	}
	return strings.ToLower(from)
}
func convertMethodParams(prefix string, m reflect.Method) (int, string, []convertFunc) {
	numOut := m.Type.NumOut()
	if numOut != 2 || m.Type.Out(1).Kind() != reflect.String {
		return -1, "", nil
	}
	numIn := m.Type.NumIn()
	path := prefix + "/" + name2route.ReplaceAllStringFunc(m.Name, nameToRoute)
	call := make([]convertFunc, numIn)
	flag := 0 // -1:"ignore current" 0:"GET" 1:"POST"
	for i := 1; i < numIn; i++ {
		switch m.Type.In(i).Kind() {
		case reflect.String:
			call[i] = newString
			path += fmt.Sprintf("/:%d", i)
		case reflect.Int:
			call[i] = newInt
			path += fmt.Sprintf("/:%d", i)
		case reflect.Int64:
			call[i] = newInt64
			path += fmt.Sprintf("/:%d", i)
		case reflect.Float32:
			call[i] = newFloat32
			path += fmt.Sprintf("/:%d", i)
		case reflect.Float64:
			call[i] = newFloat64
			path += fmt.Sprintf("/:%d", i)
		case reflect.Struct, reflect.Ptr, reflect.Slice:
			if flag == 1 {
				flag = -1
				break
			} else {
				call[i] = newJsonCall(m.Type.In(i))
				flag = 1
			}
		case reflect.Func:
			call[i] = newNilCall(m.Type.In(i))
		default:
			logger.Println("default", i, m.Type.In(i).Kind())
			flag = -1
		}
		if flag == -1 {
			break
		}
	}
	if conf.DocBind != "" {

	}
	return flag, path, call
}
