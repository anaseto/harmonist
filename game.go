package main

import (
	"log"
	"math/rand"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/paths"
	"github.com/anaseto/gruid/rl"
)

var Version string = "v0.4-dev-5"

// game contains the game logic's state, without ui stuff. Everything could be
// in the model struct instead, with only the game logic's fiend exported, as
// some game functions need the model anyway (like animations), but this allows
// to differenciate a bit things that are mainly game-logic from the stuff that
// is more about ui.
type game struct {
	Dungeon               *dungeon
	Player                *player
	Monsters              []*monster
	MonstersPosCache      []int // monster (dungeon index + 1) / no monster (0)
	Bands                 []bandInfo
	Events                *rl.EventQueue
	EventIndex            int
	Depth                 int
	ExploredLevels        int
	DepthPlayerTurn       int
	Turn                  int
	Highlight             map[gruid.Point]bool // highlighted positions (e.g. targeted ray)
	Objects               objects
	Clouds                map[gruid.Point]cloud
	MagicalBarriers       map[gruid.Point]cell
	GeneratedLore         map[int]bool
	GeneratedMagaras      []magaraKind
	GeneratedCloaks       []item
	GeneratedAmulets      []item
	GenPlan               [MaxDepth + 1]genFlavour
	TerrainKnowledge      map[gruid.Point]cell
	ExclusionsMap         map[gruid.Point]bool
	Noise                 map[gruid.Point]bool
	NoiseIllusion         map[gruid.Point]bool
	LastMonsterKnownAt    map[gruid.Point]int
	MonsterLOS            map[gruid.Point]bool
	MonsterTargLOS        map[gruid.Point]bool
	LightFOV              *rl.FOV
	RaysCache             rayMap
	Resting               bool
	RestingTurns          int
	Autoexploring         bool
	AutoexploreMapRebuild bool
	AutoTarget            gruid.Point
	AutoDir               gruid.Point
	autoDirNeighbors      dirNeighbors
	autoDirChanged        bool
	AutoHalt              bool
	Log                   []logEntry
	LogIndex              int
	LogNextTick           int
	InfoEntry             string
	Stats                 stats
	Wizard                bool
	WizardMode            wizardMode
	Version               string
	Places                places
	Params                startParams
	//Opts                startOpts
	md                *model // needed for animations and a few more cases
	LiberatedShaedra  bool
	LiberatedArtifact bool
	PlayerAgain       bool
	mfov              *rl.FOV
	PR                *paths.PathRange
	PRauto            *paths.PathRange
	autosources       []gruid.Point // cache
	nbs               paths.Neighbors
}

type specialEvent int

const (
	NormalLevel specialEvent = iota
	UnstableLevel
	EarthquakeLevel
	MistLevel
)

const spEvMax = int(MistLevel)

type startParams struct {
	Lore         map[int]bool
	Blocked      map[int]bool
	Special      []specialRoom
	Event        map[int]specialEvent
	Windows      map[int]bool
	Trees        map[int]bool
	Holes        map[int]bool
	Stones       map[int]bool
	Tables       map[int]bool
	NoMagara     map[int]bool
	FakeStair    map[int]bool
	ExtraBanana  map[int]int
	HealthPotion map[int]bool
	MappingStone map[int]bool
	CrazyImp     int
}

type wizardMode int

const (
	WizardNormal wizardMode = iota
	WizardMap
	WizardSeeAll
)

func (g *game) FreePassableCell() gruid.Point {
	d := g.Dungeon
	count := 0
	for {
		count++
		if count > maxIterations {
			panic("FreePassableCell")
		}
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		p := gruid.Point{x, y}
		c := d.Cell(p)
		if !c.IsPassable() {
			continue
		}
		if g.Player != nil && g.Player.P == p {
			continue
		}
		mons := g.MonsterAt(p)
		if mons.Exists() {
			continue
		}
		return p
	}
}

const MaxDepth = 11
const WinDepth = 8

const (
	DungeonHeight = 21
	DungeonWidth  = 80
	DungeonNCells = DungeonWidth * DungeonHeight
)

func (g *game) GenDungeon() {
	ml := AutomataCave
	switch g.Depth {
	case 2, 6, 7:
		ml = RandomWalkCave
		if RandInt(3) == 0 {
			ml = NaturalCave
		}
	case 4, 10, 11:
		ml = RandomWalkTreeCave
		if RandInt(4) == 0 && g.Depth < 11 {
			ml = RandomSmallWalkCaveUrbanised
		} else if g.Depth == 11 && RandInt(2) == 0 {
			ml = RandomSmallWalkCaveUrbanised
		}
	case 9:
		switch RandInt(4) {
		case 0:
			ml = NaturalCave
		case 1:
			ml = RandomWalkCave
		}
	default:
		if RandInt(10) == 0 {
			ml = RandomSmallWalkCaveUrbanised
		} else if RandInt(10) == 0 {
			ml = NaturalCave
		}
	}
	g.GenRoomTunnels(ml)
}

func (g *game) InitPlayer() {
	g.Player = &player{
		HP:      DefaultHealth,
		MP:      DefaultMPmax,
		Bananas: 1,
	}
	g.Player.LOS = map[gruid.Point]bool{}
	g.Player.Statuses = map[status]int{}
	g.Player.Expire = map[status]int{}
	g.Player.Magaras = []magara{
		magara{},
		magara{},
		magara{},
		magara{},
	}
	g.GeneratedMagaras = []magaraKind{}
	g.Player.Magaras[0] = g.RandomStartingMagara()
	g.GeneratedMagaras = append(g.GeneratedMagaras, g.Player.Magaras[0].Kind)
	g.Player.Inventory.Misc = MarevorMagara
	g.Player.FOV = rl.NewFOV(visionRange(g.Player.P, TreeRange))
	// Testing
	//g.Player.Magaras[1] = magara{Kind: DispersalMagara, Charges: 10}
	//g.Player.Magaras[2] = magara{Kind: DelayedOricExplosionMagara, Charges: 10}
	//g.Player.Magaras[2] = ConfusionMagara
}

type genFlavour int

const (
	GenNothing genFlavour = iota
	//GenWeapon
	GenAmulet
	GenCloak
)

func PutRandomLevels(m map[int]bool, n int) {
	for i := 0; i < n; i++ {
		j := 1 + RandInt(MaxDepth)
		if !m[j] {
			m[j] = true
		} else {
			i--
		}
	}
}

func (g *game) InitFirstLevel() {
	g.Version = Version
	g.Depth++ // start at 1
	g.InitPlayer()
	g.AutoTarget = invalidPos
	g.RaysCache = rayMap{}
	g.GeneratedLore = map[int]bool{}
	g.Stats.KilledMons = map[monsterKind]int{}
	g.Stats.UsedMagaras = map[magaraKind]int{}
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
	PutRandomLevels(g.Params.Lore, 8)
	g.Params.HealthPotion = map[int]bool{}
	PutRandomLevels(g.Params.HealthPotion, 5)
	g.Params.MappingStone = map[int]bool{}
	PutRandomLevels(g.Params.MappingStone, 3)
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
	if RandInt(2) == 0 {
		g.Params.Special[5] = roomNixes
	}
	if RandInt(4) == 0 {
		if g.Params.Special[5] == roomNixes {
			g.Params.Special[9] = roomVampires
		} else {
			g.Params.Special[9] = roomNixes
		}
	}
	if RandInt(4) == 0 {
		if RandInt(2) == 0 {
			g.Params.Special[3] = roomFrogs
		} else {
			g.Params.Special[7] = roomFrogs
		}
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
	g.Params.Event = map[int]specialEvent{}
	for i := 0; i < 2; i++ {
		g.Params.Event[2+5*i+RandInt(5)] = specialEvent(1 + RandInt(spEvMax))
	}
	g.Params.Event[2+RandInt(MaxDepth-1)] = NormalLevel
	g.Params.FakeStair = map[int]bool{}
	if RandInt(MaxDepth) > 0 {
		g.Params.FakeStair[2+RandInt(MaxDepth-2)] = true
		if RandInt(MaxDepth) > MaxDepth/2 {
			g.Params.FakeStair[2+RandInt(MaxDepth-2)] = true
			if RandInt(MaxDepth) == 0 {
				g.Params.FakeStair[2+RandInt(MaxDepth-2)] = true
			}
		}
	}
	g.Params.ExtraBanana = map[int]int{}
	for i := 0; i < 2; i++ {
		g.Params.ExtraBanana[1+5*i+RandInt(5)]++
	}
	for i := 0; i < 2; i++ {
		g.Params.ExtraBanana[1+5*i+RandInt(5)]--
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
	g.Params.NoMagara[WinDepth] = true
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
	g.Params.CrazyImp = 2 + RandInt(MaxDepth-2)
	g.PR = paths.NewPathRange(gruid.NewRange(0, 0, DungeonWidth, DungeonHeight))
	g.PRauto = paths.NewPathRange(gruid.NewRange(0, 0, DungeonWidth, DungeonHeight))
}

func (g *game) InitLevelStructures() {
	g.MonstersPosCache = make([]int, DungeonNCells)
	g.Noise = map[gruid.Point]bool{}
	g.TerrainKnowledge = map[gruid.Point]cell{}
	g.ExclusionsMap = map[gruid.Point]bool{}
	g.MagicalBarriers = map[gruid.Point]cell{}
	g.LastMonsterKnownAt = map[gruid.Point]int{}
	g.Objects.Magaras = map[gruid.Point]magara{}
	g.Objects.Lore = map[gruid.Point]int{}
	g.Objects.Items = map[gruid.Point]item{}
	g.Objects.Scrolls = map[gruid.Point]scroll{}
	g.Objects.Stairs = map[gruid.Point]stair{}
	g.Objects.Bananas = make(map[gruid.Point]bool, 2)
	g.Objects.Barrels = map[gruid.Point]bool{}
	g.Objects.Lights = map[gruid.Point]bool{}
	g.Objects.FakeStairs = map[gruid.Point]bool{}
	g.Objects.Potions = map[gruid.Point]potion{}
	g.NoiseIllusion = map[gruid.Point]bool{}
	g.Clouds = map[gruid.Point]cloud{}
	g.MonsterLOS = map[gruid.Point]bool{}
	g.Stats.AtNotablePos = map[gruid.Point]bool{}
}

var Testing = false

func (g *game) InitLevel() {
	// Starting data
	if g.Depth == 0 {
		g.InitFirstLevel()
	}

	g.InitLevelStructures()

	// Dungeon terrain
	g.GenDungeon()

	// Events
	if g.Depth == 1 {
		g.StoryPrintf("Started with %s", g.Player.Magaras[0])
		g.Events = rl.NewEventQueue()
		//g.PushEvent(&simpleEvent{ERank: 0, EAction: PlayerTurn})
	} else {
		g.CleanEvents()
		for st := range g.Player.Statuses {
			if st.Clean() {
				g.Player.Statuses[st] = 0
			}
		}
	}
	monsters := make([]*monster, len(g.Monsters))
	copy(monsters, g.Monsters)
	rand.Shuffle(len(monsters), func(i, j int) {
		monsters[i], monsters[j] = monsters[j], monsters[i]
	})
	for _, m := range monsters {
		g.PushEvent(&monsterTurnEvent{Index: m.Index}, g.Turn)
	}
	switch g.Params.Event[g.Depth] {
	case UnstableLevel:
		g.PrintStyled("Uncontrolled oric magic fills the air on this level.", logSpecial)
		g.StoryPrint("Special event: magically unstable level")
		for i := 0; i < 7; i++ {
			g.PushEvent(&posEvent{Action: ObstructionProgression},
				g.Turn+DurationObstructionProgression+RandInt(DurationObstructionProgression/2))
		}
	case MistLevel:
		g.PrintStyled("The air seems dense on this level.", logSpecial)
		g.StoryPrint("Special event: mist level")
		for i := 0; i < 20; i++ {
			g.PushEvent(&posEvent{Action: MistProgression},
				g.Turn+DurationMistProgression+RandInt(DurationMistProgression/2))
		}
	case EarthquakeLevel:
		g.PushEvent(&posEvent{P: gruid.Point{DungeonWidth/2 - 15 + RandInt(30), DungeonHeight/2 - 5 + RandInt(10)}, Action: Earthquake},
			g.Turn+10+RandInt(50))

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
	g.ComputeMonsterLOS()
	if g.md != nil { // disable when testing
		g.md.updateStatusInfo()
	}
}

func (g *game) CleanEvents() {
	g.Events.Filter(func(ev rl.Event) bool {
		switch ev.(type) {
		case *monsterTurnEvent, *posEvent, *monsterStatusEvent, *playerEvent:
			return false
		default:
			// keep player statuses events
			return true
		}
	})
	// finish current turn's other effects (like status progression)
	turn := g.Turn
	for !g.Events.Empty() {
		ev, r := g.Events.PopR()
		if r == turn {
			e := ev.(event)
			e.Handle(g)
			continue
		}
		g.Events.PushFirst(ev, r)
		break
	}
	g.Turn++
}

func (g *game) StairsSlice() []gruid.Point {
	stairs := []gruid.Point{}
	it := g.Dungeon.Grid.Iterator()
	for it.Next() {
		c := cell(it.Cell())
		if (terrain(c) != StairCell && terrain(c) != FakeStairCell) || !explored(c) {
			continue
		}
		stairs = append(stairs, it.P())
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
	c := g.Dungeon.Cell(g.Player.P)
	if terrain(c) == StairCell && g.Objects.Stairs[g.Player.P] == WinStair {
		g.StoryPrint("Escaped!")
		g.ExploredLevels = g.Depth
		g.Depth = -1
		return true
	}
	if style != DescendNormal {
		g.md.AbyssFallAnimation()
		g.PrintStyled("You fall into the abyss. It hurts!", logDamage)
		g.StoryPrint("Fell into the abyss")
	} else {
		g.Print("You descend deeper in the dungeon.")
		g.StoryPrint("Descended stairs")
	}
	g.Depth++
	g.DepthPlayerTurn = 0
	g.InitLevel()
	g.Save()
	return false
}

func (g *game) EnterWizardMode() {
	g.Wizard = true
	g.PrintStyled("Wizard mode activated: winner status disabled.", logSpecial)
	g.StoryPrint("Entered wizard mode.")
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

func (g *game) AutoPlayer() bool {
	switch {
	case g.Resting:
		const enoughRestTurns = 25
		if g.RestingTurns < enoughRestTurns {
			g.RestingTurns++
			return true
		}
		if g.RestingTurns >= enoughRestTurns {
			g.ApplyRest()
		}
		g.Resting = false
	case g.Autoexploring:
		switch {
		case g.AutoHalt:
			// stop exploring
		default:
			var n *gruid.Point
			var finished bool
			if g.AutoexploreMapRebuild {
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
				again, err := g.PlayerBump(*n)
				if err != nil {
					g.Print(err.Error())
					break
				}
				return !again
			}
		}
		g.Autoexploring = false
	case valid(g.AutoTarget):
		if g.MoveToTarget() {
			return true
		}
		g.AutoTarget = invalidPos
	case g.AutoDir != ZP:
		if g.AutoToDir() {
			return true
		}
		g.AutoDir = ZP
	}
	return false
}

func (g *game) Died() bool {
	if g.Player.HP <= 0 {
		if g.Wizard {
			g.Player.HP = g.Player.HPMax()
			g.PrintStyled("You died.", logSpecial)
			g.StoryPrint("You died (wizard mode)")
		} else {
			g.LevelStats()
			return true
		}
	}
	return false
}

type msgAuto int

func (g *game) EndTurn() {
	g.Events.Push(endTurnAction, g.Turn+DurationTurn)
	for {
		if g.Died() {
			return
		}
		if g.Events.Empty() {
			return
		}
		ev, r := g.Events.PopR()
		g.Turn = r
		switch ev := ev.(type) {
		case endTurnEvent:
			return
		case event:
			ev.Handle(g)
		default:
			log.Print("bad event: %v", ev)
		}
	}
}

func (g *game) checks() {
	if !Testing {
		return
	}
	for _, m := range g.Monsters {
		mons := g.MonsterAt(m.P)
		if !mons.Exists() && m.Exists() {
			log.Printf("does not exist")
			continue
		}
		if mons != m {
			log.Printf("bad monster: %v vs %v", mons.Index, m.Index)
		}
	}
}
