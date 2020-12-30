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
	st         *state // game state
	gd         gruid.Grid
	mode       mode
	cursor     gruid.Point
	menu       *ui.Menu
	status     *ui.Menu
	label      *ui.Label
	pager      *ui.Pager
	targeting  bool // whether currently targeting in modeNormal
	keysNormal map[gruid.Key]action
	keysTarget map[gruid.Key]action
}

func (m *model) initKeys() {
	m.keysNormal = map[gruid.Key]action{
		"h": ActionW,
		"j": ActionS,
		"k": ActionN,
		"l": ActionE,
		"a": ActionW,
		"s": ActionS,
		"w": ActionN,
		"d": ActionE,
		"4": ActionW,
		"2": ActionS,
		"8": ActionN,
		"6": ActionE,
		"H": ActionRunW,
		"J": ActionRunS,
		"K": ActionRunN,
		"L": ActionRunE,
		".": ActionWaitTurn,
		"5": ActionWaitTurn,
		"G": ActionGoToStairs,
		"o": ActionExplore,
		"x": ActionExamine,
		"v": ActionEvoke,
		"z": ActionEvoke,
		"e": ActionInteract,
		"i": ActionInventory,
		"m": ActionLogs,
		"M": ActionMenu,
		"#": ActionDump,
		"?": ActionHelp,
		"S": ActionSave,
		"Q": ActionQuit,
		"W": ActionWizard,
		"@": ActionWizardInfo,
		">": ActionWizardDescend,
		"=": ActionConfigure,
	}
	m.keysTarget = map[gruid.Key]action{
		"h":             ActionW,
		"j":             ActionS,
		"k":             ActionN,
		"l":             ActionE,
		"a":             ActionW,
		"s":             ActionS,
		"w":             ActionN,
		"d":             ActionE,
		"4":             ActionW,
		"2":             ActionS,
		"8":             ActionN,
		"6":             ActionE,
		"H":             ActionRunW,
		"J":             ActionRunS,
		"K":             ActionRunN,
		"L":             ActionRunE,
		">":             ActionNextStairs,
		"-":             ActionPreviousMonster,
		"+":             ActionNextMonster,
		"o":             ActionNextObject,
		"]":             ActionNextObject,
		")":             ActionNextObject,
		"(":             ActionNextObject,
		"[":             ActionNextObject,
		"_":             ActionNextObject,
		"=":             ActionNextObject,
		"v":             ActionDescription,
		".":             ActionTarget,
		"t":             ActionTarget,
		"g":             ActionTarget,
		"e":             ActionExclude,
		gruid.KeySpace:  ActionEscape,
		gruid.KeyEscape: ActionEscape,
		"x":             ActionEscape,
		"X":             ActionEscape,
		"?":             ActionHelp,
	}
}

func (m *model) initWidgets() {
	m.label = ui.NewLabel(ui.StyledText{}.WithStyle(gruid.Style{Fg: ColorFg, Bg: ColorBgLOS}).WithMarkup('t', gruid.Style{Fg: ColorYellow}))
}

func (m *model) Update(msg gruid.Msg) gruid.Effect {
	if _, ok := msg.(gruid.MsgInit); ok {
		SolarizedPalette()
		GameConfig.DarkLOS = true
		GameConfig.Version = Version
		LinkColors()
		ApplyConfig()
		m.initKeys()
		m.initWidgets()
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
	var eff gruid.Effect
	switch msg := msg.(type) {
	case gruid.MsgKeyDown:
		eff = m.updateKeyDown(msg)
	}
	return eff
}

func (m *model) updateKeyDown(msg gruid.MsgKeyDown) gruid.Effect {
	m.st.ComputeNoise()
	m.st.ComputeLOS()
	m.st.ComputeMonsterLOS()
	switch msg.Key {
	case gruid.KeyEscape:
		return gruid.End()
	}
	return nil
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
	//m.DrawDungeonView(NoFlushMode)
	dgd := m.gd.Slice(m.gd.Range().Shift(0, 2, 0, -1))
	for i := range m.st.Dungeon.Cells {
		p := idxtopos(i)
		r, fg, bg := m.PositionDrawing(p)
		dgd.Set(p, gruid.Cell{Rune: r, Style: gruid.Style{Fg: fg, Bg: bg}})
	}
	m.label.SetText(m.DrawLog())
	m.label.Draw(m.gd.Slice(m.gd.Range().Lines(0, 2)))

	return m.gd
}
