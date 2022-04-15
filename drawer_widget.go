package go_hugipipes_signal_drawer

type DrawerWidget interface {
	getWidgetHeight() int
	getWidgetWidth() int
	draw(y int)
}
