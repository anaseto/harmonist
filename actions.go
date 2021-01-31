package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	//"time"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/ui"
)

type action int

const (
	ActionNone action = iota
	ActionW
	ActionS
	ActionN
	ActionE
	ActionRunW
	ActionRunS
	ActionRunN
	ActionRunE
	ActionWaitTurn
	ActionDescend
	ActionGoToStairs
	ActionExplore
	ActionExamine
	ActionEvoke
	ActionInteract
	ActionInventory
	ActionLogs
	ActionDump
	ActionHelp
	ActionSave
	ActionQuit
	ActionWizard
	ActionWizardInfo
	ActionWizardDescend

	ActionPreviousMonster
	ActionNextMonster
	ActionNextObject
	ActionTarget
	ActionExclude
	ActionEscape

	ActionSettings
	ActionMenu
	ActionNextStairs
	ActionMenuCommandHelp
	ActionMenuTargetingHelp

	ActionSetKeys
	ActionInvertLOS
	ActionToggleTiles
	ActionToggleShowNumbers
)

var ConfigurableKeyActions = [...]action{
	ActionW,
	ActionS,
	ActionN,
	ActionE,
	ActionRunW,
	ActionRunS,
	ActionRunN,
	ActionRunE,
	ActionWaitTurn,
	ActionEvoke,
	ActionInteract,
	ActionInventory,
	ActionExamine,
	ActionGoToStairs,
	ActionExplore,
	ActionLogs,
	ActionDump,
	ActionSave,
	ActionQuit,
	ActionMenu,
	ActionPreviousMonster,
	ActionNextMonster,
	ActionNextObject,
	ActionNextStairs,
	ActionTarget,
	ActionExclude}

func (k action) normalModeAction() bool {
	switch k {
	case ActionW, ActionS, ActionN, ActionE,
		ActionRunW, ActionRunS, ActionRunN, ActionRunE,
		ActionWaitTurn,
		ActionDescend,
		ActionGoToStairs,
		ActionExplore,
		ActionExamine,
		ActionEvoke,
		ActionInteract,
		ActionInventory,
		ActionLogs,
		ActionDump,
		ActionHelp,
		ActionMenu,
		ActionMenuCommandHelp,
		ActionMenuTargetingHelp,
		ActionSave,
		ActionQuit,
		ActionSettings,
		ActionWizard,
		ActionWizardInfo:
		return true
	default:
		return false
	}
}

func (k action) String() (text string) {
	switch k {
	case ActionW:
		text = "Move west"
	case ActionS:
		text = "Move south"
	case ActionN:
		text = "Move north"
	case ActionE:
		text = "Move east"
	case ActionRunW:
		text = "Travel west"
	case ActionRunS:
		text = "Travel south"
	case ActionRunN:
		text = "Travel north"
	case ActionRunE:
		text = "Travel east"
	case ActionWaitTurn:
		text = "Wait a turn"
	case ActionDescend:
		text = "Descend stairs"
	case ActionGoToStairs:
		text = "Go to nearest stairs"
	case ActionExplore:
		text = "Autoexplore"
	case ActionExamine:
		text = "Examine"
	case ActionEvoke:
		text = "Evoke card"
	case ActionInteract:
		text = "Interact"
	case ActionInventory:
		text = "Inventory"
	case ActionLogs:
		text = "View previous messages"
	case ActionDump:
		text = "Write state statistics to file"
	case ActionSave:
		text = "Save and Quit"
	case ActionQuit:
		text = "Quit without saving"
	case ActionHelp:
		text = "Help (keys and mouse)"
	case ActionMenuCommandHelp:
		text = "Help (general commands)"
	case ActionMenuTargetingHelp:
		text = "Help (targeting commands)"
	case ActionSettings:
		text = "Settings and key bindings"
	case ActionWizard:
		text = "Wizard (debug) mode"
	case ActionWizardInfo:
		text = "Wizard (debug) mode information"
	case ActionMenu:
		text = "Action Menu"
	case ActionSetKeys:
		text = "Change key bindings"
	case ActionInvertLOS:
		text = "Toggle dark/light LOS"
	case ActionToggleTiles:
		text = "Toggle tiles/ascii display"
	case ActionToggleShowNumbers:
		text = "Toggle hearts/numbers"
	}
	return text
}

func (k action) targetingModeDescription() (text string) {
	switch k {
	case ActionW:
		text = "Move cursor west"
	case ActionS:
		text = "Move cursor south"
	case ActionN:
		text = "Move cursor north"
	case ActionE:
		text = "Move cursor east"
	case ActionRunW:
		text = "Big move cursor west"
	case ActionRunS:
		text = "Big move cursor south"
	case ActionRunN:
		text = "Big move north"
	case ActionRunE:
		text = "Big move east"
	case ActionDescend:
		text = "Target next stair"
	case ActionPreviousMonster:
		text = "Target previous monster"
	case ActionNextMonster:
		text = "Target next monster"
	case ActionNextObject:
		text = "Target next object"
	case ActionNextStairs:
		text = "Target next stairs"
	case ActionTarget:
		text = "Go to"
	case ActionExclude:
		text = "Toggle exclude area from auto-travel"
	case ActionEscape:
		text = "Quit targeting mode"
	}
	return text
}

func (k action) targetingModeAction() bool {
	switch k {
	case ActionW, ActionS, ActionN, ActionE,
		ActionRunW, ActionRunS, ActionRunN, ActionRunE,
		ActionPreviousMonster,
		ActionNextMonster,
		ActionNextObject,
		ActionNextStairs,
		ActionTarget,
		ActionExclude,
		ActionEscape:
		return true
	default:
		return false
	}
}

func (md *model) normalModeAction(action action) (again bool, eff gruid.Effect, err error) {
	g := md.g
	switch action {
	case ActionNone:
		// not used
		again = true
	case ActionW, ActionS, ActionN, ActionE:
		if !md.mp.kbTargeting {
			err = g.PlayerBump(To(KeyToDir(action), g.Player.Pos))
		} else {
			p := To(KeyToDir(action), md.mp.ex.pos)
			if valid(p) {
				md.Examine(p)
			}
			again = true
		}
	case ActionRunW, ActionRunS, ActionRunN, ActionRunE:
		if !md.mp.kbTargeting {
			err = g.GoToDir(KeyToDir(action))
		} else {
			q := InvalidPos
			for i := 0; i < 5; i++ {
				p := To(KeyToDir(action), md.mp.ex.pos)
				if !valid(p) {
					break
				}
				q = p
			}
			if q != InvalidPos {
				md.Examine(q)
			}
			again = true
		}
	case ActionExclude:
		again = true
		md.excludeZone(md.mp.ex.pos)
	case ActionPreviousMonster:
		again = true
		md.nextMonster("-", md.mp.ex.pos, md.mp.ex)
	case ActionNextMonster:
		again = true
		md.nextMonster("+", md.mp.ex.pos, md.mp.ex)
	case ActionNextObject:
		again = true
		md.nextObject(md.mp.ex.pos, md.mp.ex)
	case ActionTarget:
		again = true
		err = md.target()
		if err != nil {
			break
		}
		if g.MoveToTarget() {
			again = false
		}
	case ActionEscape:
		again = true
		md.CancelExamine()
	case ActionWaitTurn:
		g.WaitTurn()
	case ActionGoToStairs:
		again = true
		stairs := g.StairsSlice()
		sortedStairs := g.SortedNearestTo(stairs, g.Player.Pos)
		if len(sortedStairs) > 0 {
			stair := sortedStairs[0]
			if g.Player.Pos == stair {
				err = errors.New("You are already on the stairs.")
				break
			}
			md.mp.ex.pos = stair
			err = md.target()
			if err != nil {
				err = errors.New("There is no safe path to the nearest stairs.")
			} else if !g.MoveToTarget() {
				err = errors.New("You could not move toward stairs.")
			}
		} else {
			err = errors.New("You cannot go to any stairs.")
		}
	case ActionInteract:
		c := g.Dungeon.Cell(g.Player.Pos)
		switch terrain(c) {
		case StairCell:
			if terrain(g.Dungeon.Cell(g.Player.Pos)) == StairCell && g.Objects.Stairs[g.Player.Pos] != BlockedStair {
				// TODO: animation
				//ui.MenuSelectedAnimation(MenuInteract, true)
				strt := g.Objects.Stairs[g.Player.Pos]
				err = md.checkShaedra(strt)
				if err != nil {
					break
				}
				if g.Descend(DescendNormal) {
					md.win()
					return again, eff, err
				}
				//ui.DrawDungeonView(NormalMode)
			} else if terrain(g.Dungeon.Cell(g.Player.Pos)) == StairCell && g.Objects.Stairs[g.Player.Pos] == BlockedStair {
				err = errors.New("The stairs are blocked by a magical stone barrier energies.")
			} else {
				err = errors.New("No stairs here.")
			}
		case BarrelCell:
			//ui.MenuSelectedAnimation(MenuInteract, true)
			err = g.Rest()
			//if err != nil {
			//ui.MenuSelectedAnimation(MenuInteract, false)
			//}
		case MagaraCell:
			again = true
			md.equipMagaraMenu()
		case StoneCell:
			//ui.MenuSelectedAnimation(MenuInteract, true)
			err = g.ActivateStone()
			//if err != nil {
			//ui.MenuSelectedAnimation(MenuInteract, false)
			//}
		case ScrollCell:
			again = true
			md.readScroll()
		case ItemCell:
			err = md.g.EquipItem()
		case LightCell:
			err = g.ExtinguishFire()
		case StoryCell:
			if g.Objects.Story[g.Player.Pos] == StoryArtifact && !g.LiberatedArtifact {
				g.PushEvent(&simpleEvent{ERank: g.Ev.Rank(), EAction: ArtifactAnimation})
				g.LiberatedArtifact = true
				g.Ev.Renew(g, DurationTurn)
			} else if g.Objects.Story[g.Player.Pos] == StoryArtifactSealed {
				err = errors.New("The artifact is protected by a magical stone barrier.")
			} else {
				err = errors.New("You cannot interact with anything here.")
			}
		default:
			err = errors.New("You cannot interact with anything here.")
		}
	case ActionEvoke:
		again = true
		md.evokeMagaraMenu()
	case ActionInventory:
		again = true
		md.openIventory()
	case ActionExplore:
		err = g.Autoexplore()
	case ActionExamine:
		again = true
		md.KeyboardExamine()
	case ActionHelp, ActionMenuCommandHelp:
		if md.mp.kbTargeting {
			md.ExamineHelp()
		} else {
			md.KeysHelp()
		}
		again = true
	case ActionMenuTargetingHelp:
		md.ExamineHelp()
		again = true
	case ActionMenu:
		md.openMenu()
		again = true
	case ActionLogs:
		if len(md.logs) > 0 {
			md.logs = md.logs[:len(md.logs)-1]
		}
		for _, e := range g.Log[len(md.logs):] {
			md.logs = append(md.logs, md.pagerMarkup.WithText(e.MText))
		}
		md.pager.SetLines(md.logs)
		md.mode = modePager
		md.pagerMode = modeLogs
		again = true
	case ActionSave:
		again = true
		errsave := g.Save()
		if errsave != nil {
			g.PrintfStyled("Error: %v", logError, errsave)
			g.PrintStyled("Could not save state.", logError)
			md.mode = modeNormal
		} else {
			md.mode = modeQuit
			eff = gruid.End()
		}
	case ActionDump:
		errdump := g.WriteDump()
		if errdump != nil {
			g.PrintfStyled("Error: %v", logError, errdump)
		} else {
			dataDir, _ := DataDir()
			if dataDir != "" {
				g.Printf("Game statistics written to %s.", filepath.Join(dataDir, "dump"))
			} else {
				g.Print("Game statistics written.")
			}
		}
		again = true
	case ActionWizardInfo:
		if g.Wizard {
			err = md.HandleWizardAction()
			again = true
		} else {
			err = actionErrorUnknown
		}
	case ActionWizardDescend:
		if g.Wizard && g.Depth == WinDepth {
			g.RescuedShaedra()
		}
		if g.Wizard && g.Depth < MaxDepth {
			g.StoryPrint("Descended wizardly")
			if g.Descend(DescendNormal) {
				md.win() // TODO: win
				//quit = true
				return again, eff, err
			}
		} else {
			err = actionErrorUnknown
		}
	case ActionWizard:
		md.EnterWizard()
		return true, eff, nil
	case ActionQuit:
		md.Quit()
		again = true
	case ActionSettings:
		md.openSettings()
		again = true
	case ActionSetKeys:
		again = true
		md.openKeyBindings()
	case ActionInvertLOS:
		again = true
		GameConfig.DarkLOS = !GameConfig.DarkLOS
		err := g.SaveConfig()
		if err != nil {
			g.Print(err.Error())
		}
		if GameConfig.DarkLOS {
			ApplyDarkLOS()
		} else {
			ApplyLightLOS()
		}
		md.mode = modeNormal
	case ActionToggleTiles:
		again = true
		md.ApplyToggleTiles()
		err := g.SaveConfig()
		if err != nil {
			g.Print(err.Error())
		}
	case ActionToggleShowNumbers:
		again = true
		GameConfig.ShowNumbers = !GameConfig.ShowNumbers
		err := g.SaveConfig()
		if err != nil {
			g.Print(err.Error())
		}
		md.updateStatusInfo()
		md.mode = modeNormal
	default:
		err = actionErrorUnknown
	}
	if err != nil {
		again = true
	}
	return again, eff, err
}

func (md *model) readScroll() {
	sc, ok := md.g.Objects.Scrolls[md.g.Player.Pos]
	if !ok {
		md.g.PrintStyled("Error while reading message.", logError)
		return
	}
	md.g.Print("You read the message.")
	md.mode = modeSmallPager
	st := gruid.Style{}
	switch sc {
	case ScrollLore:
		md.smallPager.SetBox(&ui.Box{Title: ui.Text("Lore Message").WithStyle(st.WithFg(ColorCyan))})
		stts := []ui.StyledText{}
		text := ui.Text(sc.Text(md.g)).Format(56)
		for _, s := range strings.Split(text.Text(), "\n") {
			stts = append(stts, ui.Text(s))
		}
		md.smallPager.SetLines(stts)
		if !md.g.Stats.Lore[md.g.Depth] {
			md.g.StoryPrint("Read lore message")
		}
		md.g.Stats.Lore[md.g.Depth] = true
		if len(md.g.Stats.Lore) == 4 {
			AchLoreStudent.Get(md.g)
		}
		if len(md.g.Stats.Lore) == len(md.g.Params.Lore) {
			AchLoremaster.Get(md.g)
		}
	default:
		md.smallPager.SetBox(&ui.Box{Title: ui.Text("Story Message").WithStyle(st.WithFg(ColorCyan))})
		stts := []ui.StyledText{}
		text := ui.Text(sc.Text(md.g)).Format(56)
		for _, s := range strings.Split(text.Text(), "\n") {
			stts = append(stts, ui.Text(s))
		}
		md.smallPager.SetLines(stts)
	}
}

type actionError int

const (
	actionErrorUnknown actionError = iota
)

func (e actionError) Error() string {
	switch e {
	case actionErrorUnknown:
		return "unknown action"
	}
	return ""
}

func altBgEntries(entries []ui.MenuEntry) {
	for i := range entries {
		if i%2 == 1 {
			st := entries[i].Text.Style()
			entries[i].Text = entries[i].Text.WithStyle(st.WithBg(ColorBgLOS))
		}
	}
}

var settingsActions = []action{
	ActionSetKeys,
	ActionInvertLOS,
	ActionToggleShowNumbers,
}

func (md *model) openSettings() {
	entries := []ui.MenuEntry{}
	r := 'a'
	for _, it := range settingsActions {
		entries = append(entries, ui.MenuEntry{
			Text: ui.Textf("%c - %s", r, it),
			Keys: []gruid.Key{gruid.Key(r)},
		})
		r++
	}
	altBgEntries(entries)
	md.menu.SetBox(&ui.Box{Title: ui.Text("Settings").WithStyle(gruid.Style{}.WithFg(ColorYellow))})
	md.menu.SetEntries(entries)
	md.mode = modeMenu
	md.menuMode = modeSettings
}

func (md *model) keysForAction(a action) string {
	keys := []gruid.Key{}
	for k, action := range md.keysNormal {
		if a == action && !k.In(keys) {
			keys = append(keys, k)
		}
	}
	for k, action := range md.keysTarget {
		if a == action && !k.In(keys) {
			keys = append(keys, k)
		}
	}
	// TODO: sort keys
	switch len(keys) {
	case 0:
		return ""
	case 1:
		return string(keys[0])
	}
	b := strings.Builder{}
	b.WriteString(string(keys[0]))
	for _, k := range keys {
		b.WriteString(" or ")
		b.WriteString(string(k))
	}
	return b.String()
}

func (md *model) openKeyBindings() {
	entries := []ui.MenuEntry{}
	r := 'a'
	for _, it := range ConfigurableKeyActions {
		desc := it.String()
		if !it.normalModeAction() {
			desc = it.targetingModeDescription()
		}
		desc = fmt.Sprintf(" %-36s %s", desc, md.keysForAction(it))
		entries = append(entries, ui.MenuEntry{
			Text: ui.Text(desc),
		})
		r++
	}
	altBgEntries(entries)
	md.keysMenu.SetBox(&ui.Box{Title: ui.Text("Key Bindings").WithStyle(gruid.Style{}.WithFg(ColorYellow))})
	md.keysMenu.SetEntries(entries)
	md.mode = modeMenu
	md.menuMode = modeKeys
}

var menuActions = []action{
	ActionLogs,
	ActionMenuCommandHelp,
	ActionMenuTargetingHelp,
	ActionSettings,
	ActionSave,
	ActionQuit,
}

func (md *model) openMenu() {
	entries := []ui.MenuEntry{}
	r := 'a'
	for _, it := range menuActions {
		entries = append(entries, ui.MenuEntry{
			Text: ui.Textf("%c - %s", r, it),
			Keys: []gruid.Key{gruid.Key(r)},
		})
		r++
	}
	altBgEntries(entries)
	md.menu.SetBox(&ui.Box{Title: ui.Text("Menu").WithStyle(gruid.Style{}.WithFg(ColorYellow))})
	md.menu.SetEntries(entries)
	md.mode = modeMenu
	md.menuMode = modeGameMenu
}

func (md *model) openIventory() {
	entries := []ui.MenuEntry{}
	items := []item{md.g.Player.Inventory.Body, md.g.Player.Inventory.Neck, md.g.Player.Inventory.Misc}
	parts := []string{"body", "neck", "backpack"}
	r := 'a'
	for i, it := range items {
		entries = append(entries, ui.MenuEntry{
			Text: ui.Textf("%c - %s (%s)", r, it.ShortDesc(md.g), parts[i]),
			Keys: []gruid.Key{gruid.Key(r)},
		})
		r++
	}
	altBgEntries(entries)
	md.menu.SetBox(&ui.Box{Title: ui.Text("Inventory").WithStyle(gruid.Style{}.WithFg(ColorYellow))})
	md.menu.SetEntries(entries)
	md.mode = modeMenu
	md.menuMode = modeInventory
	md.description.StyledText = ui.Text(items[md.menu.Active()].Desc(md.g)).Format(UIWidth/2 - 1 - 2)
}

func (md *model) evokeMagaraMenu() {
	entries := []ui.MenuEntry{}
	items := md.g.Player.Magaras
	r := 'a'
	for _, it := range items {
		entries = append(entries, ui.MenuEntry{
			Text: ui.Textf("%c - %s ", r, it.ShortDesc()),
			Keys: []gruid.Key{gruid.Key(r)},
		})
		r++
	}
	altBgEntries(entries)
	md.menu.SetBox(&ui.Box{Title: ui.Text("Evoke Magara").WithStyle(gruid.Style{}.WithFg(ColorYellow))})
	md.menu.SetEntries(entries)
	md.mode = modeMenu
	md.menuMode = modeEvokation
	md.description.StyledText = ui.Text(items[md.menu.Active()].Desc(md.g)).Format(UIWidth/2 - 1 - 2)
}

func (md *model) equipMagaraMenu() {
	entries := []ui.MenuEntry{}
	items := md.g.Player.Magaras
	r := 'a'
	for _, it := range items {
		entries = append(entries, ui.MenuEntry{
			Text: ui.Textf("%c - %s ", r, it.ShortDesc()),
			Keys: []gruid.Key{gruid.Key(r)},
		})
		r++
	}
	altBgEntries(entries)
	md.menu.SetBox(&ui.Box{Title: ui.Text("Equip Magara").WithStyle(gruid.Style{}.WithFg(ColorYellow))})
	md.menu.SetEntries(entries)
	md.mode = modeMenu
	md.menuMode = modeEquip
	md.description.StyledText = ui.Text(items[md.menu.Active()].Desc(md.g)).Format(UIWidth/2 - 1 - 2)
}

//func (ui *model) HandleKey(rka runeKeyAction) (again bool, quit bool, err error) {
//g := ui.st
//switch rka.k {
//case ActionW, ActionS, ActionN, ActionE:
//err = g.PlayerBump(To(KeyToDir(rka.k), g.Player.Pos))
//case ActionRunW, ActionRunS, ActionRunN, ActionRunE:
//err = g.GoToDir(KeyToDir(rka.k))
//case ActionWaitTurn:
//g.WaitTurn()
//case ActionGoToStairs:
//stairs := g.StairsSlice()
//sortedStairs := g.SortedNearestTo(stairs, g.Player.Pos)
//if len(sortedStairs) > 0 {
//stair := sortedStairs[0]
//if g.Player.Pos == stair {
//err = errors.New("You are already on the stairs.")
//break
//}
//ex := &examiner{stairs: true}
//err = ex.Action(g, stair)
//if err == nil && !g.MoveToTarget() {
//err = errors.New("You could not move toward stairs.")
//}
//if ex.Done() {
//g.Targeting = InvalidPos
//}
//} else {
//err = errors.New("You cannot go to any stairs.")
//}
//case ActionInteract:
//c := g.Dungeon.Cell(g.Player.Pos)
//switch c.T {
//case StairCell:
//if g.Dungeon.Cell(g.Player.Pos).T == StairCell && g.Objects.Stairs[g.Player.Pos] != BlockedStair {
//// TODO: animation
////ui.MenuSelectedAnimation(MenuInteract, true)
//strt := g.Objects.Stairs[g.Player.Pos]
//err = ui.OptionalDescendConfirmation(strt)
//if err != nil {
//break
//}
//if g.Descend(DescendNormal) {
//ui.Win()
//quit = true
//return again, quit, err
//}
////ui.DrawDungeonView(NormalMode)
//} else if g.Dungeon.Cell(g.Player.Pos).T == StairCell && g.Objects.Stairs[g.Player.Pos] == BlockedStair {
//err = errors.New("The stairs are blocked by a magical stone barrier energies.")
//} else {
//err = errors.New("No stairs here.")
//}
//case BarrelCell:
////ui.MenuSelectedAnimation(MenuInteract, true)
//err = g.Rest()
//if err != nil {
////ui.MenuSelectedAnimation(MenuInteract, false)
//}
//case MagaraCell:
//err = ui.EquipMagara()
//err = ui.CleanError(err)
//case StoneCell:
////ui.MenuSelectedAnimation(MenuInteract, true)
//err = g.ActivateStone()
//if err != nil {
////ui.MenuSelectedAnimation(MenuInteract, false)
//}
//case ScrollCell:
//err = ui.ReadScroll()
//err = ui.CleanError(err)
//case ItemCell:
//err = ui.st.EquipItem()
//case LightCell:
//err = g.ExtinguishFire()
//case StoryCell:
//if g.Objects.Story[g.Player.Pos] == StoryArtifact && !g.LiberatedArtifact {
//g.PushEvent(&simpleEvent{ERank: g.Ev.Rank(), EAction: ArtifactAnimation})
//g.LiberatedArtifact = true
//g.Ev.Renew(g, DurationTurn)
//} else if g.Objects.Story[g.Player.Pos] == StoryArtifactSealed {
//err = errors.New("The artifact is protected by a magical stone barrier.")
//} else {
//err = errors.New("You cannot interact with anything here.")
//}
//default:
//err = errors.New("You cannot interact with anything here.")
//}
//case ActionEvoke:
//err = ui.SelectMagara()
//err = ui.CleanError(err)
//case ActionInventory:
//err = ui.SelectItem()
//err = ui.CleanError(err)
//case ActionExplore:
//err = g.Autoexplore()
//case ActionExamine:
//again, quit, err = ui.Examine(nil)
//case ActionHelp, ActionMenuCommandHelp:
//ui.KeysHelp()
//again = true
//case ActionMenuTargetingHelp:
//ui.ExamineHelp()
//again = true
//case ActionLogs:
////ui.DrawPreviousLogs()
//again = true
//case ActionSave:
//g.Ev.Renew(g, 0)
//errsave := g.Save()
//if errsave != nil {
//g.PrintfStyled("Error: %v", logError, errsave)
//g.PrintStyled("Could not save state.", logError)
//} else {
//quit = true
//}
//case ActionDump:
//errdump := g.WriteDump()
//if errdump != nil {
//g.PrintfStyled("Error: %v", logError, errdump)
//} else {
//dataDir, _ := g.DataDir()
//if dataDir != "" {
//g.Printf("Game statistics written to %s.", filepath.Join(dataDir, "dump"))
//} else {
//g.Print("Game statistics written.")
//}
//}
//again = true
//case ActionWizardInfo:
//if g.Wizard {
//err = ui.HandleWizardAction()
//again = true
//} else {
//err = errors.New("Unknown key. Type ? for help.")
//}
//case ActionWizardDescend:
//if g.Wizard && g.Depth == WinDepth {
//g.RescuedShaedra()
//}
//if g.Wizard && g.Depth < MaxDepth {
//g.StoryPrint("Descended wizardly")
//if g.Descend(DescendNormal) {
//ui.Win()
//quit = true
//return again, quit, err
//}
//} else {
//err = errors.New("Unknown key. Type ? for help.")
//}
//case ActionWizard:
//ui.EnterWizard()
//return true, false, nil
//case ActionQuit:
////if ui.Quit() {
////return false, true, nil
////}
//return true, false, nil
//case ActionConfigure:
//err = ui.HandleSettingAction()
//again = true
//case ActionDescription:
////ui.MenuSelectedAnimation(MenuView, false)
//err = fmt.Errorf("You must choose a target to describe.")
//case ActionExclude:
//err = fmt.Errorf("You must choose a target for exclusion.")
//default:
//err = fmt.Errorf("Unknown key '%c'. Type ? for help.", rka.r)
//}
//if err != nil {
//again = true
//}
//return again, quit, err
//}

type wizardAction int

const (
	WizardInfoAction wizardAction = iota
	WizardToggleMode
)

func (a wizardAction) String() (text string) {
	switch a {
	case WizardInfoAction:
		text = "Info"
	case WizardToggleMode:
		text = "toggle normal/map/all wizard mode"
	}
	return text
}

var WizardActions = []wizardAction{
	WizardInfoAction,
	WizardToggleMode,
}

func (md *model) HandleWizardAction() error {
	// TODO: rewrite
	//g := ui.st
	//s, err := ui.SelectWizardMagic(WizardActions)
	//if err != nil {
	//return err
	//}
	//switch s {
	//case WizardInfoAction:
	//ui.WizardInfo()
	//case WizardToggleMode:
	//switch g.WizardMode {
	//case WizardNormal:
	//g.WizardMode = WizardMap
	//case WizardMap:
	//g.WizardMode = WizardSeeAll
	//case WizardSeeAll:
	//g.WizardMode = WizardNormal
	//}
	//g.StoryPrint("Toggle wizard mode.")
	////ui.DrawDungeonView(NoFlushMode)
	//}
	return nil
}

func (md *model) updateKeysDescription(title string, actions []string) {
	md.pagerMode = modeHelpKeys
	md.mode = modePager
	if CustomKeys {
		title = fmt.Sprintf(" Default %s ", title)
	} else {
		title = fmt.Sprintf(" %s ", title)
	}
	md.pager.SetBox(&ui.Box{Title: ui.Text(title).WithStyle(gruid.Style{}.WithFg(ColorYellow))})
	lines := []ui.StyledText{}
	for i := 0; i < len(actions)-1; i += 2 {
		stt := ui.StyledText{}
		if actions[i+1] != "" {
			stt = stt.WithTextf(" %-36s %s", actions[i], actions[i+1])
		} else {
			stt = stt.WithTextf(" %s ", actions[i]).WithStyle(gruid.Style{}.WithFg(ColorCyan))
		}
		if i%4 == 2 {
			stt = stt.WithStyle(stt.Style().WithBg(ColorBgLOS))
		}
		lines = append(lines, stt)
	}
	md.pager.SetLines(lines)
}

func (md *model) KeysHelp() {
	md.updateKeysDescription("Basic Commands", []string{
		"Move/Jump", "arrows or wasd or hjkl or mouse left",
		"Wait a turn", "“.” or 5 or enter or mouse left on @",
		"Interact (Equip/Descend/Rest...)", "e",
		"Evoke/Zap magara", "v or z",
		"Inventory", "i",
		"Examine", "x or mouse hover",
		"Menu", "M",
		"Advanced Commands", "",
		"Save and Quit", "S",
		"View previous messages", "m",
		"Go to nearest stairs", "G",
		"Autoexplore (use with caution)", "o",
		"Write state statistics to file", "#",
		"Quit without saving", "Q",
		"Change settings and key bindings", "=",
	})
}

func (md *model) ExamineHelp() {
	md.updateKeysDescription("Examine/Travel Commands", []string{
		"Move cursor", "arrows or wasd or hjkl or mouse hover",
		"Go to/select target", "“.” or enter or mouse left",
		"View target description", "v or mouse right",
		"Cycle through monsters", "+",
		"Cycle through stairs", ">",
		"Cycle through objects", "o",
		"Toggle exclude area from auto-travel", "e or mouse middle",
	})
}

func (md *model) WizardInfo() {
	// TODO
	//g := ui.st
	//ui.Clear()
	//b := &bytes.Buffer{}
	//fmt.Fprintf(b, "Monsters: %d (%d)\n", len(g.Monsters), g.MaxMonsters())
	//fmt.Fprintf(b, "Danger: %d (%d)\n", g.Danger(), g.MaxDanger())
	//ui.DrawText(b.String(), 0, 0)
	//ui.Flush()
	//ui.WaitForContinue(-1)
}

func (md *model) EnterWizard() {
	// TODO
	//g := ui.st
	//if ui.Wizard() {
	//g.EnterWizardMode()
	//ui.DrawDungeonView(NoFlushMode)
	//} else {
	//g.Print(DoNothing)
	//}
}

func (md *model) checkShaedra(st stair) (err error) {
	g := md.g
	if g.Depth == WinDepth && st == NormalStair && terrain(g.Dungeon.Cell(g.Places.Shaedra)) == StoryCell {
		err = errors.New("You have to rescue Shaedra first!")
	}
	return err

}

// Quit enters confirmation mode for quit without saving.
func (md *model) Quit() {
	md.g.Print("Do you really want to quit without saving? [y/N]")
	md.mode = modeQuitConfirmation
}
