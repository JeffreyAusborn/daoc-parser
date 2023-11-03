package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

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
	userName              string
	armorHit              ArmorHit
	movingDamageTotal     []int // all damage
	movingDamageBaseMelee []int // melee hit without using a style
	movingDamageStyles    []int // melee hit with a style
	movingDamageSpells    []int // spells hit - this includes weapon/armor procs, style procs
	movingExtraDamage     []int // extra damage (damage add)
	movingCritDamage      []int // crit damage
	movingDamageReceived  []int // damage receieved - more so for the enemy player

	usersHit    []string // who have you hit
	usersHealed []string // who have you hit
	totalKills  int      // how many kills - pve and pvp
	totalDeaths int      // how many times you've died - pve and pvp

	experienceGained []int // experience gain

	totalSelfHeal []int // self healing - procs, styles, spells
	totalHeals    []int // healing all - procs, styles, spells,
	totalAbsorbed []int // how many absorbs a player has had

	totalStuns            int
	spellsPerformed       int
	castedSpellsPerformed int
	resistsTotal          int
	missesTotal           int
	siphonTotal           int
	blockTotal            int
	parryTotal            int
	evadeTotal            int
	overHeals             int

	startTime time.Time // First occurance of chat log opened
	endTime   time.Time // Last known occurance of chat log closed
}

func (_daocLogs *DaocLogs) findUserStats(user string) *Stats {
	for _, stats := range _daocLogs.Enemy {
		if stats.userName == user {
			return stats
		}
	}
	newUser := Stats{}
	newUser.userName = user
	_daocLogs.Enemy = append(_daocLogs.Enemy, &newUser)
	return &newUser
}

func (_daocLogs *DaocLogs) regexOffensive(line string, style bool) bool {
	match, _ := regexp.MatchString("You attack.*with your.*and hit for.*damage", line)
	if match {
		damage := strings.Split(line, "and hit for ")[1]
		damage = strings.Split(damage, " damage")[0]
		damage = strings.Split(damage, " ")[0]
		damageInt, _ := strconv.Atoi(damage)
		_daocLogs.User.movingDamageTotal = append(_daocLogs.User.movingDamageTotal, damageInt)
		if style {
			style = false
			_daocLogs.User.movingDamageStyles = append(_daocLogs.User.movingDamageStyles, damageInt)
		} else {
			_daocLogs.User.movingDamageBaseMelee = append(_daocLogs.User.movingDamageBaseMelee, damageInt)
		}

		user := strings.Split(line, "You attack ")[1]
		user = strings.Split(user, " with your")[0]
		userStats := _daocLogs.findUserStats(user)
		userStats.movingDamageReceived = append(userStats.movingDamageReceived, damageInt)
	}
	match, _ = regexp.MatchString("You prepare to perform", line)
	if match {
		style = true
	}
	match, _ = regexp.MatchString("You cast a ", line)
	if match {
		_daocLogs.User.spellsPerformed += 1
	}
	match, _ = regexp.MatchString("You begin casting a ", line)
	if match {
		_daocLogs.User.castedSpellsPerformed += 1
	}
	match, _ = regexp.MatchString("You hit.*for.*damage!", line)
	if match {
		damage := strings.Split(line, " for ")[1]
		damage = strings.Split(damage, " damage")[0]
		damage = strings.Split(damage, " ")[0]
		damageInt, _ := strconv.Atoi(damage)
		match, _ = regexp.MatchString("extra damage", line)
		if match {
			_daocLogs.User.movingExtraDamage = append(_daocLogs.User.movingExtraDamage, damageInt)
		} else {
			_daocLogs.User.movingDamageSpells = append(_daocLogs.User.movingDamageSpells, damageInt)
		}
		_daocLogs.User.movingDamageTotal = append(_daocLogs.User.movingDamageTotal, damageInt)

		user := strings.Split(line, "You hit ")[1]
		user = strings.Split(user, " for")[0]
		_daocLogs.User.usersHit = append(_daocLogs.User.usersHit, user)
	}
	match, _ = regexp.MatchString("Your.*hits.*for.*damage", line)
	if match {
		damage := strings.Split(line, " for ")[1]
		damage = strings.Split(damage, " damage")[0]
		damage = strings.Split(damage, " ")[0]
		damageInt, _ := strconv.Atoi(damage)
		_daocLogs.User.movingDamageSpells = append(_daocLogs.User.movingDamageSpells, damageInt)
		_daocLogs.User.movingDamageTotal = append(_daocLogs.User.movingDamageTotal, damageInt)

		user := strings.Split(line, "hits ")[1]
		user = strings.Split(user, " for")[0]
		userStats := _daocLogs.findUserStats(user)
		userStats.movingDamageReceived = append(userStats.movingDamageReceived, damageInt)
	}
	match, _ = regexp.MatchString("You miss!.*", line)
	if match {
		_daocLogs.User.missesTotal += 1
	}
	match, _ = regexp.MatchString("resists the effect.*", line)
	if match {
		_daocLogs.User.resistsTotal += 1
	}
	match, _ = regexp.MatchString("You gather energy from your surroundings.*", line)
	if match {
		_daocLogs.User.siphonTotal += 1
	}
	match, _ = regexp.MatchString("fully healed", line)
	if match {
		_daocLogs.User.overHeals += 1
	}
	match, _ = regexp.MatchString("You critically hit", line)
	if match {
		damage := strings.Split(line, "for an additional ")[1]
		damage = strings.Split(damage, " damage")[0]
		damageInt, _ := strconv.Atoi(damage)
		_daocLogs.User.movingDamageTotal = append(_daocLogs.User.movingDamageTotal, damageInt)
		_daocLogs.User.movingCritDamage = append(_daocLogs.User.movingCritDamage, damageInt)

		user := strings.Split(line, "hit ")[1]
		user = strings.Split(user, " for")[0]
		userStats := _daocLogs.findUserStats(user)
		userStats.movingDamageReceived = append(userStats.movingDamageReceived, damageInt)
	}
	return style
}

func (_daocLogs *DaocLogs) regexDefensives(line string) {
	match, _ := regexp.MatchString("you block", line)
	if match {
		_daocLogs.User.blockTotal += 1
	}
	match, _ = regexp.MatchString("you evade", line)
	if match {
		_daocLogs.User.evadeTotal += 1
	}
	match, _ = regexp.MatchString("you parry", line)
	if match {
		_daocLogs.User.parryTotal += 1
	}
	match, _ = regexp.MatchString("Your ablative absorbs", line)
	if match {
		absorb := strings.Split(line, "ablative absorbs ")[1]
		absorb = strings.Split(absorb, " damage")[0]
		absorbInt, _ := strconv.Atoi(absorb)
		_daocLogs.User.totalAbsorbed = append(_daocLogs.User.totalAbsorbed, absorbInt)
	}
}

func (_daocLogs *DaocLogs) regexSupport(line string) {
	match, _ := regexp.MatchString("You heal yourself for", line)
	if match {
		healing := strings.Split(line, "You heal yourself for ")[1]
		healing = strings.Split(healing, " hit points")[0]
		healingInt, _ := strconv.Atoi(healing)
		_daocLogs.User.totalSelfHeal = append(_daocLogs.User.totalSelfHeal, healingInt)
	}
	match, _ = regexp.MatchString("You heal.*for", line)
	if match {
		healing := strings.Split(line, "for ")[1]
		healing = strings.Split(healing, " hit points")[0]
		healingInt, _ := strconv.Atoi(healing)
		_daocLogs.User.totalHeals = append(_daocLogs.User.totalHeals, healingInt)
	}
	match, _ = regexp.MatchString("is stunned and cannot move", line)
	if match {
		_daocLogs.User.totalStuns += 1
		user := strings.Split(line, " is ")[0]
		user = strings.Split(user, " ")[1]
		userStats := _daocLogs.findUserStats(user)
		userStats.totalStuns += 1
	}
}

func (_daocLogs *DaocLogs) regexPets(line string) {
	match, _ := regexp.MatchString("The.*casts a spell!", line)
	if match {
	}
	match, _ = regexp.MatchString("Your.*attacks.*and hits for.*damage!", line)
	if match {
	}
	match, _ = regexp.MatchString("The.*attacks.*and misses!", line)
	if match {
	}
}

func (_daocLogs *DaocLogs) regexMisc(line string) {
	match, _ := regexp.MatchString("You gain a total of.*experience points", line)
	if match {
		exp := strings.Split(line, "You gain a total of ")[1]
		exp = strings.Split(exp, " experience")[0]
		exp = strings.ReplaceAll(exp, ",", "")
		expInt, _ := strconv.Atoi(exp)
		_daocLogs.User.experienceGained = append(_daocLogs.User.experienceGained, expInt)
	}
	match, _ = regexp.MatchString("You just killed", line)
	if match {
		_daocLogs.User.totalKills += 1
	}
}

func (_daocLogs *DaocLogs) regexEnemy(line string) {
	match, _ := regexp.MatchString("parries your attack", line)
	if match {
		user := strings.Split(line, " parries")[0]
		user = strings.Split(line, " ")[1]
		userStats := _daocLogs.findUserStats(user)
		userStats.parryTotal += 1
	}
	match, _ = regexp.MatchString("evades your attack", line)
	if match {
		user := strings.Split(line, " evades")[0]
		user = strings.Split(line, " ")[1]
		userStats := _daocLogs.findUserStats(user)
		userStats.evadeTotal += 1
	}
	match, _ = regexp.MatchString("blocks your attack", line)
	if match {
		user := strings.Split(line, " blocks")[0]
		user = strings.Split(line, " ")[1]
		userStats := _daocLogs.findUserStats(user)
		userStats.blockTotal += 1
	}
	match, _ = regexp.MatchString("hits your.*for ", line)
	if match {
		damage := strings.Split(line, "for ")[1]
		damage = strings.Split(damage, " damage")[0]
		damage = strings.Split(damage, " ")[0]
		damageInt, _ := strconv.Atoi(damage)

		user := strings.Split(line, " hits your")[0]
		user = strings.Split(line, " ")[1]
		userStats := _daocLogs.findUserStats(user)

		userStats.movingDamageTotal = append(userStats.movingDamageTotal, damageInt)

		armorPiece := strings.Split(line, " for")[0]
		armorPiece = strings.Split(armorPiece, "your ")[1]

		switch armorPiece {
		case "head":
			userStats.armorHit.Head = append(userStats.armorHit.Head, damageInt)
			userStats.movingDamageStyles = append(userStats.movingDamageStyles, damageInt)
		case "torso":
			userStats.armorHit.Torso = append(userStats.armorHit.Torso, damageInt)
			userStats.movingDamageStyles = append(userStats.movingDamageStyles, damageInt)
		case "leg":
			userStats.armorHit.Leg = append(userStats.armorHit.Leg, damageInt)
			userStats.movingDamageStyles = append(userStats.movingDamageStyles, damageInt)
		case "arm":
			userStats.armorHit.Arm = append(userStats.armorHit.Arm, damageInt)
			userStats.movingDamageStyles = append(userStats.movingDamageStyles, damageInt)
		case "hand":
			userStats.armorHit.Hand = append(userStats.armorHit.Hand, damageInt)
			userStats.movingDamageStyles = append(userStats.movingDamageStyles, damageInt)
		case "foot":
			userStats.armorHit.Foot = append(userStats.armorHit.Foot, damageInt)
			userStats.movingDamageStyles = append(userStats.movingDamageStyles, damageInt)
		default:
			break
		}
	}
	match, _ = regexp.MatchString("critically hits you", line)
	if match {
		damage := strings.Split(line, "for an additional ")[1]
		damage = strings.Split(damage, " damage")[0]
		damageInt, _ := strconv.Atoi(damage)
		playerStats := _daocLogs.findUserStats("")
		playerStats.movingDamageTotal = append(playerStats.movingDamageTotal, damageInt)
		playerStats.movingCritDamage = append(playerStats.movingCritDamage, damageInt)
	}
}
func (_daocLogs *DaocLogs) regexTime(line string) {
	if _daocLogs.User.startTime.IsZero() {
		match, _ := regexp.MatchString("Chat Log Opened", line)
		if match {
			timeObj, err := time.Parse("Mon Jan 02 15:04:05 2006", strings.TrimSuffix(strings.Split(line, ": ")[1], "\r\n"))
			if err != nil {
				fmt.Sprintln("Error:", err)
				return
			}
			_daocLogs.User.startTime = timeObj
		}
	}
	match, _ := regexp.MatchString("Chat Log Closed", line)
	if match {
		timeObj, err := time.Parse("Mon Jan 02 15:04:05 2006", strings.TrimSuffix(strings.Split(line, ": ")[1], "\r\n"))
		if err != nil {
			fmt.Sprintln("Error:", err)
			return
		}
		_daocLogs.User.endTime = timeObj
	}
}

func (_daocLogs *DaocLogs) writeLogValues() string {
	return _daocLogs.calculateArmorhits() + "\n" + _daocLogs.calculateDamageIn() + "\n" + _daocLogs.calculateEnemyDensives() + "\n" + _daocLogs.calculateDamageOut() + "\n" + _daocLogs.calculateHeal() + "\n" + _daocLogs.calculateDensives()
}

func (_daocLogs *DaocLogs) calculateDamageOut() string {
	if len(_daocLogs.User.movingDamageTotal) > 0 {
		damageIn := "----- Damage Out -----\n"
		meleeDamage := ""
		spellDamage := ""
		critDamage := ""
		spellsResists := ""
		meleeMiss := ""
		siphons := ""
		kills := ""
		if len(_daocLogs.User.movingDamageStyles)+len(_daocLogs.User.movingDamageBaseMelee) > 0 {
			meleeDamage = fmt.Sprintf("\tMelee Hit: %d\n\tMelee Damage: %d\n", len(_daocLogs.User.movingDamageStyles)+len(_daocLogs.User.movingDamageBaseMelee), sumArr(_daocLogs.User.movingDamageStyles)+sumArr(_daocLogs.User.movingDamageBaseMelee))
		}
		if len(_daocLogs.User.movingDamageSpells)+len(_daocLogs.User.movingExtraDamage) > 0 {
			spellDamage = fmt.Sprintf("\tSpell Hit: %d\n\tSpell Damage: %d\n", len(_daocLogs.User.movingDamageSpells)+len(_daocLogs.User.movingExtraDamage), sumArr(_daocLogs.User.movingDamageSpells)+sumArr(_daocLogs.User.movingExtraDamage))
		}
		if len(_daocLogs.User.movingCritDamage) > 0 {
			critDamage = fmt.Sprintf("\tCrit Hit: %d\n\tCrit Damage: %d\n", len(_daocLogs.User.movingCritDamage), sumArr(_daocLogs.User.movingCritDamage))
		}

		if _daocLogs.User.resistsTotal > 0 {
			spellsResists = fmt.Sprintf("\tResits: %d\n", _daocLogs.User.resistsTotal)
		}
		if _daocLogs.User.missesTotal > 0 {
			meleeMiss = fmt.Sprintf("\tMisses: %d\n", _daocLogs.User.missesTotal)
		}
		if _daocLogs.User.siphonTotal > 0 {
			siphons = fmt.Sprintf("\tSiphons: %d\n", _daocLogs.User.siphonTotal)
		}
		if _daocLogs.User.totalKills > 0 {
			kills = fmt.Sprintf("\tKills: %d\n", _daocLogs.User.totalKills)
		}
		return damageIn + meleeDamage + spellDamage + critDamage + spellsResists + meleeMiss + siphons + kills
	}
	return ""
}

func (_daocLogs *DaocLogs) calculateDamageIn() string {
	totalMeleeDamage := []int{}
	totalAllDamage := []int{}

	for _, user := range _daocLogs.Enemy {
		totalMeleeDamage = append(totalMeleeDamage, user.movingDamageStyles...)
		totalAllDamage = append(totalAllDamage, user.movingDamageTotal...)
	}

	totalDamage := sumArr(totalAllDamage)
	totalDamageMelee := sumArr(totalMeleeDamage)
	totalDamageSpell := totalDamage - totalDamageMelee

	if totalDamage > 0 {
		damageIn := "----- Damage In -----\n"
		meleeDamage := ""
		spellDamage := ""
		if totalDamageMelee > 0 {
			meleeDamage = fmt.Sprintf("\tMelee Hit: %d\n\tMelee Damage: %d\n", len(totalMeleeDamage), totalDamageMelee)
		}
		if totalDamageSpell > 0 {
			spellDamage = fmt.Sprintf("\tSpell Hit: %d\n\tSpell Damage: %d\n", len(totalAllDamage)-len(totalMeleeDamage), totalDamageSpell)
		}
		return damageIn + meleeDamage + spellDamage
	}
	return ""
}

func (_daocLogs *DaocLogs) calculateHeal() string {
	if len(_daocLogs.User.totalSelfHeal)+len(_daocLogs.User.totalAbsorbed)+len(_daocLogs.User.totalHeals) > 0 {
		healAndAbsorb := "----- Healing & Absorb -----\n"
		selfHeal := ""
		selfAbsorb := ""
		allHeal := ""
		overHeal := ""
		if len(_daocLogs.User.totalSelfHeal) > 0 {
			selfHeal = fmt.Sprintf("\tSelf Heals: %d\n", sumArr(_daocLogs.User.totalSelfHeal))
		}
		if len(_daocLogs.User.totalHeals) > 0 {
			allHeal = fmt.Sprintf("\tAll Heals: %d\n", sumArr(_daocLogs.User.totalHeals))
		}
		if len(_daocLogs.User.totalAbsorbed) > 0 {
			selfAbsorb = fmt.Sprintf("\tAbsorbed: %d\n", sumArr(_daocLogs.User.totalAbsorbed))
		}
		if _daocLogs.User.overHeals > 0 {
			overHeal = fmt.Sprintf("\tOverHeal Count: %d\n", _daocLogs.User.overHeals)
		}
		return healAndAbsorb + selfHeal + allHeal + overHeal + selfAbsorb
	}
	return ""
}

func (_daocLogs *DaocLogs) calculateDensives() string {
	if len(_daocLogs.User.totalSelfHeal)+len(_daocLogs.User.totalAbsorbed)+len(_daocLogs.User.totalHeals) > 0 {
		defensives := "----- Defensives -----\n"
		block := ""
		parry := ""
		evade := ""
		stuns := ""
		if _daocLogs.User.blockTotal > 0 {
			block = fmt.Sprintf("\tBlock: %d\n", _daocLogs.User.blockTotal)
		}
		if _daocLogs.User.parryTotal > 0 {
			parry = fmt.Sprintf("\tParry: %d\n", _daocLogs.User.parryTotal)
		}
		if _daocLogs.User.evadeTotal > 0 {
			evade = fmt.Sprintf("\tEvade: %d\n", _daocLogs.User.evadeTotal)
		}
		if _daocLogs.User.totalStuns > 0 {
			stuns = fmt.Sprintf("\tStuns: %d\n", _daocLogs.User.totalStuns)
		}
		return defensives + block + parry + evade + stuns
	}
	return ""
}

func (_daocLogs *DaocLogs) calculateEnemyDensives() string {
	blocks := 0
	evades := 0
	parries := 0
	for _, user := range _daocLogs.Enemy {
		blocks += user.blockTotal
		evades += user.evadeTotal
		parries += user.parryTotal
	}

	if blocks+parries+evades > 0 {
		defensives := "----- Enemy Defensives -----\n"
		block := ""
		parry := ""
		evade := ""
		if blocks > 0 {
			block = fmt.Sprintf("\tBlock: %d\n", blocks)
		}
		if parries > 0 {
			parry = fmt.Sprintf("\tParry: %d\n", parries)
		}
		if evades > 0 {
			evade = fmt.Sprintf("\tEvade: %d\n", evades)
		}
		return defensives + block + parry + evade
	}
	return ""
}

func (_daocLogs *DaocLogs) calculateTime() string {
	totalMinutes := int(_daocLogs.User.endTime.Sub(_daocLogs.User.startTime).Seconds()) / 60
	totalSeconds := int(_daocLogs.User.endTime.Sub(_daocLogs.User.startTime).Seconds()) - (60 * totalMinutes)
	return fmt.Sprintf("----- Total Time -----\n\t%d minutes and %d seconds\n", totalMinutes, totalSeconds)
}

func (_daocLogs *DaocLogs) calculateArmorhits() string {
	head := []int{}
	torso := []int{}
	arm := []int{}
	leg := []int{}
	hand := []int{}
	foot := []int{}
	for _, user := range _daocLogs.Enemy {
		head = append(head, user.armorHit.Head...)
		torso = append(torso, user.armorHit.Torso...)
		arm = append(arm, user.armorHit.Arm...)
		leg = append(leg, user.armorHit.Leg...)
		hand = append(hand, user.armorHit.Hand...)
		foot = append(foot, user.armorHit.Foot...)
	}
	if len(head)+len(torso)+len(arm)+len(leg)+len(hand)+len(foot) > 0 {
		headHitFmt := ""
		torsoHitFmt := ""
		armHitFmt := ""
		legHitFmt := ""
		handHitFmt := ""
		footHitFmt := ""
		armorHitFmt := "----- Armor Damaged -----\n"
		if len(head) > 0 {
			headHitFmt = fmt.Sprintf("\tHead Hit: %d\n\tHead Damage: %d\n", len(head), sumArr(head))
		}
		if len(torso) > 0 {
			torsoHitFmt = fmt.Sprintf("\tTorso Hit: %d\n\tTorso Damage: %d\n", len(torso), sumArr(torso))
		}
		if len(arm) > 0 {
			armHitFmt = fmt.Sprintf("\tArm Hit: %d\n\tArm Damage: %d\n", len(arm), sumArr(arm))
		}
		if len(leg) > 0 {
			legHitFmt = fmt.Sprintf("\tLeg Hit: %d\n\tLeg Damage: %d\n", len(leg), sumArr(leg))
		}
		if len(hand) > 0 {
			handHitFmt = fmt.Sprintf("\tHand Hit: %d\n\tHand Damage: %d\n", len(hand), sumArr(hand))
		}
		if len(foot) > 0 {
			footHitFmt = fmt.Sprintf("\tFoot Hit: %d\n\tFoot Damage: %d\n", len(foot), sumArr(foot))
		}

		return armorHitFmt + headHitFmt + torsoHitFmt + armHitFmt + legHitFmt + handHitFmt + footHitFmt
	}
	return ""
}
