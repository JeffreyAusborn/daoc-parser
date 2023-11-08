package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/*
	You cast
	You hit

	You prepare to perform a ___
	You perform your ___ perfectly
	You attack ___ with your ___ and hit for ___ damage! (Damage Modifier: 2668)
	You hit ___ for ___ extra damage!
		You perform your ___ perfectly
		You begin casting a ___ spell
		___ parries you attack
		___ blocks you attack
		___ evades you attack


	You begin casting a ___ spell
	You cast a ___ spell!
		You heal -> heal
		You hit -> spell
		You transfer -> heal

		You move and interrupt
		__ attacks you and your spell is interrupted
*/

var (
	styleName  string
	weaponName string
	spellName  string
	userHit    string
	growthInt  int
	spellStats Ability
)

func checkTempLines(reg string) {
	for _, line := range tempLines {
		match, _ := regexp.MatchString(reg, line)
		if match {
			return
		}
	}
	return
}

func (_daocLogs *DaocLogs) regexOffensive(line string, style bool) {

	match, _ := regexp.MatchString("@@", line)
	if match {
		return
	}

	match, _ = regexp.MatchString("You fail to execute", line)
	if match {
		spellName = ""
		styleName = ""
		return
	}

	match, _ = regexp.MatchString("You begin casting a.*spell", line)
	if match {
		spellName = strings.Split(line, "casting a ")[1]
		spellName = strings.Split(spellName, " spell")[0]
		styleName = ""
		_daocLogs.findSpellStats(spellName)
		return
	}

	match, _ = regexp.MatchString("You fire a.*!", line)
	if match {
		spellName = strings.Split(line, "fire a ")[1]
		spellName = strings.Split(spellName, "!")[0]
		styleName = ""
		_daocLogs.findSpellStats(spellName)
		return
	}

	match, _ = regexp.MatchString("You cast a.*spell", line)
	if match {
		spellName = strings.Split(line, "cast a ")[1]
		spellName = strings.Split(spellName, " spell")[0]
		styleName = ""
		_daocLogs.findSpellStats(spellName)
		return
	}

	match, _ = regexp.MatchString("You prepare to perform a.*", line)
	if match {
		styleName = strings.Split(line, "perform a ")[1]
		styleName = strings.Split(styleName, "!")[0]
		spellName = ""
		_daocLogs.findStyleStats(styleName)
		return
	}

	match, _ = regexp.MatchString("You perform your.*perfectly", line)
	if match {
		spellName = ""
		styleName = strings.Split(line, "perform your ")[1]
		styleName = strings.Split(styleName, " perfectly")[0]
		match, _ = regexp.MatchString("Growth", line)
		if match {
			growthRate := strings.Split(line, ", Growth")[0]
			growthRate = strings.Split(growthRate, "+")[1]
			growthInt, _ = strconv.Atoi(growthRate)
		}
		_daocLogs.findStyleStats(styleName)
		return
	}

	if spellName != "" {
		styleName = ""
		match, _ = regexp.MatchString("You move and interrupt", line)
		if match {
			healExist := _daocLogs.healAbilityExist(spellName)
			if healExist {
				healStats := _daocLogs.findHealStats(spellName)
				userStats := healStats.findUserStats("yourself")
				userStats.Interrupted += 1
			} else {
				spellStats := _daocLogs.findSpellStats(spellName)
				spellStats.Interupts = append(spellStats.Interupts, "self")
				userStats := spellStats.findUserStats("yourself")
				userStats.Interrupted += 1
			}
			spellName = ""
			return
		}

		match, _ = regexp.MatchString("attacks you and your spell is interrupted", line)
		if match {
			user := strings.Split(line, "attacks")[0]
			user = strings.Split(user, "] ")[1]

			healExist := _daocLogs.healAbilityExist(spellName)
			if healExist {
				healStats := _daocLogs.findHealStats(spellName)
				userStats := healStats.findUserStats(user)
				userStats.Interrupted += 1
			} else {
				spellStats := _daocLogs.findSpellStats(spellName)
				userStats := spellStats.findUserStats(user)
				userStats.Interrupted += 1
			}
			spellName = ""
			return
		}

		match, _ = regexp.MatchString("You hit.*for.*damage", line)
		if match {
			match, _ = regexp.MatchString("critically hit", line)
			if !match {
				damage := strings.Split(line, " for ")[1]
				damage = strings.Split(damage, " damage")[0]
				damage = strings.Split(damage, " ")[0]
				damageInt, _ := strconv.Atoi(damage)

				spellStats := _daocLogs.findSpellStats(spellName)
				spellStats.Output = append(spellStats.Output, damageInt)

				user := strings.Split(line, "You hit ")[1]
				user = strings.Split(user, " for")[0]
				userHit = user
				userStats := spellStats.findUserStats(user)
				userStats.MovingDamageReceived = append(userStats.MovingDamageReceived, damageInt)
			}
			return
		}
		match, _ = regexp.MatchString("You gather energy from your surroundings.*", line)
		if match {
			spellStats := _daocLogs.findSpellStats(spellName)
			spellStats.Siphon += 1
			return
		}
	}

	if styleName != "" {
		spellName = ""
		match, _ = regexp.MatchString("you hit.*extra damage", line)
		if match {
			damage := strings.Split(line, " for ")[1]
			damage = strings.Split(damage, " extra")[0]
			damage = strings.Split(damage, " ")[0]
			damageInt, _ := strconv.Atoi(damage)

			styleStats := _daocLogs.findStyleStats(styleName)
			styleStats.ExtraDamage = append(styleStats.ExtraDamage, damageInt)

			user := strings.Split(line, "You hit ")[1]
			user = strings.Split(user, " for")[0]
			userStats := styleStats.findUserStats(user)
			userStats.MovingDamageReceived = append(userStats.MovingDamageReceived, damageInt)
			if weaponName != "" {
				weaponStats := styleStats.findWeaponStats(weaponName)
				weaponStats.ExtraDamage = append(weaponStats.ExtraDamage, damageInt)
			}
			return
		}

		match, _ = regexp.MatchString("You hit.*for.*damage", line)
		if match {
			match, _ = regexp.MatchString("critically hit", line)
			if !match {
				damage := strings.Split(line, " for ")[1]
				damage = strings.Split(damage, " damage")[0]
				damage = strings.Split(damage, " ")[0]
				damageInt, _ := strconv.Atoi(damage)
				styleStats := _daocLogs.findStyleStats(styleName)
				styleStats.ExtraDamage = append(styleStats.ExtraDamage, damageInt)

				user := strings.Split(line, "You hit ")[1]
				user = strings.Split(user, " for")[0]
				userHit = user
				userStats := styleStats.findUserStats(user)
				userStats.MovingDamageReceived = append(userStats.MovingDamageReceived, damageInt)

				if weaponName != "" {
					weaponStats := styleStats.findWeaponStats(weaponName)
					weaponStats.ExtraDamage = append(weaponStats.ExtraDamage, damageInt)
				}
			}
			return
		}

		match, _ = regexp.MatchString("You attack.*with your.*and hit for.*damage", line)
		if match {
			match, _ = regexp.MatchString("critically hit", line)
			if !match {
				damage := strings.Split(line, "and hit for ")[1]
				damage = strings.Split(damage, " damage")[0]
				damage = strings.Split(damage, " ")[0]
				damageInt, _ := strconv.Atoi(damage)

				user := strings.Split(line, "You attack ")[1]
				user = strings.Split(user, " with your")[0]
				userStats := _daocLogs.findEnemyStats(user)
				userStats.MovingDamageReceived = append(userStats.MovingDamageReceived, damageInt)

				weaponName = strings.Split(line, "with your ")[1]
				weaponName = strings.Split(weaponName, " and hit")[0]

				styleStats := _daocLogs.findStyleStats(styleName)
				styleStats.Output = append(styleStats.Output, damageInt)

				weaponStats := styleStats.findWeaponStats(weaponName)
				weaponStats.Output = append(weaponStats.Output, damageInt)

				if growthInt > 0 {
					styleStats.GrowthRate = append(styleStats.GrowthRate, growthInt)
				}
				if growthInt > 0 {
					weaponStats.GrowthRate = append(weaponStats.GrowthRate, growthInt)
				}
			}
			return
		}
		match, _ = regexp.MatchString("You perform your.*perfectly", line)
		if match {
			styleName = strings.Split(line, "perform your ")[1]
			styleName = strings.Split(styleName, " perfectly")[0]
			match, _ = regexp.MatchString("Growth", line)
			if match {
				growthRate := strings.Split(line, ", Growth")[0]
				growthRate = strings.Split(growthRate, "+")[1]
				growthInt, _ = strconv.Atoi(growthRate)
			}
		}
	}

	match, _ = regexp.MatchString("You attack.*with your.*and hit for.*damage", line)
	if match {
		match, _ = regexp.MatchString("critically hit", line)
		if !match {
			damage := strings.Split(line, "and hit for ")[1]
			damage = strings.Split(damage, " damage")[0]
			damage = strings.Split(damage, " ")[0]
			damageInt, _ := strconv.Atoi(damage)

			user := strings.Split(line, "You attack ")[1]
			user = strings.Split(user, " with your")[0]
			userStats := _daocLogs.findEnemyStats(user)
			userStats.MovingDamageReceived = append(userStats.MovingDamageReceived, damageInt)

			weaponName = strings.Split(line, "with your ")[1]
			weaponName = strings.Split(weaponName, " and hit")[0]

			styleStats := _daocLogs.findStyleStats("base")
			styleStats.Output = append(styleStats.Output, damageInt)

			weaponStats := styleStats.findWeaponStats(weaponName)
			weaponStats.Output = append(weaponStats.Output, damageInt)
		}
	}

	// dots and pets (theur) have similar text, need to phase out pet ones
	match, _ = regexp.MatchString("Your.*hits.*for.*damage", line)
	if match {
		damage := ""
		damageInt := 0
		match, _ = regexp.MatchString("critically hits", line)
		if match {
			damage = strings.Split(line, "additional ")[1]
			damage = strings.Split(damage, " damage")[0]
			damage = strings.Split(damage, " ")[0]
			damageInt, _ = strconv.Atoi(damage)
		} else {
			damage = strings.Split(line, " for ")[1]
			damage = strings.Split(damage, " damage")[0]
			damage = strings.Split(damage, " ")[0]
			damageInt, _ = strconv.Atoi(damage)
		}

		spellName = strings.Split(line, "Your ")[1]
		user := ""

		match, _ = regexp.MatchString("attacks", line)
		matchShoot, _ := regexp.MatchString("Your.*shoots.*for.*damage", line)
		if match {
			spellName = strings.Split(spellName, " attacks")[0]
			user = strings.Split(line, "attacks ")[1]
			user = strings.Split(user, " and")[0]
		} else if matchShoot {
			spellName = strings.Split(spellName, " shoots")[0]
			user = strings.Split(line, " and hits")[0]
			user = strings.Split(user, " shoots ")[1]
		} else {
			spellName = strings.Split(spellName, " hits")[0]
			user = strings.Split(line, "hits ")[1]
			user = strings.Split(user, " for")[0]
		}

		spellStats := _daocLogs.findSpellStats(spellName)
		spellStats.Output = append(spellStats.Output, damageInt)

		match, _ = regexp.MatchString("critically hits", line)
		// Zombie servant crit doesn't have user name - will need to keep an eye on this for other pet crits
		if !match {
			userStats := _daocLogs.findEnemyStats(user)
			userStats.MovingDamageReceived = append(userStats.MovingDamageReceived, damageInt)
		}
	}

	match, _ = regexp.MatchString("You critically hit", line)
	if match {
		damage := strings.Split(line, "for an additional ")[1]
		damage = strings.Split(damage, " damage")[0]
		damageInt, _ := strconv.Atoi(damage)

		user := strings.Split(line, "hit ")[1]
		user = strings.Split(user, " for")[0]

		if spellName != "" {
			spellStats := _daocLogs.findSpellStats(spellName)
			spellStats.Crit = append(spellStats.Crit, damageInt)
			userStats := spellStats.findUserStats(userHit)
			userStats.MovingDamageReceived = append(userStats.MovingDamageReceived, damageInt)
		} else if styleName != "" {
			styleStats := _daocLogs.findStyleStats(styleName)
			styleStats.Crit = append(styleStats.Crit, damageInt)
			userStats := styleStats.findUserStats(userHit)
			userStats.MovingDamageReceived = append(userStats.MovingDamageReceived, damageInt)
		}

		if weaponName != "" {
			if styleName == "" {
				styleName = "base"
			}

			styleStats := _daocLogs.findStyleStats(styleName)
			weaponStats := styleStats.findWeaponStats(weaponName)
			weaponStats.Crit = append(weaponStats.Crit, damageInt)
		}
		userHit = "unknown"
	}
}

func (_daocLogs *DaocLogs) regexDefensives(line string) {
	match, _ := regexp.MatchString("@@", line)
	if match {
		return
	}
	match, _ = regexp.MatchString("you block", line)
	if match {
		_daocLogs.getUser().BlockTotal += 1
		return
	}
	match, _ = regexp.MatchString("absorbed by a magical barrier", line)
	if match {
		_daocLogs.getUser().BladeturnTotal += 1
		return
	}
	match, _ = regexp.MatchString("you evade", line)
	if match {
		_daocLogs.getUser().EvadeTotal += 1
		return
	}
	match, _ = regexp.MatchString("you parry", line)
	if match {
		_daocLogs.getUser().ParryTotal += 1
		return
	}
	match, _ = regexp.MatchString("Your ablative absorbs", line)
	if match {
		absorb := strings.Split(line, "ablative absorbs ")[1]
		absorb = strings.Split(absorb, " damage")[0]
		absorbInt, _ := strconv.Atoi(absorb)
		_daocLogs.getUser().TotalAbsorbed = append(_daocLogs.getUser().TotalAbsorbed, absorbInt)
		return
	}
	match, _ = regexp.MatchString("You resist the effect", line)
	if match {
		_daocLogs.getUser().ResistsInTotal += 1
		return
	}
}

func (_daocLogs *DaocLogs) regexSupport(line string) {
	match, _ := regexp.MatchString("@@", line)
	if match {
		return
	}

	if spellName != "" {
		match, _ = regexp.MatchString("You heal.*for.*hit points", line)
		if match {
			healing := strings.Split(line, "for ")[1]
			healing = strings.Split(healing, " hit points")[0]
			healingInt, _ := strconv.Atoi(healing)
			_daocLogs.getUser().TotalHeals = append(_daocLogs.getUser().TotalHeals, healingInt)

			user := strings.Split(line, " for")[0]
			user = strings.Split(user, "heal ")[1]

			healStats := _daocLogs.findHealStats(spellName)
			healStats.Output = append(healStats.Output, healingInt)

			userStats := healStats.findUserStats(user)
			userStats.TotalHeals = append(userStats.TotalHeals, healingInt)
			return

		}
		match, _ = regexp.MatchString("You transfer.*hit points", line)
		if match {
			healing := strings.Split(line, "transfer ")[1]
			healing = strings.Split(healing, " hit points")[0]
			healingInt, _ := strconv.Atoi(healing)
			_daocLogs.getUser().TotalHeals = append(_daocLogs.getUser().TotalHeals, healingInt)

			user := "unknown"

			healStats := _daocLogs.findHealStats(spellName)
			healStats.Output = append(healStats.Output, healingInt)

			userStats := healStats.findUserStats(user)
			userStats.TotalHeals = append(userStats.TotalHeals, healingInt)
		}
		match, _ = regexp.MatchString("heal criticals", line)
		if match {
			healing := strings.Split(line, "for an extra ")[1]
			healing = strings.Split(healing, " amount")[0]
			healingInt, _ := strconv.Atoi(healing)

			user := strings.Split(line, " for")[0]
			user = strings.Split(user, "heal ")[1]

			healStats := _daocLogs.findHealStats(spellName)
			healStats.Crit = append(healStats.Crit, healingInt)

			userStats := healStats.findUserStats(user)
			userStats.TotalHeals = append(userStats.TotalHeals, healingInt)
		}
		match, _ = regexp.MatchString("You steal.*hit points", line)
		if match {
			heal := strings.Split(line, "You steal ")[1]
			heal = strings.Split(heal, " hit points")[0]
			healInt, _ := strconv.Atoi(heal)

			healStats := _daocLogs.findHealStats(spellName)
			healStats.Output = append(healStats.Output, healInt)
			playerStats := healStats.findUserStats("yourself")
			playerStats.TotalHeals = append(playerStats.TotalHeals, healInt)
		}
	}

	// TODO
	if styleName != "" {
		match, _ = regexp.MatchString("You heal yourself for", line)
		if match {
			healing := strings.Split(line, "You heal yourself for ")[1]
			healing = strings.Split(healing, " hit points")[0]
			healingInt, _ := strconv.Atoi(healing)

			healStats := _daocLogs.findHealStats(styleName)
			healStats.Output = append(healStats.Output, healingInt)

			userStats := healStats.findUserStats("yourself")
			userStats.TotalHeals = append(userStats.TotalHeals, healingInt)
			styleName = ""
			return
		}
		match, _ = regexp.MatchString("You heal.*for.*hit points", line)
		if match {
			healing := strings.Split(line, "for ")[1]
			healing = strings.Split(healing, " hit points")[0]
			healingInt, _ := strconv.Atoi(healing)
			_daocLogs.getUser().TotalHeals = append(_daocLogs.getUser().TotalHeals, healingInt)

			user := strings.Split(line, " for")[0]
			user = strings.Split(user, "heal ")[1]

			healStats := _daocLogs.findHealStats(styleName)
			healStats.Output = append(healStats.Output, healingInt)

			userStats := healStats.findUserStats(user)
			userStats.TotalHeals = append(userStats.TotalHeals, healingInt)
			return

		}
		match, _ = regexp.MatchString("You steal.*hit points", line)
		if match {
			heal := strings.Split(line, "You steal ")[1]
			heal = strings.Split(heal, " hit points")[0]
			healInt, _ := strconv.Atoi(heal)

			healStats := _daocLogs.findHealStats(styleName)
			healStats.Output = append(healStats.Output, healInt)
			playerStats := healStats.findUserStats("yourself")
			playerStats.TotalHeals = append(playerStats.TotalHeals, healInt)
		}
	}

	match, _ = regexp.MatchString("You heal yourself for", line)
	if match {
		healing := strings.Split(line, "You heal yourself for ")[1]
		healing = strings.Split(healing, " hit points")[0]
		healingInt, _ := strconv.Atoi(healing)

		healStats := _daocLogs.findHealStats("unknown")
		healStats.Output = append(healStats.Output, healingInt)

		userStats := healStats.findUserStats("yourself")
		userStats.TotalHeals = append(userStats.TotalHeals, healingInt)
		styleName = ""
		return
	}

	match, _ = regexp.MatchString("is stunned and cannot move", line)
	if match {
		_daocLogs.getUser().TotalStuns += 1
		user := strings.Split(line, " is ")[0]
		user = strings.Split(user, "] ")[1]

		if spellName != "" {
			spellStats := _daocLogs.findSpellStats(spellName)
			spellStats.Stunned += 1
		} else if styleName != "" {
			styleStats := _daocLogs.findStyleStats(styleName)
			styleStats.Stunned += 1
		}

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
		damage := strings.Split(line, " for ")[1]
		damage = strings.Split(damage, " damage")[0]
		damage = strings.Split(damage, " ")[0]
		damageInt, _ := strconv.Atoi(damage)

		tmpSpellName := strings.Split(line, "Your ")[1]
		user := ""

		match, _ = regexp.MatchString("attacks", line)
		if match {
			tmpSpellName = strings.Split(tmpSpellName, " attacks")[0]
			user = strings.Split(line, "attacks ")[1]
			user = strings.Split(user, " and")[0]
		} else {
			tmpSpellName = strings.Split(tmpSpellName, " hits")[0]
			user = strings.Split(line, "hits ")[1]
			user = strings.Split(user, " for")[0]
		}

		spellStats := _daocLogs.findSpellStats(tmpSpellName)
		spellStats.Output = append(spellStats.Output, damageInt)
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
	}
	match, _ = regexp.MatchString("You have died", line)
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
		if styleName != "" {
			styleStats := _daocLogs.findStyleStats(styleName)
			styleStats.Parried += 1
		}
		styleName = ""
	}
	match, _ = regexp.MatchString("evades your attack", line)
	if match {
		user := strings.Split(line, " evades")[0]
		user = strings.Split(user, " ")[1]
		userStats := _daocLogs.findEnemyStats(user)
		userStats.EvadeTotal += 1
		if styleName != "" {
			styleStats := _daocLogs.findStyleStats(styleName)
			styleStats.Evaded += 1
		}
		styleName = ""
	}
	match, _ = regexp.MatchString("blocks your attack", line)
	if match {
		user := strings.Split(line, " blocks")[0]
		user = strings.Split(user, " ")[1]
		userStats := _daocLogs.findEnemyStats(user)
		userStats.BlockTotal += 1
		if styleName != "" {
			styleStats := _daocLogs.findStyleStats(styleName)
			styleStats.Blocked += 1
		}
		styleName = ""
	}

	match, _ = regexp.MatchString("resists the effect", line)
	if match {
		user := strings.Split(line, " resists")[0]
		user = strings.Split(user, "] ")[1]
		userStats := _daocLogs.findEnemyStats(user)
		userStats.ResistsInTotal += 1

		if spellName != "" {
			spellStats := _daocLogs.findSpellStats(spellName)
			spellStats.Resists += 1
		}
		return
	}
	// pets?
	match, _ = regexp.MatchString("resists your.*effect", line)
	if match {
		user := strings.Split(line, " resists")[0]
		user = strings.Split(user, "] ")[1]

		tempSpellName := strings.Split(line, "resists your ")[1]
		tempSpellName = strings.Split(tempSpellName, " effect")[0]
		tempSpellName = strings.Split(tempSpellName, "'")[0]
		spellStats := _daocLogs.findSpellStats(tempSpellName)
		spellStats.Resists += 1
		userStats := spellStats.findUserStats(user)
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

			armorPiece := strings.Split(line, " for ")[0]
			armorPiece = strings.Split(armorPiece, " your ")[1]

			switch armorPiece {
			case "head":
				userStats.ArmorHit.Head = append(userStats.ArmorHit.Head, damageInt)
				// userStats.MovingDamageStyles = append(userStats.MovingDamageStyles, damageInt)
			case "torso":
				userStats.ArmorHit.Torso = append(userStats.ArmorHit.Torso, damageInt)
				// userStats.MovingDamageStyles = append(userStats.MovingDamageStyles, damageInt)
			case "leg":
				userStats.ArmorHit.Leg = append(userStats.ArmorHit.Leg, damageInt)
				// userStats.MovingDamageStyles = append(userStats.MovingDamageStyles, damageInt)
			case "arm":
				userStats.ArmorHit.Arm = append(userStats.ArmorHit.Arm, damageInt)
				// userStats.MovingDamageStyles = append(userStats.MovingDamageStyles, damageInt)
			case "hand":
				userStats.ArmorHit.Hand = append(userStats.ArmorHit.Hand, damageInt)
				// userStats.MovingDamageStyles = append(userStats.MovingDamageStyles, damageInt)
			case "foot":
				userStats.ArmorHit.Foot = append(userStats.ArmorHit.Foot, damageInt)
				// userStats.MovingDamageStyles = append(userStats.MovingDamageStyles, damageInt)
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
