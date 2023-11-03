package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gosuri/uilive"
)

type DaocLogs struct {
	User  Stats
	Enemy Stats
}

type Stats struct {
	movingDamageTotal     []int
	movingDamageBaseMelee []int
	movingDamageStyles    []int
	movingDamageSpells    []int
	movingExtraDamage     []int
	movingCritDamage      []int
	usersHit              []string
	experienceGained      int
	totalSelfHeal         int
	totalAblative         int
	totalKills            int
	totalDeaths           int
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
	startTime             time.Time
	endTime               time.Time
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
		_daocLogs.User.usersHit = append(_daocLogs.User.usersHit, user)
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
		match, _ = regexp.MatchString("extra damage", line)
		_daocLogs.User.movingDamageTotal = append(_daocLogs.User.movingDamageTotal, damageInt)
		_daocLogs.User.movingCritDamage = append(_daocLogs.User.movingCritDamage, damageInt)
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
		ablative := strings.Split(line, "ablative absorbs ")[1]
		ablative = strings.Split(ablative, " damage")[0]
		ablativeInt, _ := strconv.Atoi(ablative)
		_daocLogs.User.totalAblative += ablativeInt
	}
}

func (_daocLogs *DaocLogs) regexSupport(line string) {
	match, _ := regexp.MatchString("You heal yourself for", line)
	if match {
		healing := strings.Split(line, "You heal yourself for ")[1]
		healing = strings.Split(healing, " hit points")[0]
		healingInt, _ := strconv.Atoi(healing)
		_daocLogs.User.totalSelfHeal += healingInt
	}
	match, _ = regexp.MatchString("is stunned and cannot move", line)
	if match {
		_daocLogs.User.totalStuns += 1
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
		_daocLogs.User.experienceGained += expInt
	}
	match, _ = regexp.MatchString("dies", line)
	if match {
		_daocLogs.User.totalKills += 1
	}
}
func (_daocLogs *DaocLogs) regexEnemy(line string) {
	match, _ := regexp.MatchString("parries your attack", line)
	if match {
		_daocLogs.Enemy.parryTotal += 1
	}
	match, _ = regexp.MatchString("evades your attack", line)
	if match {
		_daocLogs.Enemy.evadeTotal += 1
	}
	match, _ = regexp.MatchString("blocks your attack", line)
	if match {
		_daocLogs.Enemy.blockTotal += 1
	}
	match, _ = regexp.MatchString("hits your.*for ", line)
	if match {
		damage := strings.Split(line, "for ")[1]
		damage = strings.Split(damage, " damage")[0]
		damage = strings.Split(damage, " ")[0]
		damageInt, _ := strconv.Atoi(damage)
		_daocLogs.Enemy.movingDamageTotal = append(_daocLogs.Enemy.movingDamageTotal, damageInt)
	}
}
func (_daocLogs *DaocLogs) regexTime(line string) {
	if _daocLogs.User.startTime.IsZero() {
		match, _ := regexp.MatchString("Chat Log Opened", line)
		if match {
			timeObj, err := time.Parse("Mon Jan 02 15:04:05 2006", strings.TrimSuffix(strings.Split(line, ": ")[1], "\r\n"))
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			_daocLogs.User.startTime = timeObj
		}
	}
	match, _ := regexp.MatchString("Chat Log Closed", line)
	if match {
		timeObj, err := time.Parse("Mon Jan 02 15:04:05 2006", strings.TrimSuffix(strings.Split(line, ": ")[1], "\r\n"))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		_daocLogs.User.endTime = timeObj
	}
}

func (_daocLogs *DaocLogs) writeLogValues(writer *uilive.Writer) {
	fmt.Fprintf(writer, "Dark Age of Camelot - Chat Parser\nWritten by: Theorist\nIf you have any feedback, feel free to DM in Discord.\n\n")
	fmt.Fprintln(writer, "********** Melee **********")
	fmt.Fprintf(writer, "Styles Performed: %d\n", len(_daocLogs.User.movingDamageStyles))
	fmt.Fprintf(writer, "Base Hits Performed: %d\n", len(_daocLogs.User.movingDamageBaseMelee))
	fmt.Fprintf(writer, "Misses: %d\n", _daocLogs.User.missesTotal)
	fmt.Fprintf(writer, "Total Style Damage: %d\n", sumArr(_daocLogs.User.movingDamageStyles))
	fmt.Fprintf(writer, "Total Base Hit Damage: %d\n", sumArr(_daocLogs.User.movingDamageBaseMelee))
	fmt.Fprintf(writer, "Total Melee Damage: %d\n", sumArr(_daocLogs.User.movingDamageStyles)+sumArr(_daocLogs.User.movingDamageBaseMelee))
	fmt.Fprintln(writer, "********** Spells **********")
	fmt.Fprintf(writer, "Total Spells Performed: %d\n", _daocLogs.User.spellsPerformed)
	fmt.Fprintf(writer, "Casted Spells Performed: %d\n", _daocLogs.User.castedSpellsPerformed)
	fmt.Fprintf(writer, "Insta Spells Performed: %d\n", _daocLogs.User.spellsPerformed-_daocLogs.User.castedSpellsPerformed)
	fmt.Fprintf(writer, "Spells with Damage: %d\n", len(_daocLogs.User.movingDamageSpells))
	fmt.Fprintf(writer, "Total Resists: %d\n", _daocLogs.User.resistsTotal)
	fmt.Fprintf(writer, "Total Siphons: %d\n", _daocLogs.User.siphonTotal)
	fmt.Fprintf(writer, "Spell Damage: %d\n", sumArr(_daocLogs.User.movingDamageSpells))
	fmt.Fprintf(writer, "Spell Extra Damage: %d\n", sumArr(_daocLogs.User.movingExtraDamage))
	fmt.Fprintln(writer, "********** Criticals **********")
	fmt.Fprintf(writer, "Total Crits: %d\n", len(_daocLogs.User.movingCritDamage))
	fmt.Fprintf(writer, "Total Crit Damage: %d\n", sumArr(_daocLogs.User.movingCritDamage))
	fmt.Fprintln(writer, "********** Defensives **********")
	fmt.Fprintf(writer, "Total Blocks: %d\n", _daocLogs.User.blockTotal)
	fmt.Fprintf(writer, "Total Parrys: %d\n", _daocLogs.User.parryTotal)
	fmt.Fprintf(writer, "Total Evades: %d\n", _daocLogs.User.evadeTotal)
	fmt.Fprintf(writer, "Total Stuns: %d\n", _daocLogs.User.totalStuns)
	fmt.Fprintf(writer, "Total Self Heals: %d\n", _daocLogs.User.totalSelfHeal)
	fmt.Fprintf(writer, "Total Ablative Absorbs: %d\n", _daocLogs.User.totalAblative)
	fmt.Fprintln(writer, "********** Enemy **********")
	fmt.Fprintf(writer, "Enemy Total Hits: %d\n", len(_daocLogs.Enemy.movingDamageTotal))
	fmt.Fprintf(writer, "Enemy Total Parrys: %d\n", _daocLogs.Enemy.parryTotal)
	fmt.Fprintf(writer, "Enemy Total Evades: %d\n", _daocLogs.Enemy.evadeTotal)
	fmt.Fprintf(writer, "Enemy Total Blocks: %d\n", _daocLogs.Enemy.blockTotal)
	fmt.Fprintf(writer, "Total Damage Taken From Enemy: %d\n", sumArr(_daocLogs.Enemy.movingDamageTotal))
	fmt.Fprintln(writer, "********** Total **********")
	fmt.Fprintf(writer, "Total Damage: %d\n", sumArr(_daocLogs.User.movingDamageTotal))
	fmt.Fprintf(writer, "Total Experience Gained: %d\n", _daocLogs.User.experienceGained)
	fmt.Fprintf(writer, "Total Killed: %d\n", _daocLogs.User.totalKills)
	fmt.Fprintf(writer, "Users Hit: %s\n", strings.Join(dedupe(_daocLogs.User.usersHit), ","))
	totalMinutes := int(_daocLogs.User.endTime.Sub(_daocLogs.User.startTime).Seconds()) / 60
	totalSeconds := int(_daocLogs.User.endTime.Sub(_daocLogs.User.startTime).Seconds()) - (60 * totalMinutes)
	fmt.Fprintf(writer, "Total Time: %d minutes and %d seconds\n", totalMinutes, totalSeconds)
}
