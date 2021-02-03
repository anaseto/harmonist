package main

import (
	"github.com/anaseto/gruid"
)

type event interface {
	Rank() int
	Action(*game)
	Renew(*game, int)
}

type simpleAction int

const (
	PlayerTurn simpleAction = iota
	ShaedraAnimation
	ArtifactAnimation
	AbyssFall
	ExhaustionEnd
	SwiftEnd
	EvasionEnd
	LignificationEnd
	ConfusionEnd
	NauseaEnd
	DigEnd
	LevitationEnd
	ShadowsEnd
	IlluminatedEnd
	TransparentEnd
	DisguisedEnd
	DispersalEnd
)

func (g *game) PushEvent(ev event) {
	g.Events.Push(ev, ev.Rank())
}

func (g *game) PushEventFirst(ev event) {
	g.Events.PushFirst(ev, ev.Rank())
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
		g.PushEventFirst(sev)
		if sev.EAction == PlayerTurn {
			g.PlayerAgain = true
		}
	} else {
		g.PushEvent(sev)
	}
}

var StatusEndMsgs = [...]string{
	ExhaustionEnd:    "You no longer feel exhausted.",
	SwiftEnd:         "You no longer feel speedy.",
	LignificationEnd: "You no longer feel attached to the ground.",
	ConfusionEnd:     "You no longer feel confused.",
	NauseaEnd:        "You no longer feel sick.",
	DigEnd:           "You no longer feel like an earth dragon.",
	LevitationEnd:    "You no longer levitate.",
	ShadowsEnd:       "You are no longer surrounded by shadows.",
	IlluminatedEnd:   "You are no longer illuminated.",
	TransparentEnd:   "You are no longer transparent.",
	DisguisedEnd:     "You are no longer disguised.",
	DispersalEnd:     "You are no longer unstable.",
}

var EndStatuses = [...]status{
	ExhaustionEnd:    StatusExhausted,
	LignificationEnd: StatusLignification,
	ConfusionEnd:     StatusConfusion,
	DigEnd:           StatusDig,
	LevitationEnd:    StatusLevitation,
	ShadowsEnd:       StatusShadows,
	IlluminatedEnd:   StatusIlluminated,
	TransparentEnd:   StatusTransparent,
	DisguisedEnd:     StatusDisguised,
	DispersalEnd:     StatusDispersal,
}

var StatusEndActions = [...]simpleAction{
	StatusExhausted:     ExhaustionEnd,
	StatusLignification: LignificationEnd,
	StatusConfusion:     ConfusionEnd,
	StatusDig:           DigEnd,
	StatusLevitation:    LevitationEnd,
	StatusShadows:       ShadowsEnd,
	StatusIlluminated:   IlluminatedEnd,
	StatusTransparent:   TransparentEnd,
	StatusDisguised:     DisguisedEnd,
	StatusDispersal:     DispersalEnd,
}

func (sev *simpleEvent) Action(g *game) {
	switch sev.EAction {
	case PlayerTurn:
		if !g.PlayerAgain {
			g.ComputeNoise()
			g.ComputeLOS() // TODO: optimize? most of the time almost redundant (unless on a tree)
			g.ComputeMonsterLOS()
		}
		g.PlayerAgain = false
		g.LogNextTick = g.LogIndex
		g.AutoNext = g.AutoPlayer()
		g.TurnStats()
	case ShaedraAnimation:
		g.ComputeLOS()
		g.ui.FreeingShaedraAnimation()
	case ArtifactAnimation:
		g.ComputeLOS()
		g.ui.TakingArtifactAnimation()
	case AbyssFall:
		if terrain(g.Dungeon.Cell(g.Player.Pos)) == ChasmCell {
			g.FallAbyss(DescendFall)
		}
	default:
		st := EndStatuses[sev.EAction]
		g.Player.Statuses[st] -= DurationStatusStep
		if g.Player.Statuses[st] <= 0 {
			g.Player.Statuses[st] = 0
			g.PrintStyled(StatusEndMsgs[sev.EAction], logStatusEnd)
			//g.ui.StatusEndAnimation()
			switch sev.EAction {
			case LevitationEnd:
				if terrain(g.Dungeon.Cell(g.Player.Pos)) == ChasmCell {
					g.FallAbyss(DescendFall)
				}
			case LignificationEnd:
				g.Player.HPbonus -= LignificationHPbonus
				if g.Player.HPbonus < 0 {
					g.Player.HPbonus = 0
				}
			}
		} else {
			sev.Renew(g, DurationStatusStep)
		}
	}
}

type monsterAction int

const (
	MonsterTurn monsterAction = iota
	MonsConfusionEnd
	MonsExhaustionEnd
	MonsParalysedEnd
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

//var MonsStatusEndMsgs = [...]string{
//MonsConfusionEnd:     "confused",
//MonsLignificationEnd: "lignified",
//MonsParalysedEnd:     "slowed",
//MonsExhaustionEnd:    "exhausted",
//MonsSatiatedEnd:      "satiated",
//}

var MonsEndStatuses = [...]monsterStatus{
	MonsConfusionEnd:     MonsConfused,
	MonsLignificationEnd: MonsLignified,
	MonsParalysedEnd:     MonsParalysed,
	MonsExhaustionEnd:    MonsExhausted,
	MonsSatiatedEnd:      MonsSatiated,
}

var MonsStatusEndActions = [...]monsterAction{
	MonsConfused:  MonsConfusionEnd,
	MonsLignified: MonsLignificationEnd,
	MonsParalysed: MonsParalysedEnd,
	MonsExhausted: MonsExhaustionEnd,
	MonsSatiated:  MonsSatiatedEnd,
}

func (mev *monsterEvent) Action(g *game) {
	switch mev.EAction {
	case MonsterTurn:
		mons := g.Monsters[mev.NMons]
		if mons.Exists() {
			mons.HandleTurn(g)
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
	Pos     gruid.Point
	EAction posAction
	Timer   int
}

func (cev *posEvent) Rank() int {
	return cev.ERank
}

func (cev *posEvent) Action(g *game) {
	switch cev.EAction {
	case CloudEnd:
		delete(g.Clouds, cev.Pos)
		g.ComputeLOS()
	case ObstructionEnd:
		t := g.MagicalBarriers[cev.Pos]
		if !g.Player.Sees(cev.Pos) && terrain(g.Dungeon.Cell(cev.Pos)) == BarrierCell {
			// XXX does not handle all cases
			g.UpdateKnowledge(cev.Pos, BarrierCell)
		} else {
			delete(g.MagicalBarriers, cev.Pos)
			delete(g.TerrainKnowledge, cev.Pos)
		}
		if terrain(g.Dungeon.Cell(cev.Pos)) != BarrierCell {
			break
		}
		g.Dungeon.SetCell(cev.Pos, t)
	case ObstructionProgression:
		pos := g.FreePassableCell()
		g.MagicalBarrierAt(pos)
		if g.Player.Sees(pos) {
			g.Printf("You see an oric barrier appear out of thin air.")
			g.StopAuto()
		}
		g.PushEvent(&posEvent{ERank: cev.Rank() + DurationObstructionProgression + RandInt(DurationObstructionProgression/4),
			EAction: ObstructionProgression})
	case FireProgression:
		if _, ok := g.Clouds[cev.Pos]; !ok {
			break
		}
		for _, pos := range g.Dungeon.FreeNeighbors(cev.Pos) {
			if RandInt(10) == 0 {
				continue
			}
			g.Burn(pos)
		}
		delete(g.Clouds, cev.Pos)
		g.NightFog(cev.Pos, 1, &simpleEvent{ERank: cev.Rank()})
		g.ComputeLOS()
	case NightProgression:
		if _, ok := g.Clouds[cev.Pos]; !ok {
			break
		}
		if cev.Timer <= 0 {
			delete(g.Clouds, cev.Pos)
			g.ComputeLOS()
			break
		}
		g.MakeCreatureSleep(cev.Pos)
		cev.Timer--
		cev.Renew(g, DurationTurn)
	case MistProgression:
		pos := g.FreePassableCell()
		g.Fog(pos, 1)
		g.PushEvent(&posEvent{ERank: cev.Rank() + DurationMistProgression + RandInt(DurationMistProgression/4),
			EAction: MistProgression})
	case Earthquake:
		g.PrintStyled("The earth suddenly shakes with force!", logSpecial)
		g.PrintStyled("Craack!", logSpecial)
		g.StoryPrint("Special event: earthquake!")
		g.MakeNoise(EarthquakeNoise, cev.Pos)
		g.NoiseIllusion[cev.Pos] = true
		it := g.Dungeon.Grid.Iterator()
		for it.Next() {
			pos := it.P()
			c := cell(it.Cell())
			if !c.IsDiggable() || !g.Dungeon.HasFreeNeighbor(pos) {
				continue
			}
			if Distance(cev.Pos, pos) > RandInt(35) || RandInt(2) == 0 {
				continue
			}
			g.Dungeon.SetCell(pos, RubbleCell)
			g.UpdateKnowledge(pos, terrain(c))
			g.Fog(pos, 1)
		}
	case DelayedHarmonicNoiseEvent:
		if cev.Timer <= 1 {
			g.Player.Statuses[StatusDelay] = 0
			g.Print("Pop!")
			g.NoiseIllusion[cev.Pos] = true
			g.MakeNoise(DelayedHarmonicNoise, cev.Pos)
		} else {
			cev.Timer--
			g.Player.Statuses[StatusDelay] = cev.Timer
			cev.Renew(g, DurationTurn)
		}
	case DelayedOricExplosionEvent:
		if cev.Timer <= 1 {
			g.Player.Statuses[StatusDelay] = 0
			g.Print(g.CrackSound())
			g.NoiseIllusion[cev.Pos] = true
			dij := &gridPath{dungeon: g.Dungeon}
			g.MakeNoise(OricExplosionNoise, cev.Pos)
			nodes := g.PR.DijkstraMap(dij, []gruid.Point{cev.Pos}, 7)
			fogs := []gruid.Point{}
			terrains := []cell{}
			for _, n := range nodes {
				c := g.Dungeon.Cell(n.P)
				if !c.IsDiggable() {
					continue
				}
				g.Dungeon.SetCell(n.P, RubbleCell)
				g.Stats.Digs++
				if g.Player.Sees(n.P) {
					g.ui.WallExplosionAnimation(n.P)
				}
				fogs = append(fogs, n.P)
				terrains = append(terrains, terrain(c))
			}
			for _, pos := range fogs {
				g.Fog(pos, 1)
			}
			g.ComputeLOS()
			for i, pos := range fogs {
				g.UpdateKnowledge(pos, terrains[i])
			}
		} else {
			cev.Timer--
			g.Player.Statuses[StatusDelay] = cev.Timer
			cev.Renew(g, DurationTurn)
		}
	}
}

func (g *game) NightFog(at gruid.Point, radius int, ev event) {
	dij := &noisePath{state: g}
	nodes := g.PR.DijkstraMap(dij, []gruid.Point{at}, radius)
	for _, n := range nodes {
		_, ok := g.Clouds[n.P]
		if !ok {
			g.Clouds[n.P] = CloudNight
			g.PushEvent(&posEvent{ERank: ev.Rank() + DurationCloudProgression, EAction: NightProgression,
				Pos: n.P, Timer: DurationNightFog})
			g.MakeCreatureSleep(n.P)
		}
	}
	g.ComputeLOS()
}

func (g *game) MakeCreatureSleep(pos gruid.Point) {
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

func (g *game) Burn(pos gruid.Point) {
	if _, ok := g.Clouds[pos]; ok {
		return
	}
	c := g.Dungeon.Cell(pos)
	if !c.Flammable() {
		return
	}
	g.Stats.Burns++
	switch terrain(c) {
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
		g.UpdateKnowledge(pos, terrain(c))
	} else {
		g.ComputeLOS()
	}
	g.PushEvent(&posEvent{ERank: g.Ev.Rank() + DurationCloudProgression, EAction: FireProgression, Pos: pos})
}

func (cev *posEvent) Renew(g *game, delay int) {
	cev.ERank += delay
	g.PushEvent(cev)
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

func (g *game) RenewEvent(delay int) {
	g.Ev.Renew(g, delay)
}
