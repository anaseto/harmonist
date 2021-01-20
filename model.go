package main

import (
	"fmt"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/ui"
)

type mode int

const (
	modeNormal mode = iota
	modePager
	modeMenu
	modeQuit
	modeQuitConfirmation
)

type pagerMode int

const (
	modeLogs pagerMode = iota
	modeHelpKeys
)

type menuMode int

const (
	modeInventory menuMode = iota
	modeSettings
	modeGameMenu
	modeEvokation
	modeEquip
)

type model struct {
	g           *game // game state
	gd          gruid.Grid
	mode        mode
	menuMode    menuMode
	pagerMode   pagerMode
	menu        *ui.Menu
	help        *ui.Menu
	status      *ui.Menu
	log         *ui.Label
	description *ui.Label
	pager       *ui.Pager
	pagerMarkup ui.StyledText
	mp          mapUI
	logs        []ui.StyledText
	keysNormal  map[gruid.Key]action
	keysTarget  map[gruid.Key]action
	quit        bool
	finished    bool
}

type mapUI struct {
	kbTargeting bool
	ex          *examination
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
		Grid: gruid.NewGrid(UIWidth, UIHeight-1),
		Box:  &ui.Box{},
	})
	md.pagerMarkup = ui.StyledText{}.WithMarkups(logStyles)
	style := ui.MenuStyle{
		Active: gruid.Style{}.WithFg(ColorYellow),
	}
	md.menu = ui.NewMenu(ui.MenuConfig{
		Grid:  gruid.NewGrid(UIWidth/2, UIHeight-1),
		Box:   &ui.Box{},
		Style: style,
	})
	md.status = ui.NewMenu(ui.MenuConfig{
		Grid:  gruid.NewGrid(UIWidth, 1),
		Style: ui.MenuStyle{Layout: gruid.Point{0, 1}},
	})
}

func (md *model) init() gruid.Effect {
	SolarizedPalette()
	GameConfig.DarkLOS = true
	GameConfig.Version = Version
	GameConfig.Tiles = true
	LinkColors()
	//ApplyConfig()
	md.initKeys()
	md.initWidgets()

	g := md.g

	load, err := g.LoadConfig()
	var cfgerrstr string
	var cfgreseterr string
	if load && err != nil {
		cfgerrstr = fmt.Sprintf("Error loading config: %s", err.Error())
		err = g.SaveConfig()
		if err != nil {
			cfgreseterr = fmt.Sprintf("Error resetting config: %s", err.Error())
		}
	} else if load {
		CustomKeys = true
	}
	ApplyConfig()
	//ui.DrawWelcome()
	load, err = g.Load()
	if !load {
		g.InitLevel()
	} else if err != nil {
		g.InitLevel()
		g.PrintfStyled("Error: %v", logError, err)
		g.PrintStyled("Could not load saved state… starting new state.", logError)
	}
	if cfgerrstr != "" {
		g.PrintStyled(cfgerrstr, logError)
	}
	if cfgreseterr != "" {
		g.PrintStyled(cfgreseterr, logError)
	}

	//md.g.InitLevel()
	md.g.ComputeNoise()
	md.g.ComputeLOS()
	md.g.ComputeMonsterLOS()
	md.updateStatus()
	md.mp.ex = &examination{}
	return nil
}

func (md *model) Update(msg gruid.Msg) gruid.Effect {
	if _, ok := msg.(gruid.MsgInit); ok {
		return md.init()
	}
	if md.finished {
		switch msg.(type) {
		case gruid.MsgKeyDown:
			return gruid.End()
		default:
			return nil
		}
	}
	switch md.mode {
	case modeQuit:
		return nil
	case modeQuitConfirmation:
		eff := md.updateQuitConfirmation(msg)
		if md.mode == modeQuit {
			err := md.g.RemoveSaveFile()
			if err != nil {
				md.g.PrintfStyled("Error removing save file: %v", logError, err)
			}
		}
		return eff
	}
	if _, ok := msg.(gruid.MsgQuit); ok {
		md.mode = modeQuit
		md.g.Save() // TODO: log error ?
		return gruid.End()
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

func (md *model) updateQuitConfirmation(msg gruid.Msg) gruid.Effect {
	switch msg := msg.(type) {
	case gruid.MsgKeyDown:
		if msg.Key == "y" || msg.Key == "Y" {
			md.mode = modeQuit
			return gruid.End()
		} else {
			md.mode = modeNormal
		}
	}
	return nil
}

func (md *model) updateNormal(msg gruid.Msg) gruid.Effect {
	var eff gruid.Effect
	switch msg := msg.(type) {
	case gruid.MsgKeyDown:
		eff = md.updateKeyDown(msg)
	case gruid.MsgMouse:
		eff = md.updateMouse(msg)
	}
	return eff
}

func (md *model) updateKeyDown(msg gruid.MsgKeyDown) gruid.Effect {
	md.g.Ev = &simpleEvent{EAction: PlayerTurn, ERank: md.g.Turn}
	if !md.mp.kbTargeting && valid(md.mp.ex.pos) {
		md.CancelExamine()
	}
	again, eff, err := md.normalModeKeyDown(msg.Key)
	if again {
		return eff
	}
	if err != nil {
		md.g.Print(err.Error())
		return eff
	}
	md.EndTurn()
	return nil
}

func (md *model) updateMouse(msg gruid.MsgMouse) gruid.Effect {
	p := msg.P.Add(gruid.Point{0, -2}) // relative position ignoring log
	switch msg.Action {
	case gruid.MouseMove:
		if valid(p) {
			md.Examine(p)
		} else {
			md.CancelExamine()
		}
	}
	return nil
}

func (md *model) EndTurn() {
	md.mode = modeNormal
	md.g.EndTurn()
	if md.g.Player.HP <= 0 {
		md.Death()
		md.finished = true
		return
	}
	md.g.ComputeNoise()
	md.g.ComputeLOS()
	md.g.ComputeMonsterLOS()
	md.updateStatus()
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
	switch act := md.menu.Action(); act {
	case ui.MenuQuit:
		md.mode = modeNormal
	case ui.MenuMove, ui.MenuInvoke:
		switch md.menuMode {
		case modeInventory:
			items := []item{md.g.Player.Inventory.Body, md.g.Player.Inventory.Neck, md.g.Player.Inventory.Misc}
			it := items[md.menu.Active()]
			md.description.StyledText = ui.Text(it.Desc(md.g)).Format(UIWidth/2 - 1 - 2)
		case modeEvokation:
			items := md.g.Player.Magaras
			it := items[md.menu.Active()]
			md.description.StyledText = ui.Text(it.Desc(md.g)).Format(UIWidth/2 - 1 - 2)
			if act != ui.MenuInvoke {
				break
			}
			err := md.g.UseMagara(md.menu.Active())
			if err != nil {
				md.g.Printf("%v", err)
				md.mode = modeNormal
				break
			}
			md.EndTurn()
		case modeEquip:
			items := md.g.Player.Magaras
			it := items[md.menu.Active()]
			md.description.StyledText = ui.Text(it.Desc(md.g)).Format(UIWidth/2 - 1 - 2)
			if act != ui.MenuInvoke {
				break
			}
			err := md.g.EquipMagara(md.menu.Active())
			if err != nil {
				md.g.Printf("%v", err)
				md.mode = modeNormal
				break
			}
			md.EndTurn()
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
		if md.g.Highlight[p] || p == md.mp.ex.pos {
			attrs |= AttrReverse
		}
		dgd.Set(p, gruid.Cell{Rune: r, Style: gruid.Style{Fg: fg, Bg: bg, Attrs: attrs}})
	}
	md.log.StyledText = md.DrawLog()
	md.log.Draw(md.gd.Slice(md.gd.Range().Lines(0, 2)))
	if md.mp.ex.pos != InvalidPos {
		md.DrawPosInfo()
	}
	switch md.mode {
	case modePager:
		md.gd.Copy(md.pager.Draw())
	case modeMenu:
		switch md.menuMode {
		case modeInventory, modeEquip, modeEvokation:
			md.gd.Copy(md.menu.Draw())
			md.description.Box = &ui.Box{Title: ui.Text("Description")}
			md.description.Draw(md.gd.Slice(md.gd.Range().Columns(UIWidth/2+1, UIWidth)))
		}
	}
	md.gd.Slice(md.gd.Range().Line(UIHeight - 1)).Copy(md.status.Draw())
	return md.gd
}
