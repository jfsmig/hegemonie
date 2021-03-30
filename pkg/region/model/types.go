// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"sync"
)

const (
	ResourceActions   = iota
	ResourceGold      = iota
	ResourceCereal    = iota
	ResourceLivestock = iota
	ResourceStone     = iota
	ResourceWood      = iota
	ResourceMax       = iota
)

const (
	// Disband the Army and transfer its whole content to the local City
	// If there is no local City at the current position, the content of
	// the Army is lost.
	// Argument: ignored
	CmdCityDisband = "disband"

	// Like CmdMove but the command doesn't expire
	// Argument: ignored
	CmdWait = "wait"

	// Do nothing. Useful for waypoints
	// Argument: ActionArgMove
	CmdMove = "move"

	// Start a fight or join a running fight on the side of the attackers
	// Argument: ActionArgAssault
	CmdCityAttack = "attack"

	// Join a running fight on the side of the defenders, or Watch the City if
	// Argument: ignored
	CmdCityDefend = "defend"
)

const (
	CityStateIdle = iota
	CityStateAssigned
	CityStateConfigured
	CityStateMgmtConfigured
	CityStateReady
	CityStateActive
	CityStateAuto
	CityStateSuspended
)

type World struct {
	// Core configuration common to all the Regions of the current World.
	Config Configuration

	// Data that define the current World, shared among all the Regions.
	Definitions DefinitionsBase

	// All the regions that have been instantiated into the current World.
	Regions SetOfRegions

	// Interface to the notification system
	notifier Notifier

	// Interface to the map
	mapView MapClient

	rw sync.RWMutex
}

type Configuration struct {
	// Ratio applied to the production of resources that is applied for each
	// Massacre underwent by any city. It only impacts the production of the City itself.
	MassacreImpact float64

	// Should resource transfers happen instantly or should an actual transport
	// be emitted by the sender? Set to `true` for an instant transfer or to
	// `false` for a transport.
	InstantTransfers bool

	// Permanent bonus to the Popularity when a City creates an Army
	PopBonusArmyCreate int64

	// Permanent bonus to the Popularity when a City disband an Army
	PopBonusArmyDisband int64

	// Transient bonus to the Popularity of a City for each of its live Army
	PopBonusArmyAlive int64

	// Default Overlord rate: percentage of the production of a City that is
	// taxed by its Overlord
	RateOverlord float64

	// A city pattern is picked randomly among this set when a city is created.
	// So the configuration of the world may introduce a variation between
	// Cities
	CityPatterns []City
}

type DefinitionsBase struct {
	// All the possible Units that can be trained or hired in a World
	// IMMUTABLE: Only read accesses allowed.
	Units SetOfUnitTypes
	// All the possible Buildings that can be spanwed in Cities of the current World
	// IMMUTABLE: Only read accesses allowed.
	Buildings SetOfBuildingTypes
	// All the possible Knowledge that can be learned in Cities of the current World
	// IMMUTABLE: Only read accesses allowed.
	Knowledges SetOfKnowledgeTypes
}

type Region struct {
	// Unique name of the region
	Name string

	// Identifier of the map in use for the current Region
	MapName string

	// All the cities present on the Region
	Cities SetOfCities

	// Fights currently happening. The armies involved in the Fight are owned
	// By the Fight and do not appear in the "Armies" field.
	Fights SetOfFights

	Models SetOfTemplates

	// Back-pointer to the World the current Region belongs to.
	world *World
}

type Resources [ResourceMax]uint64

type ResourcesIncrement [ResourceMax]int64

type ResourcesMultiplier [ResourceMax]float64

type Artifact struct {
	// UUID
	ID string `json:"id"`
	// UUID
	Type string `json:"type"`
	// UUID
	Name    string `json:"name"`
	Visible bool   `json:"visible,omitempty"`
}

type ResourceModifiers struct {
	Mult ResourcesMultiplier
	Plus ResourcesIncrement
}

type CityProduction struct {
	Base      Resources
	Knowledge ResourceModifiers
	Buildings ResourceModifiers
	Actual    Resources
}

type CityStock struct {
	Base      Resources
	Knowledge ResourceModifiers
	Buildings ResourceModifiers
	Actual    Resources

	Usage Resources
}

// CityActivityCounters gathers counters that depict the activity of the City
// and that are continuously updated by the regions service.
type CityActivityCounters struct {
	ResourceProduced Resources
	ResourceSent     Resources
	ResourceReceived Resources
	TaxSent          Resources
	TaxReceived      Resources

	Moves        uint64
	FightsJoined uint64
	FightsLeft   uint64
	FightsWon    uint64
	FightsLost   uint64
	UnitsRaised  uint64
	UnitsLost    uint64
}

// CityStats gathers gauges and counters that give a hint on the activity of
// a City. The counters, that are centralized in the Activity field, are
// extracted as-is from a struct that is continuously updated while the
// gauges are mostly computed.
type CityStats struct {
	// Gauges
	StockCapacity  Resources
	StockUsage     Resources
	ScoreBuildings uint64
	ScoreKnowledge uint64
	ScoreMilitary  uint64
	// Counters
	Activity CityActivityCounters
}

type KnowledgeType struct {
	ID    uint64 `json:"Id"`
	Name  string `json:"Name"`
	Ticks uint32 `json:",omitempty"`

	// Transient bonus of Popularity, when the Knowledge is present
	PopBonus int64

	// Permanent bonus of Popularity when the Knowledge is achieved
	PopBonusLearn int64

	// Permanent bonus of Popularity (to the owner) when the Knowledge is stolen
	PopBonusStealVictim int64

	// Permanent bonus of Popularity (to the robber) when the Knowledge is stolen
	PopBonusStealActor int64

	// Impat of the current Building on the total storage capacity of the City.
	Stock ResourceModifiers

	// Increment of resources produced by this building.
	Prod ResourceModifiers

	// Amount of resources spent when the City starts learning this knowledge
	Cost0 Resources

	// Amount of resources spent to advance of one tick
	Cost Resources

	// All the knowledges that are required to start the current Knowledge
	// (this is an AND, not an OR)
	Requires []uint64

	// All the knowledge that are forbidden by the current knowledge
	Conflicts []uint64
}

type Knowledge struct {
	ID    string `json:"Id"`
	Type  uint64
	Ticks uint32 `json:",omitempty"`
}

type BuildingType struct {
	// Unique ID of the BuildingType
	ID uint64 `json:"Id"`

	// Display name of the current BuildingType
	Name string

	// How many ticks for the construction
	Ticks uint32 `json:",omitempty"`

	// How much does the production cost to start the the building process.
	Cost0 Resources

	// How much does the production cost at each tick.
	Cost Resources

	// Has the building to be unique a the City
	MultipleAllowed bool `json:",omitempty"`

	// Amount of total popularity required to start the construction of the building
	PopRequired int64

	// Transient bonus of Popularity, when the Building is alive
	PopBonus int64

	// Permanent bonus of Popularity given when the Building is achieved
	PopBonusBuild int64

	// Permanent bonus of Popularity given to the owner of the Building when it is destroyed.
	PopBonusFall int64

	// Permanent bonus of Popularity given to the destroyer of the Building
	PopBonusDestroy int64

	// Permanent bonus of Popularity given to the owner of the Building when it is dismantled.
	PopBonusDismantle int64

	// Impat of the current Building on the total storage capacity of the City.
	Stock ResourceModifiers

	// Increment of resources produced by this building.
	Prod ResourceModifiers

	// A set of KnowledgeType ID that must all be present in a City to let that City start
	// this kind of building.
	Requires []uint64

	// A set of KnowledgeType ID that must all be absent in a City to let that City start
	// this kind of building.
	Conflicts []uint64
}

type Building struct {
	// The unique ID of the current Building
	ID string `json:"Id"`

	// The unique ID of the BuildingType associated to the current Building
	Type uint64

	// How many construction rounds remain before the building's achievement
	Ticks uint32 `json:",omitempty"`
}

// City is the central point of a game instance.
// It is also the locking granularity for most actions.
type City struct {
	// The unique ID of the current City/
	// It is identical to the ID of the location (Vertex) on the Map.
	ID uint64 `json:"Id"`

	// The current FSM state, whose value is valid if chosen among the macros in CityState*
	State uint32 `json:"state"`

	// The unique ID of the User in charge of the City.
	Owner string

	// The unique ID of a City who is the boss of the current City.
	// Used for resources production computations.
	Overlord uint64

	// Ratio of the produced resources automatically sent to the Overlord City.
	TaxRate ResourcesMultiplier

	Assault *Fight `json:",omitempty"`

	// The display name of the current City
	Name string

	// From Lawful (<0) to Chaotic (>0) (0 for neutral)
	Chaotic int32

	// From Bad (<0) to Good (>0) (0 for neutral)
	Alignment int32

	// Race, Tribe, whatever (0 for unset)
	EthnicGroup uint32

	// Major political orientation (0 for none)
	PoliticalGroup uint32

	// God, Pantheon, Philosophy (0 for unset)
	Cult uint32

	// Permanent Popularity of the current City
	// The total value is the permanent value plus several "transient" bonus
	PermanentPopularity int64

	// Permanent Health of the current City. In other words, how it resists
	// to diseases, propagates pandemies, etc.
	// Higher is better.
	// The total value is the permanent value plus several "transient" bonus
	PermanentHealth int64

	// Permanent Intelligence of the current City. In other words, how to is
	// able to spy other and resist to their intelligence assaults.
	// Higher is better.
	// The total value is the permanent value plus several "transient" bonus
	PermanentIntelligence int64

	// Number of massacres the current City undergo.
	// It takes one production turn to recover one Massacre.
	TicksMassacres uint32 `json:",omitempty"`

	// Resources stock owned by the current City
	Stock Resources

	// Maximum amounts of each resources that might be stored in the town hall
	// of the city. That limit doesn't consider the modifiers.
	StockCapacity Resources

	// Resources produced each round by the City, before the enforcing of
	// Production Boosts ans Production Multipliers
	Production Resources

	Knowledges SetOfKnowledges

	Buildings SetOfBuildings

	// Units directly defending the current City
	Units SetOfUnits

	// Armies under the responsibility of the current City
	Armies SetOfArmies

	// Artifacts currently placed in the City.
	Artifacts SetOfArtifacts

	// Stats has a self-explanatory name
	Counters CityActivityCounters

	// PRIVATE
	// Pointer to the current Overlord of the current City
	pOverlord *City

	// PRIVATE
	// Pointer to cities we currently are the overlord of
	lieges SetOfCities
}

// UnitType gathers the core statistics of a kind of Unit. It actually dictates the
// behavior of the Unit.
type UnitType struct {
	// Unique ID of the Unit Type
	ID uint64 `json:"Id"`

	// The number of health point for that type of unit.
	// A health equal to 0 means the death of the unit.
	Health uint32

	// How affected is that type of unit by a loss of Health.
	// Must be between 0 and 1.
	// 0 means that the capacity of the Unit isn't affected by a health reduction.
	// 1 means that the capacity of the Unit loses an equal percentage of its capacity
	// for a loss of health (in other words, a HealthFactor of 1 means that the Unit
	// will hit at 90% of its maximal power if it has 90% of its health points).
	HealthFactor float64

	// The display name of the Unit Type
	Name string

	// How many ticks
	Ticks uint32

	// Transient bonus of Popularity, when the Unit is alive
	PopBonus int64

	// Permanent bonus of Popularity given when the Unit's training is done
	PopBonusTrain int64

	// Permanent bonus of Popularity given to the owner of the Unit when it dies
	PopBonusDeath int64

	// Permanent bonus of Popularity given to the killer of the Unit
	PopBonusKill int64

	// Permanent bonus of Popularity given to the ownerof the Unit when it is disbanded.
	PopBonusDisband int64

	// Instantiation cost of the current UnitType at the beginning of the process
	Cost0 Resources

	// Instantiation cost of the current UnitType, at each step of the process.
	Cost Resources

	// Might positive (resource boost) or more commonly negative (maintenance cost)
	Prod ResourceModifiers

	// Required Popularity to start training this type of troop
	ReqPop int64

	// A UnitType is only dependant on the presence of a Building of that BuildingType.
	RequiredBuilding uint64
}

// Unit is the part of an Army that can participate to a Fight.
// A Unit is owned after a training period. It has a regular cost.
type Unit struct {
	// Unique ID of the Unit
	ID string `json:"Id"`

	// A copy of the definition for the current Unit.
	Type uint64

	// How many ticks remain before the Troop training is finished
	Ticks uint32

	// The number of health points of the unit, Health should be less or equal to HealthMax
	Health uint32 `json:"H,omitempty"`
}

// Command is a queued order for a specific Army.
type Command struct {
	// The unique ID of the Cell to target
	Cell uint64

	// What to do once arrived at the given Cell.
	Action string `json:"action"`

	// Json
	Args string `json:"args"`
}

// ActionArgMove gives additional actions to be performed when the Army reaches its destination.
type ActionArgMove struct {
	// Resources to be transferred
	Amount Resources `json:"amount,omitempty"`

	// Artifact to be transferred
	Artifact uint64 `json:"artifact,omitempty"`

	// Troops to be transferred
	Units []uint64 `json:"units,omitempty"`
}

// ActionArgAssault tells what to do if the army is victorious
type ActionArgAssault struct {
	// Become the overlord of the City.
	Overlord bool `json:"overlord,omitempty"`

	// Break a random building
	Break bool `json:"break,omitempty"`

	// Impact the production of the subbsequent seasons
	Massacre bool `json:"massacre,omitempty"`
}

// Army is the entity able to act on a map.
// The actions are produced in a queued (fielf :Targets:) and the queue
// is consumed by a periodical task.
type Army struct {
	// The unique ID of the current Army
	ID string `json:"Id"`

	// A display name for the current City
	Name string

	// The City that controls the current Army
	City *City `json:"-"`

	// The ID of the Fight this Army is involved in.
	FightID string `json:"fight,omitempty"`

	// The ID of the Cell the Army is on
	Cell uint64 `json:",omitempty"`

	// How many resources are carried by that Army
	// The set may be empty.
	Stock Resources `json:",omitempty"`

	// Units that compose the current Army.
	// The set may be empty.
	Units SetOfUnits `json:",omitempty"`

	// Artifacts carried by the current Army
	// The set may be empty
	Artifacts SetOfArtifacts `json:",omitempty"`

	// The IS of a Cell of the Map that is a goal of the current movement of the Army
	Targets []Command `json:",omitempty"`

	// An array of Postures against armies of other cities.
	// A positive value means "defend"
	// A negative value means "assault"
	Postures []int64 `json:",omitempty"`
}

type Fight struct {
	// The unique ID of the
	ID string `json:"Id"`

	// The unique ID of the MapVertex the current Fight is happening on.
	Cell uint64

	// The set of ID of armies involved in the current Fight on the "attack" side
	// (the side that initiated the fight)
	Attack SetOfArmies

	/// The set of ID of armies involved in the current Fight on the "defence" side
	// (the side that has been force-pulled).
	Defense SetOfArmies
}

type SetOfFights []*Fight

//go:generate go run github.com/jfsmig/hegemonie/pkg/gen-set ./world_auto.go region:SetOfArtifacts:*Artifact ID:string
//go:generate go run github.com/jfsmig/hegemonie/pkg/gen-set ./world_auto.go region:SetOfArmies:*Army ID:string
//go:generate go run github.com/jfsmig/hegemonie/pkg/gen-set ./world_auto.go region:SetOfBuildings:*Building ID:string
//go:generate go run github.com/jfsmig/hegemonie/pkg/gen-set ./world_auto.go region:SetOfBuildingTypes:*BuildingType
//go:generate go run github.com/jfsmig/hegemonie/pkg/gen-set ./world_auto.go region:SetOfCities:*City
//go:generate go run github.com/jfsmig/hegemonie/pkg/gen-set ./world_auto.go region:SetOfTemplates:*City Name:string
//go:generate go run github.com/jfsmig/hegemonie/pkg/gen-set ./world_auto.go region:SetOfId:uint64 :uint64
//go:generate go run github.com/jfsmig/hegemonie/pkg/gen-set ./world_auto.go region:SetOfKnowledges:*Knowledge ID:string
//go:generate go run github.com/jfsmig/hegemonie/pkg/gen-set ./world_auto.go region:SetOfKnowledgeTypes:*KnowledgeType
//go:generate go run github.com/jfsmig/hegemonie/pkg/gen-set ./world_auto.go region:SetOfUnits:*Unit ID:string
//go:generate go run github.com/jfsmig/hegemonie/pkg/gen-set ./world_auto.go region:SetOfUnitTypes:*UnitType
//go:generate go run github.com/jfsmig/hegemonie/pkg/gen-set ./world_auto.go region:SetOfRegions:*Region Name:string
