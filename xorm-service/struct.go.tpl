package {{.Model}}

import (
	"errors"
	"github.com/inu1255/gemrest"
	"{{.Path}}model"
)

{{range .Tables}}
type {{Mapper .Name}}Service struct {
	gemrest.ModelService
}

func (this *{{Mapper .Name}}Service) Insert(src *model.{{Mapper .Name}}) (*model.{{Mapper .Name}}, error) {
	_, err := this.Db.Insert(src)
	return src, err
}

func (this *{{Mapper .Name}}Service) Update(id string,src *model.{{Mapper .Name}}) (*model.{{Mapper .Name}}, error) {
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
	inst.SetTable(&model.{{Mapper .Name}}{})
	return inst
}

var _ = gemrest.Reg(New{{Mapper .Name}}Service())
{{end}}