package main

import (
	"fmt"
	"log"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/ui"
)

type logStyle int

const (
	logNormal logStyle = iota
	logCritic
	logNotable
	logDamage
	logSpecial
	logStatusEnd
	logError
	logConfirm
)

type logEntry struct {
	Text  string
	MText string
	Index int
	Tick  bool
	Style logStyle
	Dups  int
}

func (e logEntry) String() string {
	tick := ""
	if e.Tick {
		tick = "@t•@N "
	}
	s := e.Text
	if e.Dups > 0 {
		s += fmt.Sprintf(" (%d×)", e.Dups+1)
	}
	r := e.Style.Rune()
	if r != 0 {
		s = fmt.Sprintf("%s@%c%s@N", tick, r, s)
	} else {
		s = fmt.Sprintf("%s%s", tick, s)
	}
	return s
}

func (e logEntry) dumpString() string {
	tick := ""
	if e.Tick {
		tick = "• "
	}
	s := e.Text
	if e.Dups > 0 {
		s += fmt.Sprintf(" (%d×)", e.Dups+1)
	}
	s = fmt.Sprintf("%s%s", tick, s)
	return s
}

func (g *game) Print(s string) {
	e := logEntry{Text: s, Index: g.LogIndex}
	g.PrintEntry(e)
}

func (g *game) PrintStyled(s string, style logStyle) {
	e := logEntry{Text: s, Index: g.LogIndex, Style: style}
	g.PrintEntry(e)
}

func (g *game) Printf(format string, a ...interface{}) {
	e := logEntry{Text: fmt.Sprintf(format, a...), Index: g.LogIndex}
	g.PrintEntry(e)
}

func (g *game) PrintfStyled(format string, style logStyle, a ...interface{}) {
	e := logEntry{Text: fmt.Sprintf(format, a...), Index: g.LogIndex, Style: style}
	g.PrintEntry(e)
}

func (g *game) PrintEntry(e logEntry) {
	if e.Index == g.LogNextTick {
		e.Tick = true
	}
	if !e.Tick && len(g.Log) > 0 {
		le := g.Log[len(g.Log)-1]
		if le.Text == e.Text {
			le.Dups++
			le.MText = le.String()
			g.Log[len(g.Log)-1] = le
			return
		}
	}
	e.MText = e.String()
	if LogGame {
		log.Printf("Depth %d:Turn %d:%v", g.Depth, g.Turn, e.dumpString())
	}
	g.Log = append(g.Log, e)
	g.LogIndex++
	if len(g.Log) > 100000 {
		g.Log = g.Log[10000:]
	}
}

func (g *game) StoryPrint(s string) {
	g.Stats.Story = append(g.Stats.Story, fmt.Sprintf("Depth %2d|Turn %5d| %s", g.Depth, g.Turn, s))
}

func (g *game) StoryPrintf(format string, a ...interface{}) {
	g.Stats.Story = append(g.Stats.Story, fmt.Sprintf("Depth %2d|Turn %5d| %s", g.Depth, g.Turn, fmt.Sprintf(format, a...)))
}

func (g *game) CrackSound() (text string) {
	switch RandInt(4) {
	case 0:
		text = "Crack!"
	case 1:
		text = "Crash!"
	case 2:
		text = "Crunch!"
	case 3:
		text = "Creak!"
	}
	return text
}

func (g *game) ExplosionSound() (text string) {
	switch RandInt(3) {
	case 0:
		text = "Bang!"
	case 1:
		text = "Pop!"
	case 2:
		text = "Boom!"
	}
	return text
}

func (st logStyle) Rune() rune {
	var r rune
	switch st {
	case logCritic:
		r = 'r'
	case logNotable:
		r = 'g'
	case logDamage:
		r = 'o'
	case logSpecial:
		r = 'm'
	case logStatusEnd:
		r = 'v'
	case logError:
		r = 'e'
	case logConfirm:
		r = 'c'
	default:
		r = 'N'
	}
	return r
}

var logStyles = map[rune]gruid.Style{
	'r': gruid.Style{}.WithFg(ColorRed),
	't': gruid.Style{}.WithFg(ColorYellow),
	'g': gruid.Style{}.WithFg(ColorGreen),
	'o': gruid.Style{}.WithFg(ColorOrange),
	'm': gruid.Style{}.WithFg(ColorMagenta),
	'v': gruid.Style{}.WithFg(ColorViolet),
	'e': gruid.Style{}.WithFg(ColorRed),
	'c': gruid.Style{}.WithFg(ColorCyan),
}

// DrawLog draws 2 compacted lines of log.
func (md *model) DrawLog() ui.StyledText {
	g := md.g
	stt := ui.StyledText{}.WithMarkups(logStyles)
	tick := false
	for i := len(g.Log) - 1; i >= 0; i-- {
		e := g.Log[i]
		s := e.MText
		if stt.Text() != "" {
			if tick {
				s = s + "\n"
				tick = false
			} else {
				s = s + " "
			}
		}
		if e.Tick {
			tick = true
		}
		if stt.WithText(s+stt.Text()).Format(79).Size().Y > 2 {
			break
		}
		stt = stt.WithText(s + stt.Text()).Format(79)
	}
	return stt
}
