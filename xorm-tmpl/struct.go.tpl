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

func (this *{{Mapper .Name}})TableName()string{
	return "{{.Name}}"
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

func New{{Mapper .Name}}Service() *{{Mapper .Name}}Service {
	inst:= &{{Mapper .Name}}Service{}
	inst.SetTable(&{{Mapper .Name}}{})
	return inst
}

{{end}}