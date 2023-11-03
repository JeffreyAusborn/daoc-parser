package main

import "time"

type DaocLogs struct {
	User  Stats
	Enemy []*Stats
}

/*
how often a piece of gear was hit
what it was hit for.
*/
type ArmorHit struct {
	Head  []int
	Torso []int
	Arm   []int
	Leg   []int
	Hand  []int
	Foot  []int
}

/*
	Reason for slice int versus int
	- We can use it to extract more information such as
		- min
		- max
		- average
		- etc
*/

type Stats struct {
	UserName              string
	ArmorHit              ArmorHit
	MovingDamageTotal     []int // all damage
	MovingDamageBaseMelee []int // melee hit without using a style
	MovingDamageStyles    []int // melee hit with a style
	MovingDamageSpells    []int // spells hit - this includes weapon/armor procs, style procs
	MovingExtraDamage     []int // extra damage (damage add)
	MovingCritDamage      []int // crit damage
	MovingDamageReceived  []int // damage receieved - more so for the enemy player

	UsersHit    []string // who have you hit
	UsersHealed []string // who have you hit
	TotalKills  int      // how many kills - pve and pvp
	TotalDeaths int      // how many times you've died - pve and pvp

	ExperienceGained []int // experience gain

	TotalSelfHeal []int // self healing - procs, styles, spells
	TotalHeals    []int // healing all - procs, styles, spells,
	TotalAbsorbed []int // how many absorbs a player has had

	TotalStuns            int
	SpellsPerformed       int
	CastedSpellsPerformed int
	ResistsTotal          int
	MissesTotal           int
	SiphonTotal           int
	BlockTotal            int
	ParryTotal            int
	EvadeTotal            int
	OverHeals             int

	StartTime time.Time // First occurance of chat log opened
	EndTime   time.Time // Last known occurance of chat log closed
}

/*
	Create getters and setters for the stats object
*/

func (_daocLogs *DaocLogs) writeLogValues() string {
	return _daocLogs.calculateArmorhits() + "\n" + _daocLogs.calculateDamageIn() + "\n" + _daocLogs.calculateEnemyDensives() + "\n" + _daocLogs.calculateDamageOut() + "\n" + _daocLogs.calculateHeal() + "\n" + _daocLogs.calculateDensives()
}

func (_daocLogs *DaocLogs) getUser() *Stats {
	if _daocLogs != nil {
		return &_daocLogs.User
	}
	return &Stats{}
}

func (_daocLogs *DaocLogs) findEnemyStats(user string) *Stats {
	for _, stats := range _daocLogs.Enemy {
		if stats.UserName == user {
			return stats
		}
	}
	newUser := Stats{}
	newUser.UserName = user
	_daocLogs.Enemy = append(_daocLogs.Enemy, &newUser)
	return &newUser
}
