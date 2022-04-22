package api

type Api interface {
	StartServing(port int) error
}
