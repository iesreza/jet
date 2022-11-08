package jet

import (
	"fmt"
	"io/ioutil"
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
	Janitor(10 * time.Millisecond)

	AddGlobal("global", " Global!!")

	var testSize = 1000
	for i := 0; i < testSize; i++ {
		rendered, err := Execute(template, VarMap{}, struct {
			Name string
			I    int
		}{Name: "John", I: i})
		if err != nil {
			t.Errorf("executing template: %v", err)
		}
		fmt.Println(rendered)

	}
}
