package main

import (
	"errors"
	"fmt"
	"path/filepath"
	//"time"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/ui"
)

func (md *model) HideCursor() {
	md.mp.ex.pos = InvalidPos
}

func (md *model) SetCursor(pos gruid.Point) {
	md.mp.ex.pos = pos
}

func (md *model) GetPos(i int) (int, int) {
	return i - (i/UIWidth)*UIWidth, i / UIWidth
}

type uiMode int

const (
	NormalMode uiMode = iota
	TargetingMode
	NoFlushMode
	AnimationMode
)

const DoNothing = "Do nothing, then."

func (md *model) EnterWizard() {
	//g := ui.st
	//if ui.Wizard() {
	//g.EnterWizardMode()
	//ui.DrawDungeonView(NoFlushMode)
	//} else {
	//g.Print(DoNothing)
	//}
}

func (md *model) CleanError(err error) error {
	if err != nil && err.Error() == DoNothing {
		err = errors.New("")
	}
	return err
}

type action int

const (
	ActionNothing action = iota
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
	ActionDescription
	ActionTarget
	ActionExclude
	ActionEscape

	ActionConfigure
	ActionMenu
	ActionNextStairs
	ActionMenuCommandHelp
	ActionMenuTargetingHelp
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
	ActionDescription,
	ActionTarget,
	ActionExclude}

var CustomKeys bool

func (k action) NormalModeAction() bool {
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
		ActionConfigure,
		ActionWizard,
		ActionWizardInfo:
		return true
	default:
		return false
	}
}

func (k action) NormalModeDescription() (text string) {
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
	case ActionConfigure:
		text = "Settings and key bindings"
	case ActionWizard:
		text = "Wizard (debug) mode"
	case ActionWizardInfo:
		text = "Wizard (debug) mode information"
	case ActionMenu:
		text = "Action Menu"
	}
	return text
}

func (k action) TargetingModeDescription() (text string) {
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
	case ActionDescription:
		text = "View target description"
	case ActionTarget:
		text = "Go to"
	case ActionExclude:
		text = "Toggle exclude area from auto-travel"
	case ActionEscape:
		text = "Quit targeting mode"
	case ActionMenu:
		text = "Action Menu"
	}
	return text
}

func (k action) TargetingModeAction() bool {
	switch k {
	case ActionW, ActionS, ActionN, ActionE,
		ActionRunW, ActionRunS, ActionRunN, ActionRunE,
		ActionDescend,
		ActionPreviousMonster,
		ActionNextMonster,
		ActionNextObject,
		ActionNextStairs,
		ActionDescription,
		ActionTarget,
		ActionExclude,
		ActionEscape:
		return true
	default:
		return false
	}
}

var GameConfig config

func ApplyDefaultKeyBindings() {
	// TODO: rewrite with gruid.Key
	GameConfig.RuneNormalModeKeys = map[rune]action{
		'h': ActionW,
		'j': ActionS,
		'k': ActionN,
		'l': ActionE,
		'a': ActionW,
		's': ActionS,
		'w': ActionN,
		'd': ActionE,
		'4': ActionW,
		'2': ActionS,
		'8': ActionN,
		'6': ActionE,
		'H': ActionRunW,
		'J': ActionRunS,
		'K': ActionRunN,
		'L': ActionRunE,
		'.': ActionWaitTurn,
		'5': ActionWaitTurn,
		'G': ActionGoToStairs,
		'o': ActionExplore,
		'x': ActionExamine,
		'v': ActionEvoke,
		'z': ActionEvoke,
		'e': ActionInteract,
		'i': ActionInventory,
		'm': ActionLogs,
		'M': ActionMenu,
		'#': ActionDump,
		'?': ActionHelp,
		'S': ActionSave,
		'Q': ActionQuit,
		'W': ActionWizard,
		'@': ActionWizardInfo,
		'>': ActionWizardDescend,
		'=': ActionConfigure,
	}
	GameConfig.RuneTargetModeKeys = map[rune]action{
		'h':    ActionW,
		'j':    ActionS,
		'k':    ActionN,
		'l':    ActionE,
		'a':    ActionW,
		's':    ActionS,
		'w':    ActionN,
		'd':    ActionE,
		'4':    ActionW,
		'2':    ActionS,
		'8':    ActionN,
		'6':    ActionE,
		'H':    ActionRunW,
		'J':    ActionRunS,
		'K':    ActionRunN,
		'L':    ActionRunE,
		'>':    ActionNextStairs,
		'-':    ActionPreviousMonster,
		'+':    ActionNextMonster,
		'o':    ActionNextObject,
		']':    ActionNextObject,
		')':    ActionNextObject,
		'(':    ActionNextObject,
		'[':    ActionNextObject,
		'_':    ActionNextObject,
		'=':    ActionNextObject,
		'v':    ActionDescription,
		'.':    ActionTarget,
		't':    ActionTarget,
		'g':    ActionTarget,
		'e':    ActionExclude,
		' ':    ActionEscape,
		'\x1b': ActionEscape,
		'x':    ActionEscape,
		'X':    ActionEscape,
		'?':    ActionHelp,
	}
	CustomKeys = false
}

func (md *model) OptionalDescendConfirmation(st stair) (err error) {
	g := md.g
	if g.Depth == WinDepth && st == NormalStair && g.Dungeon.Cell(g.Places.Shaedra).T == StoryCell {
		err = errors.New("You have to rescue Shaedra first!")
	}
	return err

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

func (md *model) normalModeAction(action action) (again bool, eff gruid.Effect, err error) {
	g := md.g
	switch action {
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
		md.ExcludeZone(md.mp.ex.pos)
	case ActionPreviousMonster:
		again = true
		md.NextMonster("-", md.mp.ex.pos, md.mp.ex)
	case ActionNextMonster:
		again = true
		md.NextMonster("+", md.mp.ex.pos, md.mp.ex)
	case ActionNextObject:
		again = true
		md.NextObject(md.mp.ex.pos, md.mp.ex)
	case ActionTarget:
		again = true
		err = md.Target()
		if err != nil {
			break
		}
		g.Targeting = InvalidPos
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
			err = md.Target()
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
		switch c.T {
		case StairCell:
			if g.Dungeon.Cell(g.Player.Pos).T == StairCell && g.Objects.Stairs[g.Player.Pos] != BlockedStair {
				// TODO: animation
				//ui.MenuSelectedAnimation(MenuInteract, true)
				strt := g.Objects.Stairs[g.Player.Pos]
				err = md.OptionalDescendConfirmation(strt)
				if err != nil {
					break
				}
				if g.Descend(DescendNormal) {
					md.Win()
					// TODO: win
					return again, eff, err
				}
				//ui.DrawDungeonView(NormalMode)
			} else if g.Dungeon.Cell(g.Player.Pos).T == StairCell && g.Objects.Stairs[g.Player.Pos] == BlockedStair {
				err = errors.New("The stairs are blocked by a magical stone barrier energies.")
			} else {
				err = errors.New("No stairs here.")
			}
		case BarrelCell:
			//ui.MenuSelectedAnimation(MenuInteract, true)
			err = g.Rest()
			if err != nil {
				//ui.MenuSelectedAnimation(MenuInteract, false)
			}
		case MagaraCell:
			again = true
			md.EquipMagaraMenu()
		case StoneCell:
			//ui.MenuSelectedAnimation(MenuInteract, true)
			err = g.ActivateStone()
			if err != nil {
				//ui.MenuSelectedAnimation(MenuInteract, false)
			}
		case ScrollCell:
			err = md.ReadScroll()
			err = md.CleanError(err)
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
		md.EvokeMagaraMenu()
	case ActionInventory:
		again = true
		md.OpenIventory()
	case ActionExplore:
		err = g.Autoexplore()
	case ActionExamine:
		again = true
		md.KeyboardExamine()
	case ActionHelp, ActionMenuCommandHelp:
		md.KeysHelp()
		again = true
	case ActionMenuTargetingHelp:
		md.ExamineHelp()
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
				md.Win() // TODO: win
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
	case ActionConfigure:
		err = md.HandleSettingAction()
		again = true
	case ActionDescription:
		again = true
		//ui.MenuSelectedAnimation(MenuView, false)
		err = fmt.Errorf("You must choose a target to describe.")
	default:
		err = actionErrorUnknown
	}
	if err != nil {
		again = true
	}
	return again, eff, err
}

func altBgEntries(entries []ui.MenuEntry) {
	for i := range entries {
		if i%2 == 1 {
			st := entries[i].Text.Style()
			entries[i].Text = entries[i].Text.WithStyle(st.WithBg(ColorBgLOS))
		}
	}
}

func (md *model) OpenIventory() {
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

func (md *model) EvokeMagaraMenu() {
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

func (md *model) EquipMagaraMenu() {
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

func (md *model) Death() {
	// TODO: fix this
	g := md.g
	if len(g.Stats.Achievements) == 0 {
		NoAchievement.Get(g)
	}
	g.Print("You die... [(x) to continue]")
	err := g.WriteDump()
	md.Dump(err)
}

func (md *model) Win() {
	// TODO: rewrite
	g := md.g
	err := g.RemoveSaveFile()
	if err != nil {
		g.PrintfStyled("Error removing save file: %v", logError, err)
	}
	if g.Wizard {
		g.Print("You escape by the magic portal! **WIZARD** [(x) to continue]")
	} else {
		g.Print("You escape by the magic portal! [(x) to continue]")
	}
	//ui.DrawDungeonView(NormalMode)
	//ui.WaitForContinue(-1)
	err = g.WriteDump()
	md.Dump(err)
	//ui.WaitForContinue(-1)
}

func (md *model) Dump(err error) {
	//g := ui.st
	//ui.DrawText(g.SimplifedDump(err), 0, 0)
}

func (md *model) CriticalHPWarning() {
	// TODO
	//g := ui.st
	//g.PrintStyled("*** CRITICAL HP WARNING *** [(x) to continue]", logCritic)
	//ui.DrawDungeonView(NormalMode)
	//ui.WaitForContinue(ui.MapHeight())
	//g.Print("Ok. Be careful, then.")
}

func (md *model) Quit() {
	md.g.Print("Do you really want to quit without saving? [y/N]")
	md.mode = modeQuitConfirmation
}

func (md *model) Clear() {
	c := gruid.Cell{}
	c.Rune = ' '
	c.Style = gruid.Style{Fg: ColorFg, Bg: ColorBg}
	md.gd.Fill(c)
}

func ApplyConfig() {
	if GameConfig.RuneNormalModeKeys == nil || GameConfig.RuneTargetModeKeys == nil {
		ApplyDefaultKeyBindings()
	}
	if GameConfig.DarkLOS {
		ApplyDarkLOS()
	} else {
		ApplyLightLOS()
	}
}
