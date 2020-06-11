package main

import "container/heap"

type event interface {
	Rank() int
	Renew(*game, int)
	Action(*game)
}

type iEvent struct {
	Event event
	Index int
}

func (g *game) PushEvent(ev event) {
	iev := iEvent{Event: ev, Index: g.EventIndex}
	g.EventIndex++
	heap.Push(g.Events, iev)
}

// PushEventRandomIndex pushes a new even to the heap, with randomised Index.
// Used so that monster turn order is not predictable.
func (g *game) PushEventRandomIndex(ev event) {
	iev := iEvent{Event: ev, Index: RandInt(10)}
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

func (g *game) RenewEvent(delay int) {
	g.Ev.Renew(g, delay)
}

const (
	DurationSwiftness              = 4
	DurationShadows                = 15
	DurationLevitation             = 18
	DurationShortSwiftness         = 3
	DurationDigging                = 8
	DurationParalysisMonster       = 6
	DurationCloudProgression       = 1
	DurationFog                    = 15
	DurationExhaustion             = 6
	DurationConfusionMonster       = 12
	DurationConfusionPlayer        = 5
	DurationLignificationMonster   = 15
	DurationLignificationPlayer    = 4
	DurationMagicalBarrier         = 15
	DurationObstructionProgression = 15
	DurationMistProgression        = 12
	DurationSmokingCloakFog        = 2
	DurationExhaustionMonster      = 10
	DurationSatiationMonster       = 40
	DurationIlluminated            = 7
	DurationTransparency           = 15
	DurationDisguise               = 9
	DurationDispersal              = 12
	DurationHarmonicNoiseDelay     = 13
	DurationOricExplosionDelay     = 8
	DurationTurn                   = 1
	DurationStatusStep             = 1
	DurationNightFog               = 15
)

type simpleAction int

const (
	PlayerTurn simpleAction = iota
	AbyssFall
)

type simpleEvent struct {
	ERank   int
	EAction simpleAction
}

func (ev *simpleEvent) Rank() int {
	return ev.ERank
}

func (ev *simpleEvent) Renew(g *game, delay int) {
	ev.ERank += delay
	if delay == 0 {
		g.PushAgainEvent(ev)
		if ev.EAction == PlayerTurn {
			g.PlayerAgain = true
		}
	} else {
		g.PushEvent(ev)
	}
}

func (ev *simpleEvent) Action(g *game) {
	switch ev.EAction {
	case PlayerTurn:
		if !g.PlayerAgain {
			g.ComputeNoise()
			g.ComputeLOS() // TODO: optimize? most of the time almost redundant (unless on a tree)
			g.ComputeMonsterLOS()
		}
		g.PlayerAgain = false
		g.LogNextTick = g.LogIndex
		g.AutoNext = g.AutoPlayer(ev)
		if g.AutoNext {
			g.TurnStats()
			return
		}
		g.Quit = g.ui.HandlePlayerTurn()
		if g.Quit {
			return
		}
		g.TurnStats()
	case AbyssFall:
		if g.Dungeon.Cell(g.Player.Pos).T == ChasmCell {
			g.FallAbyss(DescendFall)
		}
	}
}

type statusEvent struct {
	ERank  int
	Status status
}

var StatusEndMsgs = [...]string{
	StatusExhausted:     "You no longer feel exhausted.",
	StatusSwift:         "You no longer feel speedy.",
	StatusLignification: "You no longer feel attached to the ground.",
	StatusConfusion:     "You no longer feel confused.",
	StatusNausea:        "You no longer feel sick.",
	StatusDig:           "You no longer feel like an earth dragon.",
	StatusLevitation:    "You no longer levitate.",
	StatusShadows:       "You are no longer surrounded by shadows.",
	StatusIlluminated:   "You are no longer illuminated.",
	StatusTransparent:   "You are no longer transparent.",
	StatusDisguised:     "You are no longer disguised.",
	StatusDispersal:     "You are no longer unstable.",
}

func (ev *statusEvent) Rank() int {
	return ev.ERank
}

func (ev *statusEvent) Renew(g *game, delay int) {
	ev.ERank += delay
	g.PushEvent(ev)
}

func (ev *statusEvent) Action(g *game) {
	st := ev.Status
	g.Player.Statuses[st] -= DurationStatusStep
	if g.Player.Statuses[st] <= 0 {
		g.Player.Statuses[st] = 0
		g.PrintStyled(StatusEndMsgs[st], logStatusEnd)
		g.ui.StatusEndAnimation()
		switch st {
		case StatusLevitation:
			if g.Dungeon.Cell(g.Player.Pos).T == ChasmCell {
				g.FallAbyss(DescendFall)
			}
		case StatusLignification:
			g.Player.HPbonus -= LignificationHPbonus
			if g.Player.HPbonus < 0 {
				g.Player.HPbonus = 0
			}
		}
	} else {
		ev.Renew(g, DurationStatusStep)
	}
}

type monsterEvent struct {
	ERank int
	NMons int
}

func (mev *monsterEvent) Rank() int {
	return mev.ERank
}

func (mev *monsterEvent) Renew(g *game, delay int) {
	mev.ERank += delay
	g.PushEvent(mev)
}

func (mev *monsterEvent) Action(g *game) {
	mons := g.Monsters[mev.NMons]
	if mons.Exists() {
		mons.HandleTurn(g)
		if mons.Exists() {
			mev.Renew(g, DurationTurn)
		}
	}
}

type monsterStatusEvent struct {
	ERank  int
	NMons  int
	Status monsterStatus
}

var MonsStatusEndMsgs = [...]string{
	MonsConfused:  "confused",
	MonsLignified: "lignified",
	MonsParalysed: "slowed",
	MonsExhausted: "exhausted",
	MonsSatiated:  "satiated",
}

func (mev *monsterStatusEvent) Rank() int {
	return mev.ERank
}

func (mev *monsterStatusEvent) Renew(g *game, delay int) {
	mev.ERank += delay
	g.PushEvent(mev)
}

func (mev *monsterStatusEvent) Action(g *game) {
	mons := g.Monsters[mev.NMons]
	mons.Statuses[mev.Status] -= DurationStatusStep
	if mons.Statuses[mev.Status] <= 0 {
		mons.Statuses[mev.Status] = 0
		if g.Player.Sees(mons.Pos) {
			g.Printf("%s is no longer %s.", mons.Kind.Definite(true), StatusEndMsgs[mev.Status])
		}
		switch mev.Status {
		case MonsConfused, MonsLignified:
			mons.Path = mons.APath(g, mons.Pos, mons.Target)
		}
	} else {
		g.PushEvent(&monsterStatusEvent{NMons: mev.NMons, ERank: mev.Rank() + DurationStatusStep, Status: mev.Status})
	}
}

type posAction int

const (
	CloudEnd posAction = iota
	ObstructionEnd
	ObstructionProgression
	FireProgression
	NightProgression
	MistProgression
	Earthquake
	DelayedHarmonicNoiseEvent
	DelayedOricExplosionEvent
)

type posEvent struct {
	ERank   int
	Pos     position
	EAction posAction
	Timer   int
}

func (ev *posEvent) Rank() int {
	return ev.ERank
}

func (ev *posEvent) Renew(g *game, delay int) {
	ev.ERank += delay
	g.PushEvent(ev)
}

func (ev *posEvent) Action(g *game) {
	switch ev.EAction {
	case CloudEnd:
		delete(g.Clouds, ev.Pos)
		g.ComputeLOS()
	case ObstructionEnd:
		t := g.MagicalBarriers[ev.Pos]
		if !g.Player.Sees(ev.Pos) && g.Dungeon.Cell(ev.Pos).T == BarrierCell {
			// XXX does not handle all cases
			g.UpdateKnowledge(ev.Pos, BarrierCell)
		} else {
			delete(g.MagicalBarriers, ev.Pos)
			delete(g.TerrainKnowledge, ev.Pos)
		}
		if g.Dungeon.Cell(ev.Pos).T != BarrierCell {
			break
		}
		g.Dungeon.SetCell(ev.Pos, t)
	case ObstructionProgression:
		pos := g.FreePassableCell()
		g.MagicalBarrierAt(pos)
		if g.Player.Sees(pos) {
			g.Printf("You see an oric barrier appear out of thin air.")
			g.StopAuto()
		}
		g.PushEvent(&posEvent{ERank: ev.Rank() + DurationObstructionProgression + RandInt(DurationObstructionProgression/4),
			EAction: ObstructionProgression})
	case FireProgression:
		if _, ok := g.Clouds[ev.Pos]; !ok {
			break
		}
		for _, pos := range g.Dungeon.FreeNeighbors(ev.Pos) {
			if RandInt(10) == 0 {
				continue
			}
			g.Burn(pos)
		}
		delete(g.Clouds, ev.Pos)
		g.NightFog(ev.Pos, 1, &simpleEvent{ERank: ev.Rank()})
		g.ComputeLOS()
	case NightProgression:
		if _, ok := g.Clouds[ev.Pos]; !ok {
			break
		}
		if ev.Timer <= 0 {
			delete(g.Clouds, ev.Pos)
			g.ComputeLOS()
			break
		}
		g.MakeCreatureSleep(ev.Pos)
		ev.Timer--
		ev.Renew(g, DurationTurn)
	case MistProgression:
		pos := g.FreePassableCell()
		g.Fog(pos, 1)
		g.PushEvent(&posEvent{ERank: ev.Rank() + DurationMistProgression + RandInt(DurationMistProgression/4),
			EAction: MistProgression})
	case Earthquake:
		g.PrintStyled("The earth suddenly shakes with force!", logSpecial)
		g.PrintStyled("Craack!", logSpecial)
		g.StoryPrint("Special event: earthquake!")
		g.MakeNoise(EarthquakeNoise, ev.Pos)
		g.NoiseIllusion[ev.Pos] = true
		for i, c := range g.Dungeon.Cells {
			pos := idxtopos(i)
			if !c.T.IsDiggable() || !g.Dungeon.HasFreeNeighbor(pos) {
				continue
			}
			if ev.Pos.Distance(pos) > RandInt(35) || RandInt(2) == 0 {
				continue
			}
			g.Dungeon.SetCell(pos, RubbleCell)
			g.UpdateKnowledge(pos, c.T)
			g.Fog(pos, 1)
		}
	case DelayedHarmonicNoiseEvent:
		if ev.Timer <= 1 {
			g.Player.Statuses[StatusDelay] = 0
			g.Print("Pop!")
			g.NoiseIllusion[ev.Pos] = true
			g.MakeNoise(DelayedHarmonicNoise, ev.Pos)
		} else {
			ev.Timer--
			g.Player.Statuses[StatusDelay] = ev.Timer
			ev.Renew(g, DurationTurn)
		}
	case DelayedOricExplosionEvent:
		if ev.Timer <= 1 {
			g.Player.Statuses[StatusDelay] = 0
			g.Print(g.CrackSound())
			g.NoiseIllusion[ev.Pos] = true
			dij := &gridPath{dungeon: g.Dungeon}
			g.MakeNoise(OricExplosionNoise, ev.Pos)
			nm := Dijkstra(dij, []position{ev.Pos}, 7)
			fogs := []position{}
			terrains := []terrain{}
			nm.iter(ev.Pos, func(n *node) {
				c := g.Dungeon.Cell(n.Pos)
				if !c.T.IsDiggable() {
					return
				}
				g.Dungeon.SetCell(n.Pos, RubbleCell)
				g.Stats.Digs++
				if g.Player.Sees(n.Pos) {
					g.ui.WallExplosionAnimation(n.Pos)
				}
				fogs = append(fogs, n.Pos)
				terrains = append(terrains, c.T)
			})
			for _, pos := range fogs {
				g.Fog(pos, 1)
			}
			g.ComputeLOS()
			for i, pos := range fogs {
				g.UpdateKnowledge(pos, terrains[i])
			}
		} else {
			ev.Timer--
			g.Player.Statuses[StatusDelay] = ev.Timer
			ev.Renew(g, DurationTurn)
		}
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
			g.PushEvent(&posEvent{ERank: ev.Rank() + DurationCloudProgression, EAction: NightProgression,
				Pos: pos, Timer: DurationNightFog})
			g.MakeCreatureSleep(pos)
		}
	})
	g.ComputeLOS()
}

func (g *game) MakeCreatureSleep(pos position) {
	if pos == g.Player.Pos {
		if g.PutStatus(StatusConfusion, DurationConfusionPlayer) {
			g.Print("The clouds of night confuse you.")
		}
		return
	}
	mons := g.MonsterAt(pos)
	if !mons.Exists() || (RandInt(2) == 0 && mons.Status(MonsExhausted)) {
		// do not always make already exhausted monsters sleep (they were probably awaken)
		return
	}
	mons.EnterConfusion(g)
	if mons.State != Resting && g.Player.Sees(mons.Pos) {
		g.Printf("%s falls asleep.", mons.Kind.Definite(true))
	}
	mons.State = Resting
	mons.Dir = NoDir
	mons.ExhaustTime(g, 4+RandInt(2))
}

func (g *game) Burn(pos position) {
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
	case EssenciaticSourceCell:
		g.Print("The barrel vanishes in magical flames.")
		delete(g.Objects.EssenciaticSources, pos)
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
	g.PushEvent(&posEvent{ERank: g.Ev.Rank() + DurationCloudProgression, EAction: FireProgression, Pos: pos})
}

// eventQueue datastructure for the queue
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
