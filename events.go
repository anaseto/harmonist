package main

import (
	"github.com/anaseto/gruid"
)

type event interface {
	Action(*game)
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
	EAction playerEventAction
}

func (sev *playerEvent) Action(g *game) {
	switch sev.EAction {
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
		if terrain(g.Dungeon.Cell(g.Player.Pos)) == ChasmCell {
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

func (sev *statusEvent) Action(g *game) {
	st := sev.Status
	g.Player.Statuses[st] -= DurationStatusStep
	if g.Player.Statuses[st] <= 0 {
		g.Player.Statuses[st] = 0
		g.PrintStyled(StatusEndMsgs[st], logStatusEnd)
		g.md.StatusEndAnimation()
		switch st {
		case StatusLevitation:
			if terrain(g.Dungeon.Cell(g.Player.Pos)) == ChasmCell {
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
	Mons *monster
}

func (mev *monsterTurnEvent) Action(g *game) {
	mons := mev.Mons
	if mons.Exists() {
		mons.HandleTurn(g)
	}
}

type monsterStatusEvent struct {
	Mons   *monster
	Status monsterStatus
}

func (mev *monsterStatusEvent) Action(g *game) {
	mons := mev.Mons
	st := mev.Status
	mons.Statuses[st] -= DurationStatusStep
	if mons.Statuses[st] <= 0 {
		mons.Statuses[st] = 0
		if g.Player.Sees(mons.Pos) {
			g.Printf("%s is no longer %s.", mons.Kind.Definite(true), st)
		}
		switch st {
		case MonsConfused, MonsLignified:
			mons.Path = mons.APath(g, mons.Pos, mons.Target)
		}
	} else {
		g.PushEventD(&monsterStatusEvent{Mons: mev.Mons, Status: st}, DurationStatusStep)
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
	Pos     gruid.Point
	EAction posAction
	Timer   int
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
		g.PushEvent(&posEvent{EAction: ObstructionProgression},
			g.Turn+DurationObstructionProgression+RandInt(DurationObstructionProgression/4))
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
		g.NightFog(cev.Pos, 1)
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
		g.PushEventD(cev, DurationTurn)
	case MistProgression:
		pos := g.FreePassableCell()
		g.Fog(pos, 1)
		g.PushEvent(&posEvent{EAction: MistProgression},
			g.Turn+DurationMistProgression+RandInt(DurationMistProgression/4))
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
			g.PushEventD(cev, DurationTurn)
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
					g.md.WallExplosionAnimation(n.P)
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
			g.PushEventD(&posEvent{EAction: NightProgression,
				Pos: n.P, Timer: DurationNightFog}, DurationCloudProgression)

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
	g.PushEventD(&posEvent{Pos: pos, EAction: FireProgression}, DurationCloudProgression)
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
