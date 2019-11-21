package main

import "container/heap"

type event interface {
	Rank() int
	Action(*game)
	Renew(*game, int)
}

type iEvent struct {
	Event event
	Index int
}

type eventQueue []iEvent

func (evq eventQueue) Len() int {
	return len(evq)
}

func (evq eventQueue) Less(i, j int) bool {
	return evq[i].Event.Rank() < evq[j].Event.Rank() ||
		evq[i].Event.Rank() == evq[j].Event.Rank() && evq[i].Index < evq[j].Index
}

func (evq eventQueue) Swap(i, j int) {
	evq[i], evq[j] = evq[j], evq[i]
}

func (evq *eventQueue) Push(x interface{}) {
	no := x.(iEvent)
	*evq = append(*evq, no)
}

func (evq *eventQueue) Pop() interface{} {
	old := *evq
	n := len(old)
	no := old[n-1]
	*evq = old[0 : n-1]
	return no
}

type simpleAction int

const (
	PlayerTurn simpleAction = iota
	SlowEnd
	ExhaustionEnd
	HasteEnd
	EvasionEnd
	LignificationEnd
	ConfusionEnd
	NauseaEnd
	DigEnd
	LevitationEnd
	ShadowsEnd
	IlluminatedEnd
	ShaedraAnimation
	ArtifactAnimation
)

func (g *game) PushEvent(ev event) {
	iev := iEvent{Event: ev, Index: g.EventIndex}
	g.EventIndex++
	heap.Push(g.Events, iev)
}

func (g *game) PushAgainEvent(ev event) {
	iev := iEvent{Event: ev, Index: 0}
	heap.Push(g.Events, iev)
}

func (g *game) PopIEvent() iEvent {
	iev := heap.Pop(g.Events).(iEvent)
	return iev
}

type simpleEvent struct {
	ERank   int
	EAction simpleAction
}

func (sev *simpleEvent) Rank() int {
	return sev.ERank
}

func (sev *simpleEvent) Renew(g *game, delay int) {
	sev.ERank += delay
	if delay == 0 {
		g.PushAgainEvent(sev)
	} else {
		g.PushEvent(sev)
	}
}

const DurationStatusStep = 10

var StatusEndMsgs = [...]string{
	SlowEnd:          "You no longer feel slow.",
	ExhaustionEnd:    "You no longer feel exhausted.",
	HasteEnd:         "You no longer feel speedy.",
	LignificationEnd: "You no longer feel attached to the ground.",
	ConfusionEnd:     "You no longer feel confused.",
	NauseaEnd:        "You no longer feel sick.",
	DigEnd:           "You no longer feel like an earth dragon.",
	LevitationEnd:    "You no longer levitate.",
	ShadowsEnd:       "You are no longer surrounded by shadows.",
	IlluminatedEnd:   "You are no longer illuminated.",
}

var EndStatuses = [...]status{
	SlowEnd:          StatusSlow,
	ExhaustionEnd:    StatusExhausted,
	HasteEnd:         StatusSwift,
	LignificationEnd: StatusLignification,
	ConfusionEnd:     StatusConfusion,
	NauseaEnd:        StatusNausea,
	DigEnd:           StatusDig,
	LevitationEnd:    StatusLevitation,
	ShadowsEnd:       StatusShadows,
	IlluminatedEnd:   StatusIlluminated,
}

var StatusEndActions = [...]simpleAction{
	StatusSlow:          SlowEnd,
	StatusExhausted:     ExhaustionEnd,
	StatusSwift:         HasteEnd,
	StatusLignification: LignificationEnd,
	StatusConfusion:     ConfusionEnd,
	StatusNausea:        NauseaEnd,
	StatusDig:           DigEnd,
	StatusLevitation:    LevitationEnd,
	StatusShadows:       ShadowsEnd,
	StatusIlluminated:   IlluminatedEnd,
}

func (sev *simpleEvent) Action(g *game) {
	switch sev.EAction {
	case PlayerTurn:
		g.ComputeNoise()
		g.ComputeLOS() // TODO: optimize? most of the time almost redundant (unless on a tree)
		g.ComputeMonsterLOS()
		g.LogNextTick = g.LogIndex
		g.AutoNext = g.AutoPlayer(sev)
		if g.AutoNext {
			g.TurnStats()
			return
		}
		g.Quit = g.ui.HandlePlayerTurn(sev)
		if g.Quit {
			return
		}
		g.TurnStats()
	case ShaedraAnimation:
		g.ComputeLOS()
		g.ui.FreeingShaedraAnimation()
	case ArtifactAnimation:
		g.ComputeLOS()
		g.ui.TakingArtifactAnimation()
	case SlowEnd, ExhaustionEnd, HasteEnd, LignificationEnd, ConfusionEnd, NauseaEnd, DigEnd, LevitationEnd, ShadowsEnd, IlluminatedEnd:
		g.Player.Statuses[EndStatuses[sev.EAction]] -= DurationStatusStep
		if g.Player.Statuses[EndStatuses[sev.EAction]] <= 0 {
			g.Player.Statuses[EndStatuses[sev.EAction]] = 0
			g.PrintStyled(StatusEndMsgs[sev.EAction], logStatusEnd)
			g.ui.StatusEndAnimation()
			switch sev.EAction {
			case LevitationEnd:
				if g.Dungeon.Cell(g.Player.Pos).T == ChasmCell {
					g.FallAbyss(DescendFall)
				}
			case LignificationEnd:
				g.Player.HPbonus -= LignificationHPbonus
				if g.Player.HPbonus < 0 {
					g.Player.HPbonus = 0
				}
			}
		} else {
			g.PushEvent(&simpleEvent{ERank: sev.Rank() + DurationStatusStep, EAction: sev.EAction})
		}
	}
}

type monsterAction int

const (
	MonsterTurn monsterAction = iota
	MonsConfusionEnd
	MonsExhaustionEnd
	MonsSlowEnd
	MonsSatiatedEnd
	MonsLignificationEnd
)

type monsterEvent struct {
	ERank   int
	NMons   int
	EAction monsterAction
}

func (mev *monsterEvent) Rank() int {
	return mev.ERank
}

var MonsStatusEndMsgs = [...]string{
	MonsConfusionEnd:     "confused",
	MonsLignificationEnd: "lignified",
	MonsSlowEnd:          "slowed",
	MonsExhaustionEnd:    "exhausted",
	MonsSatiatedEnd:      "satiated",
}

var MonsEndStatuses = [...]monsterStatus{
	MonsConfusionEnd:     MonsConfused,
	MonsLignificationEnd: MonsLignified,
	MonsSlowEnd:          MonsSlow,
	MonsExhaustionEnd:    MonsExhausted,
	MonsSatiatedEnd:      MonsSatiated,
}

var MonsStatusEndActions = [...]monsterAction{
	MonsConfused:  MonsConfusionEnd,
	MonsLignified: MonsLignificationEnd,
	MonsSlow:      MonsSlowEnd,
	MonsExhausted: MonsExhaustionEnd,
	MonsSatiated:  MonsSatiatedEnd,
}

func (mev *monsterEvent) Action(g *game) {
	switch mev.EAction {
	case MonsterTurn:
		mons := g.Monsters[mev.NMons]
		if mons.Exists() {
			mons.HandleTurn(g, mev)
		}
	default:
		mons := g.Monsters[mev.NMons]
		mons.Statuses[MonsEndStatuses[mev.EAction]] -= DurationStatusStep
		if mons.Statuses[MonsEndStatuses[mev.EAction]] <= 0 {
			mons.Statuses[MonsEndStatuses[mev.EAction]] = 0
			if g.Player.Sees(mons.Pos) {
				g.Printf("%s is no longer %s.", mons.Kind.Definite(true), StatusEndMsgs[mev.EAction])
			}
			switch mev.EAction {
			case MonsConfusionEnd, MonsLignificationEnd:
				mons.Path = mons.APath(g, mons.Pos, mons.Target)
			}
		} else {
			g.PushEvent(&monsterEvent{NMons: mev.NMons, ERank: mev.Rank() + DurationStatusStep, EAction: mev.EAction})
		}
	}
}

func (mev *monsterEvent) Renew(g *game, delay int) {
	mev.ERank += delay
	g.PushEvent(mev)
}

type cloudAction int

const (
	CloudEnd cloudAction = iota
	ObstructionEnd
	ObstructionProgression
	FireProgression
	NightProgression
)

type cloudEvent struct {
	ERank   int
	Pos     position
	EAction cloudAction
}

func (cev *cloudEvent) Rank() int {
	return cev.ERank
}

func (cev *cloudEvent) Action(g *game) {
	switch cev.EAction {
	case CloudEnd:
		delete(g.Clouds, cev.Pos)
		g.ComputeLOS()
	case ObstructionEnd:
		t := g.MagicalBarriers[cev.Pos]
		if !g.Player.Sees(cev.Pos) && g.Dungeon.Cell(cev.Pos).T == BarrierCell {
			// XXX does not handle all cases
			g.UpdateKnowledge(cev.Pos, BarrierCell)
		} else {
			delete(g.MagicalBarriers, cev.Pos)
			delete(g.TerrainKnowledge, cev.Pos)
		}
		if g.Dungeon.Cell(cev.Pos).T != BarrierCell {
			break
		}
		g.Dungeon.SetCell(cev.Pos, t)
	case ObstructionProgression:
		pos := g.FreePassableCell()
		g.MagicalBarrierAt(pos, cev)
		if g.Player.Sees(pos) {
			g.Printf("You see an oric barrier appear out of thin air.")
			g.StopAuto()
		}
		g.PushEvent(&cloudEvent{ERank: cev.Rank() + DurationObstructionProgression + RandInt(DurationObstructionProgression/4),
			EAction: ObstructionProgression})
	case FireProgression:
		if _, ok := g.Clouds[cev.Pos]; !ok {
			break
		}
		for _, pos := range g.Dungeon.FreeNeighbors(cev.Pos) {
			if RandInt(4) == 0 {
				continue
			}
			g.Burn(pos, cev)
		}
		delete(g.Clouds, cev.Pos)
		g.NightFog(cev.Pos, 1, &simpleEvent{ERank: cev.Rank()})
		g.ComputeLOS()
	case NightProgression:
		if _, ok := g.Clouds[cev.Pos]; !ok {
			break
		}
		g.MakeCreatureSleep(cev.Pos, cev)
		if RandInt(20) == 0 {
			delete(g.Clouds, cev.Pos)
			g.ComputeLOS()
			break
		}
		cev.Renew(g, 10)
	}
}

func (g *game) NightFog(at position, radius int, ev event) {
	dij := &noisePath{game: g}
	nm := Dijkstra(dij, []position{at}, radius)
	nm.iter(at, func(n *node) {
		pos := n.Pos
		_, ok := g.Clouds[pos]
		if !ok {
			g.Clouds[pos] = CloudNight
			g.PushEvent(&cloudEvent{ERank: ev.Rank() + DurationCloudProgression, EAction: NightProgression, Pos: pos})
			g.MakeCreatureSleep(pos, ev)
		}
	})
	g.ComputeLOS()
}

func (g *game) MakeCreatureSleep(pos position, ev event) {
	if pos == g.Player.Pos {
		if g.PutStatus(StatusSlow, DurationSleepSlow) {
			g.Print("The clouds of night make you sleepy.")
		}
		return
	}
	mons := g.MonsterAt(pos)
	if !mons.Exists() || (RandInt(2) == 0 && mons.Status(MonsExhausted)) {
		// do not always make already exhausted monsters sleep (they were probably awaken)
		return
	}
	if mons.State != Resting && g.Player.Sees(mons.Pos) {
		g.Printf("%s falls asleep.", mons.Kind.Definite(true))
	}
	mons.State = Resting
	mons.Dir = NoDir
	mons.ExhaustTime(g, 40+RandInt(10))
}

func (g *game) Burn(pos position, ev event) {
	if _, ok := g.Clouds[pos]; ok {
		return
	}
	c := g.Dungeon.Cell(pos)
	if !c.Flammable() {
		return
	}
	g.Stats.Burns++
	switch c.T {
	case DoorCell:
		g.Print("The door vanishes in magical flames.")
	case TableCell:
		g.Print("The table vanishes in magical flames.")
	case BarrelCell:
		g.Print("The barrel vanishes in magical flames.")
		delete(g.Objects.Barrels, pos)
	case TreeCell:
		g.Print("The tree vanishes in magical flames.")
	}
	g.Dungeon.SetCell(pos, GroundCell)
	g.Clouds[pos] = CloudFire
	if !g.Player.Sees(pos) {
		g.UpdateKnowledge(pos, c.T)
	} else {
		g.ComputeLOS()
	}
	g.PushEvent(&cloudEvent{ERank: ev.Rank() + DurationCloudProgression, EAction: FireProgression, Pos: pos})
}

func (cev *cloudEvent) Renew(g *game, delay int) {
	cev.ERank += delay
	g.PushEvent(cev)
}

const (
	DurationSwiftness              = 70
	DurationShadows                = 150
	DurationLevitation             = 180
	DurationShortSwiftness         = 30
	DurationDigging                = 80
	DurationSlowMonster            = 170
	DurationSleepSlow              = 40
	DurationCloudProgression       = 10
	DurationFog                    = 150
	DurationExhaustion             = 60
	DurationConfusionMonster       = 130
	DurationConfusionPlayer        = 50
	DurationLignificationMonster   = 150
	DurationLignificationPlayer    = 40
	DurationMagicalBarrier         = 150
	DurationObstructionProgression = 150
	DurationSmokingCloakFog        = 20
	DurationExhaustionMonster      = 100
	DurationSatiationMonster       = 400
	DurationIlluminated            = 70
)
