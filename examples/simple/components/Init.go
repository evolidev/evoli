package components

type Init struct {
	Name string
}

func (i *Init) Mount() {
	//todo do stuff here on mount
}

func (i *Init) Update(p any) {
	i.Name = p.(string)
}
