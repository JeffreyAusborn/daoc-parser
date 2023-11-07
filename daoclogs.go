package main

import (
	"strings"
	"time"
)

type DaocLogs struct {
	User     Stats
	Enemy    []*Stats
	Friendly []*Stats
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

	UsersHit    []string // who you hit
	UsersHealed []string // who you healed
	UsersKilled []string // who you killed

	TotalKills  int // how many kills - pve and pvp
	TotalDeaths int // how many times you've died - pve and pvp

	Spells    []*Ability // This will only work if the ability with damage exist, like it does for dots and pets
	DotsNPets []*Ability // This will only work if the ability with damage exist, like it does for dots and pets
	Heals     []*Ability
	Styles    []*Ability

	ExperienceGained []int // experience gain

	TotalSelfHeal   []int // self healing
	TotalHeals      []int // healing all
	TotalHealsCrits []int // healing crits
	TotalAbsorbed   []int // how many absorbs a player has had

	TotalStuns            int
	SpellsPerformed       int
	CastedSpellsPerformed int
	ResistsOutTotal       int
	ResistsInTotal        int
	MissesTotal           int
	SiphonTotal           int
	BlockTotal            int
	ParryTotal            int
	EvadeTotal            int
	OverHeals             int
	Interrupted           int

	StartTime time.Time // First occurance of chat log opened
	EndTime   time.Time // Last known occurance of chat log closed
}

/*
Breakdown of each ability that we'll aggregate later
*/
type Ability struct {
	Name string
	Type string // heal, spell, melee

	Output      []int
	ExtraDamage []int
	Crit        []int
	GrowthRate  []int
	Interupts   []string
	Users       []*Stats
	Weapons     []*Ability

	Stunned   int
	Resists   int
	Blocked   int
	Parried   int
	Evaded    int
	Siphon    int
	OverHeals int
}

type Weapon struct {
	Name   string
	Output []int
}

/*
	Create getters and setters for the stats object
*/

func (_daocLogs *DaocLogs) getUser() *Stats {
	if _daocLogs != nil {
		return &_daocLogs.User
	}
	return &Stats{}
}

func (_daocLogs *DaocLogs) findEnemyStats(user string) *Stats {
	user = strings.TrimSpace(strings.ToLower(user))
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

func (_daocLogs *DaocLogs) findSpellStats(ability string) *Ability {
	ability = strings.TrimSpace(strings.ToLower(ability))
	if strings.Contains(ability, "call ") {
		ability = strings.Split(ability, "call ")[1]
	}
	if strings.Contains(ability, "summon ") {
		ability = strings.Split(ability, "summon ")[1]
	}
	for _, stats := range _daocLogs.User.Spells {
		if stats.Name == ability {
			return stats
		}
	}
	newAbility := Ability{}
	newAbility.Name = ability
	_daocLogs.getUser().Spells = append(_daocLogs.getUser().Spells, &newAbility)
	return &newAbility
}

func (_daocLogs *DaocLogs) findDotsNPetsStats(ability string) *Ability {
	ability = strings.TrimSpace(strings.ToLower(ability))
	for _, stats := range _daocLogs.User.DotsNPets {
		if stats.Name == ability {
			return stats
		}
	}
	newAbility := Ability{}
	newAbility.Name = ability
	_daocLogs.getUser().DotsNPets = append(_daocLogs.getUser().DotsNPets, &newAbility)
	return &newAbility
}

func (_spell *Ability) findUserStats(userName string) *Stats {
	userName = strings.TrimSpace(strings.ToLower(userName))
	for _, user := range _spell.Users {
		if user.UserName == userName {
			return user
		}
	}
	newUser := Stats{}
	newUser.UserName = userName
	_spell.Users = append(_spell.Users, &newUser)
	return &newUser
}

func (_daocLogs *DaocLogs) findStyleStats(ability string) *Ability {
	ability = strings.TrimSpace(strings.ToLower(ability))
	for _, stats := range _daocLogs.getUser().Styles {
		if stats.Name == ability {
			return stats
		}
	}
	newAbility := Ability{}
	newAbility.Name = ability
	_daocLogs.getUser().Styles = append(_daocLogs.getUser().Styles, &newAbility)
	return &newAbility
}

func (_ability *Ability) findWeaponStats(weaponName string) *Ability {
	weaponName = strings.TrimSpace(strings.ToLower(weaponName))
	for _, weapon := range _ability.Weapons {
		if weapon.Name == weaponName {
			return weapon
		}
	}
	newWeapon := Ability{}
	newWeapon.Name = weaponName
	_ability.Weapons = append(_ability.Weapons, &newWeapon)
	return &newWeapon
}

func (_daocLogs *DaocLogs) findHealStats(ability string) *Ability {
	ability = strings.TrimSpace(strings.ToLower(ability))
	for _, stats := range _daocLogs.getUser().Heals {
		if stats.Name == ability {
			return stats
		}
	}
	newAbility := Ability{}
	newAbility.Name = ability
	_daocLogs.getUser().Heals = append(_daocLogs.getUser().Heals, &newAbility)
	return &newAbility
}

func (_daocLogs *DaocLogs) healAbilityExist(ability string) bool {
	ability = strings.TrimSpace(strings.ToLower(ability))
	for _, stats := range _daocLogs.User.Heals {
		if stats.Name == ability {
			return true
		}
	}
	return false
}

func (_daocLogs *DaocLogs) findFriendlyStats(user string) *Stats {
	user = strings.TrimSpace(strings.ToLower(user))
	for _, stats := range _daocLogs.Friendly {
		if stats.UserName == user {
			return stats
		}
	}
	newUser := Stats{}
	newUser.UserName = user
	_daocLogs.Friendly = append(_daocLogs.Friendly, &newUser)
	return &newUser
}
