package esrc

type Store interface {
	Load(id interface{})
}
