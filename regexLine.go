package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (_daocLogs *DaocLogs) regexOffensive(line string, style bool) bool {
	match, _ := regexp.MatchString("@@", line)
	if match {
		return false
	}
	match, _ = regexp.MatchString("You attack.*with your.*and hit for.*damage", line)
	if match {
		match, _ = regexp.MatchString("critically hit", line)
		if !match {
			damage := strings.Split(line, "and hit for ")[1]
			damage = strings.Split(damage, " damage")[0]
			damage = strings.Split(damage, " ")[0]
			damageInt, _ := strconv.Atoi(damage)
			_daocLogs.getUser().MovingDamageTotal = append(_daocLogs.getUser().MovingDamageTotal, damageInt)
			if style {
				style = false
				_daocLogs.getUser().MovingDamageStyles = append(_daocLogs.getUser().MovingDamageStyles, damageInt)
			} else {
				_daocLogs.getUser().MovingDamageBaseMelee = append(_daocLogs.getUser().MovingDamageBaseMelee, damageInt)
			}

			user := strings.Split(line, "You attack ")[1]
			user = strings.Split(user, " with your")[0]
			userStats := _daocLogs.findEnemyStats(user)
			userStats.MovingDamageReceived = append(userStats.MovingDamageReceived, damageInt)
		}
	}
	match, _ = regexp.MatchString("You prepare to perform", line)
	if match {
		style = true
	}
	match, _ = regexp.MatchString("You cast a ", line)
	if match {
		_daocLogs.getUser().SpellsPerformed += 1
	}
	match, _ = regexp.MatchString("You begin casting a ", line)
	if match {
		_daocLogs.getUser().CastedSpellsPerformed += 1
	}
	match, _ = regexp.MatchString("You hit.*for.*damage", line)
	if match {
		match, _ = regexp.MatchString("critically hit", line)
		if !match {
			damage := strings.Split(line, " for ")[1]
			damage = strings.Split(damage, " damage")[0]
			damage = strings.Split(damage, " ")[0]
			damageInt, _ := strconv.Atoi(damage)
			match, _ = regexp.MatchString("extra damage", line)
			if match {
				_daocLogs.getUser().MovingExtraDamage = append(_daocLogs.getUser().MovingExtraDamage, damageInt)
			} else {
				_daocLogs.getUser().MovingDamageSpells = append(_daocLogs.getUser().MovingDamageSpells, damageInt)
			}
			_daocLogs.getUser().MovingDamageTotal = append(_daocLogs.getUser().MovingDamageTotal, damageInt)

			user := strings.Split(line, "You hit ")[1]
			user = strings.Split(user, " for")[0]
			userStats := _daocLogs.findEnemyStats(user)
			userStats.MovingDamageReceived = append(userStats.MovingDamageReceived, damageInt)
		}
	}
	match, _ = regexp.MatchString("Your.*hits.*for.*damage", line)
	if match {
		damage := strings.Split(line, " for ")[1]
		damage = strings.Split(damage, " damage")[0]
		damage = strings.Split(damage, " ")[0]
		damageInt, _ := strconv.Atoi(damage)
		_daocLogs.getUser().MovingDamageSpells = append(_daocLogs.getUser().MovingDamageSpells, damageInt)
		_daocLogs.getUser().MovingDamageTotal = append(_daocLogs.getUser().MovingDamageTotal, damageInt)

		user := strings.Split(line, "hits ")[1]
		user = strings.Split(user, " for")[0]
		userStats := _daocLogs.findEnemyStats(user)
		userStats.MovingDamageReceived = append(userStats.MovingDamageReceived, damageInt)
	}
	match, _ = regexp.MatchString("You miss!.*", line)
	if match {
		_daocLogs.getUser().MissesTotal += 1
	}
	match, _ = regexp.MatchString("resists the effect.*", line)
	if match {
		_daocLogs.getUser().ResistsOutTotal += 1
	}
	match, _ = regexp.MatchString("You gather energy from your surroundings.*", line)
	if match {
		_daocLogs.getUser().SiphonTotal += 1
	}
	match, _ = regexp.MatchString("fully healed", line)
	if match {
		_daocLogs.getUser().OverHeals += 1
	}
	match, _ = regexp.MatchString("You critically hit", line)
	if match {
		damage := strings.Split(line, "for an additional ")[1]
		damage = strings.Split(damage, " damage")[0]
		damageInt, _ := strconv.Atoi(damage)
		_daocLogs.getUser().MovingDamageTotal = append(_daocLogs.getUser().MovingDamageTotal, damageInt)
		_daocLogs.getUser().MovingCritDamage = append(_daocLogs.getUser().MovingCritDamage, damageInt)

		user := strings.Split(line, "hit ")[1]
		user = strings.Split(user, " for")[0]
		userCheck := strings.Split(user, " ")
		if len(userCheck) == 0 {
			userStats := _daocLogs.findEnemyStats(user)
			userStats.MovingDamageReceived = append(userStats.MovingDamageReceived, damageInt)
		}

	}
	return style
}

func (_daocLogs *DaocLogs) regexDefensives(line string) {
	match, _ := regexp.MatchString("@@", line)
	if match {
		return
	}
	match, _ = regexp.MatchString("you block", line)
	if match {
		_daocLogs.getUser().BlockTotal += 1
	}
	match, _ = regexp.MatchString("you evade", line)
	if match {
		_daocLogs.getUser().EvadeTotal += 1
	}
	match, _ = regexp.MatchString("you parry", line)
	if match {
		_daocLogs.getUser().ParryTotal += 1
	}
	match, _ = regexp.MatchString("Your ablative absorbs", line)
	if match {
		absorb := strings.Split(line, "ablative absorbs ")[1]
		absorb = strings.Split(absorb, " damage")[0]
		absorbInt, _ := strconv.Atoi(absorb)
		_daocLogs.getUser().TotalAbsorbed = append(_daocLogs.getUser().TotalAbsorbed, absorbInt)
	}
	match, _ = regexp.MatchString("You resist the effect", line)
	if match {
		_daocLogs.getUser().ResistsInTotal += 1
	}
}

func (_daocLogs *DaocLogs) regexSupport(line string) {
	match, _ := regexp.MatchString("@@", line)
	if match {
		return
	}
	match, _ = regexp.MatchString("You heal yourself for", line)
	if match {
		healing := strings.Split(line, "You heal yourself for ")[1]
		healing = strings.Split(healing, " hit points")[0]
		healingInt, _ := strconv.Atoi(healing)
		_daocLogs.getUser().TotalSelfHeal = append(_daocLogs.getUser().TotalSelfHeal, healingInt)
		_daocLogs.getUser().TotalHeals = append(_daocLogs.getUser().TotalHeals, healingInt)
	}
	match, _ = regexp.MatchString("You transfer.*hit points", line)
	if match {
		healing := strings.Split(line, "transfer ")[1]
		healing = strings.Split(healing, " hit points")[0]
		healingInt, _ := strconv.Atoi(healing)
		_daocLogs.getUser().TotalHeals = append(_daocLogs.getUser().TotalHeals, healingInt)
		_daocLogs.getUser().TotalHeals = append(_daocLogs.getUser().TotalHeals, healingInt)

		user := "unknown"
		userStats := _daocLogs.findFriendlyStats(user)
		userStats.TotalHeals = append(userStats.TotalHeals, healingInt)
	}
	match, _ = regexp.MatchString("heal criticals", line)
	if match {
		healing := strings.Split(line, "for an extra ")[1]
		healing = strings.Split(healing, " amount")[0]
		healingInt, _ := strconv.Atoi(healing)
		_daocLogs.getUser().TotalHealsCrits = append(_daocLogs.getUser().TotalHealsCrits, healingInt)
		_daocLogs.getUser().TotalHeals = append(_daocLogs.getUser().TotalHeals, healingInt)
	}
	match, _ = regexp.MatchString("is stunned and cannot move", line)
	if match {
		_daocLogs.getUser().TotalStuns += 1
		user := strings.Split(line, " is ")[0]
		user = strings.Split(user, "] ")[1]
		userStats := _daocLogs.findEnemyStats(user)
		userStats.TotalStuns += 1
	}
}

func (_daocLogs *DaocLogs) regexPets(line string) {
	match, _ := regexp.MatchString("@@", line)
	if match {
		return
	}
	match, _ = regexp.MatchString("The.*casts a spell!", line)
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
	match, _ := regexp.MatchString("@@", line)
	if match {
		return
	}
	match, _ = regexp.MatchString("You gain a total of.*experience points", line)
	if match {
		exp := strings.Split(line, "You gain a total of ")[1]
		exp = strings.Split(exp, " experience")[0]
		exp = strings.ReplaceAll(exp, ",", "")
		expInt, _ := strconv.Atoi(exp)
		_daocLogs.getUser().ExperienceGained = append(_daocLogs.getUser().ExperienceGained, expInt)
	}
	match, _ = regexp.MatchString("You just killed", line)
	if match {
		_daocLogs.getUser().TotalKills += 1
		user := strings.Split(line, "You just killed ")[1]
		user = strings.Split(user, "!")[0]
		userStats := _daocLogs.findEnemyStats(user)
		userStats.TotalDeaths += 1
		// _daocLogs.getUser().UsersKilled = append(_daocLogs.getUser().UsersKilled, user)
	}
	match, _ = regexp.MatchString("You have dued", line)
	if match {
		_daocLogs.getUser().TotalDeaths += 1
	}
}

func (_daocLogs *DaocLogs) regexEnemy(line string) {
	match, _ := regexp.MatchString("@@", line)
	if match {
		return
	}
	match, _ = regexp.MatchString("parries your attack", line)
	if match {
		user := strings.Split(line, " parries")[0]
		user = strings.Split(user, " ")[1]
		userStats := _daocLogs.findEnemyStats(user)
		userStats.ParryTotal += 1
	}
	match, _ = regexp.MatchString("evades your attack", line)
	if match {
		user := strings.Split(line, " evades")[0]
		user = strings.Split(user, " ")[1]
		userStats := _daocLogs.findEnemyStats(user)
		userStats.EvadeTotal += 1
	}
	match, _ = regexp.MatchString("blocks your attack", line)
	if match {
		user := strings.Split(line, " blocks")[0]
		user = strings.Split(user, " ")[1]
		userStats := _daocLogs.findEnemyStats(user)
		userStats.BlockTotal += 1
	}

	match, _ = regexp.MatchString("resists the effect", line)
	if match {
		user := strings.Split(line, " resists")[0]
		user = strings.Split(user, "] ")[1]
		userStats := _daocLogs.findEnemyStats(user)
		userStats.ResistsInTotal += 1
	}
	match, _ = regexp.MatchString("hits you for.*damage ", line)
	if match {
		match, _ = regexp.MatchString("critically hit", line)
		if !match {
			damage := strings.Split(line, "for ")[1]
			damage = strings.Split(damage, " damage")[0]
			damage = strings.Split(damage, " ")[0]
			damageInt, _ := strconv.Atoi(damage)

			user := strings.Split(line, " hits you")[0]
			user = strings.Split(user, "] ")[1]
			userStats := _daocLogs.findEnemyStats(user)

			userStats.MovingDamageTotal = append(userStats.MovingDamageTotal, damageInt)
		}
	}
	match, _ = regexp.MatchString("hits your.*for ", line)
	if match {
		match, _ = regexp.MatchString("critically hits you", line)
		if !match {
			damage := strings.Split(line, "for ")[1]
			damage = strings.Split(damage, " damage")[0]
			damage = strings.Split(damage, " ")[0]
			damageInt, _ := strconv.Atoi(damage)

			user := strings.Split(line, " hits your")[0]
			user = strings.Split(user, "] ")[1]
			userStats := _daocLogs.findEnemyStats(user)

			userStats.MovingDamageTotal = append(userStats.MovingDamageTotal, damageInt)

			armorPiece := strings.Split(line, " for")[0]
			armorPiece = strings.Split(armorPiece, "your ")[1]

			switch armorPiece {
			case "head":
				userStats.ArmorHit.Head = append(userStats.ArmorHit.Head, damageInt)
				userStats.MovingDamageStyles = append(userStats.MovingDamageStyles, damageInt)
			case "torso":
				userStats.ArmorHit.Torso = append(userStats.ArmorHit.Torso, damageInt)
				userStats.MovingDamageStyles = append(userStats.MovingDamageStyles, damageInt)
			case "leg":
				userStats.ArmorHit.Leg = append(userStats.ArmorHit.Leg, damageInt)
				userStats.MovingDamageStyles = append(userStats.MovingDamageStyles, damageInt)
			case "arm":
				userStats.ArmorHit.Arm = append(userStats.ArmorHit.Arm, damageInt)
				userStats.MovingDamageStyles = append(userStats.MovingDamageStyles, damageInt)
			case "hand":
				userStats.ArmorHit.Hand = append(userStats.ArmorHit.Hand, damageInt)
				userStats.MovingDamageStyles = append(userStats.MovingDamageStyles, damageInt)
			case "foot":
				userStats.ArmorHit.Foot = append(userStats.ArmorHit.Foot, damageInt)
				userStats.MovingDamageStyles = append(userStats.MovingDamageStyles, damageInt)
			default:
				break
			}
		}
	}
	match, _ = regexp.MatchString("critically hits you", line)
	if match {
		damage := strings.Split(line, "for an additional ")[1]
		damage = strings.Split(damage, " damage")[0]
		damageInt, _ := strconv.Atoi(damage)
		user := strings.Split(line, " critically")[0]
		user = strings.Split(user, " ")[0]
		playerStats := _daocLogs.findEnemyStats(user)
		playerStats.MovingDamageTotal = append(playerStats.MovingDamageTotal, damageInt)
		playerStats.MovingCritDamage = append(playerStats.MovingCritDamage, damageInt)
	}
}
func (_daocLogs *DaocLogs) regexTime(line string) {
	match, _ := regexp.MatchString("@@", line)
	if match {
		return
	}
	if _daocLogs.getUser().StartTime.IsZero() {
		match, _ = regexp.MatchString("Chat Log Opened", line)
		if match {
			timeObj, err := time.Parse("Mon Jan 02 15:04:05 2006", strings.TrimSuffix(strings.Split(line, ": ")[1], "\r\n"))
			if err != nil {
				fmt.Sprintln("Error:", err)
				return
			}
			_daocLogs.getUser().StartTime = timeObj
		}
	}
	match, _ = regexp.MatchString("Chat Log Closed", line)
	if match {
		timeObj, err := time.Parse("Mon Jan 02 15:04:05 2006", strings.TrimSuffix(strings.Split(line, ": ")[1], "\r\n"))
		if err != nil {
			fmt.Sprintln("Error:", err)
			return
		}
		_daocLogs.getUser().EndTime = timeObj
	}
}
