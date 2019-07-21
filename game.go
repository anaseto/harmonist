package main

import "container/heap"

var Version string = "v0.2"

type game struct {
	Dungeon            *dungeon
	Player             *player
	Monsters           []*monster
	MonstersPosCache   []int // monster (dungeon index + 1) / no monster (0)
	Bands              []bandInfo
	Events             *eventQueue
	Ev                 event
	EventIndex         int
	Depth              int
	ExploredLevels     int
	DepthPlayerTurn    int
	Turn               int
	Highlight          map[position]bool // highlighted positions (e.g. targeted ray)
	Objects            objects
	Clouds             map[position]cloud
	MagicalBarriers    map[position]terrain
	GeneratedUniques   map[monsterBand]int
	GeneratedLore      map[int]bool
	GeneratedMagaras   []magara
	GeneratedCloaks    []item
	GeneratedAmulets   []item
	GenPlan            [MaxDepth + 1]genFlavour
	TerrainKnowledge   map[position]terrain
	ExclusionsMap      map[position]bool
	Noise              map[position]bool
	NoiseIllusion      map[position]bool
	LastMonsterKnownAt map[position]*monster
	MonsterLOS         map[position]bool
	MonsterTargLOS     map[position]bool
	Illuminated        []bool
	RaysCache          rayMap
	Resting            bool
	RestingTurns       int
	Autoexploring      bool
	DijkstraMapRebuild bool
	Targeting          position
	AutoTarget         position
	AutoDir            direction
	AutoHalt           bool
	AutoNext           bool
	DrawBuffer         []UICell
	drawBackBuffer     []UICell
	DrawLog            []drawFrame
	Log                []logEntry
	LogIndex           int
	LogNextTick        int
	InfoEntry          string
	Stats              stats
	Boredom            int
	Quit               bool
	Wizard             bool
	WizardMap          bool
	Version            string
	Places             places
	Params             startParams
	//Opts                startOpts
	ui                *gameui
	LiberatedShaedra  bool
	LiberatedArtifact bool
}

type startParams struct {
	Lore     map[int]bool
	Blocked  map[int]bool
	Special  []specialRoom
	Unstable map[int]bool
	Windows  map[int]bool
	Trees    map[int]bool
	Holes    map[int]bool
	Stones   map[int]bool
	Tables   map[int]bool
	NoMagara map[int]bool
}

type places struct {
	Shaedra  position
	Monolith position
	Marevor  position
	Artifact position
}

func (g *game) FreePassableCell() position {
	d := g.Dungeon
	count := 0
	for {
		count++
		if count > 1000 {
			panic("FreeCell")
		}
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		pos := position{x, y}
		c := d.Cell(pos)
		if !c.IsPassable() {
			continue
		}
		if g.Player != nil && g.Player.Pos == pos {
			continue
		}
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		return pos
	}
}

func (g *game) FreeCellForPlayer() position {
	// TODO: not used now, but could be for cases when you fall into the abyss
	center := position{DungeonWidth / 2, DungeonHeight / 2}
	bestpos := g.FreePassableCell()
	for i := 0; i < 5; i++ {
		pos := g.FreePassableCell()
		if pos.Distance(center) > bestpos.Distance(center) {
			bestpos = pos
		}
	}
	return bestpos
}

func (g *game) FreeCellForMonster() position {
	d := g.Dungeon
	count := 0
	for {
		count++
		if count > 1000 {
			panic("FreeCellForMonster")
		}
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		pos := position{x, y}
		c := d.Cell(pos)
		if !c.IsPassable() {
			continue
		}
		if g.Player != nil && g.Player.Pos.Distance(pos) < 8 {
			continue
		}
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		return pos
	}
}

func (g *game) FreeCellForBandMonster(pos position) position {
	count := 0
	for {
		count++
		if count > 1000 {
			return g.FreeCellForMonster()
		}
		neighbors := g.Dungeon.FreeNeighbors(pos)
		r := RandInt(len(neighbors))
		pos = neighbors[r]
		if g.Player != nil && g.Player.Pos.Distance(pos) < 8 {
			continue
		}
		mons := g.MonsterAt(pos)
		if mons.Exists() || !g.Dungeon.Cell(pos).IsPassable() {
			continue
		}
		return pos
	}
}

const MaxDepth = 11
const WinDepth = 8

const (
	DungeonHeight = 21
	DungeonWidth  = 79
	DungeonNCells = DungeonWidth * DungeonHeight
)

func (g *game) GenDungeon() {
	ml := AutomataCave
	//ml := RandomWalkTreeCave
	switch g.Depth {
	case 2, 6, 7:
		ml = RandomWalkCave
	case 4, 10, 11:
		ml = RandomWalkTreeCave
	}
	g.GenRoomTunnels(ml)
}

func (g *game) InitPlayer() {
	g.Player = &player{
		HP:      DefaultHealth,
		MP:      DefaultMPmax,
		Bananas: 1,
	}
	g.Player.Rays = rayMap{}
	g.Player.LOS = map[position]bool{}
	g.Player.Statuses = map[status]int{}
	g.Player.Expire = map[status]int{}
	g.Player.Magaras = []magara{
		NoMagara,
		NoMagara,
		NoMagara,
		NoMagara,
	}
	g.GeneratedMagaras = []magara{}
	g.Player.Magaras[0] = g.RandomMagara()
	g.GeneratedMagaras = append(g.GeneratedMagaras, g.Player.Magaras[0])
	g.Player.Inventory.Misc = MarevorMagara
	// Testing
	//g.Player.Magaras[2] = NoiseMagara
	//g.Player.Magaras[2] = SlowingMagara
	//g.Player.Magaras[2] = ConfusionMagara
}

type genFlavour int

const (
	GenNothing genFlavour = iota
	//GenWeapon
	GenAmulet
	GenCloak
)

func (g *game) InitFirstLevel() {
	g.Version = Version
	g.Depth++ // start at 1
	g.InitPlayer()
	g.AutoTarget = InvalidPos
	g.Targeting = InvalidPos
	g.Illuminated = make([]bool, DungeonNCells)
	g.RaysCache = rayMap{}
	g.GeneratedUniques = map[monsterBand]int{}
	g.GeneratedLore = map[int]bool{}
	g.Stats.KilledMons = map[monsterKind]int{}
	g.Stats.UsedMagaras = map[magara]int{}
	g.Stats.Achievements = map[achievement]int{}
	g.Stats.Lore = map[int]bool{}
	g.Stats.Statuses = map[status]int{}
	g.GenPlan = [MaxDepth + 1]genFlavour{
		1:  GenNothing,
		2:  GenCloak,
		3:  GenNothing,
		4:  GenAmulet,
		5:  GenNothing,
		6:  GenCloak,
		7:  GenNothing,
		8:  GenAmulet,
		9:  GenNothing,
		10: GenCloak,
		11: GenNothing,
	}
	g.Params.Lore = map[int]bool{}
	const nlore = 8
	for i := 0; i < nlore; i++ {
		j := RandInt(MaxDepth)
		if !g.Params.Lore[j] {
			g.Params.Lore[j] = true
		} else {
			i--
		}
	}
	g.Params.Blocked = map[int]bool{}
	if RandInt(10) > 0 {
		g.Params.Blocked[2+RandInt(WinDepth-2)] = true
	}
	if RandInt(10) == 0 {
		// a second one sometimes!
		g.Params.Blocked[2+RandInt(WinDepth-2)] = true
	}
	g.Params.Special = []specialRoom{
		noSpecialRoom, // unused (depth 0)
		noSpecialRoom,
		noSpecialRoom,
		roomMilfids,
		roomCelmists,
		roomVampires,
		roomHarpies,
		roomTreeMushrooms,
		roomShaedra,
		roomCelmists,
		roomMirrorSpecters,
		roomArtifact,
	}
	if RandInt(4) == 0 {
		g.Params.Special[10], g.Params.Special[5] = g.Params.Special[5], g.Params.Special[10]
	}
	if RandInt(4) == 0 {
		g.Params.Special[6], g.Params.Special[7] = g.Params.Special[7], g.Params.Special[6]
	}
	if RandInt(4) == 0 {
		g.Params.Special[3], g.Params.Special[4] = g.Params.Special[4], g.Params.Special[3]
	}
	g.Params.Unstable = map[int]bool{}
	if RandInt(MaxDepth) > MaxDepth/2 {
		g.Params.Unstable[2+RandInt(MaxDepth-1)] = true
		if RandInt(MaxDepth) == 0 {
			g.Params.Unstable[2+RandInt(MaxDepth-1)] = true
		}
	}
	g.Params.Windows = map[int]bool{}
	if RandInt(MaxDepth) > MaxDepth/2 {
		g.Params.Windows[2+RandInt(MaxDepth-1)] = true
		if RandInt(MaxDepth) == 0 {
			g.Params.Windows[2+RandInt(MaxDepth-1)] = true
		}
	}
	g.Params.Holes = map[int]bool{}
	if RandInt(MaxDepth) > MaxDepth/2 {
		g.Params.Holes[2+RandInt(MaxDepth-1)] = true
		if RandInt(MaxDepth) == 0 {
			g.Params.Holes[2+RandInt(MaxDepth-1)] = true
		}
	}
	g.Params.Trees = map[int]bool{}
	if RandInt(MaxDepth) > MaxDepth/2 {
		g.Params.Trees[2+RandInt(MaxDepth-1)] = true
		if RandInt(MaxDepth) == 0 {
			g.Params.Trees[2+RandInt(MaxDepth-1)] = true
		}
	}
	g.Params.Tables = map[int]bool{}
	if RandInt(MaxDepth) > MaxDepth/2 {
		g.Params.Tables[2+RandInt(MaxDepth-1)] = true
		if RandInt(MaxDepth) == 0 {
			g.Params.Tables[2+RandInt(MaxDepth-1)] = true
		}
	}
	g.Params.NoMagara = map[int]bool{}
	switch RandInt(3) {
	case 0:
		g.Params.NoMagara[4] = true
	case 1:
		g.Params.NoMagara[6] = true
	case 2:
		g.Params.NoMagara[8] = true
	}
	g.Params.Stones = map[int]bool{}
	if RandInt(MaxDepth) > MaxDepth/2 {
		g.Params.Stones[2+RandInt(MaxDepth-1)] = true
		if RandInt(MaxDepth) == 0 {
			g.Params.Stones[2+RandInt(MaxDepth-1)] = true
		}
	}
	permi := RandInt(WinDepth - 1)
	switch permi {
	case 0, 1, 2, 3:
		g.GenPlan[permi+1], g.GenPlan[permi+2] = g.GenPlan[permi+2], g.GenPlan[permi+1]
	}
	if RandInt(4) == 0 {
		g.GenPlan[6], g.GenPlan[7] = g.GenPlan[7], g.GenPlan[6]
	}
	if RandInt(4) == 0 {
		g.GenPlan[MaxDepth-1], g.GenPlan[MaxDepth] = g.GenPlan[MaxDepth], g.GenPlan[MaxDepth-1]
	}
}

func (g *game) InitLevelStructures() {
	g.MonstersPosCache = make([]int, DungeonNCells)
	g.Noise = map[position]bool{}
	g.TerrainKnowledge = map[position]terrain{}
	g.ExclusionsMap = map[position]bool{}
	g.MagicalBarriers = map[position]terrain{}
	g.LastMonsterKnownAt = map[position]*monster{}
	g.Objects.Magaras = map[position]magara{}
	g.Objects.Lore = map[position]int{}
	g.Objects.Items = map[position]item{}
	g.Objects.Scrolls = map[position]scroll{}
	g.Objects.Stairs = map[position]stair{}
	g.Objects.Bananas = make(map[position]bool, 2)
	g.Objects.Barrels = map[position]bool{}
	g.Objects.Lights = map[position]bool{}
	g.NoiseIllusion = map[position]bool{}
	g.Clouds = map[position]cloud{}
	g.MonsterLOS = map[position]bool{}
	g.Stats.AtNotablePos = map[position]bool{}
}

var Testing = false

func (g *game) InitLevel() {
	// Starting data
	if g.Depth == 0 {
		g.InitFirstLevel()
	} else if !Testing {
		g.ui.DrawLoading()
	}

	g.InitLevelStructures()

	// Dungeon terrain
	g.GenDungeon()

	// Events
	if g.Depth == 1 {
		g.StoryPrintf("Started with %s", g.Player.Magaras[0])
		g.Events = &eventQueue{}
		heap.Init(g.Events)
		g.PushEvent(&simpleEvent{ERank: 0, EAction: PlayerTurn})
	} else {
		g.CleanEvents()
	}
	for i := range g.Monsters {
		g.PushEvent(&monsterEvent{ERank: g.Turn + RandInt(10), EAction: MonsterTurn, NMons: i})
	}
	if g.Params.Unstable[g.Depth] {
		g.PrintStyled("Uncontrolled oric magic fills the air on this level.", logSpecial)
		for i := 0; i < 5; i++ {
			g.PushEvent(&cloudEvent{ERank: g.Turn + 50 + RandInt(100), EAction: ObstructionProgression})
		}
	}

	// initialize LOS
	if g.Depth == 1 {
		g.PrintStyled("► Press ? for help on keys or use the mouse and [buttons].", logSpecial)
	}
	if g.Depth == WinDepth {
		g.PrintStyled("Finally! Shaedra should be imprisoned somewhere around here.", logSpecial)
	} else if g.Depth == MaxDepth {
		g.PrintStyled("This the bottom floor, you now have to look for the artifact.", logSpecial)
	}
	g.ComputeLOS()
	g.MakeMonstersAware()
}

func (g *game) CleanEvents() {
	evq := &eventQueue{}
	for g.Events.Len() > 0 {
		iev := g.PopIEvent()
		switch iev.Event.(type) {
		case *monsterEvent:
		case *cloudEvent:
		default:
			heap.Push(evq, iev)
		}
	}
	g.Events = evq
}

func (g *game) StairsSlice() []position {
	// TODO: use cache?
	stairs := []position{}
	for i, c := range g.Dungeon.Cells {
		if c.T != StairCell || !c.Explored {
			continue
		}
		pos := idxtopos(i)
		stairs = append(stairs, pos)
	}
	return stairs
}

type descendstyle int

const (
	DescendNormal descendstyle = iota
	DescendJump
	DescendFall
)

func (g *game) Descend(style descendstyle) bool {
	g.LevelStats()
	if g.Stats.DUSpotted[g.Depth] < 3 {
		AchStealthNovice.Get(g)
	}
	if g.Depth >= 3 {
		if g.Stats.DRests[g.Depth] == 0 && g.Stats.DRests[g.Depth-1] == 0 {
			AchInsomniaNovice.Get(g)
		}
	}
	if g.Depth >= 5 {
		if g.Stats.DRests[g.Depth] == 0 && g.Stats.DRests[g.Depth-1] == 0 && g.Stats.DRests[g.Depth-2] == 0 &&
			g.Stats.DRests[g.Depth-3] == 0 {
			AchInsomniaInitiate.Get(g)
		}
	}
	if g.Depth >= 8 {
		if g.Stats.DRests[g.Depth] == 0 && g.Stats.DRests[g.Depth-1] == 0 && g.Stats.DRests[g.Depth-2] == 0 &&
			g.Stats.DRests[g.Depth-3] == 0 && g.Stats.DRests[g.Depth-4] == 0 && g.Stats.DRests[g.Depth-5] == 0 {
			AchInsomniaMaster.Get(g)
		}
	}
	if g.Depth >= 3 {
		if g.Stats.DMagaraUses[g.Depth] == 0 && g.Stats.DMagaraUses[g.Depth-1] == 0 {
			AchAntimagicNovice.Get(g)
		}
	}
	if g.Depth >= 5 {
		if g.Stats.DMagaraUses[g.Depth] == 0 && g.Stats.DMagaraUses[g.Depth-1] == 0 && g.Stats.DMagaraUses[g.Depth-2] == 0 &&
			g.Stats.DMagaraUses[g.Depth-3] == 0 {
			AchAntimagicInitiate.Get(g)
		}
	}
	if g.Depth >= 8 {
		if g.Stats.DMagaraUses[g.Depth] == 0 && g.Stats.DMagaraUses[g.Depth-1] == 0 && g.Stats.DMagaraUses[g.Depth-2] == 0 &&
			g.Stats.DMagaraUses[g.Depth-3] == 0 && g.Stats.DMagaraUses[g.Depth-4] == 0 && g.Stats.DMagaraUses[g.Depth-5] == 0 {
			AchAntimagicMaster.Get(g)
		}
	}
	if g.Depth >= 5 {
		if g.Stats.DUSpotted[g.Depth] < 3 && g.Stats.DSpotted[g.Depth-1] < 3 && g.Stats.DSpotted[g.Depth-2] < 3 {
			AchStealthInitiate.Get(g)
		}
	}
	if g.Depth >= 8 {
		if g.Stats.DUSpotted[g.Depth] < 3 && g.Stats.DUSpotted[g.Depth-1] < 3 && g.Stats.DSpotted[g.Depth-2] < 3 &&
			g.Stats.DSpotted[g.Depth-3] < 3 {
			AchStealthMaster.Get(g)
		}
	}
	c := g.Dungeon.Cell(g.Player.Pos)
	if c.T == StairCell && g.Objects.Stairs[g.Player.Pos] == WinStair {
		g.StoryPrint("Escaped!")
		g.ExploredLevels = g.Depth
		g.Depth = -1
		return true
	}
	if style != DescendNormal {
		// TODO: add animation?
		g.Print("You fall into the abyss. It hurts!")
		g.StoryPrint("Fell into the abyss")
	} else {
		g.Print("You descend deeper in the dungeon.")
		g.StoryPrint("Descended stairs")
	}
	g.Depth++
	g.DepthPlayerTurn = 0
	g.Boredom = 0
	if style != DescendFall {
		g.PushEvent(&simpleEvent{ERank: g.Ev.Rank(), EAction: PlayerTurn})
	}
	g.InitLevel()
	g.Save()
	return false
}

func (g *game) WizardMode() {
	g.Wizard = true
	g.PrintStyled("You are now in wizard mode and cannot obtain winner status.", logSpecial)
}

func (g *game) ApplyRest() {
	g.Player.HP = g.Player.HPMax()
	g.Player.HPbonus = 0
	g.Player.MP = g.Player.MPMax()
	g.Stats.Rest++
	g.Stats.DRests[g.Depth]++
	g.PrintStyled("You feel fresh again after eating banana and sleeping.", logStatusEnd)
	g.StoryPrintf("Rested in barrel (bananas: %d)", g.Player.Bananas)
	if g.Stats.Rest == 10 {
		AchSleepy.Get(g)
	}
}

func (g *game) AutoPlayer(ev event) bool {
	if g.Resting {
		const enoughRestTurns = 25
		if g.RestingTurns < enoughRestTurns {
			g.WaitTurn(ev)
			g.RestingTurns++
			return true
		}
		if g.RestingTurns >= enoughRestTurns {
			g.ApplyRest()
		}
		g.Resting = false
	} else if g.Autoexploring {
		if g.ui.ExploreStep() {
			g.AutoHalt = true
			g.Print("Stopping, then.")
		}
		switch {
		case g.AutoHalt:
			// stop exploring
		default:
			var n *position
			var finished bool
			if g.DijkstraMapRebuild {
				if g.AllExplored() {
					g.Print("You finished exploring.")
					break
				}
				sources := g.AutoexploreSources()
				g.BuildAutoexploreMap(sources)
			}
			n, finished = g.NextAuto()
			if finished {
				n = nil
			}
			if finished && g.AllExplored() {
				g.Print("You finished exploring.")
			} else if n == nil {
				g.Print("You could not safely reach some places.")
			}
			if n != nil {
				err := g.MovePlayer(*n, ev)
				if err != nil {
					g.Print(err.Error())
					break
				}
				return true
			}
		}
		g.Autoexploring = false
	} else if g.AutoTarget.valid() {
		if !g.ui.ExploreStep() && g.MoveToTarget(ev) {
			return true
		} else {
			g.AutoTarget = InvalidPos
		}
	} else if g.AutoDir != NoDir {
		if !g.ui.ExploreStep() && g.AutoToDir(ev) {
			return true
		} else {
			g.AutoDir = NoDir
		}
	}
	return false
}

func (g *game) EventLoop() {
loop:
	for {
		if g.Player.HP <= 0 {
			if g.Wizard {
				g.Player.HP = g.Player.HPMax()
			} else {
				g.LevelStats()
				err := g.RemoveSaveFile()
				if err != nil {
					g.PrintfStyled("Error removing save file: %v", logError, err.Error())
				}
				g.ui.Death()
				break loop
			}
		}
		if g.Events.Len() == 0 {
			break loop
		}
		ev := g.PopIEvent().Event
		g.Turn = ev.Rank()
		g.Ev = ev
		ev.Action(g)
		if g.AutoNext {
			continue loop
		}
		if g.Quit {
			break loop
		}
	}
}
