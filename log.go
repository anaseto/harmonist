package main

import "fmt"

type logStyle int

const (
	logNormal logStyle = iota
	logCritic
	logPlayerHit
	logMonsterHit
	logSpecial
	logStatusEnd
	logError
)

type logEntry struct {
	Text  string
	Index int
	Tick  bool
	Style logStyle
	Dups  int
}

func (e logEntry) String() string {
	if e.Dups > 0 {
		return fmt.Sprintf("%s (%d×)", e.Text, e.Dups+1)
	}
	return e.Text
}

func (g *state) Print(s string) {
	e := logEntry{Text: s, Index: g.LogIndex}
	g.PrintEntry(e)
}

func (g *state) PrintStyled(s string, style logStyle) {
	e := logEntry{Text: s, Index: g.LogIndex, Style: style}
	g.PrintEntry(e)
}

func (g *state) Printf(format string, a ...interface{}) {
	e := logEntry{Text: fmt.Sprintf(format, a...), Index: g.LogIndex}
	g.PrintEntry(e)
}

func (g *state) PrintfStyled(format string, style logStyle, a ...interface{}) {
	e := logEntry{Text: fmt.Sprintf(format, a...), Index: g.LogIndex, Style: style}
	g.PrintEntry(e)
}

func (g *state) PrintEntry(e logEntry) {
	if e.Index == g.LogNextTick {
		e.Tick = true
	}
	if !e.Tick && len(g.Log) > 0 {
		le := g.Log[len(g.Log)-1]
		if le.Text == e.Text {
			le.Dups++
			g.Log[len(g.Log)-1] = le
			return
		}
	}
	g.Log = append(g.Log, e)
	g.LogIndex++
	if len(g.Log) > 100000 {
		g.Log = g.Log[10000:]
	}
}

func (g *state) StoryPrint(s string) {
	g.Stats.Story = append(g.Stats.Story, fmt.Sprintf("Depth %2d|Turn %5d| %s", g.Depth, g.Turn, s))
}

func (g *state) StoryPrintf(format string, a ...interface{}) {
	g.Stats.Story = append(g.Stats.Story, fmt.Sprintf("Depth %2d|Turn %5d| %s", g.Depth, g.Turn, fmt.Sprintf(format, a...)))
}

func (g *state) CrackSound() (text string) {
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

func (g *state) ExplosionSound() (text string) {
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
