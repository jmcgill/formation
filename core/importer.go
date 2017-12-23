package core

type KnownResource struct {
	ID   string
	Name string
}

type Importer interface {
	List() []*KnownResource
}
