package use

type Facade struct {
	values *Collection[string, any]
}

var facades = NewFacade()

func GetFacade(name string) any {
	return facades.Get(name)
}

func AddFacade(name string, facade any) {
	facades.Add(name, facade)
}

func NewFacade() *Facade {
	return &Facade{
		values: NewCollection[string, any](),
	}
}

func (f *Facade) Add(name string, facade any) {
	f.values.Add(name, facade)
}

func (f *Facade) Get(name string) any {
	return f.values.Get(name)
}
