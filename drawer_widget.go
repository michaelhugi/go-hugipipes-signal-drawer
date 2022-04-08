package go_hugipipes_signal_drawer

type DrawerWidget interface {
	GetHeight() int
	GetWidth() int
	Draw(y int)
}
