package go_hugipipes_signal_drawer

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"time"
)

//WaveDrawerItems contains a list of plot points to draw in the spectrum. Multiple items can be plotted to one
//plot (like left and right channel)
type WaveDrawerItems struct {
	points []float64
	color  color.Color
}

//NewWaveDrawerItems is the constructor for WaveDrawerItems
//points are all the data-points. It will be automatically scaled to the plot
//color is the color the plot should have
func NewWaveDrawerItems(points []float64, color color.Color) *WaveDrawerItems {
	return &WaveDrawerItems{
		points,
		color,
	}
}

//WaveDrawer is a widget that can be used in drawer to draw a time-based wave signal
type WaveDrawer struct {
	*DrawerBuilder
	cache           *waveDrawerCache
	title           string
	times           []time.Duration
	items           []WaveDrawerItems
	backgroundColor color.Color
	dividerColor    color.Color
	axisColor       color.Color
	titleColor      color.Color
	startTime       time.Duration
	endTime         time.Duration
}

//NewWaveDrawer is the constructor for WaveDrawer
func NewWaveDrawer(drawer *DrawerBuilder, times []time.Duration, title string) *WaveDrawer {
	return &WaveDrawer{
		DrawerBuilder:   drawer,
		title:           title,
		times:           times,
		backgroundColor: image.Black.C,
		axisColor:       image.White.C,
		titleColor:      image.White.C,
		items:           make([]WaveDrawerItems, 0),
		dividerColor:    gray,
		startTime:       times[0],
		endTime:         times[len(times)-1],
	}
}

//waveDrawerCache contains data that would be recalculated often during drawing
type waveDrawerCache struct {
	timeFactor       float64
	calculatedWidth  int
	calculatedHeight int
}

//SetItems adds a data-set to be plotted
func (s *WaveDrawer) SetItems(items *WaveDrawerItems) *WaveDrawer {
	s.items = append(s.items, *items)
	return s
}

//BackgroundColor sets the background-color of the plot. Default is black
func (s *WaveDrawer) BackgroundColor(backgroundColor color.Color) *WaveDrawer {
	s.backgroundColor = backgroundColor
	return s
}

//DividerColor sets the divider-color of the plot. Default is gray
func (s *WaveDrawer) DividerColor(dividerColor color.Color) *WaveDrawer {
	s.dividerColor = dividerColor
	return s
}

//AxisColor sets the color of the axis. Default is white
func (s *WaveDrawer) AxisColor(axisColor color.Color) *WaveDrawer {
	s.axisColor = axisColor
	return s
}

//TitleColor sets the color of the title of the Spectrum.
func (s *WaveDrawer) TitleColor(titleColor color.Color) *WaveDrawer {
	s.titleColor = titleColor
	return s
}

//StartTime sets the start time for the plot. Default is 0
func (s *WaveDrawer) StartTime(startTime time.Duration) *WaveDrawer {
	if startTime.Milliseconds() >= s.endTime.Milliseconds() {
		return s
	}
	s.startTime = startTime
	return s
}

//EndTime sets the highest shown time in the plot. Default is the latest plot point provided
func (s *WaveDrawer) EndTime(endTime time.Duration) *WaveDrawer {
	if s.startTime.Milliseconds() >= endTime.Milliseconds() {
		return s
	}
	s.endTime = endTime
	return s
}

//newSpectrumDrawerCache creates a new cache with pre-calculated values for plotting to avoid executing the same operation multiple times
func (s *WaveDrawer) newWaveDrawerCache() *waveDrawerCache {
	return &waveDrawerCache{
		timeFactor:       float64(s.plotWidth) / float64(s.endTime.Nanoseconds()-s.startTime.Nanoseconds()),
		calculatedWidth:  s.plotWidth + 2*s.labelSpace,
		calculatedHeight: s.plotHeight + 2*s.labelSpace,
	}
}

//freqToX recalculates a frequency to the x-coordinates
func (s *WaveDrawer) timeToX(time time.Duration) int {
	if time.Milliseconds() < s.startTime.Milliseconds() || time.Milliseconds() > s.endTime.Milliseconds() {
		return -1000
	}
	t := time - s.startTime
	return int(float64(t.Nanoseconds())*s.cache.timeFactor) + s.labelSpace
}

//drawBackground plots the background
func (s *WaveDrawer) drawBackground(y int) {
	top := y
	bottom := top + s.cache.calculatedHeight
	for x := 0; x <= s.cache.calculatedWidth; x++ {
		for y := top; y <= bottom; y++ {
			s.drawable.Set(x, y, s.backgroundColor)
		}
	}
}

//drawXAxis draws the x-axis of the plot
func (s *WaveDrawer) drawXAxis(y int) {
	y += s.labelSpace + (s.plotHeight / 2)
	maxX := s.cache.calculatedWidth - s.labelSpace
	for x := s.labelSpace - s.spacePart; x <= maxX; x++ {
		s.drawable.Set(x, y, s.axisColor)
	}
	y += s.plotHeight / 2
	dt := s.endTime - s.startTime
	dt = dt / 5
	tt := s.startTime

	for tt.Milliseconds() <= s.endTime.Milliseconds() {
		s.drawTime(tt, y)
		tt += dt
	}
}

//drawTime draws a time-label to the x-axis
func (s *WaveDrawer) drawTime(t time.Duration, lineY int) {
	x := s.timeToX(t)
	bottom := lineY + s.spacePart*3
	for y := lineY; y <= bottom; y++ {
		s.drawable.Set(x, y, s.axisColor)
	}
	s.drawable.DrawString(x, bottom+s.spacePart*2, fmt.Sprintf("%dms", t.Milliseconds()), s.axisColor)
}

//draw draws all content to the drawable
func (s *WaveDrawer) draw(y int) {
	s.cache = s.newWaveDrawerCache()
	s.drawBackground(y)
	s.drawPlotTitle(s.title, s.spacePart*3+y)
	for _, item := range s.items {
		s.drawItem(item, y)
	}
	s.drawXAxis(y)
	s.drawYAxis(y)
	if y > 0 {
		s.drawDivider(y)
	}

}

//drawItem draws the plot-points of a points set to the wave
func (s *WaveDrawer) drawItem(item WaveDrawerItems, y int) {
	maxValue := item.points[0]
	minValue := item.points[0]
	for _, v := range item.points {
		maxValue = math.Max(maxValue, v)
		minValue = math.Min(minValue, v)
	}
	offset := -minValue
	maxValue += offset

	factor := float64(s.plotHeight) / maxValue

	bottom := y + s.labelSpace + s.plotHeight
	for i, t := range s.times {
		if t.Milliseconds() >= s.startTime.Milliseconds() && t.Milliseconds() <= s.endTime.Milliseconds() {
			it := item.points[i]
			x := s.timeToX(t)
			if x > 0 {
				yPoint := it + offset
				yPoint = yPoint * factor
				YPoint := bottom - int(yPoint)
				s.drawable.Set(x, YPoint, item.color)
			}
		}
	}
}

//drawDivider draws a horizontal line a the end of the plot
func (s *WaveDrawer) drawDivider(y int) {
	for x := 0; x <= s.cache.calculatedWidth; x++ {
		s.drawable.Set(x, y, s.dividerColor)
	}
}

//Draws the plot title
func (s *WaveDrawer) drawPlotTitle(title string, lineTop int) {
	x := s.labelSpace
	y := lineTop + 3*s.spacePart
	s.drawable.DrawString(x, y, title, s.titleColor)
}

//getWidgetWidth implements Widget interface
func (s *WaveDrawer) getWidgetWidth() int {
	s.cache = s.newWaveDrawerCache()
	return s.cache.calculatedWidth
}

//getWidgetHeight implements Widget interface
func (s *WaveDrawer) getWidgetHeight() int {
	s.cache = s.newWaveDrawerCache()
	return s.cache.calculatedHeight
}

//Draws the y axis
func (s *WaveDrawer) drawYAxis(top int) {
	top += s.labelSpace
	bottom := top + s.plotHeight + s.spacePart
	x := s.labelSpace

	for y := top; y <= bottom; y++ {
		s.drawable.Set(x, y, s.axisColor)
	}

}
