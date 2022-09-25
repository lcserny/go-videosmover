package core

type Action interface {
	Execute([]byte) (string, error)
}

type ActionRepository interface {
	Register(key string, a Action)
	Retrieve(key string) Action
}
