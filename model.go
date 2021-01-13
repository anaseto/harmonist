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

type menuMode int

const (
	modeInventory menuMode = iota
	modeSettings
	modeGameMenu
)

type model struct {
	g           *game // game state
	gd          gruid.Grid
	mode        mode
	menuMode    menuMode
	menu        *ui.Menu
	status      *ui.Menu
	log         *ui.Label
	description *ui.Label
	pager       *ui.Pager
	mp          mapUI
	logs        []string
	keysNormal  map[gruid.Key]action
	keysTarget  map[gruid.Key]action
}

type mapUI struct {
	targeting bool
	ex        *examination
}

func (md *model) initKeys() {
	md.keysNormal = map[gruid.Key]action{
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
	md.keysTarget = map[gruid.Key]action{
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

func (md *model) initWidgets() {
	md.log = ui.NewLabel(ui.StyledText{}.WithStyle(gruid.Style{}).WithMarkup('t', gruid.Style{Fg: ColorYellow}))
	md.description = ui.NewLabel(ui.StyledText{}.WithStyle(gruid.Style{}).WithMarkup('t', gruid.Style{Fg: ColorYellow}))
	md.description.AdjustWidth = false
	md.pager = ui.NewPager(ui.PagerConfig{
		Grid:       gruid.NewGrid(UIWidth, UIHeight),
		Box:        &ui.Box{},
		StyledText: ui.StyledText{}.WithMarkups(logStyles),
	})
	style := ui.MenuStyle{
		BgAlt:  ColorBgLOS,
		Active: ColorYellow,
	}
	md.menu = ui.NewMenu(ui.MenuConfig{
		Grid:  gruid.NewGrid(UIWidth/2, UIHeight-1),
		Box:   &ui.Box{},
		Style: style,
	})
	st := gruid.Style{}
	md.status = ui.NewMenu(ui.MenuConfig{
		Grid: gruid.NewGrid(UIWidth, 1),
		StyledText: ui.StyledText{}.WithMarkups(map[rune]gruid.Style{
			'G': st.WithFg(ColorFgHPok),
			'g': st.WithFg(ColorFgMPok),
			'W': st.WithFg(ColorFgHPwounded),
			'w': st.WithFg(ColorFgMPpartial),
			'C': st.WithFg(ColorFgHPcritical),
			'c': st.WithFg(ColorFgMPcritical),
			'B': st.WithFg(ColorCyan),
			'M': st.WithFg(ColorYellow).WithAttrs(AttrInMap),
		}),
		Style: ui.MenuStyle{Layout: gruid.Point{0, 1}},
	})
}

func (md *model) Update(msg gruid.Msg) gruid.Effect {
	if _, ok := msg.(gruid.MsgInit); ok {
		SolarizedPalette()
		GameConfig.DarkLOS = true
		GameConfig.Version = Version
		GameConfig.Tiles = true
		LinkColors()
		ApplyConfig()
		md.initKeys()
		md.initWidgets()
		md.g.InitLevel()
		md.g.ComputeNoise()
		md.g.ComputeLOS()
		md.g.ComputeMonsterLOS()
		md.updateStatus()
		return nil
	}
	var eff gruid.Effect
	switch md.mode {
	case modeNormal:
		eff = md.updateNormal(msg)
	case modePager:
		eff = md.updatePager(msg)
	case modeMenu:
		eff = md.updateMenu(msg)
	}
	return eff
}

func (md *model) updateNormal(msg gruid.Msg) gruid.Effect {
	var eff gruid.Effect
	switch msg := msg.(type) {
	case gruid.MsgKeyDown:
		eff = md.updateKeyDown(msg)
	}
	return eff
}

func (md *model) updateKeyDown(msg gruid.MsgKeyDown) gruid.Effect {
	switch msg.Key {
	case gruid.KeyEscape:
		return gruid.End()
	default:
		md.g.Ev = &simpleEvent{EAction: PlayerTurn, ERank: md.g.Turn}
		again, err := md.normalModeKeyDown(msg.Key)
		if again {
			break
		}
		if err != nil {
			md.g.Print(err.Error())
			break
		}
		md.g.EndTurn()
		md.g.ComputeNoise()
		md.g.ComputeLOS()
		md.g.ComputeMonsterLOS()
		md.updateStatus()
	}
	return nil
}

func (md *model) updatePager(msg gruid.Msg) gruid.Effect {
	md.pager.Update(msg)
	if md.pager.Action() == ui.PagerQuit {
		md.mode = modeNormal
	}
	return nil
}

func (md *model) updateMenu(msg gruid.Msg) gruid.Effect {
	md.menu.Update(msg)
	switch md.menu.Action() {
	case ui.MenuQuit:
		md.mode = modeNormal
	case ui.MenuMove, ui.MenuInvoke:
		switch md.menuMode {
		case modeInventory:
			items := []item{md.g.Player.Inventory.Body, md.g.Player.Inventory.Neck, md.g.Player.Inventory.Misc}
			it := items[md.menu.Active()]
			md.description.StyledText = ui.NewStyledText(it.Desc(md.g)).Format(UIWidth/2 - 1 - 2)
		}
	}
	return nil
}

func (md *model) Draw() gruid.Grid {
	md.gd.Fill(gruid.Cell{Rune: ' '})
	dgd := md.gd.Slice(md.gd.Range().Shift(0, 2, 0, -1))
	for i := range md.g.Dungeon.Cells {
		p := idxtopos(i)
		r, fg, bg := md.PositionDrawing(p)
		attrs := AttrInMap
		if md.g.Highlight[p] {
			attrs |= AttrReverse
		}
		dgd.Set(p, gruid.Cell{Rune: r, Style: gruid.Style{Fg: fg, Bg: bg, Attrs: attrs}})
	}
	md.log.StyledText = md.DrawLog()
	md.log.Draw(md.gd.Slice(md.gd.Range().Lines(0, 2)))
	if md.mp.targeting {
		md.DrawPosInfo()
	}
	switch md.mode {
	case modePager:
		md.gd.Copy(md.pager.Draw())
	case modeMenu:
		md.gd.Copy(md.menu.Draw())
		switch md.menuMode {
		case modeInventory:
			md.description.Box = &ui.Box{Title: ui.NewStyledText("Description")}
			md.description.Draw(md.gd.Slice(md.gd.Range().Columns(UIWidth/2+1, UIWidth)))
		}
	}
	md.gd.Slice(md.gd.Range().Line(UIHeight - 1)).Copy(md.status.Draw())
	return md.gd
}
