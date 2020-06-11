package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type UICell struct {
	Fg    uicolor
	Bg    uicolor
	R     rune
	InMap bool
}

type uiInput struct {
	key       string
	mouse     bool
	mouseX    int
	mouseY    int
	button    int
	interrupt bool
}

func (ui *gameui) HideCursor() {
	ui.cursor = InvalidPos
}

func (ui *gameui) SetCursor(pos position) {
	ui.cursor = pos
}

func (ui *gameui) KeyToRuneKeyAction(in uiInput) rune {
	if utf8.RuneCountInString(in.key) != 1 {
		return 0
	}
	return ui.ReadKey(in.key)
}

func (ui *gameui) WaitForContinue(line int) {
loop:
	for {
		in := ui.PollEvent()
		r := ui.KeyToRuneKeyAction(in)
		switch r {
		case '\x1b', ' ', 'x', 'X':
			break loop
		}
		if in.mouse && in.button == -1 {
			continue
		}
		if in.mouse && line >= 0 {
			if in.mouseY > line || in.mouseX > DungeonWidth {
				break loop
			}
		} else if in.mouse {
			break loop
		}
	}
}

func (ui *gameui) PromptConfirmation() bool {
	// TODO: this cannot be done with the mouse
	for {
		in := ui.PollEvent()
		switch in.key {
		case "Y", "y":
			return true
		case "":
		default:
			return false
		}
	}
}

func (ui *gameui) PressAnyKey() error {
	for {
		e := ui.PollEvent()
		if e.interrupt {
			return errors.New("interrupted")
		}
		if e.key != "" || (e.mouse && e.button != -1) {
			return nil
		}
	}
}

type startAction int

const (
	StartPlay startAction = iota
	StartWatchReplay
)

func (ui *gameui) StartMenu(l int) startAction {
	for {
		in := ui.PollEvent()
		switch in.key {
		case "P", "p":
			ui.ColorLine(l, ColorYellow)
			ui.Flush()
			Sleep(AnimDurShort)
			return StartPlay
		case "W", "w":
			ui.ColorLine(l+1, ColorYellow)
			ui.Flush()
			Sleep(AnimDurShort)
			return StartWatchReplay
		}
		if in.key != "" && !in.mouse {
			continue
		}
		y := in.mouseY
		switch in.button {
		case -1:
			oih := ui.itemHover
			if y < l || y >= l+2 {
				ui.itemHover = -1
				if oih != -1 {
					ui.ColorLine(oih, ColorFg)
					ui.Flush()
				}
				break
			}
			if y == oih {
				break
			}
			ui.itemHover = y
			ui.ColorLine(y, ColorYellow)
			if oih != -1 {
				ui.ColorLine(oih, ColorFg)
			}
			ui.Flush()
		case 0:
			if y < l || y >= l+2 {
				ui.itemHover = -1
				break
			}
			ui.itemHover = -1
			switch y - l {
			case 0:
				return StartPlay
			case 1:
				return StartWatchReplay
			}
		}
	}
}

func (ui *gameui) PlayerTurnEvent() (again, quit bool, err error) {
	g := ui.g
	again = true
	in := ui.PollEvent()
	switch in.key {
	case "":
		if in.mouse {
			var mpos position
			if CenteredCamera {
				mpos = ui.CameraTargetPosition(in.mouseX, in.mouseY, true)
			} else {
				mpos = position{in.mouseX, in.mouseY}
			}
			switch in.button {
			case -1:
				if in.mouseY == ui.MapHeight() {
					m, ok := ui.WhichButton(in.mouseX)
					omh := ui.menuHover
					if ok {
						ui.menuHover = m
					} else {
						ui.menuHover = -1
					}
					if ui.menuHover != omh {
						ui.DrawMenus()
						ui.Flush()
					}
					break
				}
				ui.menuHover = -1
				if in.mouseX >= ui.MapWidth() || in.mouseY >= ui.MapHeight() {
					again = true
					break
				}
				fallthrough
			case 0:
				if in.mouseY == ui.MapHeight() {
					m, ok := ui.WhichButton(in.mouseX)
					if !ok {
						again = true
						break
					}
					again, quit, err = ui.HandleKeyAction(runeKeyAction{k: m.Key(g)})
					if err != nil {
						again = true
					}
					return again, quit, err
				} else if in.mouseX >= ui.MapWidth() || in.mouseY >= ui.MapHeight() {
					again = true
				} else {
					again, quit, err = ui.ExaminePos(mpos)
				}
			case 2:
				again, quit, err = ui.HandleKeyAction(runeKeyAction{k: ActionMenu})
				if err != nil {
					again = true
				}
				return again, quit, err
			}
		}
	default:
		r := ui.KeyToRuneKeyAction(in)
		if r == 0 {
			again = true
		} else {
			again, quit, err = ui.HandleKeyAction(runeKeyAction{r: r})
		}
	}
	if err != nil {
		again = true
	}
	return again, quit, err
}

func (ui *gameui) Scroll(n int) (m int, quit bool) {
	in := ui.PollEvent()
	switch in.key {
	case "Escape", "\x1b", " ", "x", "X":
		quit = true
	case "u", "9", "b":
		n -= 12
	case "d", "3", "f":
		n += 12
	case "j", "2", ".":
		n++
	case "k", "8":
		n--
	case "":
		if in.mouse {
			switch in.button {
			case 0:
				y := in.mouseY
				x := in.mouseX
				if x >= DungeonWidth {
					quit = true
					break
				}
				if y > UIHeight {
					break
				}
				n += y - (ui.MapHeight()+3)/2
			}
		}
	}
	return n, quit
}

func (ui *gameui) GetIndex(x, y int) int {
	return y*UIWidth + x
}

func (ui *gameui) GetPos(i int) (int, int) {
	return i - (i/UIWidth)*UIWidth, i / UIWidth
}

func (ui *gameui) Select(l int) (index int, alternate bool, err error) {
	if ui.itemHover >= 1 && ui.itemHover <= l {
		ui.ColorLine(ui.itemHover, ColorYellow)
		ui.Flush()
	} else {
		ui.itemHover = -1
	}
	for {
		in := ui.PollEvent()
		r := ui.ReadKey(in.key)
		switch {
		case in.key == "\x1b" || in.key == "Escape" || in.key == " " || in.key == "x" || in.key == "X":
			return -1, false, errors.New(DoNothing)
		case in.key == "?":
			return -1, true, nil
		case 97 <= r && int(r) < 97+l:
			if ui.itemHover >= 1 && ui.itemHover <= l {
				ui.ColorLine(ui.itemHover, ColorFg)
			}
			ui.itemHover = int(r-97) + 1
			return int(r - 97), false, nil
		case in.key == "2":
			oih := ui.itemHover
			ui.itemHover++
			if ui.itemHover < 1 || ui.itemHover > l {
				ui.itemHover = 1
			}
			if oih > 0 && oih <= l {
				ui.ColorLine(oih, ColorFg)
			}
			ui.ColorLine(ui.itemHover, ColorYellow)
			ui.Flush()
		case in.key == "8":
			oih := ui.itemHover
			ui.itemHover--
			if ui.itemHover < 1 {
				ui.itemHover = l
			}
			if oih > 0 && oih <= l {
				ui.ColorLine(oih, ColorFg)
			}
			ui.ColorLine(ui.itemHover, ColorYellow)
			ui.Flush()
		case in.key == "." && ui.itemHover >= 1 && ui.itemHover <= l:
			if ui.itemHover >= 1 && ui.itemHover <= l {
				ui.ColorLine(ui.itemHover, ColorFg)
			}
			return ui.itemHover - 1, false, nil
		case in.key == "" && in.mouse:
			y := in.mouseY
			x := in.mouseX
			switch in.button {
			case -1:
				oih := ui.itemHover
				if y <= 0 || y > l || x >= DungeonWidth {
					ui.itemHover = -1
					if oih > 0 {
						ui.ColorLine(oih, ColorFg)
						ui.Flush()
					}
					break
				}
				if y == oih {
					break
				}
				ui.itemHover = y
				ui.ColorLine(y, ColorYellow)
				if oih > 0 {
					ui.ColorLine(oih, ColorFg)
				}
				ui.Flush()
			case 0:
				if y < 0 || y > l || x >= DungeonWidth {
					ui.itemHover = -1
					return -1, false, errors.New(DoNothing)
				}
				if y == 0 {
					ui.itemHover = -1
					return -1, true, nil
				}
				ui.itemHover = -1
				return y - 1, false, nil
			case 2:
				ui.itemHover = -1
				return -1, true, nil
			case 1:
				ui.itemHover = -1
				return -1, false, errors.New(DoNothing)
			}
		}
	}
}

func (ui *gameui) KeyMenuAction(n int) (m int, action keyConfigAction) {
	in := ui.PollEvent()
	r := ui.KeyToRuneKeyAction(in)
	switch string(r) {
	case "a":
		action = ChangeKeys
	case "\x1b", " ", "x", "X":
		action = QuitKeyConfig
	case "u":
		n -= ui.MapHeight() / 2
	case "d":
		n += ui.MapHeight() / 2
	case "j", "2", "ArrowDown":
		n++
	case "k", "8", "ArrowUp":
		n--
	case "R":
		action = ResetKeys
	default:
		if r == 0 && in.mouse {
			y := in.mouseY
			x := in.mouseX
			switch in.button {
			case 0:
				if x > DungeonWidth || y > ui.MapHeight() {
					action = QuitKeyConfig
				}
			case 1:
				action = QuitKeyConfig
			}
		}
	}
	return n, action
}

func (ui *gameui) TargetModeEvent(targ Targeter, data *examineData) (again, quit, notarg bool, err error) {
	g := ui.g
	again = true
	in := ui.PollEvent()
	switch in.key {
	case "\x1b", "Escape", " ", "x", "X":
		g.Targeting = InvalidPos
		notarg = true
	case "":
		if !in.mouse {
			return
		}
		switch in.button {
		case -1:
			if in.mouseY == ui.MapHeight() {
				m, ok := ui.WhichButton(in.mouseX)
				omh := ui.menuHover
				if ok {
					ui.menuHover = m
				} else {
					ui.menuHover = -1
				}
				if ui.menuHover != omh {
					ui.DrawMenus()
					ui.Flush()
				}
				g.Targeting = InvalidPos
				notarg = true
				err = errors.New(DoNothing)
				break
			}
			ui.menuHover = -1
			if in.mouseY >= ui.MapHeight() || in.mouseX >= ui.MapWidth() {
				g.Targeting = InvalidPos
				notarg = true
				err = errors.New(DoNothing)
				break
			}
			var mpos position
			if CenteredCamera {
				mpos = ui.CameraTargetPosition(in.mouseX, in.mouseY, true)
			} else {
				mpos = position{in.mouseX, in.mouseY}
			}
			if g.Targeting == mpos {
				break
			}
			g.Targeting = InvalidPos
			fallthrough
		case 0:
			if in.mouseY == ui.MapHeight() {
				m, ok := ui.WhichButton(in.mouseX)
				if !ok {
					g.Targeting = InvalidPos
					notarg = true
					err = errors.New(DoNothing)
					break
				}
				again, quit, notarg, err = ui.CursorKeyAction(targ, runeKeyAction{k: m.Key(g)}, data)
			} else if in.mouseX >= ui.MapWidth() || in.mouseY >= ui.MapHeight() {
				g.Targeting = InvalidPos
				notarg = true
				err = errors.New(DoNothing)
			} else {
				var mpos position
				if CenteredCamera {
					mpos = ui.CameraTargetPosition(in.mouseX, in.mouseY, true)
				} else {
					mpos = position{in.mouseX, in.mouseY}
				}
				again, notarg = ui.CursorMouseLeft(targ, mpos, data)
			}
		case 2:
			if in.mouseY >= ui.MapHeight() || in.mouseX >= ui.MapWidth() {
				again, quit, notarg, err = ui.CursorKeyAction(targ, runeKeyAction{k: ActionMenu}, data)
			} else {
				again, quit, notarg, err = ui.CursorKeyAction(targ, runeKeyAction{k: ActionDescription}, data)
			}
		case 1:
			again, quit, notarg, err = ui.CursorKeyAction(targ, runeKeyAction{k: ActionExclude}, data)
		}
	default:
		r := ui.KeyToRuneKeyAction(in)
		if r != 0 {
			return ui.CursorKeyAction(targ, runeKeyAction{r: r}, data)
		}
		again = true
	}
	return
}

func (ui *gameui) ReadRuneKey() rune {
	for {
		in := ui.PollEvent()
		switch in.key {
		case "\x1b", "Escape", " ", "x", "X":
			return 0
		case "Enter":
			return '.'
		}
		r := ui.ReadKey(in.key)
		if unicode.IsPrint(r) {
			return r
		}
	}
}

func (ui *gameui) ReadKey(s string) (r rune) {
	bs := strings.NewReader(s)
	r, _, _ = bs.ReadRune()
	return r
}

type uiMode int

const (
	NormalMode uiMode = iota
	TargetingMode
	NoFlushMode
	AnimationMode
)

const DoNothing = "Do nothing, then."

func (ui *gameui) EnterWizard() {
	g := ui.g
	if ui.Wizard() {
		g.EnterWizardMode()
		ui.DrawDungeonView(NoFlushMode)
	} else {
		g.Print(DoNothing)
	}
}

func (ui *gameui) CleanError(err error) error {
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

func FixedRuneKey(r rune) bool {
	switch r {
	case ' ', '?', '=', '2', '4', '8', '6', '.', '5', '\x1b', 'x', 'X':
		return true
	default:
		return false
	}
}

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
		text = "Write game statistics to file"
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

type runeKeyAction struct {
	r rune
	k action
}

func (ui *gameui) HandleKeyAction(rka runeKeyAction) (again bool, quit bool, err error) {
	if rka.r != 0 {
		var ok bool
		rka.k, ok = GameConfig.RuneNormalModeKeys[rka.r]
		if !ok {
			switch rka.r {
			case 's':
				err = errors.New("Unknown key. Did you mean capital S for save and quit?")
			case 'q':
				err = errors.New("Unknown key. Did you mean capital Q for quit without saving?")
			default:
				err = fmt.Errorf("Unknown key '%c'. Type ? for help.", rka.r)
			}
			return again, quit, err
		}
	}
	if rka.k == ActionMenu {
		rka.k, err = ui.SelectAction(menuActions)
		if err != nil {
			err = ui.CleanError(err)
			return again, quit, err
		}
	}
	return ui.HandleKey(rka)
}

func (ui *gameui) OptionalDescendConfirmation(st stair) (err error) {
	g := ui.g
	if g.Depth == WinDepth && st == NormalStair && g.Dungeon.Cell(g.Places.Shaedra).T == StoryCell {
		err = errors.New("You have to rescue Shaedra first!")
	}
	return err

}

func (ui *gameui) HandleKey(rka runeKeyAction) (again bool, quit bool, err error) {
	g := ui.g
	switch rka.k {
	case ActionW, ActionS, ActionN, ActionE:
		err = g.PlayerBump(g.Player.Pos.To(KeyToDir(rka.k)))
	case ActionRunW, ActionRunS, ActionRunN, ActionRunE:
		err = g.GoToDir(KeyToDir(rka.k))
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
			ex := &examiner{stairs: true}
			err = ex.Action(g, stair)
			if err == nil && !g.MoveToTarget() {
				err = errors.New("You could not move toward stairs.")
			}
			if ex.Done() {
				g.Targeting = InvalidPos
			}
		} else {
			err = errors.New("You cannot go to any stairs.")
		}
	case ActionInteract:
		c := g.Dungeon.Cell(g.Player.Pos)
		switch c.T {
		case StairCell:
			if g.Dungeon.Cell(g.Player.Pos).T == StairCell && g.Objects.Stairs[g.Player.Pos] != BlockedStair {
				ui.MenuSelectedAnimation(MenuInteract, true)
				strt := g.Objects.Stairs[g.Player.Pos]
				err = ui.OptionalDescendConfirmation(strt)
				if err != nil {
					break
				}
				if g.Descend(DescendNormal) {
					ui.Win()
					quit = true
					return again, quit, err
				}
				ui.DrawDungeonView(NormalMode)
			} else if g.Dungeon.Cell(g.Player.Pos).T == StairCell && g.Objects.Stairs[g.Player.Pos] == BlockedStair {
				err = errors.New("The stairs are blocked by a magical stone barrier energies.")
			} else {
				err = errors.New("No stairs here.")
			}
		case EssenciaticSourceCell:
			ui.MenuSelectedAnimation(MenuInteract, true)
			err = g.Rest()
			if err != nil {
				ui.MenuSelectedAnimation(MenuInteract, false)
			}
		case StoneCell:
			ui.MenuSelectedAnimation(MenuInteract, true)
			err = g.ActivateStone()
			if err != nil {
				ui.MenuSelectedAnimation(MenuInteract, false)
			}
		case ScrollCell:
			err = ui.ReadScroll()
			err = ui.CleanError(err)
		case ItemCell:
			err = ui.g.EquipItem()
		case LightCell:
			err = g.ExtinguishFire()
		case StoryCell:
			//if g.Objects.Story[g.Player.Pos] == StoryArtifact && !g.LiberatedArtifact {
			//g.PushEvent(&simpleEvent{ERank: g.Ev.Rank(), EAction: ArtifactAnimation})
			//g.LiberatedArtifact = true
			//g.Ev.Renew(g, DurationTurn)
			//} else if g.Objects.Story[g.Player.Pos] == StoryArtifactSealed {
			//err = errors.New("The artifact is protected by a magical stone barrier.")
			//} else {
			//err = errors.New("You cannot interact with anything here.")
			//}
		default:
			err = errors.New("You cannot interact with anything here.")
		}
	//case ActionEvoke:
	//err = ui.SelectMagara()
	//err = ui.CleanError(err)
	case ActionInventory:
		err = ui.SelectItem()
		err = ui.CleanError(err)
	case ActionExplore:
		err = g.Autoexplore()
	case ActionExamine:
		again, quit, err = ui.Examine(nil)
	case ActionHelp, ActionMenuCommandHelp:
		ui.KeysHelp()
		again = true
	case ActionMenuTargetingHelp:
		ui.ExamineHelp()
		again = true
	case ActionLogs:
		ui.DrawPreviousLogs()
		again = true
	case ActionSave:
		g.RenewEvent(0)
		errsave := g.Save()
		if errsave != nil {
			g.PrintfStyled("Error: %v", logError, errsave)
			g.PrintStyled("Could not save game.", logError)
		} else {
			quit = true
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
		//if g.Wizard && g.Depth == WinDepth {
		//g.RescuedShaedra()
		//}
		if g.Wizard && g.Depth < MaxDepth {
			g.StoryPrint("Descended wizardly")
			if g.Descend(DescendNormal) {
				ui.Win()
				quit = true
				return again, quit, err
			}
		} else {
			err = errors.New("Unknown key. Type ? for help.")
		}
	case ActionWizard:
		ui.EnterWizard()
		return true, false, nil
	case ActionQuit:
		if ui.Quit() {
			return false, true, nil
		}
		return true, false, nil
	case ActionConfigure:
		err = ui.HandleSettingAction()
		again = true
	case ActionDescription:
		//ui.MenuSelectedAnimation(MenuView, false)
		err = fmt.Errorf("You must choose a target to describe.")
	case ActionExclude:
		err = fmt.Errorf("You must choose a target for exclusion.")
	default:
		err = fmt.Errorf("Unknown key '%c'. Type ? for help.", rka.r)
	}
	if err != nil {
		again = true
	}
	return again, quit, err
}

func (ui *gameui) ExaminePos(pos position) (again, quit bool, err error) {
	var start *position
	if pos.valid() {
		start = &pos
	}
	again, quit, err = ui.Examine(start)
	return again, quit, err
}

func (ui *gameui) Examine(start *position) (again, quit bool, err error) {
	ex := &examiner{}
	again, quit, err = ui.CursorAction(ex, start)
	return again, quit, err
}

func (ui *gameui) ChooseTarget(targ Targeter) error {
	_, _, err := ui.CursorAction(targ, nil)
	if err != nil {
		return err
	}
	if !targ.Done() {
		return errors.New(DoNothing)
	}
	return nil
}

func (ui *gameui) NextMonster(r rune, pos position, data *examineData) {
	g := ui.g
	nmonster := data.nmonster
	for i := 0; i < len(g.Monsters); i++ {
		if r == '+' {
			nmonster++
		} else {
			nmonster--
		}
		if nmonster > len(g.Monsters)-1 {
			nmonster = 0
		} else if nmonster < 0 {
			nmonster = len(g.Monsters) - 1
		}
		mons := g.Monsters[nmonster]
		if mons.Exists() && g.Player.LOS[mons.Pos] && pos != mons.Pos {
			pos = mons.Pos
			break
		}
	}
	data.npos = pos
	data.nmonster = nmonster
}

func (ui *gameui) NextStair(data *examineData) {
	g := ui.g
	if data.sortedStairs == nil {
		stairs := g.StairsSlice()
		data.sortedStairs = g.SortedNearestTo(stairs, g.Player.Pos)
	}
	if data.stairIndex >= len(data.sortedStairs) {
		data.stairIndex = 0
	}
	if len(data.sortedStairs) > 0 {
		data.npos = data.sortedStairs[data.stairIndex]
		data.stairIndex++
	}
}

func (ui *gameui) NextObject(pos position, data *examineData) {
	g := ui.g
	nobject := data.nobject
	if len(data.objects) == 0 {
		for p := range g.Objects.Stairs {
			data.objects = append(data.objects, p)
		}
		for p := range g.Objects.FakeStairs {
			data.objects = append(data.objects, p)
		}
		for p := range g.Objects.Stones {
			data.objects = append(data.objects, p)
		}
		for p := range g.Objects.EssenciaticSources {
			data.objects = append(data.objects, p)
		}
		for p := range g.Objects.Bananas {
			data.objects = append(data.objects, p)
		}
		for p := range g.Objects.Items {
			data.objects = append(data.objects, p)
		}
		for p := range g.Objects.Scrolls {
			data.objects = append(data.objects, p)
		}
		for p := range g.Objects.Potions {
			data.objects = append(data.objects, p)
		}
		data.objects = g.SortedNearestTo(data.objects, g.Player.Pos)
	}
	for i := 0; i < len(data.objects); i++ {
		p := data.objects[nobject]
		nobject++
		if nobject > len(data.objects)-1 {
			nobject = 0
		}
		if g.Dungeon.Cell(p).Explored {
			pos = p
			break
		}
	}
	data.npos = pos
	data.nobject = nobject
}

func (ui *gameui) ExcludeZone(pos position) {
	g := ui.g
	if !g.Dungeon.Cell(pos).Explored {
		g.Print("You cannot choose an unexplored cell for exclusion.")
	} else {
		toggle := !g.ExclusionsMap[pos]
		g.ComputeExclusion(pos, toggle)
	}
}

func (ui *gameui) CursorMouseLeft(targ Targeter, pos position, data *examineData) (again, notarg bool) {
	g := ui.g
	again = true
	if data.npos == pos {
		err := targ.Action(g, pos)
		if err != nil {
			g.Print(err.Error())
		} else {
			if g.MoveToTarget() {
				again = false
			}
			if targ.Done() {
				notarg = true
			}
		}
	} else {
		data.npos = pos
	}
	return again, notarg
}

func (ui *gameui) CursorKeyAction(targ Targeter, rka runeKeyAction, data *examineData) (again, quit, notarg bool, err error) {
	g := ui.g
	pos := data.npos
	again = true
	if rka.r != 0 {
		var ok bool
		rka.k, ok = GameConfig.RuneTargetModeKeys[rka.r]
		if !ok {
			err = fmt.Errorf("Invalid targeting mode key '%c'. Type ? for help.", rka.r)
			return again, quit, notarg, err
		}
	}
	if rka.k == ActionMenu {
		rka.k, err = ui.SelectAction(menuActions)
		if err != nil {
			err = ui.CleanError(err)
			return again, quit, notarg, err
		}
	}
	switch rka.k {
	case ActionW, ActionS, ActionN, ActionE:
		data.npos = pos.To(KeyToDir(rka.k))
	case ActionRunW, ActionRunS, ActionRunN, ActionRunE:
		for i := 0; i < 5; i++ {
			p := data.npos.To(KeyToDir(rka.k))
			if !p.valid() {
				break
			}
			data.npos = p
		}
	case ActionNextStairs:
		ui.NextStair(data)
	case ActionDescend:
		if g.Dungeon.Cell(g.Player.Pos).T == StairCell && g.Objects.Stairs[g.Player.Pos] != BlockedStair {
			ui.MenuSelectedAnimation(MenuInteract, true)
			strt := g.Objects.Stairs[g.Player.Pos]
			err = ui.OptionalDescendConfirmation(strt)
			if err != nil {
				break
			}
			again = false
			g.Targeting = InvalidPos
			notarg = true
			if g.Descend(DescendNormal) {
				ui.Win()
				quit = true
				return again, quit, notarg, err
			}
		} else if g.Dungeon.Cell(g.Player.Pos).T == StairCell && g.Objects.Stairs[g.Player.Pos] == BlockedStair {
			err = errors.New("The stairs are blocked by a magical stone barrier energies.")
		} else {
			err = errors.New("No stairs here.")
		}
	case ActionPreviousMonster, ActionNextMonster:
		ui.NextMonster(rka.r, pos, data)
	case ActionNextObject:
		ui.NextObject(pos, data)
	case ActionHelp, ActionMenuTargetingHelp:
		ui.HideCursor()
		ui.ExamineHelp()
		ui.SetCursor(pos)
	case ActionMenuCommandHelp:
		ui.HideCursor()
		ui.KeysHelp()
		ui.SetCursor(pos)
	case ActionTarget:
		err = targ.Action(g, pos)
		if err != nil {
			break
		}
		g.Targeting = InvalidPos
		if g.MoveToTarget() {
			again = false
		}
		if targ.Done() {
			notarg = true
		}
	case ActionDescription:
		ui.HideCursor()
		ui.ViewPositionDescription(pos)
		ui.SetCursor(pos)
	case ActionExclude:
		ui.ExcludeZone(pos)
	case ActionEscape:
		g.Targeting = InvalidPos
		notarg = true
		err = errors.New(DoNothing)
	case ActionExplore, ActionLogs, ActionEvoke, ActionInventory:
		// XXX: hm, this is only useful with mouse in terminal, rarely tested.
		if _, ok := targ.(*examiner); !ok {
			break
		}
		again, quit, err = ui.HandleKey(rka)
		if err != nil {
			notarg = true
		}
		g.Targeting = InvalidPos
	case ActionConfigure:
		err = ui.HandleSettingAction()
	case ActionSave:
		g.RenewEvent(0)
		g.Highlight = nil
		g.Targeting = InvalidPos
		errsave := g.Save()
		if errsave != nil {
			g.PrintfStyled("Error: %v", logError, errsave)
			g.PrintStyled("Could not save game.", logError)
		} else {
			notarg = true
			again = false
			quit = true
		}
	case ActionQuit:
		if ui.Quit() {
			quit = true
			again = false
		}
	default:
		err = fmt.Errorf("Invalid targeting mode key '%c'. Type ? for help.", rka.r)
	}
	return again, quit, notarg, err
}

type examineData struct {
	npos         position
	nmonster     int
	objects      []position
	nobject      int
	sortedStairs []position
	stairIndex   int
}

var InvalidPos = position{-1, -1}

func (ui *gameui) CursorAction(targ Targeter, start *position) (again, quit bool, err error) {
	g := ui.g
	pos := g.Player.Pos
	if start != nil {
		pos = *start
	} else {
		minDist := 999
		for _, mons := range g.Monsters {
			if mons.Exists() && g.Player.LOS[mons.Pos] {
				dist := mons.Pos.Distance(g.Player.Pos)
				if minDist > dist {
					minDist = dist
					pos = mons.Pos
				}
			}
		}
	}
	data := &examineData{
		npos:    pos,
		objects: []position{},
	}
	if _, ok := targ.(*examiner); ok && pos == g.Player.Pos && start == nil {
		ui.NextObject(InvalidPos, data)
		if !data.npos.valid() {
			ui.NextStair(data)
		}
		if data.npos.valid() && pos.Distance(data.npos) < DefaultLOSRange+5 {
			pos = data.npos
		}
	}
	opos := InvalidPos
loop:
	for {
		err = nil
		if pos != opos {
			ui.DescribePosition(pos, targ)
		}
		opos = pos
		targ.ComputeHighlight(g, pos)
		ui.SetCursor(pos)
		m := g.MonsterAt(pos)
		if m.Exists() && g.Player.Sees(pos) {
			g.ComputeMonsterCone(m)
		} else {
			g.MonsterTargLOS = nil
		}
		ui.DrawDungeonView(TargetingMode)
		ui.DrawInfoLine(g.InfoEntry)
		if !ui.Small() {
			st := " Examine/Travel mode "
			if _, ok := targ.(*examiner); !ok {
				st = " Targeting mode "
			}
			ui.DrawStyledTextLine(st, ui.MapHeight()+2, FooterLine)
		}
		ui.SetCell(DungeonWidth, ui.MapHeight(), '┤', ColorFg, ColorBg)
		ui.Flush()
		data.npos = pos
		var notarg bool
		again, quit, notarg, err = ui.TargetModeEvent(targ, data)
		if err != nil {
			err = ui.CleanError(err)
		}
		if !again || notarg {
			break loop
		}
		if err != nil {
			g.Print(err.Error())
		}
		if data.npos.valid() {
			pos = data.npos
		}
	}
	g.Highlight = nil
	g.MonsterTargLOS = nil
	ui.HideCursor()
	return again, quit, err
}

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

func (m menu) Key(g *game) (key action) {
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

var MenuCols = [][2]int{
	//MenuExplore:  {0, 0},
	MenuOther:     {0, 0},
	MenuInventory: {0, 0},
	MenuEvoke:     {0, 0},
	MenuInteract:  {0, 0}}

func init() {
	for i := range MenuCols {
		runes := utf8.RuneCountInString(menu(i).String())
		if i == 0 {
			MenuCols[0] = [2]int{25, 25 + runes}
			continue
		}
		MenuCols[i] = [2]int{MenuCols[i-1][1] + 2, MenuCols[i-1][1] + 2 + runes}
	}
}

func (ui *gameui) WhichButton(col int) (menu, bool) {
	g := ui.g
	if ui.Small() {
		return MenuOther, false
	}
	end := len(MenuCols) - 1
	switch g.Dungeon.Cell(g.Player.Pos).T {
	case StairCell, EssenciaticSourceCell, ScrollCell, StoneCell, LightCell:
		end++
	case StoryCell:
		if g.Objects.Story[g.Player.Pos] == StoryArtifactSealed || g.Objects.Story[g.Player.Pos] == StoryArtifact {
			end++
		}
	}
	for i, cols := range MenuCols[0:end] {
		if cols[0] >= 0 && col >= cols[0] && col < cols[1] {
			return menu(i), true
		}
	}
	return MenuOther, false
}

func (ui *gameui) UpdateInteractButton() string {
	g := ui.g
	var interactMenu string
	var show bool
	switch g.Dungeon.Cell(g.Player.Pos).T {
	case StairCell:
		interactMenu = "[descend]"
		if g.Objects.Stairs[g.Player.Pos] == WinStair {
			interactMenu = "[escape]"
		}
		show = true
	case EssenciaticSourceCell:
		interactMenu = "[rest]"
		show = true
	case StoneCell:
		interactMenu = "[activate]"
		show = true
	case ScrollCell:
		interactMenu = "[read]"
		show = true
	case ItemCell:
		interactMenu = "[equip]"
		show = true
	case LightCell:
		interactMenu = "[extinguish]"
		show = true
	case StoryCell:
		if g.Objects.Story[g.Player.Pos] == StoryArtifactSealed || g.Objects.Story[g.Player.Pos] == StoryArtifact {
			interactMenu = "[take]"
			show = true
		}
	}
	if !show {
		return ""
	}
	i := len(MenuCols) - 1
	runes := utf8.RuneCountInString(interactMenu)
	MenuCols[i][1] = MenuCols[i][0] + runes
	return interactMenu
}

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

func (ui *gameui) HandleWizardAction() error {
	g := ui.g
	s, err := ui.SelectWizardMagic(WizardActions)
	if err != nil {
		return err
	}
	switch s {
	case WizardInfoAction:
		ui.WizardInfo()
	case WizardToggleMode:
		switch g.WizardMode {
		case WizardNormal:
			g.WizardMode = WizardMap
		case WizardMap:
			g.WizardMode = WizardSeeAll
		case WizardSeeAll:
			g.WizardMode = WizardNormal
		}
		g.StoryPrint("Toggle wizard mode.")
		ui.DrawDungeonView(NoFlushMode)
	}
	return nil
}

func (ui *gameui) Death() {
	g := ui.g
	if len(g.Stats.Achievements) == 0 {
		NoAchievement.Get(g)
	}
	g.Print("You die... [(x) to continue]")
	ui.DrawDungeonView(NormalMode)
	ui.WaitForContinue(-1)
	err := g.WriteDump()
	ui.Dump(err)
	ui.WaitForContinue(-1)
}

func (ui *gameui) Win() {
	g := ui.g
	err := g.RemoveSaveFile()
	if err != nil {
		g.PrintfStyled("Error removing save file: %v", logError, err)
	}
	if g.Wizard {
		g.Print("You escape by the magic portal! **WIZARD** [(x) to continue]")
	} else {
		g.Print("You escape by the magic portal! [(x) to continue]")
	}
	ui.DrawDungeonView(NormalMode)
	ui.WaitForContinue(-1)
	err = g.WriteDump()
	ui.Dump(err)
	ui.WaitForContinue(-1)
}

func (ui *gameui) Dump(err error) {
	g := ui.g
	ui.Clear()
	ui.DrawText(g.SimplifedDump(err), 0, 0)
	ui.Flush()
}

func (ui *gameui) CriticalHPWarning() {
	g := ui.g
	g.PrintStyled("*** CRITICAL HP WARNING *** [(x) to continue]", logCritic)
	ui.DrawDungeonView(NormalMode)
	ui.WaitForContinue(ui.MapHeight())
	g.Print("Ok. Be careful, then.")
}

func (ui *gameui) Quit() bool {
	g := ui.g
	g.Print("Do you really want to quit without saving? [y/N]")
	ui.DrawDungeonView(NormalMode)
	quit := ui.PromptConfirmation()
	if quit {
		err := g.RemoveSaveFile()
		if err != nil {
			g.PrintfStyled("Error removing save file: %v ——press any key to quit——", logError, err)
			ui.DrawDungeonView(NormalMode)
			ui.PressAnyKey()
		}
	} else {
		g.Print(DoNothing)
	}
	return quit
}

func (ui *gameui) Wizard() bool {
	g := ui.g
	g.Print("Do you really want to enter wizard mode (no return)? [y/N]")
	ui.DrawDungeonView(NormalMode)
	return ui.PromptConfirmation()
}

func (ui *gameui) HandlePlayerTurn() bool {
	g := ui.g
getKey:
	for {
		var err error
		var again, quit bool
		if g.Targeting.valid() {
			again, quit, err = ui.ExaminePos(g.Targeting)
		} else {
			ui.DrawDungeonView(NormalMode)
			again, quit, err = ui.PlayerTurnEvent()
		}
		if err != nil && err.Error() != "" {
			g.Print(err.Error())
		}
		if again {
			continue getKey
		}
		return quit
	}
}

func (ui *gameui) ExploreStep() bool {
	next := make(chan bool)
	var stop bool
	go func() {
		Sleep(10)
		ui.Interrupt()
	}()
	go func() {
		err := ui.PressAnyKey()
		interrupted := err != nil
		next <- !interrupted
	}()
	stop = <-next
	ui.DrawDungeonView(NormalMode)
	return stop
}

func (ui *gameui) Clear() {
	for i := 0; i < UIHeight*UIWidth; i++ {
		x, y := ui.GetPos(i)
		ui.SetCell(x, y, ' ', ColorFg, ColorBg)
	}
}

func (ui *gameui) DrawBufferInit() {
	if len(ui.g.DrawBuffer) == 0 {
		ui.g.DrawBuffer = make([]UICell, UIHeight*UIWidth)
	} else if len(ui.g.DrawBuffer) != UIHeight*UIWidth {
		ui.g.DrawBuffer = make([]UICell, UIHeight*UIWidth)
	}
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

func (ui *gameui) ColorLine(y int, fg uicolor) {
	for x := 0; x < DungeonWidth; x++ {
		i := ui.GetIndex(x, y)
		c := ui.g.DrawBuffer[i]
		ui.SetCell(x, y, c.R, fg, c.Bg)
	}
}

func Sleep(d time.Duration) {
	time.Sleep(d * time.Millisecond)
}
