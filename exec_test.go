package jet

import (
	"fmt"
	"io/ioutil"
	"sync"
	"testing"
	"time"
)

func TestExecuteConcurrency(t *testing.T) {
	l := NewInMemLoader()
	l.Set("foo", "{{if true}}Hi {{ .Name }}!{{end}}")

	set := NewSet(l)

	tpl, err := set.GetTemplate("foo")
	if err != nil {
		t.Errorf("getting template from set: %v", err)
	}

	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("CC_%d", i), func(t *testing.T) {
			t.Parallel()

			err := tpl.Execute(ioutil.Discard, nil, struct{ Name string }{Name: "John"})
			if err != nil {
				t.Errorf("executing template: %v", err)
			}
		})
	}
}

func TestSetExecute(t *testing.T) {
	l := NewInMemLoader()
	var template = "{{if true}}Hi {{ .Name }} {{ .I }}!{{end}}"
	l.Set("foo", template)

	set := NewSet(l)

	tpl, err := set.GetTemplate("foo")
	if err != nil {
		t.Errorf("getting template from set: %v", err)
	}

	for i := 0; i < 100000000; i++ {

		err := tpl.Execute(ioutil.Discard, nil, struct {
			Name string
			I    int
		}{Name: "John", I: i})
		if err != nil {
			t.Errorf("executing template: %v", err)
		}

	}
}

func TestExecute(t *testing.T) {
	var template = "{{if true}}Hi {{ .Name }} {{ .I }}!{{end}}"
	Janitor(10 * time.Millisecond)
	for i := 0; i < 100000000; i++ {
		_, err := Execute(template, VarMap{}, struct {
			Name string
			I    int
		}{Name: "John", I: i})
		if err != nil {
			t.Errorf("executing template: %v", err)
		}
		//fmt.Println(rendered)

	}
}

func TestRandomExecute(t *testing.T) {
	var template = "{{if true}}Hi {{ .Name }} {{ .I }}!{{end}} {{global}}"
	Janitor(1 * time.Second)

	AddGlobal("global", " Global!!")

	var testSize = 10000000
	var wg = sync.WaitGroup{}
	for i := 0; i < testSize; i++ {
		wg.Add(1)
		go func() {
			var x = ""
			if i < 10000000/2 {
				x = fmt.Sprint(i % 100)
			}
			_, err := Execute(template+x, VarMap{}, struct {
				Name string
				I    int
			}{Name: "John", I: i})
			if err != nil {
				t.Errorf("executing template: %v", err)
			}
			wg.Done()
			time.Sleep(10 * time.Second)
		}()
	}
	wg.Wait()
}
