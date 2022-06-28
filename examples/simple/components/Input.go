package components

import "fmt"

type Input struct {
	Type  string
	Name  string
	Field string
}

func (i *Input) Mount(Type string, name string, field string) {
	i.Type = Type
	i.Name = name
	i.Field = field
}

func (i *Input) Update(field string) {
	fmt.Println("Update", i.Type, i.Name, i.Field)
}
