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
	menu       *ui.Menu
	status     *ui.Menu
	label      *ui.Label
	pager      *ui.Pager
	mp         mapUI
	keysNormal map[gruid.Key]action
	keysTarget map[gruid.Key]action
}

type mapUI struct {
	targeting bool
	ex        *examination
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
	m.label = ui.NewLabel(ui.StyledText{}.WithStyle(gruid.Style{}).WithMarkup('t', gruid.Style{Fg: ColorYellow}))
	m.pager = ui.NewPager(ui.PagerConfig{
		Grid: gruid.NewGrid(UIWidth, UIHeight),
		Box:  &ui.Box{},
	})
	m.menu = ui.NewMenu(ui.MenuConfig{
		Grid: gruid.NewGrid(UIWidth/2, UIHeight-1),
		Box:  &ui.Box{},
	})
	m.status = ui.NewMenu(ui.MenuConfig{
		Grid: m.gd.Slice(gruid.NewRange(0, DungeonHeight, UIWidth, UIHeight)),
		Box:  &ui.Box{},
	})
}

func (m *model) Update(msg gruid.Msg) gruid.Effect {
	if _, ok := msg.(gruid.MsgInit); ok {
		SolarizedPalette()
		GameConfig.DarkLOS = true
		GameConfig.Version = Version
		GameConfig.Tiles = true
		LinkColors()
		ApplyConfig()
		m.initKeys()
		m.initWidgets()
		m.st.InitLevel()
		m.st.ComputeNoise()
		m.st.ComputeLOS()
		m.st.ComputeMonsterLOS()
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
	switch msg.Key {
	case gruid.KeyEscape:
		return gruid.End()
	default:
		m.st.Ev = &simpleEvent{EAction: PlayerTurn, ERank: m.st.Turn}
		again, err := m.normalModeKeyDown(msg.Key)
		if again {
			break
		}
		if err != nil {
			m.st.Print(err.Error())
			break
		}
		m.st.EndTurn()
		m.st.ComputeNoise()
		m.st.ComputeLOS()
		m.st.ComputeMonsterLOS()
	}
	return nil
}

func (m *model) updatePager(msg gruid.Msg) gruid.Effect {
	m.pager.Update(msg)
	if m.pager.Action() == ui.PagerQuit {
		m.mode = modeNormal
	}
	return nil
}

func (m *model) updateMenu(msg gruid.Msg) gruid.Effect {
	m.menu.Update(msg)
	if m.menu.Action() == ui.MenuQuit {
		m.mode = modeNormal
	}
	return nil
}

func (m *model) Draw() gruid.Grid {
	dgd := m.gd.Slice(m.gd.Range().Shift(0, 2, 0, -1))
	for i := range m.st.Dungeon.Cells {
		p := idxtopos(i)
		r, fg, bg := m.PositionDrawing(p)
		attrs := AttrInMap
		if m.st.Highlight[p] {
			attrs |= AttrReverse
		}
		dgd.Set(p, gruid.Cell{Rune: r, Style: gruid.Style{Fg: fg, Bg: bg, Attrs: attrs}})
	}
	m.label.AdjustWidth = false
	m.label.Box = nil
	m.label.SetText(m.DrawLog())
	m.label.Draw(m.gd.Slice(m.gd.Range().Lines(0, 2)))
	if m.mp.targeting {
		m.DrawPosInfo()
	}
	switch m.mode {
	case modePager:
		m.gd.Copy(m.pager.Draw())
	case modeMenu:
		m.gd.Copy(m.menu.Draw())
	}
	return m.gd
}
