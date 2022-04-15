package go_hugipipes_signal_drawer

import (
	"fmt"
	mn "github.com/michaelhugi/go-hugipipes-musical-notes"
	"image"
	"image/color"
	"math"
	"strings"
)

//SpectrumDrawerMark can be used to highlight a frequency in a spectrum
type SpectrumDrawerMark struct {
	frequency float64
	color     color.Color
}

//NewSpectrumDrawerMark is the constructor for SpectrumDrawerMark
func NewSpectrumDrawerMark(frequency float64, color color.Color) *SpectrumDrawerMark {
	return &SpectrumDrawerMark{
		frequency: frequency,
		color:     color,
	}
}

//SpectrumDrawerItems contains a list of plot points to draw in the spectrum. Multiple items can be plotted to one
//spectrum (like amplitude and phase)
type SpectrumDrawerItems struct {
	points   []float64
	drawLine bool
	color    color.Color
}

//NewSpectrumDrawerItems is the constructor for SpectrumDrawerItems
//points are all the data-points. It will be automatically scaled to the plot
//drawLine says, if the items should be plotted as lines from the bottom of the plot (amplitudes) or single points (phases)
//color is the color the plot should have
func NewSpectrumDrawerItems(points []float64, drawLine bool, color color.Color) *SpectrumDrawerItems {
	return &SpectrumDrawerItems{
		points,
		drawLine,
		color,
	}
}

//SpectrumDrawer is a widget that can be used in drawer to draw a Frequency-Spectrum
type SpectrumDrawer struct {
	*DrawerBuilder
	cache           *spectrumDrawerCache
	title           string
	frequencies     []float64
	items           []SpectrumDrawerItems
	marks           []SpectrumDrawerMark
	backgroundColor color.Color
	dividerColor    color.Color
	axisColor       color.Color
	titleColor      color.Color
	temp            mn.MTemperament
	startFreq       float64
	endFreq         float64
}

//NewSpectrumDrawer is the constructor for SpectrumDrawer
func NewSpectrumDrawer(drawer *DrawerBuilder, frequencies []float64, title string) *SpectrumDrawer {
	return &SpectrumDrawer{
		DrawerBuilder:   drawer,
		title:           title,
		frequencies:     frequencies,
		backgroundColor: image.Black.C,
		axisColor:       image.White.C,
		titleColor:      image.White.C,
		items:           make([]SpectrumDrawerItems, 0),
		marks:           make([]SpectrumDrawerMark, 0),
		dividerColor:    gray,
		temp:            mn.NewMTemperamentEqual(440),
		startFreq:       20,
		endFreq:         20000,
	}
}

//spectrumDrawerCache contains data that is recalculated often during drawing
type spectrumDrawerCache struct {
	freqFactor       float64
	calculatedWidth  int
	calculatedHeight int
}

//Temperament sets the temperament of the musical notes for the x-axis. Default is equal at A4=440Hz
func (s *SpectrumDrawer) Temperament(temp mn.MTemperament) *SpectrumDrawer {
	s.temp = temp
	return s
}

//SetItems adds a data-set to be plotted
func (s *SpectrumDrawer) SetItems(items *SpectrumDrawerItems) *SpectrumDrawer {
	s.items = append(s.items, *items)
	return s
}

//SetMark adds a mark to highlight a frequency
func (s *SpectrumDrawer) SetMark(mark *SpectrumDrawerMark) *SpectrumDrawer {
	s.marks = append(s.marks, *mark)
	return s
}

//BackgroundColor sets the background-color of the plot. Default is black
func (s *SpectrumDrawer) BackgroundColor(backgroundColor color.Color) *SpectrumDrawer {
	s.backgroundColor = backgroundColor
	return s
}

//DividerColor sets the divider-color of the plot. Default is gray
func (s *SpectrumDrawer) DividerColor(dividerColor color.Color) *SpectrumDrawer {
	s.dividerColor = dividerColor
	return s
}

//AxisColor sets the color of the axis. Default is white
func (s *SpectrumDrawer) AxisColor(axisColor color.Color) *SpectrumDrawer {
	s.axisColor = axisColor
	return s
}

//TitleColor sets the color of the title of the Spectrum.
func (s *SpectrumDrawer) TitleColor(titleColor color.Color) *SpectrumDrawer {
	s.titleColor = titleColor
	return s
}

//StartFreq sets the lowest shown frequency in the plot. Default is 20Hz
func (s *SpectrumDrawer) StartFreq(startFreq float64) *SpectrumDrawer {
	if startFreq >= s.endFreq {
		return s
	}
	s.startFreq = startFreq
	return s
}

//EndFreq sets the highest shown frequency in the plot. Default is 20kHz
func (s *SpectrumDrawer) EndFreq(endFreq float64) *SpectrumDrawer {
	if s.startFreq >= endFreq {
		return s
	}
	s.endFreq = endFreq
	return s
}

//StartNote sets the lowest shown frequency in the plot. Default is 20Hz
func (s *SpectrumDrawer) StartNote(note mn.MNote) *SpectrumDrawer {
	return s.StartFreq(note.LowerFrequency())
}

//EndNote sets the highest shown frequency in the plot. Default is 20kHz
func (s *SpectrumDrawer) EndNote(note mn.MNote) *SpectrumDrawer {
	return s.EndFreq(note.UpperFrequency())
}

//newSpectrumDrawerCache creates a new cache with pre-calculated values for plotting to avoid executing the same operation multiple times
func (s *SpectrumDrawer) newSpectrumDrawerCache() *spectrumDrawerCache {
	return &spectrumDrawerCache{
		freqFactor:       float64(s.plotWidth) / (s.endFreq - s.startFreq),
		calculatedWidth:  s.plotWidth + 2*s.labelSpace,
		calculatedHeight: s.plotHeight + 2*s.labelSpace,
	}
}

//freqToX recalculates a frequency to the x-coordinates
func (s *SpectrumDrawer) freqToX(freq float64) int {
	if freq < s.startFreq || freq > s.endFreq {
		return -1000
	}

	return int((freq-s.startFreq)*s.cache.freqFactor) + s.labelSpace
}

//drawBackground plots the background
func (s *SpectrumDrawer) drawBackground(y int) {
	top := y
	bottom := top + s.cache.calculatedHeight
	for x := 0; x <= s.cache.calculatedWidth; x++ {
		for y := top; y <= bottom; y++ {
			s.drawable.Set(x, y, s.backgroundColor)
		}
	}
}

//drawXAxis draws the x-axis of the plot
func (s *SpectrumDrawer) drawXAxis(y int) {
	y += s.labelSpace + s.plotHeight
	maxX := s.cache.calculatedWidth - s.labelSpace
	for x := s.labelSpace - s.spacePart; x <= maxX; x++ {
		s.drawable.Set(x, y, s.axisColor)
	}

	//s.drawXAxisOctave(s.temp.Octave(mn.OctaveMinus1), y)
	s.drawXAxisOctave(s.temp.Octave(mn.Octave0), y)
	s.drawXAxisOctave(s.temp.Octave(mn.Octave1), y)
	s.drawXAxisOctave(s.temp.Octave(mn.Octave2), y)
	s.drawXAxisOctave(s.temp.Octave(mn.Octave3), y)
	s.drawXAxisOctave(s.temp.Octave(mn.Octave4), y)
	s.drawXAxisOctave(s.temp.Octave(mn.Octave5), y)
	s.drawXAxisOctave(s.temp.Octave(mn.Octave6), y)
	s.drawXAxisOctave(s.temp.Octave(mn.Octave7), y)
	s.drawXAxisOctave(s.temp.Octave(mn.Octave8), y)
	s.drawXAxisOctave(s.temp.Octave(mn.Octave9), y)
	s.drawStartAndEndFreq(y)
}

//drawStartAndEndFreq draws the frequencies in the x-axis of the lowest and highest frequencies
func (s *SpectrumDrawer) drawStartAndEndFreq(lineTop int) {
	lineBottom := lineTop + 5*s.spacePart
	x1 := s.freqToX(s.startFreq)
	yFreq := lineTop + s.spacePart*7
	x2 := s.freqToX(s.endFreq)

	for y := lineTop; y <= lineBottom; y++ {
		s.drawable.Set(x1, y, s.axisColor)
		s.drawable.Set(x2, y, s.axisColor)
	}
	x1 = x1 + 5
	x2 = x2 - 100
	s.drawable.DrawString(x1, yFreq, fmt.Sprintf("%fHz", s.startFreq), s.axisColor)
	s.drawable.DrawString(x2, yFreq, fmt.Sprintf("%fHz", s.endFreq), s.axisColor)

}

//drawXAxisOctave draws one musical octave in the x-axis
func (s *SpectrumDrawer) drawXAxisOctave(oct mn.MOctave, lineTop int) {
	//Draw line
	x1 := s.freqToX(oct.Note(mn.C).ExactFrequency())
	if x1 < 0 {
		return
	}
	x2 := s.freqToX(oct.Note(mn.C).ExactFrequency() * 2)
	lineBottom := lineTop + 4*s.spacePart
	for y := lineTop; y <= lineBottom; y++ {
		s.drawable.Set(x1, y, s.axisColor)
		s.drawable.Set(x2, y, s.axisColor)
	}
	notes := oct.AllNotes()
	for _, note := range notes {
		s.drawXAxisNote(note, lineTop)
	}
}

//draw draws all content to the drawable
func (s *SpectrumDrawer) draw(y int) {
	s.cache = s.newSpectrumDrawerCache()
	s.drawBackground(y)
	s.drawPlotTitle(s.title, s.spacePart*3+y)
	for _, mark := range s.marks {
		s.drawMark(mark, y)
	}
	for _, item := range s.items {
		s.drawItem(item, y)
	}
	s.drawXAxis(y)
	s.drawYAxis(y)
	if y > 0 {
		s.drawDivider(y)
	}

}

//drawMark draws a line to highlight a special frequency
func (s *SpectrumDrawer) drawMark(mark SpectrumDrawerMark, y int) {
	x := s.freqToX(mark.frequency)
	bottom := y + s.plotHeight + s.labelSpace
	top := y + s.labelSpace
	for y := top; y <= bottom; y++ {
		s.drawable.Set(x, y, mark.color)
	}
}

//drawItem draws the plot-points of a points set to the spectrum
func (s *SpectrumDrawer) drawItem(item SpectrumDrawerItems, y int) {
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
	for i, f := range s.frequencies {
		it := item.points[i]
		x := s.freqToX(f)
		if x > 0 {
			yPoint := it + offset
			yPoint = yPoint * factor
			YPoint := bottom - int(yPoint)
			if item.drawLine {
				for y := bottom; y >= YPoint; y-- {
					s.drawable.Set(x, y, item.color)
				}
			} else {
				s.drawable.Set(x, YPoint, item.color)
			}
		}
	}
}

//drawDivider draws a horizontal line a the end of the plot
func (s *SpectrumDrawer) drawDivider(y int) {
	for x := 0; x <= s.cache.calculatedWidth; x++ {
		s.drawable.Set(x, y, s.dividerColor)
	}
}

//Draws a musical note to the x-axis
func (s *SpectrumDrawer) drawXAxisNote(n mn.MNote, lineTop int) {
	x1 := s.freqToX(n.ExactFrequency())
	lineBottom := lineTop + s.spacePart
	for y := lineTop; y <= lineBottom; y++ {
		s.drawable.Set(x1, y, s.axisColor)
	}
	y := lineBottom + s.spacePart + 3
	x := x1 - 2
	if !strings.Contains(n.String(), "#") {
		s.drawable.DrawString(x+5, y, n.String(), s.axisColor)
		y += s.spacePart*2 + 3
		s.drawable.DrawString(x+5, y, fmt.Sprintf("%d", n.MidiNoteNumber()), s.axisColor)
	}

}

//Draws the plot title
func (s *SpectrumDrawer) drawPlotTitle(title string, lineTop int) {
	x := s.labelSpace
	y := lineTop + 3*s.spacePart
	s.drawable.DrawString(x, y, title, s.titleColor)
}

//getWidgetWidth implements Widget interface
func (s *SpectrumDrawer) getWidgetWidth() int {
	s.cache = s.newSpectrumDrawerCache()
	return s.cache.calculatedWidth
}

//getWidgetHeight implements Widget interface
func (s *SpectrumDrawer) getWidgetHeight() int {
	s.cache = s.newSpectrumDrawerCache()
	return s.cache.calculatedHeight
}

//Draws the y axis
func (s *SpectrumDrawer) drawYAxis(top int) {
	top += s.labelSpace
	bottom := top + s.plotHeight + s.spacePart
	x := s.labelSpace

	for y := top; y <= bottom; y++ {
		s.drawable.Set(x, y, s.axisColor)
	}

}
