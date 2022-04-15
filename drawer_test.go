package go_hugipipes_signal_drawer

import (
	"testing"
	"time"
)

func TestWidgets(t *testing.T) {
	defer func() {
		recover()
	}()
	spec := NewSpectrumDrawer(nil, make([]float64, 0), "")
	checkDrawerWidgetInterface(spec)
	tim := NewWaveDrawer(nil, make([]time.Duration, 0), "")
	checkDrawerWidgetInterface(tim)
}

func checkDrawerWidgetInterface(i DrawerWidget) {

}
