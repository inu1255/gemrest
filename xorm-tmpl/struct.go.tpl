package {{.Model}}

import (
	"errors"
	"github.com/inu1255/gemrest"
	{{range .Imports}}"{{.}}"{{end}}
)

{{range .Tables}}
type {{Mapper .Name}}Search struct{
{{$table := .}}
{{range .ColumnsSeq}}{{$col := $table.GetColumn .}}	{{Mapper $col.Name}}	{{Type $col}} {{Tag $table $col}}
{{end}}
}
type {{Mapper .Name}} struct {
{{Mapper .Name}}Search `xorm:"extends"`
}
type {{Mapper .Name}}Detail struct{
    {{Mapper .Name}}
}

func (this *{{Mapper .Name}})GetDetail()interface{}{
	return {{Mapper .Name}}Detail{ {{Mapper .Name}}:*this}
}

func (this *{{Mapper .Name}})GetSearch()interface{}{
    return this.{{Mapper .Name}}Search
}

type {{Mapper .Name}}Service struct {
	gemrest.ModelService
}

func (this *{{Mapper .Name}}Service) Insert(src *{{Mapper .Name}}) (*{{Mapper .Name}}, error) {
	_, err := this.Db.Insert(src)
	return src, err
}

func (this *{{Mapper .Name}}Service) Update(id string,src *{{Mapper .Name}}) (*{{Mapper .Name}}, error) {
	n, err := this.Db.Id(id).Update(src)	
	if err != nil {
		return nil, err
	}
	if n <1 {
		return nil,errors.New("没有变化")
	}
	return src, nil
}

func New{{Mapper .Name}}Service() *{{Mapper .Name}}Service {
	inst:= &{{Mapper .Name}}Service{}
	inst.SetTable(&{{Mapper .Name}}{})
	return inst
}

{{end}}