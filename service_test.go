package gemrest

import (
	"encoding/json"
	"testing"

	"github.com/go-gem/gem"
	"github.com/go-gem/tests"
	"github.com/valyala/fasthttp"
)

type Node struct {
	Value int
}

type Service struct {
	DefaultService
	Ta Node
	Tb *Node
}

func (s *Service) A() (int, error) {
	s.Ta.Value++
	return s.Ta.Value, nil
}
func (s *Service) B() (int, error) {
	s.Tb.Value++
	return s.Tb.Value, nil
}

func TestServiceRouter(t *testing.T) {
	var err error
	SetDb("mysql", "root:199337@tcp(114.215.86.245:3306)/movie_baidu")
	Db.ShowSQL(true)
	Bind("/test", &Service{Ta: Node{}, Tb: &Node{}})
	srv := gem.New("", Router.Handler())

	test1 := tests.New(srv)
	test1.Timeout *= 1000
	test1.Method = gem.MethodGet
	test1.Url = "/test/a"
	test1.Expect().Custom(func(r fasthttp.Response) (err error) {
		var v ApiMsg
		err = json.Unmarshal(r.Body(), &v)
		if err != nil {
			t.Error(err)
		} else if v.Data.(float64) != 1 {
			t.Error("/test/a should be 1 but get", v.Data)
		}
		return
	})
	if err = test1.Run(); err != nil {
		t.Error(err)
	}

	count := 0
	test2 := tests.New(srv)
	test2.Timeout *= 1000
	test2.Method = gem.MethodGet
	test2.Url = "/test/b"
	test2.Expect().Custom(func(r fasthttp.Response) (err error) {
		var v ApiMsg
		err = json.Unmarshal(r.Body(), &v)
		count++
		if err != nil {
			t.Error(err)
		} else if int(v.Data.(float64)) != count {
			t.Error("/test/a should be", count, "but get", v.Data)
		}
		return
	})
	if err = test2.Run(); err != nil {
		t.Error(err)
	}
	if err = test2.Run(); err != nil {
		t.Error(err)
	}
}
