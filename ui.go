package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/anaseto/gruid"
)

func (ui *model) HideCursor() {
	ui.mp.cursor = InvalidPos
}

func (ui *model) SetCursor(pos gruid.Point) {
	ui.mp.cursor = pos
}

func (ui *model) GetPos(i int) (int, int) {
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

func (ui *model) EnterWizard() {
	//g := ui.st
	//if ui.Wizard() {
	//g.EnterWizardMode()
	//ui.DrawDungeonView(NoFlushMode)
	//} else {
	//g.Print(DoNothing)
	//}
}

func (ui *model) CleanError(err error) error {
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

//type runeKeyAction struct {
//r rune
//k action
//}

//func (ui *model) HandleKeyAction(rka runeKeyAction) (again bool, quit bool, err error) {
//if rka.r != 0 {
//var ok bool
//rka.k, ok = GameConfig.RuneNormalModeKeys[rka.r]
//if !ok {
//switch rka.r {
//case 's':
//err = errors.New("Unknown key. Did you mean capital S for save and quit?")
//case 'q':
//err = errors.New("Unknown key. Did you mean capital Q for quit without saving?")
//default:
//err = fmt.Errorf("Unknown key '%c'. Type ? for help.", rka.r)
//}
//return again, quit, err
//}
//}
//if rka.k == ActionMenu {
//rka.k, err = ui.SelectAction(menuActions)
//if err != nil {
//err = ui.CleanError(err)
//return again, quit, err
//}
//}
//return ui.HandleKey(rka)
//}

func (ui *model) OptionalDescendConfirmation(st stair) (err error) {
	g := ui.st
	if g.Depth == WinDepth && st == NormalStair && g.Dungeon.Cell(g.Places.Shaedra).T == StoryCell {
		err = errors.New("You have to rescue Shaedra first!")
	}
	return err

}

func (ui *model) normalModeKeyDown(key gruid.Key) (again bool, err error) {
	g := ui.st
	action := ui.keysNormal[key]
	if ui.mp.targeting {
		again = true
		action = ui.keysTarget[key]
	}
	switch action {
	case ActionW, ActionS, ActionN, ActionE:
		if !ui.mp.targeting {
			err = g.PlayerBump(To(KeyToDir(action), g.Player.Pos))
		} else {
			p := To(KeyToDir(action), ui.mp.ex.pos)
			if valid(p) {
				ui.Examine(p)
			}
		}
	case ActionRunW, ActionRunS, ActionRunN, ActionRunE:
		if !ui.mp.targeting {
			err = g.GoToDir(KeyToDir(action))
		} else {
			q := InvalidPos
			for i := 0; i < 5; i++ {
				p := To(KeyToDir(action), ui.mp.ex.pos)
				if !valid(p) {
					break
				}
				q = p
			}
			if q != InvalidPos {
				ui.Examine(q)
			}
		}
	case ActionExclude:
		ui.ExcludeZone(ui.mp.ex.pos)
	case ActionPreviousMonster, ActionNextMonster:
		ui.NextMonster(key, ui.mp.ex.pos, ui.mp.ex)
	case ActionNextObject:
		ui.NextObject(ui.mp.ex.pos, ui.mp.ex)
	case ActionTarget:
		err = ui.Target()
		if err != nil {
			break
		}
		g.Targeting = InvalidPos
		if g.MoveToTarget() {
			again = false
		}
	case ActionEscape:
		ui.CancelExamine()
	case ActionWaitTurn:
		g.WaitTurn()
	case ActionGoToStairs:
		stairs := g.StairsSlice()
		sortedStairs := g.SortedNearestTo(stairs, g.Player.Pos)
		if len(sortedStairs) > 0 {
			stair := sortedStairs[0]
			if g.Player.Pos == stair {
				err = errors.New("You are already on the stairs.")
				break
			}
			ui.mp.ex.pos = stair
			err = ui.Target()
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
				err = ui.OptionalDescendConfirmation(strt)
				if err != nil {
					break
				}
				if g.Descend(DescendNormal) {
					ui.Win()
					// TODO: win
					return again, err
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
			err = ui.EquipMagara()
			err = ui.CleanError(err)
		case StoneCell:
			//ui.MenuSelectedAnimation(MenuInteract, true)
			err = g.ActivateStone()
			if err != nil {
				//ui.MenuSelectedAnimation(MenuInteract, false)
			}
		case ScrollCell:
			err = ui.ReadScroll()
			err = ui.CleanError(err)
		case ItemCell:
			err = ui.st.EquipItem()
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
		err = ui.SelectMagara()
		err = ui.CleanError(err)
	case ActionInventory:
		err = ui.SelectItem()
		err = ui.CleanError(err)
	case ActionExplore:
		err = g.Autoexplore()
	case ActionExamine:
		again = true
		ui.StartExamine()
	case ActionHelp, ActionMenuCommandHelp:
		ui.KeysHelp()
		again = true
	case ActionMenuTargetingHelp:
		ui.ExamineHelp()
		again = true
	case ActionLogs:
		//ui.DrawPreviousLogs()
		again = true
	case ActionSave:
		g.Ev.Renew(g, 0)
		errsave := g.Save()
		if errsave != nil {
			g.PrintfStyled("Error: %v", logError, errsave)
			g.PrintStyled("Could not save state.", logError)
		} else {
			// TODO quit on save
			again = true
		}
	case ActionDump:
		errdump := g.WriteDump()
		if errdump != nil {
			g.PrintfStyled("Error: %v", logError, errdump)
		} else {
			dataDir, _ := g.DataDir()
			if dataDir != "" {
				g.Printf("Game statistics written to %s.", filepath.Join(dataDir, "dump"))
			} else {
				g.Print("Game statistics written.")
			}
		}
		again = true
	case ActionWizardInfo:
		if g.Wizard {
			err = ui.HandleWizardAction()
			again = true
		} else {
			err = errors.New("Unknown key. Type ? for help.")
		}
	case ActionWizardDescend:
		if g.Wizard && g.Depth == WinDepth {
			g.RescuedShaedra()
		}
		if g.Wizard && g.Depth < MaxDepth {
			g.StoryPrint("Descended wizardly")
			if g.Descend(DescendNormal) {
				ui.Win() // TODO: win
				//quit = true
				return again, err
			}
		} else {
			err = errors.New("Unknown key. Type ? for help.")
		}
	case ActionWizard:
		ui.EnterWizard()
		return true, nil
	case ActionQuit:
		//if ui.Quit() {
		//return false, true, nil
		//}
		return true, nil
	case ActionConfigure:
		err = ui.HandleSettingAction()
		again = true
	case ActionDescription:
		//ui.MenuSelectedAnimation(MenuView, false)
		err = fmt.Errorf("You must choose a target to describe.")
	default:
		err = fmt.Errorf("Unknown key '%s'. Type ? for help.", key)
	}
	if err != nil {
		again = true
	}
	return again, err
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

//func (ui *model) CursorMouseLeft(targ Targeter, pos gruid.Point, data *examineData) (again, notarg bool) {
//// TODO: rewrite
//g := ui.st
//again = true
//if data.pos == pos {
//err := targ.Action(g, pos)
//if err != nil {
//g.Print(err.Error())
//} else {
//if g.MoveToTarget() {
//again = false
//}
//}
//} else {
//data.pos = pos
//}
//return again, notarg
//}

//func (ui *model) CursorKeyAction(targ Targeter, rka runeKeyAction, data *examineData) (again, quit, notarg bool, err error) {
//// TODO: rewrite
//g := ui.st
//pos := data.npos
//again = true
//if rka.r != 0 {
//var ok bool
//rka.k, ok = GameConfig.RuneTargetModeKeys[rka.r]
//if !ok {
//err = fmt.Errorf("Invalid targeting mode key '%c'. Type ? for help.", rka.r)
//return again, quit, notarg, err
//}
//}
//if rka.k == ActionMenu {
//rka.k, err = ui.SelectAction(menuActions)
//if err != nil {
//err = ui.CleanError(err)
//return again, quit, notarg, err
//}
//}
//switch rka.k {
//case ActionW, ActionS, ActionN, ActionE:
//data.npos = To(KeyToDir(rka.k), pos)
//case ActionRunW, ActionRunS, ActionRunN, ActionRunE:
//for i := 0; i < 5; i++ {
//p := To(KeyToDir(rka.k), data.npos)
//if !valid(p) {
//break
//}
//data.npos = p
//}
//case ActionNextStairs:
//ui.NextStair(data)
//case ActionDescend:
//if g.Dungeon.Cell(g.Player.Pos).T == StairCell && g.Objects.Stairs[g.Player.Pos] != BlockedStair {
////ui.MenuSelectedAnimation(MenuInteract, true)
//strt := g.Objects.Stairs[g.Player.Pos]
//err = ui.OptionalDescendConfirmation(strt)
//if err != nil {
//break
//}
//again = false
//g.Targeting = InvalidPos
//notarg = true
//if g.Descend(DescendNormal) {
//ui.Win()
//quit = true
//return again, quit, notarg, err
//}
//} else if g.Dungeon.Cell(g.Player.Pos).T == StairCell && g.Objects.Stairs[g.Player.Pos] == BlockedStair {
//err = errors.New("The stairs are blocked by a magical stone barrier energies.")
//} else {
//err = errors.New("No stairs here.")
//}
//case ActionPreviousMonster, ActionNextMonster:
//ui.NextMonster(rka.r, pos, data)
//case ActionNextObject:
//ui.NextObject(pos, data)
//case ActionHelp, ActionMenuTargetingHelp:
//ui.HideCursor()
//ui.ExamineHelp()
//ui.SetCursor(pos)
//case ActionMenuCommandHelp:
//ui.HideCursor()
//ui.KeysHelp()
//ui.SetCursor(pos)
//case ActionTarget:
//err = targ.Action(g, pos)
//if err != nil {
//break
//}
//g.Targeting = InvalidPos
//if g.MoveToTarget() {
//again = false
//}
//if targ.Done() {
//notarg = true
//}
//case ActionDescription:
//ui.HideCursor()
////ui.ViewPositionDescription(pos)
//ui.SetCursor(pos)
//case ActionExclude:
//ui.ExcludeZone(pos)
//case ActionEscape:
//g.Targeting = InvalidPos
//notarg = true
//err = errors.New(DoNothing)
//case ActionExplore, ActionLogs, ActionEvoke, ActionInventory:
//// XXX: hm, this is only useful with mouse in terminal, rarely tested.
//if _, ok := targ.(*examiner); !ok {
//break
//}
//again, quit, err = ui.HandleKey(rka)
//if err != nil {
//notarg = true
//}
//g.Targeting = InvalidPos
//case ActionConfigure:
//err = ui.HandleSettingAction()
//case ActionSave:
//g.Ev.Renew(g, 0)
//g.Highlight = nil
//g.Targeting = InvalidPos
//errsave := g.Save()
//if errsave != nil {
//g.PrintfStyled("Error: %v", logError, errsave)
//g.PrintStyled("Could not save state.", logError)
//} else {
//notarg = true
//again = false
//quit = true
//}
//case ActionQuit:
////if ui.Quit() {
////quit = true
////again = false
////}
//default:
//err = fmt.Errorf("Invalid targeting mode key '%c'. Type ? for help.", rka.r)
//}
//return again, quit, notarg, err
//}

//func (ui *model) CursorAction(targ Targeter, start *gruid.Point) (again, quit bool, err error) {
//// TODO: rewrite
//g := ui.st
//pos := g.Player.Pos
//if start != nil {
//pos = *start
//} else {
//minDist := 999
//for _, mons := range g.Monsters {
//if mons.Exists() && g.Player.LOS[mons.Pos] {
//dist := Distance(mons.Pos, g.Player.Pos)
//if minDist > dist {
//minDist = dist
//pos = mons.Pos
//}
//}
//}
//}
//data := &examineData{
//pos:     pos,
//objects: []gruid.Point{},
//}
//if _, ok := targ.(*examiner); ok && pos == g.Player.Pos && start == nil {
//ui.NextObject(InvalidPos, data)
//if !valid(data.pos) {
//ui.NextStair(data)
//}
//if valid(data.pos) && Distance(pos, data.pos) < DefaultLOSRange+5 {
//pos = data.pos
//}
//}
//opos := InvalidPos
//loop:
//for {
//err = nil
//if pos != opos {
////ui.DescribePosition(pos, targ)
//}
//opos = pos
////targ.ComputeHighlight(g, pos)
//ui.SetCursor(pos)
//m := g.MonsterAt(pos)
//if m.Exists() && g.Player.Sees(pos) {
//g.ComputeMonsterCone(m)
//} else {
//g.MonsterTargLOS = nil
//}
////ui.DrawDungeonView(TargetingMode)
////ui.DrawInfoLine(g.InfoEntry)
////ui.SetCell(DungeonWidth, ui.MapHeight(), '┤', ColorFg, ColorBg)
////ui.Flush()
//data.pos = pos
//var notarg bool
////again, quit, notarg, err = ui.TargetModeEvent(targ, data)
//if err != nil {
//err = ui.CleanError(err)
//}
//if !again || notarg {
//break loop
//}
//if err != nil {
//g.Print(err.Error())
//}
//if valid(data.pos) {
//pos = data.pos
//}
//}
//g.Highlight = nil
//g.MonsterTargLOS = nil
//ui.HideCursor()
//return again, quit, err
//}

type menu int

const (
	//MenuExplore menu = iota
	MenuOther menu = iota
	MenuInventory
	MenuEvoke
	MenuInteract
)

func (m menu) String() (text string) {
	switch m {
	//case MenuExplore:
	//text = "explore"
	case MenuEvoke:
		text = "evoke"
	case MenuInventory:
		text = "inventory"
	case MenuOther:
		text = "menu"
	case MenuInteract:
		text = "interact"
	}
	return "[" + text + "]"
}

func (m menu) Key(g *state) (key action) {
	switch m {
	//case MenuExplore:
	//key = KeyExplore
	case MenuOther:
		key = ActionMenu
	case MenuInventory:
		key = ActionInventory
	case MenuEvoke:
		key = ActionEvoke
	case MenuInteract:
		key = ActionInteract
	}
	return key
}

//func (ui *model) UpdateInteractButton() string {
//g := ui.st
//var interactMenu string
//var show bool
//switch g.Dungeon.Cell(g.Player.Pos).T {
//case StairCell:
//interactMenu = "[descend]"
//if g.Objects.Stairs[g.Player.Pos] == WinStair {
//interactMenu = "[escape]"
//}
//show = true
//case BarrelCell:
//interactMenu = "[rest]"
//show = true
//case MagaraCell:
//interactMenu = "[equip]"
//show = true
//case StoneCell:
//interactMenu = "[activate]"
//show = true
//case ScrollCell:
//interactMenu = "[read]"
//show = true
//case ItemCell:
//interactMenu = "[equip]"
//show = true
//case LightCell:
//interactMenu = "[extinguish]"
//show = true
//case StoryCell:
//if g.Objects.Story[g.Player.Pos] == StoryArtifactSealed || g.Objects.Story[g.Player.Pos] == StoryArtifact {
//interactMenu = "[take]"
//show = true
//}
//}
//if !show {
//return ""
//}
//i := len(MenuCols) - 1
//runes := utf8.RuneCountInString(interactMenu)
//MenuCols[i][1] = MenuCols[i][0] + runes
//return interactMenu
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

func (ui *model) HandleWizardAction() error {
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

func (ui *model) Death() {
	g := ui.st
	if len(g.Stats.Achievements) == 0 {
		NoAchievement.Get(g)
	}
	g.Print("You die... [(x) to continue]")
	//ui.DrawDungeonView(NormalMode)
	//ui.WaitForContinue(-1)
	err := g.WriteDump()
	ui.Dump(err)
	//ui.WaitForContinue(-1)
}

func (ui *model) Win() {
	// TODO: rewrite
	g := ui.st
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
	ui.Dump(err)
	//ui.WaitForContinue(-1)
}

func (ui *model) Dump(err error) {
	//g := ui.st
	//ui.DrawText(g.SimplifedDump(err), 0, 0)
}

func (ui *model) CriticalHPWarning() {
	// TODO
	//g := ui.st
	//g.PrintStyled("*** CRITICAL HP WARNING *** [(x) to continue]", logCritic)
	//ui.DrawDungeonView(NormalMode)
	//ui.WaitForContinue(ui.MapHeight())
	//g.Print("Ok. Be careful, then.")
}

//func (ui *model) Quit() bool {
//g := ui.st
//g.Print("Do you really want to quit without saving? [y/N]")
//ui.DrawDungeonView(NormalMode)
//quit := ui.PromptConfirmation()
//if quit {
//err := g.RemoveSaveFile()
//if err != nil {
//g.PrintfStyled("Error removing save file: %v ——press any key to quit——", logError, err)
//ui.DrawDungeonView(NormalMode)
//ui.PressAnyKey()
//}
//} else {
//g.Print(DoNothing)
//}
//return quit
//}

func (ui *model) Clear() {
	c := gruid.Cell{}
	c.Rune = ' '
	c.Style = gruid.Style{Fg: ColorFg, Bg: ColorBg}
	ui.gd.Fill(c)
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

func Sleep(d time.Duration) {
	// TODO: fix animations
	//time.Sleep(d * time.Millisecond)
}

type posInfo struct {
	Pos         gruid.Point
	Unknown     bool
	Noise       bool
	Unreachable bool
	Sees        bool
	Player      bool
	Monster     *monster
	Cell        string
	Cloud       string
	Lighted     bool
}

func (pi *posInfo) Draw(ui *model) {
	//ui.DrawDungeonView(TargetingMode)
	g := ui.st
	ie := pi // XXX
	y := 3
	//var x int
	//if ui.cursor.X > DungeonWidth/2 {
	//x = 2
	//} else {
	//x = DungeonWidth/2 + 3
	//}
	//width := 35

	// TODO: replace this with ui.Label
	formatBox := func(title, s string, fg gruid.Color) {
		//h := ui.DrawLabel(title, s, x, y, width, fg)
		h := 1 // TODO
		y += h + 2
	}

	features := []string{}
	if !ie.Unknown {
		features = append(features, ie.Cell)
		if ie.Cloud != "" && ie.Sees {
			features = append(features, ie.Cloud)
		}
		if ie.Lighted && ie.Sees {
			features = append(features, "(lighted)")
		}
	} else {
		features = append(features, "unknown place")
	}
	if ie.Noise {
		features = append(features, "noise")
	}
	if ie.Unreachable {
		features = append(features, "unreachable")
	}
	t := " Terrain Features"
	if !ie.Sees && !ie.Unknown {
		t += " (seen)"
	} else if ie.Unknown {
		t += " (unknown)"
	}
	fg := ColorFg
	if ie.Unreachable {
		t += " - Unreachable"
		fg = ColorOrange
	}
	formatBox(t+" ", strings.Join(features, ", "), fg)

	if ie.Player {
		formatBox("Player", "This is you.", ColorBlue)
	}

	mons := ie.Monster
	if mons == nil {
		return
	}
	title := fmt.Sprintf(" %s %s (%s %s) ", mons.Kind, mons.State, mons.Dir.String())
	fg = mons.Color(g)
	var mdesc []string

	statuses := mons.StatusesText()
	if statuses != "" {
		mdesc = append(mdesc, "Statuses: %s", statuses)
	}
	mdesc = append(mdesc, "Traits: "+mons.Traits())
	// TODO
	formatBox(title, strings.Join(mdesc, "\n"), fg)
}

func (m *monster) Color(gs *state) gruid.Color {
	var fg gruid.Color
	if m.Status(MonsLignified) {
		fg = ColorFgLignifiedMonster
	} else if m.Status(MonsConfused) {
		fg = ColorFgConfusedMonster
	} else if m.Status(MonsParalysed) {
		fg = ColorFgParalysedMonster
	} else if m.State == Resting {
		fg = ColorFgSleepingMonster
	} else if m.State == Hunting {
		fg = ColorFgMonster
	} else if m.Peaceful(gs) {
		fg = ColorFgPlayer
	} else {
		fg = ColorFgWanderingMonster
	}
	return fg
}

func (m *monster) StatusesText() string {
	infos := []string{}
	for st, i := range m.Statuses {
		if i > 0 {
			infos = append(infos, fmt.Sprintf("%s %d", monsterStatus(st), m.Statuses[monsterStatus(st)]))
		}
	}
	return strings.Join(infos, ", ")
}

func (m *monster) Traits() string {
	var info string
	info += fmt.Sprintf("Their size is %s.", m.Kind.Size())
	if m.Kind.Peaceful() {
		info += " " + fmt.Sprint("They are peaceful.")
	}
	if m.Kind.CanOpenDoors() {
		info += " " + fmt.Sprint("They can open doors.")
	}
	if m.Kind.CanFly() {
		info += " " + fmt.Sprint("They can fly.")
	}
	if m.Kind.CanSwim() {
		info += " " + fmt.Sprint("They can swim.")
	}
	if m.Kind.ShallowSleep() {
		info += " " + fmt.Sprint("They have very shallow sleep.")
	}
	if m.Kind.ResistsLignification() {
		info += " " + fmt.Sprint("They are unaffected by lignification.")
	}
	if m.Kind.ReflectsTeleport() {
		info += " " + fmt.Sprint("They partially reflect back oric teleport magic.")
	}
	if m.Kind.GoodFlair() {
		info += " " + fmt.Sprint("They have good flair.")
	}
	return info
}

func (pi *posInfo) Update(g *state, pos gruid.Point) {
	*pi = posInfo{}
	pi.Pos = pos
	switch {
	case !g.Dungeon.Cell(pos).Explored:
		pi.Unknown = true
		if g.Noise[pos] || g.NoiseIllusion[pos] {
			pi.Noise = true
		}
		return
		//case !targ.Reachable(g, pos):
		//pi.Unreachable = true
		//return
	}
	mons := g.MonsterAt(pos)
	if pos == g.Player.Pos {
		pi.Player = true
	}
	if g.Player.Sees(pos) {
		pi.Sees = true
	}
	c := g.Dungeon.Cell(pos)
	if t, ok := g.TerrainKnowledge[pos]; ok {
		c.T = t
	}
	if mons.Exists() && g.Player.Sees(pos) {
		pi.Monster = mons
	}
	if cld, ok := g.Clouds[pos]; ok && g.Player.Sees(pos) {
		pi.Cloud = cld.String()
	}
	pi.Cell = c.ShortDesc(g, pos)
	if g.Illuminated[idx(pos)] && c.IsIlluminable() && g.Player.Sees(pos) {
		pi.Lighted = true
	}
	if g.Noise[pos] || g.NoiseIllusion[pos] {
		pi.Noise = true
	}
}
