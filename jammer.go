package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var notes = map[string]int{
	"B0":  31,
	"C1":  33,
	"CS1": 35,
	"D1":  37,
	"DS1": 39,
	"E1":  41,
	"F1":  44,
	"FS1": 46,
	"G1":  49,
	"GS1": 52,
	"A1":  55,
	"AS1": 58,
	"B1":  62,
	"C2":  65,
	"CS2": 69,
	"D2":  73,
	"DS2": 78,
	"E2":  82,
	"F2":  87,
	"FS2": 93,
	"G2":  98,
	"GS2": 104,
	"A2":  110,
	"AS2": 117,
	"B2":  123,
	"C3":  131,
	"CS3": 139,
	"D3":  147,
	"DS3": 156,
	"E3":  165,
	"F3":  175,
	"FS3": 185,
	"G3":  196,
	"GS3": 208,
	"A3":  220,
	"AS3": 233,
	"B3":  247,
	"C4":  262,
	"CS4": 277,
	"D4":  294,
	"DS4": 311,
	"E4":  330,
	"F4":  349,
	"FS4": 370,
	"G4":  392,
	"GS4": 415,
	"A4":  440,
	"AS4": 466,
	"B4":  494,
	"C5":  523,
	"CS5": 554,
	"D5":  587,
	"DS5": 622,
	"E5":  659,
	"F5":  698,
	"FS5": 740,
	"G5":  784,
	"GS5": 831,
	"A5":  880,
	"AS5": 932,
	"B5":  988,
	"C6":  1047,
	"CS6": 1109,
	"D6":  1175,
	"DS6": 1245,
	"E6":  1319,
	"F6":  1397,
	"FS6": 1480,
	"G6":  1568,
	"GS6": 1661,
	"A6":  1760,
	"AS6": 1865,
	"B6":  1976,
	"C7":  2093,
	"CS7": 2217,
	"D7":  2349,
	"DS7": 2489,
	"E7":  2637,
	"F7":  2794,
	"FS7": 2960,
	"G7":  3136,
	"GS7": 3322,
	"A7":  3520,
	"AS7": 3729,
	"B7":  3951,
	"C8":  4186,
	"CS8": 4435,
	"D8":  4699,
	"DS8": 4978,
}

// Jammer is an object which can play sounds on your motherboard speaker
type Jammer struct {
	beeper      *Beeper
	tempo       float64
	Lines       []string
	CurrentLine int
}

// NewJammer creates a new instance of the Jammer object
func NewJammer(music string) (*Jammer, error) {
	j := &Jammer{}

	b, err := NewBeeper()
	if err != nil {
		return nil, err
	}

	j.beeper = b
	j.tempo = 1
	j.Lines = strings.Split(music, "\n")
	return j, nil
}

func (j *Jammer) Play() {
	for i := 0; i < len(j.Lines); i++ {
		j.PlayNext()
		j.CurrentLine++
	}
}

func (j *Jammer) PlayNext() {
	line := strings.TrimSpace(j.Lines[j.CurrentLine])

	if line == "" {
		return
	}

	if line[0] == ';' {
		return
	}

	parts := strings.Split(line, " ")
	switch parts[0] {
	case "PAUSE":
		dur, err := strconv.ParseFloat(parts[1], 64)
		j.check(err)
		time.Sleep(time.Duration(int(dur*j.tempo)) * time.Millisecond)
	case "FREQ":
		freq, err := strconv.Atoi(parts[1])
		j.check(err)

		dur, err := strconv.ParseFloat(parts[2], 64)
		j.check(err)

		j.PlayFreq(float32(freq), dur)
	case "TEMPO":
		tempo, err := strconv.Atoi(parts[1])
		j.check(err)

		j.tempo = bpmToMs(tempo)
	case ";":
		// Comment
	default:
		note := parts[0]
		dur, err := strconv.ParseFloat(parts[1], 64)
		j.check(err)

		j.PlayNote(note, dur)
	}
}

func (j *Jammer) check(err error) {
	if err != nil {
		fmt.Println("Error line", j.CurrentLine, ":", j.Lines[j.CurrentLine])
		panic(err)
	}
}

func bpmToMs(bpm int) float64 {
	bps := float64(bpm) / 60.0
	seconds_per_beat := 1.0 / bps
	return seconds_per_beat * 1000
}

// PlayFreq plays a specific frequency for a certain duration on the internal speaker
func (j *Jammer) PlayFreq(freq float32, dur float64) {
	j.beeper.Beep(freq, int(dur*j.tempo))
}

// PlayNote plays a specific note
func (j *Jammer) PlayNote(note string, dur float64) {
	freq, ok := notes[note]
	if !ok {
		panic("Unknown note " + note + " on line " + strconv.Itoa(j.CurrentLine))
	}

	j.PlayFreq(float32(freq), dur)
}
