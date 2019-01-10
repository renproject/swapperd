package server

type Server interface {
	Run(doneCh <-chan struct{})
}
