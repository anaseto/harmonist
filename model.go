package main

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/ui"
)

const (
	UIWidth  = 80
	UIHeight = 24
)

var (
	DisableAnimations bool = false
	Xterm256Color          = false
	Terminal               = false
)

// CustomKeys tracks whether we're using custom key bindings.
var CustomKeys bool
var GameConfig config

type mode int

const (
	modeNormal mode = iota
	modePager
	modeSmallPager
	modeMenu
	modeQuit
	modeQuitConfirmation
	modeJumpConfirmation
	modeWizardConfirmation
	modeDump // simplified dump visualization after end
	modeEnd  // win or death
	modeHPCritical
	modeWelcome
	modeStory
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
	modeKeys
	modeKeysChange
	modeGameMenu
	modeEvokation
	modeEquip
	modeWizard
)

type model struct {
	g           *game // game state
	gd          gruid.Grid
	mode        mode
	menuMode    menuMode
	pagerMode   pagerMode
	menu        *ui.Menu
	keysMenu    *ui.Menu
	status      *ui.Menu
	log         *ui.Label
	description *ui.Label
	statusDesc  *ui.Label
	pager       *ui.Pager
	smallPager  *ui.Pager
	pagerMarkup ui.StyledText
	mp          mapUI
	logs        []ui.StyledText
	keysNormal  map[gruid.Key]action
	keysTarget  map[gruid.Key]action
	finished    bool
	statusFocus bool
	anims       Animations
	story       int
}

type mapUI struct {
	kbTargeting bool
	ex          *examination
}

func (md *model) initKeys() {
	md.keysNormal = map[gruid.Key]action{
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
		".":             ActionWaitTurn,
		"5":             ActionWaitTurn,
		"G":             ActionGoToStairs,
		"o":             ActionExplore,
		"x":             ActionExamine,
		"v":             ActionEvoke,
		"z":             ActionEvoke,
		"e":             ActionInteract,
		"E":             ActionInteract,
		"i":             ActionInventory,
		"I":             ActionInventory,
		"m":             ActionLogs,
		"M":             ActionMenu,
		"#":             ActionDump,
		"?":             ActionHelp,
		"S":             ActionSave,
		"Q":             ActionQuit,
		"W":             ActionWizard,
		"@":             ActionWizardMenu,
		">":             ActionWizardDescend,
		"=":             ActionSettings,
		gruid.KeyEscape: ActionEscape,
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
	CustomKeys = false
}

func (md *model) initWidgets() {
	md.log = ui.NewLabel(ui.StyledText{}.WithStyle(gruid.Style{}).WithMarkup('t', gruid.Style{Fg: ColorYellow}))
	md.description = ui.NewLabel(ui.StyledText{}.WithStyle(gruid.Style{}).WithMarkup('t', gruid.Style{Fg: ColorYellow}))
	md.description.AdjustWidth = false
	md.statusDesc = ui.NewLabel(ui.StyledText{}.WithStyle(gruid.Style{}))
	md.pager = ui.NewPager(ui.PagerConfig{
		Grid: gruid.NewGrid(UIWidth, UIHeight-1),
		Box:  &ui.Box{},
		Keys: ui.PagerKeys{Quit: []gruid.Key{gruid.KeySpace, "x", "X", gruid.KeyEscape}},
	})
	md.smallPager = ui.NewPager(ui.PagerConfig{
		Grid: gruid.NewGrid(60, UIHeight-1),
		Box:  &ui.Box{},
		Keys: ui.PagerKeys{Quit: []gruid.Key{gruid.KeySpace, "x", "X", gruid.KeyEscape}},
	})
	md.pagerMarkup = ui.StyledText{}.WithMarkups(logStyles)
	style := ui.MenuStyle{
		Active: gruid.Style{}.WithFg(ColorYellow),
	}
	md.menu = ui.NewMenu(ui.MenuConfig{
		Grid:  gruid.NewGrid(UIWidth/2, UIHeight-1),
		Box:   &ui.Box{},
		Style: style,
		Keys:  ui.MenuKeys{Quit: []gruid.Key{gruid.KeySpace, "x", "X", gruid.KeyEscape}},
	})
	md.keysMenu = ui.NewMenu(ui.MenuConfig{
		Grid:  gruid.NewGrid(UIWidth, UIHeight-1),
		Box:   &ui.Box{},
		Style: style,
		Keys:  ui.MenuKeys{Quit: []gruid.Key{gruid.KeySpace, "x", "X", gruid.KeyEscape}},
	})
	md.status = ui.NewMenu(ui.MenuConfig{
		Grid:  gruid.NewGrid(UIWidth, 1),
		Style: ui.MenuStyle{Layout: gruid.Point{0, 1}, Active: style.Active},
	})
}

func initConfig() (string, string) {
	GameConfig.DarkLOS = true
	GameConfig.Version = Version
	GameConfig.Tiles = true
	load, err := LoadConfig()
	var cfgerrstr string
	var cfgreseterr string
	if load && err != nil {
		cfgerrstr = fmt.Sprintf("Error loading config: %s", err.Error())
		err = SaveConfig()
		if err != nil {
			cfgreseterr = fmt.Sprintf("Error resetting config: %s", err.Error())
		}
	} else if load {
		CustomKeys = true
	}
	return cfgerrstr, cfgreseterr
}

func (md *model) init() gruid.Effect {
	if runtime.GOOS != "js" {
		md.mode = modeWelcome
	}
	LinkColors()
	md.initKeys()
	md.initWidgets()

	g := md.g

	cfgerrstr, cfgreseterr := initConfig()
	md.applyConfig()
	load, err := g.Load()
	md.g.md = md // TODO: avoid this? (though it's handy)
	if !load {
		g.InitLevel()
	} else if err != nil {
		g.InitLevel()
		g.PrintfStyled("Error: %v", logError, err)
		g.PrintStyled("Could not load saved state… starting new state.", logError)
	} else {
		g.checks()
	}
	if cfgerrstr != "" {
		g.PrintStyled(cfgerrstr, logError)
	}
	if cfgreseterr != "" {
		g.PrintStyled(cfgreseterr, logError)
	}

	md.g.ComputeNoise()
	md.g.ComputeLOS()
	md.g.ComputeMonsterLOS()
	md.updateStatusInfo()
	md.mp.ex = &examination{}
	md.CancelExamine()
	md.initAnimations()
	return nil
}

func (md *model) more(msg gruid.Msg) bool {
	switch msg := msg.(type) {
	case gruid.MsgKeyDown:
		switch msg.Key {
		case "x", "X", gruid.KeyEscape, gruid.KeySpace:
			return true
		}
	}
	return false
}

func (md *model) interrupt(msg gruid.Msg) bool {
	switch msg := msg.(type) {
	case gruid.MsgKeyDown:
		return true
	case gruid.MsgMouse:
		return msg.Action != gruid.MouseMove
	}
	return false
}

func (md *model) Update(msg gruid.Msg) gruid.Effect {
	if _, ok := msg.(gruid.MsgInit); ok {
		return md.init()
	}
	if _, ok := msg.(gruid.MsgQuit); ok {
		md.mode = modeQuit
		md.g.Save() // TODO: log error ?
		return gruid.End()
	}
	//log.Printf("msg %+v", msg)
	md.anims.draw = false
	if msg, ok := msg.(msgAnim); ok {
		if int(msg) != md.anims.idx {
			return nil
		}
		if !md.anims.Done() {
			md.anims.draw = true
			return md.animNext()
		}
		md.anims.Finish()
		return nil
	}
	anims := !md.anims.Done()
	if anims && md.interrupt(msg) {
		md.anims.Finish()
	}
	eff := md.update(msg)
	cmd := md.animCmd()
	if !anims {
		if cmd != nil {
			return gruid.Batch(eff, cmd)
		}
	}
	return eff
}

func (md *model) update(msg gruid.Msg) gruid.Effect {
	var eff gruid.Effect
	switch md.mode {
	case modeWelcome:
		switch msg := msg.(type) {
		case gruid.MsgKeyDown:
			md.mode = modeNormal
		case gruid.MsgMouse:
			if msg.Action != gruid.MouseMove {
				md.mode = modeNormal
			}
		}
		return nil
	case modeQuit:
		return nil
	case modeEnd:
		if md.more(msg) {
			md.finished = true
			md.mode = modeDump
			md.dump(md.g.WriteDump())
		}
		return nil
	case modeDump:
		return md.updateDump(msg)
	case modeQuitConfirmation:
		eff := md.updateQuitConfirmation(msg)
		if md.mode == modeQuit {
			err := RemoveSaveFile()
			if err != nil {
				md.g.PrintfStyled("Error removing save file: %v", logError, err)
			}
			RemoveDataFile("replay.part")
			return eff
		}
		return eff
	case modeJumpConfirmation:
		md.updateJumpConfirmation(msg)
		return nil
	case modeWizardConfirmation:
		md.updateWizardConfirmation(msg)
		return nil
	case modeHPCritical:
		if md.more(msg) {
			md.mode = modeNormal
			md.g.Print("Ok. Be careful, then.")
		}
		return nil
	case modeNormal:
		eff = md.updateNormal(msg)
	case modeStory:
		if md.more(msg) {
			md.Story()
		}
	case modePager:
		eff = md.updatePager(msg)
	case modeSmallPager:
		eff = md.updateSmallPager(msg)
	case modeMenu:
		switch md.menuMode {
		case modeKeys:
			eff = md.updateKeysMenu(msg)
		case modeKeysChange:
			eff = md.updateKeysChange(msg)
		default:
			eff = md.updateMenu(msg)
		}
	}
	return eff
}

func (md *model) updateQuitConfirmation(msg gruid.Msg) gruid.Effect {
	switch msg := msg.(type) {
	case gruid.MsgKeyDown:
		if msg.Key == "y" || msg.Key == "Y" {
			md.mode = modeQuit
			return gruid.End()
		}
		md.mode = modeNormal
	}
	return nil
}

func (md *model) updateJumpConfirmation(msg gruid.Msg) {
	switch msg := msg.(type) {
	case gruid.MsgKeyDown:
		md.mode = modeNormal
		if msg.Key == "y" || msg.Key == "Y" {
			md.g.FallAbyss(DescendFall)
		} else {
			md.g.Print("No jump, then.")
		}
	}
}

func (md *model) updateWizardConfirmation(msg gruid.Msg) {
	switch msg := msg.(type) {
	case gruid.MsgKeyDown:
		md.mode = modeNormal
		if msg.Key == "y" || msg.Key == "Y" {
			md.g.EnterWizardMode()
		} else {
			md.g.Print("Continuing normally, then.")
		}
	}
}

func (md *model) updateNormal(msg gruid.Msg) gruid.Effect {
	var eff gruid.Effect
	switch msg := msg.(type) {
	case msgAuto:
		if int(msg) == md.g.Turn && md.g.AutoNext {
			return md.EndTurn()
		}
	case gruid.MsgKeyDown:
		eff = md.updateKeyDown(msg)
	case gruid.MsgMouse:
		eff = md.updateMouse(msg)
	}
	return eff
}

func (md *model) updateKeyDown(msg gruid.MsgKeyDown) gruid.Effect {
	md.statusFocus = false
	if !md.mp.kbTargeting && valid(md.mp.ex.p) {
		md.CancelExamine()
	}
	if md.mp.ex.p != InvalidPos {
		switch msg.Key {
		case gruid.KeyPageDown:
			md.mp.ex.scroll = true
			return nil
		case gruid.KeyPageUp:
			md.mp.ex.scroll = false
			return nil
		}
	}
	again, eff, err := md.normalModeKeyDown(msg.Key)
	if err != nil {
		md.g.Print(err.Error())
	}
	if again {
		return eff
	}
	return md.EndTurn()
}

func (md *model) updateMouse(msg gruid.MsgMouse) gruid.Effect {
	if msg.P.Y == UIHeight-1 {
		return md.updateStatusMouse(msg)
	}
	md.statusFocus = false
	p := msg.P.Add(gruid.Point{0, -2}) // relative position ignoring log
	switch msg.Action {
	case gruid.MouseWheelUp:
		if md.mp.ex.p != InvalidPos {
			md.mp.ex.scroll = true
		}
	case gruid.MouseWheelDown:
		if md.mp.ex.p != InvalidPos {
			md.mp.ex.scroll = false
		}
	case gruid.MouseMove:
		if valid(p) {
			md.Examine(p)
		} else {
			md.CancelExamine()
		}
	case gruid.MouseMain:
		if !valid(p) {
			return nil
		}
		again, eff, err := md.normalModeAction(ActionTarget)
		if err != nil {
			md.g.Print(err.Error())
		}
		if again {
			return eff
		}
		return md.EndTurn()
	}
	return nil
}

type statusItem int

const (
	statusDepth statusItem = iota
	statusTurns
	statusHP
	statusMP
	statusBananas
	statusMenu
	statusInventory
	statusInteract
)

func (md *model) updateStatusMouse(msg gruid.MsgMouse) gruid.Effect {
	msg.P.Y = 0
	if !msg.P.In(md.status.Bounds()) {
		md.statusFocus = false
		return nil
	}
	md.CancelExamine()
	md.status.Update(msg)
	update := !md.statusFocus
	switch md.status.Action() {
	case ui.MenuMove:
		update = true
	case ui.MenuInvoke:
		i := statusItem(md.status.Active())
		var action action
		switch i {
		case statusMenu:
			action = ActionMenu
		case statusInventory:
			action = ActionInventory
		case statusInteract:
			action = ActionInteract
			md.statusFocus = false
		}
		again, eff, err := md.normalModeAction(action)
		if err != nil {
			md.g.Printf("%v", err)
		}
		if again {
			return eff
		}
		return md.EndTurn()
	}
	if update {
		const statusIndex = statusInteract + 2
		i := statusItem(md.status.Active())
		md.statusFocus = false
		switch {
		case i == statusDepth:
			md.statusDesc.Box = &ui.Box{Title: ui.Text("Depth")}
			md.statusDesc.SetText("Dungeon depth.")
			md.statusFocus = true
		case i == statusTurns:
			md.statusDesc.Box = &ui.Box{Title: ui.Text("Turns")}
			md.statusDesc.SetText("Number of turns since the beginning.")
			md.statusFocus = true
		case i == statusHP:
			md.statusDesc.Box = &ui.Box{Title: ui.Text("Health")}
			md.statusDesc.SetText("Your hit points.")
			md.statusFocus = true
		case i == statusMP:
			md.statusDesc.Box = &ui.Box{Title: ui.Text("Magic Points")}
			md.statusDesc.SetText("Your magic points. Needed for evoking magaras.")
			md.statusFocus = true
		case i == statusBananas:
			md.statusDesc.Box = &ui.Box{Title: ui.Text("Bananas")}
			md.statusDesc.SetText("Need to eat one before sleeping in barrels.")
			md.statusFocus = true
		case i == statusMenu:
			md.statusDesc.Box = &ui.Box{Title: ui.Text("Menu")}
			md.statusDesc.SetText("Click to open menu.")
			md.statusFocus = true
		case i == statusInventory:
			md.statusDesc.Box = &ui.Box{Title: ui.Text("Inventory")}
			md.statusDesc.SetText("Click to open inventory.")
			md.statusFocus = true
		case i == statusInteract:
			s, ok := md.interact()
			if !ok {
				break
			}
			md.statusDesc.Box = &ui.Box{Title: ui.Text("Interact")}
			md.statusDesc.SetText(fmt.Sprintf("Click to %v.", s))
			md.statusFocus = true
		case i >= statusIndex:
			i := md.status.Active() - int(statusIndex)
			sts := md.sortedStatuses()
			if i > len(sts)-1 {
				break
			}
			md.statusDesc.Box = &ui.Box{Title: ui.Text(sts[i].String())}
			md.statusDesc.SetText(sts[i].Desc())
			md.statusFocus = true
		}
	}
	return nil
}

// EndTurn finalizes player's turn and runs other events until next player
// turn.
func (md *model) EndTurn() gruid.Effect {
	md.mode = modeNormal
	md.g.PushPlayerTurn()
	eff := md.g.EndTurn()
	if md.g.Player.HP <= 0 {
		md.death()
		return eff
	}
	md.updateMapInfo()
	return eff
}

func (md *model) updateMapInfo() {
	md.g.ComputeNoise()
	md.g.ComputeLOS()
	md.g.ComputeMonsterLOS()
	md.updateStatusInfo()
	if md.g.Highlight != nil {
		md.examine(md.mp.ex.p)
	}
}

func (md *model) updatePager(msg gruid.Msg) gruid.Effect {
	md.pager.Update(msg)
	if md.pager.Action() == ui.PagerQuit {
		md.mode = modeNormal
	}
	return nil
}

func (md *model) updateSmallPager(msg gruid.Msg) gruid.Effect {
	md.smallPager.Update(msg)
	if md.smallPager.Action() == ui.PagerQuit {
		md.mode = modeNormal
	}
	return nil
}

func (md *model) updateDump(msg gruid.Msg) gruid.Effect {
	md.pager.Update(msg)
	if md.pager.Action() == ui.PagerQuit {
		md.mode = modeQuit
		return gruid.End()
	}
	return nil
}

func (md *model) updateKeysChange(msg gruid.Msg) gruid.Effect {
	mkd, ok := msg.(gruid.MsgKeyDown)
	if !ok {
		return nil
	}
	key := mkd.Key
	switch {
	case key == gruid.KeyEscape:
		md.openKeyBindings()
		return nil
	case key.IsRune():
		action := ConfigurableKeyActions[md.keysMenu.Active()]
		//log.Printf("active %v, action %v", md.keysMenu.Active(), action)
		if action.normalModeAction() {
			md.keysNormal[key] = action
		}
		if action.targetingModeAction() {
			md.keysTarget[key] = action
		}
		GameConfig.NormalModeKeys = md.keysNormal
		GameConfig.TargetModeKeys = md.keysTarget
		err := SaveConfig()
		if err != nil {
			md.g.PrintStyled("Error while saving config changes.", logCritic)
		}
		md.openKeyBindings()
		return nil
	}
	return nil
}

func (md *model) updateKeysMenu(msg gruid.Msg) gruid.Effect {
	md.keysMenu.Update(msg)
	switch act := md.keysMenu.Action(); act {
	case ui.MenuQuit:
		md.mode = modeNormal
	case ui.MenuInvoke:
		md.menuMode = modeKeysChange
	case ui.MenuPass:
		msg, ok := msg.(gruid.MsgKeyDown)
		if !ok {
			return nil
		}
		if msg.Key == "R" {
			md.initKeys()
			GameConfig.NormalModeKeys = md.keysNormal
			GameConfig.TargetModeKeys = md.keysTarget
			err := SaveConfig()
			if err != nil {
				md.g.PrintStyled("Error while resetting config changes.", logCritic)
			}
			md.openKeyBindings()
		}
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
			md.description.Content = ui.Text(it.Desc(md.g)).Format(UIWidth/2 - 1 - 2)
		case modeEvokation:
			items := md.g.Player.Magaras
			it := items[md.menu.Active()]
			md.description.Content = ui.Text(it.Desc(md.g)).Format(UIWidth/2 - 1 - 2)
			if act != ui.MenuInvoke {
				break
			}
			err := md.g.UseMagara(md.menu.Active())
			if err != nil {
				md.g.Printf("%v", err)
				md.mode = modeNormal
				break
			}
			return md.EndTurn()
		case modeEquip:
			items := md.g.Player.Magaras
			it := items[md.menu.Active()]
			md.description.Content = ui.Text(it.Desc(md.g)).Format(UIWidth/2 - 1 - 2)
			if act != ui.MenuInvoke {
				break
			}
			err := md.g.EquipMagara(md.menu.Active())
			if err != nil {
				md.g.Printf("%v", err)
				md.mode = modeNormal
				break
			}
			return md.EndTurn()
		case modeGameMenu:
			if act != ui.MenuInvoke {
				break
			}
			_, eff, err := md.normalModeAction(menuActions[md.menu.Active()])
			if err != nil {
				// should not happen
				md.g.Printf("%v", err)
			}
			return eff
		case modeSettings:
			if act != ui.MenuInvoke {
				break
			}
			md.g.Printf("%v", settingsActions[md.menu.Active()])
			_, eff, err := md.normalModeAction(settingsActions[md.menu.Active()])
			if err != nil {
				// should not happen
				md.g.Printf("%v", err)
			}
			return eff
		case modeWizard:
			if act != ui.MenuInvoke {
				break
			}
			md.g.Printf("%v", wizardActions[md.menu.Active()])
			_, eff, err := md.normalModeAction(wizardActions[md.menu.Active()])
			if err != nil {
				// should not happen
				md.g.Printf("%v", err)
			}
			return eff
		}
	}
	return nil
}

func (md *model) normalModeKeyDown(key gruid.Key) (again bool, eff gruid.Effect, err error) {
	action := md.keysNormal[key]
	if md.mp.kbTargeting {
		action = md.keysTarget[key]
	}
	again, eff, err = md.normalModeAction(action)
	if _, ok := err.(actionError); ok {
		err = fmt.Errorf("Unknown key '%s'. Type ? for help.", key)
	}
	return again, eff, err
}

func (md *model) death() {
	g := md.g
	g.LevelStats()
	if len(g.Stats.Achievements) == 0 {
		NoAchievement.Get(g)
	}
	g.Print("You die... [(x) to continue]")
	md.mode = modeEnd
}

func (md *model) win() {
	g := md.g
	if g.Wizard {
		g.Print("You escape by the magic portal! **WIZARD** [(x) to continue]")
	} else {
		g.Print("You escape by the magic portal! [(x) to continue]")
	}
	md.mode = modeEnd
}

func (md *model) dump(err error) {
	s := md.g.SimplifedDump(err)
	lines := strings.Split(s, "\n")
	stts := []ui.StyledText{}
	for _, l := range lines {
		stts = append(stts, ui.Text(l))
	}
	md.pager.SetLines(stts)
	//log.Printf("%v", s)
	//log.Printf("%v", stts)
}

func (md *model) applyConfig() {
	if GameConfig.NormalModeKeys != nil {
		md.keysNormal = GameConfig.NormalModeKeys
	}
	if GameConfig.TargetModeKeys != nil {
		md.keysTarget = GameConfig.TargetModeKeys
	}
	if GameConfig.DarkLOS {
		ApplyDarkLOS()
	} else {
		ApplyLightLOS()
	}
}
