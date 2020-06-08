package main

type stats struct {
	Story             []string
	Killed            int
	KilledMons        map[monsterKind]int
	Moves             int
	Waits             int
	Jumps             int
	WallJumps         int
	ReceivedHits      int
	Dodges            int
	MagarasUsed       int
	DMagaraUses       [MaxDepth + 1]int
	UsedStones        int
	Damage            int
	DDamage           [MaxDepth + 1]int
	DExplPerc         [MaxDepth + 1]int
	DSleepingPerc     [MaxDepth + 1]int
	DKilledPerc       [MaxDepth + 1]int
	Burns             int
	Digs              int
	Rest              int
	DRests            [MaxDepth + 1]int
	Turns             int
	TWounded          int
	TMWounded         int
	TMonsLOS          int
	NSpotted          int
	NUSpotted         int
	DSpotted          [MaxDepth + 1]int
	DUSpotted         [MaxDepth + 1]int
	DUSpottedPerc     [MaxDepth + 1]int
	Achievements      map[achievement]int
	AtNotablePos      map[position]bool
	HarmonicMagUse    int
	OricMagUse        int
	FireUse           int
	DestructionUse    int
	OricTelUse        int
	ClimbedTree       int
	TableHides        int
	HoledWallsCrawled int
	DoorsOpened       int
	BarrelHides       int
	Extinguishments   int
	Lore              map[int]bool
	Statuses          map[status]int
	StolenBananas     int
	TimesPushed       int
	TimesBlinked      int
	TimesBlocked      int
}

func (g *game) TurnStats() {
	g.Stats.Turns++
	g.DepthPlayerTurn++
	if g.Player.HP < g.Player.HPMax() {
		g.Stats.TWounded++
	}
	if g.MonsterInLOS() != nil {
		g.Stats.TMonsLOS++
		if g.Player.HP < g.Player.HPMax() {
			g.Stats.TMWounded++
		}
	}
}

func (g *game) LevelStats() {
	free := 0
	exp := 0
	for _, c := range g.Dungeon.Cells {
		if c.IsWall() || c.T == ChasmCell {
			continue
		}
		free++
		if c.Explored {
			exp++
		}
	}
	g.Stats.DExplPerc[g.Depth] = exp * 100 / free
	//g.Stats.DBurns[g.Depth] = g.Stats.CurBurns // XXX to avoid little dump info leak
	nmons := len(g.Monsters)
	kmons := 0
	smons := 0
	for _, mons := range g.Monsters {
		if !mons.Exists() {
			kmons++
			continue
		}
		if mons.State == Resting {
			smons++
		}
	}
	g.Stats.DSleepingPerc[g.Depth] = smons * 100 / nmons
	g.Stats.DKilledPerc[g.Depth] = kmons * 100 / nmons
	g.Stats.DUSpottedPerc[g.Depth] = g.Stats.DUSpotted[g.Depth] * 100 / nmons
}

type achievement string

// Achievements.
const (
	NoAchievement achievement = "Pitiful Death"
)

func (ach achievement) Get(g *game) {
	if g.Stats.Achievements[ach] == 0 {
		g.Stats.Achievements[ach] = g.Turn
		g.PrintfStyled("Achievement: %s.", logSpecial, ach)
		g.StoryPrintf("Achievement: %s", ach)
	}
}

func (t terrain) ReachNotable() bool {
	switch t {
	case TreeCell, TableCell, HoledWallCell, DoorCell, BarrelCell:
		return true
	default:
		return false
	}
}

func (pos position) Reach(g *game) {
	if g.Stats.AtNotablePos[pos] {
		return
	}
	g.Stats.AtNotablePos[pos] = true
	c := g.Dungeon.Cell(pos)
	switch c.T {
	case TreeCell:
		g.Stats.ClimbedTree++
	case TableCell:
		g.Stats.TableHides++
	case HoledWallCell:
		g.Stats.HoledWallsCrawled++
	case DoorCell:
		g.Stats.DoorsOpened++
	case BarrelCell:
		g.Stats.BarrelHides++
	}
}
