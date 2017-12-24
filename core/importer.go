package core

type ResourceDescription struct {
	Links     map[string]string
	Defaults  map[string]Default
	Instances []*Instance
}

type Default struct {
	Value  string
	IsBool bool
}

type Instance struct {
	ID   string
	Name string
}

type Importer interface {
	Describe() ResourceDescription
}
