package main

import (
	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/ui"
)

type mode int

const (
	modeNormal mode = iota
	modePager
	modeMenu
)

type model struct {
	st     *state // game state
	gd     gruid.Grid
	mode   mode
	cursor gruid.Point
	menu   *ui.Menu
	label  *ui.Label
	pager  *ui.Pager
}

func (m *model) Update(msg gruid.Msg) gruid.Effect {
	if _, ok := msg.(gruid.MsgInit); ok {
		m.st.InitLevel()
		return nil
	}
	var eff gruid.Effect
	switch m.mode {
	case modeNormal:
		eff = m.updateNormal(msg)
	case modePager:
		eff = m.updatePager(msg)
	case modeMenu:
		eff = m.updateMenu(msg)
	}
	return eff
}

func (m *model) updateNormal(msg gruid.Msg) gruid.Effect {
	return nil
	// TODO
}

func (m *model) updatePager(msg gruid.Msg) gruid.Effect {
	return nil
	// TODO
}

func (m *model) updateMenu(msg gruid.Msg) gruid.Effect {
	return nil
	// TODO
}

func (m *model) Draw() gruid.Grid {
	m.DrawDungeonView(NoFlushMode)

	return m.gd
}
