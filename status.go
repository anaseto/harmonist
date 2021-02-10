package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/ui"
)

type status int

const (
	StatusExhausted status = iota
	StatusSwift
	StatusLignification
	StatusConfusion
	StatusFlames // fake status
	StatusHidden
	StatusUnhidden
	StatusLight
	StatusDig
	StatusLevitation
	StatusShadows
	StatusIlluminated
	StatusTransparent
	StatusDisguised
	StatusDelay
	StatusDispersal
)

func (st status) Flag() bool {
	switch st {
	case StatusFlames, StatusHidden, StatusUnhidden, StatusLight:
		return true
	default:
		return false
	}
}

func (st status) Clean() bool {
	switch st {
	case StatusDelay:
		return true
	default:
		return false
	}
}

func (st status) Info() bool {
	switch st {
	case StatusFlames, StatusHidden, StatusUnhidden, StatusLight, StatusDelay:
		return true
	}
	return false
}

func (st status) Good() bool {
	switch st {
	case StatusSwift, StatusDig, StatusHidden, StatusLevitation, StatusShadows, StatusTransparent, StatusDisguised, StatusDispersal:
		return true
	default:
		return false
	}
}

func (st status) Bad() bool {
	switch st {
	case StatusConfusion, StatusUnhidden, StatusIlluminated:
		return true
	default:
		return false
	}
}

func (st status) String() string {
	switch st {
	case StatusExhausted:
		return "Exhaustion"
	case StatusSwift:
		return "Swiftness"
	case StatusLignification:
		return "Lignification"
	case StatusConfusion:
		return "Confusion"
	case StatusFlames:
		return "Flames"
	case StatusHidden:
		return "Hidden"
	case StatusUnhidden:
		return "Unhidden"
	case StatusDig:
		return "Dig"
	case StatusLight:
		return "Light"
	case StatusLevitation:
		return "Levitation"
	case StatusShadows:
		return "Shadows"
	case StatusIlluminated:
		return "Illumination"
	case StatusTransparent:
		return "Transparency"
	case StatusDisguised:
		return "Disguise"
	case StatusDelay:
		return "Delay"
	case StatusDispersal:
		return "Dispersal"
	default:
		// should not happen
		return "unknown"
	}
}

func (st status) Desc() string {
	switch st {
	case StatusExhausted:
		return "Forbids some actions, such as jumping."
	case StatusSwift:
		return "Allows for moving several times in a row."
	case StatusLignification:
		return "Makes it impossible to move."
	case StatusConfusion:
		return "Makes it impossible to use magaras."
	case StatusFlames:
		return "Surrounded by magical flames."
	case StatusHidden:
		return "No monster can see you."
	case StatusUnhidden:
		return "Some monsters can see you."
	case StatusDig:
		return "Allows to walk into walls."
	case StatusLight:
		return "You are in a lighted cell."
	case StatusLevitation:
		return "Allows to fly over chasm and oric barriers."
	case StatusShadows:
		return "Only adjacent monsters see you on dark cells."
	case StatusIlluminated:
		return "You are always in a lighted cell."
	case StatusTransparent:
		return "Only adjacent monsters see you on lighted cells."
	case StatusDisguised:
		return "Most monsters will ignore you, except those with good flair."
	case StatusDelay:
		return "Time remaining before the trigger."
	case StatusDispersal:
		return "Monsters that attempt to hit you will blink away."
	default:
		// should not happen
		return "unknown"
	}
}

func (st status) Short() string {
	switch st {
	case StatusExhausted:
		return "Ex"
	case StatusSwift:
		return "Sw"
	case StatusLignification:
		return "Lg"
	case StatusConfusion:
		return "Co"
	case StatusFlames:
		return "Fl"
	case StatusHidden:
		return "H+"
	case StatusUnhidden:
		return "H-"
	case StatusDig:
		return "Di"
	case StatusLight:
		return "Li"
	case StatusLevitation:
		return "Le"
	case StatusShadows:
		return "Sh"
	case StatusIlluminated:
		return "Il"
	case StatusTransparent:
		return "Tr"
	case StatusDisguised:
		return "Dg"
	case StatusDelay:
		return "De"
	case StatusDispersal:
		return "Dp"
	default:
		// should not happen
		return "?"
	}
}

func (md *model) statusHPColor() rune {
	g := md.g
	hpColor := 'G'
	switch g.Player.HP + g.Player.HPbonus {
	case 1, 2:
		hpColor = 'C'
	case 3, 4:
		hpColor = 'W'
	}
	return hpColor
}

func (md *model) statusMPColor() rune {
	g := md.g
	mpColor := 'g'
	switch g.Player.MP {
	case 1, 2:
		mpColor = 'c'
	case 3, 4:
		mpColor = 'w'
	}
	return mpColor
}

func (md *model) sortedStatuses() statusSlice {
	g := md.g
	sts := statusSlice{}
	if cld, ok := g.Clouds[g.Player.Pos]; ok && cld == CloudFire {
		g.Player.Statuses[StatusFlames] = 1
		defer func() {
			g.Player.Statuses[StatusFlames] = 0
		}()
	}
	for st, c := range g.Player.Statuses {
		if c > 0 {
			sts = append(sts, st)
		}
	}
	sort.Sort(sts)
	return sts
}

func (md *model) updateStatusInfo() {
	g := md.g
	var entries []ui.MenuEntry

	st := gruid.Style{}
	stt := ui.StyledText{}.WithMarkups(map[rune]gruid.Style{
		'G': st.WithFg(ColorFgHPok),
		'g': st.WithFg(ColorFgMPok),
		'W': st.WithFg(ColorFgHPwounded),
		'w': st.WithFg(ColorFgMPpartial),
		'C': st.WithFg(ColorFgHPcritical),
		'c': st.WithFg(ColorFgMPcritical),
		'x': st.WithFg(ColorFgStatusExpire),
		's': st.WithFg(ColorFgStatusGood),
		'o': st.WithFg(ColorFgStatusOther),
		'b': st.WithFg(ColorFgStatusBad),
		'B': st.WithFg(ColorCyan),
		'M': st.WithFg(ColorYellow).WithAttrs(AttrInMap),
	})
	// depth
	var depth string
	if g.Depth == -1 {
		depth = "D: Out! "
	} else {
		depth = fmt.Sprintf(" D:%d ", g.Depth)
	}
	entries = append(entries, ui.MenuEntry{Text: stt.WithText(depth), Disabled: true})

	// turns
	entries = append(entries, ui.MenuEntry{Text: stt.WithTextf("T: %d ", g.Turn), Disabled: true})

	// HP
	nWounds := g.Player.HPMax() - g.Player.HP - g.Player.HPbonus
	if nWounds <= 0 {
		nWounds = 0
	}
	hpColor := md.statusHPColor()
	hps := "HP:"
	hp := g.Player.HP
	if hp < 0 {
		hp = 0
	}
	if !GameConfig.ShowNumbers {
		hps = fmt.Sprintf("%s@%c%s@B%s@N%s ",
			hps,
			hpColor,
			strings.Repeat("♥", hp),
			strings.Repeat("♥", g.Player.HPbonus),
			strings.Repeat("♥", nWounds),
		)
	} else {
		if g.Player.HPbonus > 0 {
			hps = fmt.Sprintf("@%c%d+%d/%d@N ", hpColor, hp, g.Player.HPbonus, g.Player.HPMax())
		} else {
			hps = fmt.Sprintf("@%c%d/%d@N ", hpColor, hp, g.Player.HPMax())
		}
	}
	entries = append(entries, ui.MenuEntry{Text: stt.WithText(hps), Disabled: true})

	// MP
	MPspent := g.Player.MPMax() - g.Player.MP
	if MPspent <= 0 {
		MPspent = 0
	}
	mpColor := md.statusMPColor()
	mps := "MP:"
	if !GameConfig.ShowNumbers {
		mps = fmt.Sprintf("%s@%c%s@N%s ",
			mps,
			mpColor,
			strings.Repeat("♥", g.Player.MP),
			strings.Repeat("♥", MPspent),
		)
	} else {
		mps = fmt.Sprintf("@%c%d/%d@N ", mpColor, g.Player.MP, g.Player.MPMax())
	}
	entries = append(entries, ui.MenuEntry{Text: stt.WithText(mps), Disabled: true})

	// bananas
	bananas := fmt.Sprintf("@M)@N:%1d/%1d ", g.Player.Bananas, MaxBananas)
	entries = append(entries, ui.MenuEntry{Text: stt.WithText(bananas), Disabled: true})

	// statuses (TODO: add description)
	sts := md.sortedStatuses()

	if len(sts) > 0 {
		// entries[5]
		entries = append(entries, ui.MenuEntry{Text: stt.WithText("| "), Disabled: true})
	}
	for _, st := range sts {
		r := 'o'
		if st.Good() {
			r = 's'
			t := DurationTurn
			if g.Ev != nil && g.Player.Expire[st] >= g.Ev.Rank() && g.Player.Expire[st]-g.Ev.Rank() <= t {
				r = 'x'
			}
		} else if st.Bad() {
			r = 'b'
		}
		var sttext string
		if !st.Flag() {
			sttext = fmt.Sprintf("@%c%s(%d)@N ", r, st.Short(), g.Player.Statuses[st]/DurationStatusStep)
		} else {
			sttext = fmt.Sprintf("@%c%s@N ", r, st.Short())
		}
		entries = append(entries, ui.MenuEntry{Text: stt.WithText(sttext), Disabled: true})
	}

	//altBgEntries(entries)
	md.status.SetEntries(entries)
}
