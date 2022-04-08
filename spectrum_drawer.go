package go_hugipipes_signal_drawer

import (
	"fmt"
	mn "github.com/michaelhugi/go-hugipipes-musical-notes"
	"image"
	"image/color"
	"math"
	"strings"
)

const xLogScaleFactor = 600
const freqLeftOffset = 2500

type SpectrumDrawerItems struct {
	points   []float64
	drawLine bool
	color    color.RGBA
}

func NewSpectrumDrawerItems(points []float64, drawLine bool, color color.RGBA) *SpectrumDrawerItems {
	return &SpectrumDrawerItems{
		points,
		drawLine,
		color,
	}
}

type SpectrumDrawerBuilder struct {
	*DrawerBuilder
	title           string
	frequencies     []float64
	items           []SpectrumDrawerItems
	freqLogarithmic bool
	backgroundColor color.Color
	axisColor       color.Color
	titleColor      color.Color
	temp            mn.MTemperament
}

func NewSpectrumDrawer(drawer *DrawerBuilder, frequencies []float64) *SpectrumDrawerBuilder {
	return &SpectrumDrawerBuilder{
		DrawerBuilder:   drawer,
		title:           "No Title set",
		frequencies:     frequencies,
		backgroundColor: image.Black.C,
		axisColor:       image.White.C,
		titleColor:      image.White.C,
		items:           make([]SpectrumDrawerItems, 0),
		temp:            mn.NewMTemperamentEqual(440),
	}
}
func (s *SpectrumDrawerBuilder) Temperament(temp mn.MTemperament) *SpectrumDrawerBuilder {
	s.temp = temp
	return s
}
func (s *SpectrumDrawer) SetItems(items SpectrumDrawerItems) *SpectrumDrawer {
	s.items = append(s.items, items)
	return s
}
func (s *SpectrumDrawerBuilder) BackgroundColor(backgroundColor color.Color) *SpectrumDrawerBuilder {
	s.backgroundColor = backgroundColor
	return s
}
func (s *SpectrumDrawerBuilder) AxisColor(axisColor color.Color) *SpectrumDrawerBuilder {
	s.axisColor = axisColor
	return s
}
func (s *SpectrumDrawerBuilder) TitleColor(titleColor color.Color) *SpectrumDrawerBuilder {
	s.titleColor = titleColor
	return s
}
func (s *SpectrumDrawerBuilder) LogarithmicFreq() *SpectrumDrawerBuilder {
	s.freqLogarithmic = true
	return s
}

func (s *SpectrumDrawerBuilder) Title(title string) *SpectrumDrawerBuilder {
	s.title = title
	return s
}

func (s *SpectrumDrawerBuilder) Build() *SpectrumDrawer {
	plotWidth := len(s.frequencies)
	if s.freqLogarithmic {
		plotWidth = int(math.Log2(float64(len(s.frequencies)))*xLogScaleFactor) - freqLeftOffset
	}

	calculatedWidth := plotWidth + 2*s.labelSpace

	calculatedHeight := s.plotHeight + 2*s.labelSpace
	maxFrequency := s.frequencies[len(s.frequencies)-1]
	freqFactor := float64(plotWidth) / maxFrequency

	return &SpectrumDrawer{
		calculatedWidth:  calculatedWidth,
		calculatedHeight: calculatedHeight,
		freqFactor:       freqFactor,
		maxFrequency:     maxFrequency,
	}
}

type SpectrumDrawer struct {
	*SpectrumDrawerBuilder
	calculatedWidth  int
	calculatedHeight int
	freqFactor       float64
	maxFrequency     float64
}

func (s *SpectrumDrawer) freqToX(freq float64) int {
	if !s.freqLogarithmic {
		return int(s.freqFactor*freq) + s.labelSpace
	}
	return int(math.Log2(float64(int(s.freqFactor*freq)+s.labelSpace))*xLogScaleFactor) - freqLeftOffset
}

func (s *SpectrumDrawer) drawBackground() {
	for x := 0; x <= s.calculatedWidth; x++ {
		for y := 0; y <= s.calculatedHeight; y++ {
			s.drawable.Set(x, y, s.backgroundColor)
		}
	}
}
func (s *SpectrumDrawer) drawXAxis(y int) {
	y += s.labelSpace + s.plotHeight
	for x := s.labelSpace - s.spacePart; x <= s.calculatedWidth; x++ {
		s.drawable.Set(x, y, s.axisColor)
	}

	s.drawXAxisOctave(s.temp.Octave(mn.OctaveMinus1), y)
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
}
func (s *SpectrumDrawer) drawXAxisOctave(oct mn.MOctave, lineTop int) {
	//Draw line
	x1 := s.freqToX(oct.LowerFrequency())
	x2 := s.freqToX(oct.UpperFrequency())
	lineBottom := lineTop + 5*s.spacePart
	for y := lineTop; y <= lineBottom; y++ {
		s.drawable.Set(x1, y, s.axisColor)
		s.drawable.Set(x2, y, s.axisColor)
	}
	notes := oct.AllNotes()
	if len(notes) == 0 {
		y := lineTop + 2*s.spacePart
		xCenter := s.freqToX((oct.LowerFrequency()+oct.UpperFrequency())/2) - 3
		s.drawable.DrawString(xCenter, y, oct.String(), s.axisColor)
	}
	if oct.Octave() > mn.Octave1 {
		xFreq := s.freqToX(oct.LowerFrequency()) + 5
		yFreq := lineTop + s.spacePart*5
		s.drawable.DrawString(xFreq, yFreq, fmt.Sprintf("%fHz", oct.LowerFrequency()), blue)
	}
	for _, note := range notes {
		s.drawXAxisNote(note, lineTop)
	}
}

func (s *SpectrumDrawer) drawDivider(y int) {
	for x := 0; x <= s.calculatedWidth; x++ {
		s.drawable.Set(x, y, gray)
	}
}

func (s *SpectrumDrawer) drawXAxisNote(n mn.MNote, lineTop int) {
	x1 := s.freqToX(n.ExactFrequency())
	lineBottom := lineTop + s.spacePart
	for y := lineTop; y <= lineBottom; y++ {
		s.drawable.Set(x1, y, s.axisColor)
	}
	y := lineBottom + s.spacePart + 3
	x := x1 - 2
	if !strings.Contains(n.String(), "#") {
		s.drawable.DrawString(x, y, n.String(), s.axisColor)
	}

	y += s.spacePart + 3
	s.drawable.DrawString(x, y, fmt.Sprintf("%d", n.MidiNoteNumber()), s.axisColor)
}

func (s *SpectrumDrawer) drawPlotTitle(title string, lineTop int) {
	x := s.labelSpace
	y := lineTop + 3*s.spacePart
	s.drawable.DrawString(x, y, title, s.titleColor)
}

func (s *SpectrumDrawer) drawYAxis(top int, labels []xlabel) {
	top += s.labelSpace
	bottom := top + s.plotHeight + s.spacePart
	x := s.labelSpace

	for y := top; y <= bottom; y++ {
		s.drawable.Set(x, y, s.axisColor)
	}
	if labels != nil {
		for _, label := range labels {
			for x := s.labelSpace - s.spacePart; x < s.labelSpace; x++ {
				s.drawable.Set(x, label.Y+top, image.White)
			}
			s.drawable.DrawString(s.labelSpace-4*s.spacePart, label.Y+top-3, label.Text, image.White)
		}
	}
}

type xlabel struct {
	Y    int
	Text string
}
