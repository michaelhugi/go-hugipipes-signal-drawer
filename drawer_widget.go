package go_hugipipes_signal_drawer

type DrawerWidget interface {
	GetWidgetHeight() int
	GetWidgetWidth() int
	Draw(y int)
}
