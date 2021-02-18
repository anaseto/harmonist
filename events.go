package main

import (
	"github.com/anaseto/gruid"
)

type event interface {
	Handle(*game)
}

func (g *game) PushEvent(ev event, r int) {
	g.Events.Push(ev, r)
}

func (g *game) PushEventD(ev event, delay int) {
	g.Events.Push(ev, g.Turn+delay)
}

func (g *game) PushEventFirst(ev event, r int) {
	g.Events.PushFirst(ev, r)
}

type playerEventAction int

const (
	PlayerTurn playerEventAction = iota
	StorySequence
	AbyssFall
)

type playerEvent struct {
	Action playerEventAction
}

func (sev *playerEvent) Handle(g *game) {
	switch sev.Action {
	case PlayerTurn:
		g.ComputeNoise()
		g.ComputeLOS() // TODO: optimize? most of the time almost redundant (unless on a tree)
		g.ComputeMonsterLOS()
		g.LogNextTick = g.LogIndex
		g.AutoNext = g.AutoPlayer()
		g.TurnStats()
	case StorySequence:
		g.ComputeLOS()
		g.md.Story()
	case AbyssFall:
		if terrain(g.Dungeon.Cell(g.Player.P)) == ChasmCell {
			g.FallAbyss(DescendFall)
		}
	}
}

type statusEvent struct {
	Status status
}

var StatusEndMsgs = [...]string{
	StatusExhausted:     "You no longer feel exhausted.",
	StatusSwift:         "You no longer feel speedy.",
	StatusLignification: "You no longer feel attached to the ground.",
	StatusConfusion:     "You no longer feel confused.",
	StatusDig:           "You no longer feel like an earth dragon.",
	StatusLevitation:    "You no longer levitate.",
	StatusShadows:       "You are no longer surrounded by shadows.",
	StatusIlluminated:   "You are no longer illuminated.",
	StatusTransparent:   "You are no longer transparent.",
	StatusDisguised:     "You are no longer disguised.",
	StatusDispersal:     "You are no longer unstable.",
}

func (sev *statusEvent) Handle(g *game) {
	st := sev.Status
	g.Player.Statuses[st] -= DurationStatusStep
	if g.Player.Statuses[st] <= 0 {
		g.Player.Statuses[st] = 0
		g.PrintStyled(StatusEndMsgs[st], logStatusEnd)
		g.md.StatusEndAnimation()
		switch st {
		case StatusLevitation:
			if terrain(g.Dungeon.Cell(g.Player.P)) == ChasmCell {
				g.FallAbyss(DescendFall)
			}
		case StatusLignification:
			g.Player.HPbonus -= LignificationHPbonus
			if g.Player.HPbonus < 0 {
				g.Player.HPbonus = 0
			}
		}
	} else {
		g.PushEventD(sev, DurationStatusStep)
	}
}

type monsterTurnEvent struct {
	Index int
}

func (mev *monsterTurnEvent) Handle(g *game) {
	mons := g.Monsters[mev.Index]
	if mons.Exists() {
		mons.HandleTurn(g)
	}
}

type monsterStatusEvent struct {
	Index  int
	Status monsterStatus
}

func (mev *monsterStatusEvent) Handle(g *game) {
	mons := g.Monsters[mev.Index]
	st := mev.Status
	mons.Statuses[st] -= DurationStatusStep
	if mons.Statuses[st] <= 0 {
		mons.Statuses[st] = 0
		if g.Player.Sees(mons.P) {
			g.Printf("%s is no longer %s.", mons.Kind.Definite(true), st)
		}
		switch st {
		case MonsConfused, MonsLignified:
			mons.Path = mons.APath(g, mons.P, mons.Target)
		}
	} else {
		g.PushEventD(&monsterStatusEvent{Index: mev.Index, Status: st}, DurationStatusStep)
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
	P      gruid.Point
	Action posAction
	Timer  int
}

func (cev *posEvent) Handle(g *game) {
	switch cev.Action {
	case CloudEnd:
		delete(g.Clouds, cev.P)
		g.ComputeLOS()
	case ObstructionEnd:
		t := g.MagicalBarriers[cev.P]
		if terrain(g.Dungeon.Cell(cev.P)) == BarrierCell {
			_, ok := g.TerrainKnowledge[cev.P]
			if !ok {
				g.UpdateKnowledge(cev.P, BarrierCell)
			}
		}
		if g.Player.Sees(cev.P) {
			delete(g.MagicalBarriers, cev.P)
			delete(g.TerrainKnowledge, cev.P)
		}
		if terrain(g.Dungeon.Cell(cev.P)) != BarrierCell {
			break
		}
		g.Dungeon.SetCell(cev.P, t)
	case ObstructionProgression:
		p := g.FreePassableCell()
		g.MagicalBarrierAt(p)
		if g.Player.Sees(p) {
			g.Printf("You see an oric barrier appear out of thin air.")
			g.StopAuto()
		}
		g.PushEvent(&posEvent{Action: ObstructionProgression},
			g.Turn+DurationObstructionProgression+RandInt(DurationObstructionProgression/4))
	case FireProgression:
		if _, ok := g.Clouds[cev.P]; !ok {
			break
		}
		for _, p := range g.playerPassableNeighbors(cev.P) {
			if RandInt(10) == 0 {
				continue
			}
			g.Burn(p)
		}
		delete(g.Clouds, cev.P)
		g.NightFog(cev.P, 1)
		g.ComputeLOS()
	case NightProgression:
		if _, ok := g.Clouds[cev.P]; !ok {
			break
		}
		if cev.Timer <= 0 {
			delete(g.Clouds, cev.P)
			g.ComputeLOS()
			break
		}
		g.MakeCreatureSleep(cev.P)
		cev.Timer--
		g.PushEventD(cev, DurationTurn)
	case MistProgression:
		p := g.FreePassableCell()
		g.Fog(p, 1)
		g.PushEvent(&posEvent{Action: MistProgression},
			g.Turn+DurationMistProgression+RandInt(DurationMistProgression/4))
	case Earthquake:
		g.PrintStyled("The earth suddenly shakes with force!", logSpecial)
		g.PrintStyled("Craack!", logSpecial)
		g.StoryPrint("Special event: earthquake!")
		g.MakeNoise(EarthquakeNoise, cev.P)
		g.NoiseIllusion[cev.P] = true
		it := g.Dungeon.Grid.Iterator()
		for it.Next() {
			p := it.P()
			c := cell(it.Cell())
			if !c.IsDiggable() || !g.Dungeon.HasFreeNeighbor(&g.nbs, p) {
				continue
			}
			if distance(cev.P, p) > RandInt(35) || RandInt(2) == 0 {
				continue
			}
			g.Dungeon.SetCell(p, RubbleCell)
			g.UpdateKnowledge(p, terrain(c))
			g.Fog(p, 1)
		}
	case DelayedHarmonicNoiseEvent:
		if cev.Timer <= 1 {
			g.Player.Statuses[StatusDelay] = 0
			g.Print("Pop!")
			g.NoiseIllusion[cev.P] = true
			g.MakeNoise(DelayedHarmonicNoise, cev.P)
		} else {
			cev.Timer--
			g.Player.Statuses[StatusDelay] = cev.Timer
			g.PushEventD(cev, DurationTurn)
		}
	case DelayedOricExplosionEvent:
		if cev.Timer <= 1 {
			g.Player.Statuses[StatusDelay] = 0
			g.Print(g.CrackSound())
			g.NoiseIllusion[cev.P] = true
			dij := &gridPath{dungeon: g.Dungeon}
			g.MakeNoise(OricExplosionNoise, cev.P)
			nodes := g.PR.DijkstraMap(dij, []gruid.Point{cev.P}, 7)
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
					g.md.WallExplosionAnimation(n.P)
				}
				fogs = append(fogs, n.P)
				terrains = append(terrains, terrain(c))
			}
			for _, p := range fogs {
				g.Fog(p, 1)
			}
			g.ComputeLOS()
			for i, p := range fogs {
				g.UpdateKnowledge(p, terrains[i])
			}
		} else {
			cev.Timer--
			g.Player.Statuses[StatusDelay] = cev.Timer
			g.PushEventD(cev, DurationTurn)
		}
	}
}

func (g *game) NightFog(at gruid.Point, radius int) {
	dij := &noisePath{state: g}
	nodes := g.PR.DijkstraMap(dij, []gruid.Point{at}, radius)
	for _, n := range nodes {
		_, ok := g.Clouds[n.P]
		if !ok {
			g.Clouds[n.P] = CloudNight
			g.PushEventD(&posEvent{Action: NightProgression,
				P: n.P, Timer: DurationNightFog}, DurationCloudProgression)

			g.MakeCreatureSleep(n.P)
		}
	}
	g.ComputeLOS()
}

func (g *game) MakeCreatureSleep(p gruid.Point) {
	if p == g.Player.P {
		if g.PutStatus(StatusConfusion, DurationConfusionPlayer) {
			g.Print("The clouds of night confuse you.")
		}
		return
	}
	mons := g.MonsterAt(p)
	if !mons.Exists() || (RandInt(2) == 0 && mons.Status(MonsExhausted)) {
		// do not always make already exhausted monsters sleep (they were probably awaken)
		return
	}
	mons.EnterConfusion(g)
	if mons.State != Resting && g.Player.Sees(mons.P) {
		g.Printf("%s falls asleep.", mons.Kind.Definite(true))
	}
	mons.State = Resting
	mons.Dir = ZP
	mons.ExhaustTime(g, 4+RandInt(2))
}

func (g *game) Burn(p gruid.Point) {
	if _, ok := g.Clouds[p]; ok {
		return
	}
	c := g.Dungeon.Cell(p)
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
		delete(g.Objects.Barrels, p)
	case TreeCell:
		g.Print("The tree vanishes in magical flames.")
	}
	g.Dungeon.SetCell(p, GroundCell)
	g.Clouds[p] = CloudFire
	if !g.Player.Sees(p) {
		g.UpdateKnowledge(p, terrain(c))
	} else {
		g.ComputeLOS()
	}
	g.PushEventD(&posEvent{P: p, Action: FireProgression}, DurationCloudProgression)
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
	DurationExhaustion             = 5
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
