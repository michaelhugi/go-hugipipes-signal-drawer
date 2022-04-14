package go_hugipipes_signal_drawer

type DrawerBuilder struct {
	plotHeight int
	plotWidth  int
	labelSpace int
	spacePart  int
	plots      []DrawerWidget
	drawable   Drawable
}

func NewDrawer() *DrawerBuilder {
	return &DrawerBuilder{
		plotHeight: 300,
		plotWidth:  2000,
		labelSpace: 80,
		spacePart:  10,
		plots:      make([]DrawerWidget, 0),
	}

}

func (s *DrawerBuilder) Build() *Drawer {
	self := newDrawer(s)
	return self
}

func (s *DrawerBuilder) LabelSpace(labelSpace int) *DrawerBuilder {
	s.labelSpace = labelSpace
	s.spacePart = labelSpace / 8
	return s
}
func (s *DrawerBuilder) PlotHeight(plotHeight int) *DrawerBuilder {
	s.plotHeight = plotHeight
	return s
}

func (s *DrawerBuilder) AddPlot(plot DrawerWidget) *DrawerBuilder {
	s.plots = append(s.plots, plot)
	return s
}

func (s *DrawerBuilder) GetHeight() int {
	h := 0
	for _, p := range s.plots {
		h += p.GetWidgetHeight()
	}
	return h
}
func (s *DrawerBuilder) GetWidth() int {
	w := 0
	for _, p := range s.plots {
		w = max(w, p.GetWidgetWidth())
	}
	return w
}

func (s *DrawerBuilder) SetDrawable(drawable Drawable) *DrawerBuilder {
	s.drawable = drawable
	return s
}

type Drawer struct {
	*DrawerBuilder
}

func newDrawer(builder *DrawerBuilder) *Drawer {
	return &Drawer{
		DrawerBuilder: builder,
	}
}

func (s *Drawer) Draw() {
	y := 0
	for _, p := range s.plots {
		p.Draw(y)
		y += p.GetWidgetHeight()
	}
}

func max(one int, two int) int {
	if one > two {
		return one
	}
	return two
}
