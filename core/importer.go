package core

type Instance struct {
	ID   string
	Name string
}

type Importer interface {
	Describe(meta interface{}) ([]*Instance, error)
	Links() map[string]string
}