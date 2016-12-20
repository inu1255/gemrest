package {{.Model}}

import (
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

func (this *{{Mapper .Name}}Service) Insert(src *{{Mapper .Name}}) (*{{Mapper .Name}}, string) {
	_, err := this.Db.Insert(src)
	if err != nil {
		return nil, err.Error()
	}
	return src, ""
}

func (this *{{Mapper .Name}}Service) Update(id string,src *{{Mapper .Name}}) (*{{Mapper .Name}}, string) {
	n, err := this.Db.Id(id).Update(src)	
	if err != nil {
		return nil, err.Error()
	}
	if n <1 {
		return nil,"没有变化"
	}
	return src, ""
}

func New{{Mapper .Name}}Service() *{{Mapper .Name}}Service {
	inst:= &{{Mapper .Name}}Service{}
	inst.SetTable(&{{Mapper .Name}}{})
	return inst
}

{{end}}