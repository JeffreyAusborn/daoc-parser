package main

import (
	"fmt"
	"strings"
)

func (_daocLogs *DaocLogs) calculateDamageOut() string {
	if len(_daocLogs.User.MovingDamageTotal) > 0 {
		damageIn := "----- Damage Out -----\n"
		meleeDamage := ""
		spellDamage := ""
		critDamage := ""
		spellsResists := ""
		meleeMiss := ""
		siphons := ""
		kills := ""
		if len(_daocLogs.User.MovingDamageStyles)+len(_daocLogs.User.MovingDamageBaseMelee) > 0 {
			meleeDamage = fmt.Sprintf("Melee Hit: %d\nMelee Damage: %d\n", len(_daocLogs.User.MovingDamageStyles)+len(_daocLogs.User.MovingDamageBaseMelee), sumArr(_daocLogs.User.MovingDamageStyles)+sumArr(_daocLogs.User.MovingDamageBaseMelee))
		}
		if len(_daocLogs.User.MovingDamageSpells)+len(_daocLogs.User.MovingExtraDamage) > 0 {
			spellDamage = fmt.Sprintf("Spell Hit: %d\nSpell Damage: %d\n", len(_daocLogs.User.MovingDamageSpells)+len(_daocLogs.User.MovingExtraDamage), sumArr(_daocLogs.User.MovingDamageSpells)+sumArr(_daocLogs.User.MovingExtraDamage))
		}
		if len(_daocLogs.User.MovingCritDamage) > 0 {
			critDamage = fmt.Sprintf("Crit Hit: %d\nCrit Damage: %d\n", len(_daocLogs.User.MovingCritDamage), sumArr(_daocLogs.User.MovingCritDamage))
		}

		if _daocLogs.User.ResistsOutTotal > 0 {
			spellsResists = fmt.Sprintf("Resits: %d\n", _daocLogs.User.ResistsOutTotal)
		}
		if _daocLogs.User.MissesTotal > 0 {
			meleeMiss = fmt.Sprintf("Misses: %d\n", _daocLogs.User.MissesTotal)
		}
		if _daocLogs.User.SiphonTotal > 0 {
			siphons = fmt.Sprintf("Siphons: %d\n", _daocLogs.User.SiphonTotal)
		}
		if _daocLogs.User.TotalKills > 0 {
			kills = fmt.Sprintf("Kills: %d\n", _daocLogs.User.TotalKills)
		}
		return damageIn + meleeDamage + spellDamage + critDamage + spellsResists + meleeMiss + siphons + kills
	}
	return ""
}

func (_daocLogs *DaocLogs) calculateDamageIn() string {
	totalMeleeDamage := []int{}
	totalAllDamage := []int{}

	for _, user := range _daocLogs.Enemy {
		totalMeleeDamage = append(totalMeleeDamage, user.MovingDamageStyles...)
		totalAllDamage = append(totalAllDamage, user.MovingDamageTotal...)
	}

	totalDamage := sumArr(totalAllDamage)
	totalDamageMelee := sumArr(totalMeleeDamage)
	totalDamageSpell := totalDamage - totalDamageMelee

	if totalDamage > 0 {
		damageIn := "----- Damage In -----\n"
		meleeDamage := ""
		spellDamage := ""
		if totalDamageMelee > 0 {
			meleeDamage = fmt.Sprintf("Melee Hit: %d\nMelee Damage: %d\n", len(totalMeleeDamage), totalDamageMelee)
		}
		if totalDamageSpell > 0 {
			spellDamage = fmt.Sprintf("Spell Hit: %d\nSpell Damage: %d\n", len(totalAllDamage)-len(totalMeleeDamage), totalDamageSpell)
		}
		if _daocLogs.getUser().TotalDeaths > 0 {
			spellDamage = fmt.Sprintf("Deaths: %d\n", daocLogs.getUser().TotalDeaths)
		}
		return damageIn + meleeDamage + spellDamage
	}
	return ""
}

func (_daocLogs *DaocLogs) calculateHeal() string {
	if len(_daocLogs.User.TotalSelfHeal)+len(_daocLogs.User.TotalAbsorbed)+len(_daocLogs.User.TotalHeals) > 0 {
		healAndAbsorb := "----- Healing & Absorb -----\n"
		selfHeal := ""
		selfAbsorb := ""
		allHeal := ""
		overHeal := ""
		healCrits := ""
		if len(_daocLogs.User.TotalSelfHeal) > 0 {
			selfHeal = fmt.Sprintf("Self Heals: %d\n", sumArr(_daocLogs.User.TotalSelfHeal))
		}
		if len(_daocLogs.User.TotalHeals) > 0 {
			allHeal = fmt.Sprintf("All Heals: %d\n", sumArr(_daocLogs.User.TotalHeals))
		}
		if len(_daocLogs.User.TotalHealsCrits) > 0 {
			healCrits = fmt.Sprintf("Heal Crits: %d\n", sumArr(_daocLogs.User.TotalHealsCrits))
		}
		if len(_daocLogs.User.TotalAbsorbed) > 0 {
			selfAbsorb = fmt.Sprintf("Absorbed: %d\n", sumArr(_daocLogs.User.TotalAbsorbed))
		}
		if _daocLogs.User.OverHeals > 0 {
			overHeal = fmt.Sprintf("OverHeal Count: %d\n", _daocLogs.User.OverHeals)
		}
		return healAndAbsorb + selfHeal + allHeal + healCrits + overHeal + selfAbsorb
	}
	return ""
}

func (_daocLogs *DaocLogs) calculateDensives() string {
	if _daocLogs.User.BlockTotal+_daocLogs.User.ParryTotal+_daocLogs.User.EvadeTotal+_daocLogs.User.TotalStuns+_daocLogs.User.ResistsInTotal > 0 {
		defensives := "----- Defensives -----\n"
		block := ""
		parry := ""
		evade := ""
		stuns := ""
		resists := ""
		if _daocLogs.User.BlockTotal > 0 {
			block = fmt.Sprintf("Block: %d\n", _daocLogs.User.BlockTotal)
		}
		if _daocLogs.User.ParryTotal > 0 {
			parry = fmt.Sprintf("Parry: %d\n", _daocLogs.User.ParryTotal)
		}
		if _daocLogs.User.EvadeTotal > 0 {
			evade = fmt.Sprintf("Evade: %d\n", _daocLogs.User.EvadeTotal)
		}
		if _daocLogs.User.TotalStuns > 0 {
			stuns = fmt.Sprintf("Stuns: %d\n", _daocLogs.User.TotalStuns)
		}
		if _daocLogs.User.TotalStuns > 0 {
			resists = fmt.Sprintf("Resists: %d\n", _daocLogs.User.ResistsInTotal)
		}
		return defensives + block + parry + evade + stuns + resists
	}
	return ""
}

func (_daocLogs *DaocLogs) calculateEnemyDensives() string {
	blocks := 0
	evades := 0
	parries := 0
	for _, user := range _daocLogs.Enemy {
		blocks += user.BlockTotal
		evades += user.EvadeTotal
		parries += user.ParryTotal
	}

	if blocks+parries+evades > 0 {
		defensives := "----- Enemy Defensives -----\n"
		block := ""
		parry := ""
		evade := ""
		if blocks > 0 {
			block = fmt.Sprintf("Block: %d\n", blocks)
		}
		if parries > 0 {
			parry = fmt.Sprintf("Parry: %d\n", parries)
		}
		if evades > 0 {
			evade = fmt.Sprintf("Evade: %d\n", evades)
		}
		return defensives + block + parry + evade
	}
	return ""
}

func (_daocLogs *DaocLogs) getCombativeUsers() string {
	users := []string{}
	for _, user := range _daocLogs.Enemy {
		users = append(users, user.UserName)
	}

	if len(users) > 0 {
		combative := "----- Combatives -----\n"
		combativeUsers := fmt.Sprintf("%s\n", strings.Join(dedupe(users), "\n"))
		return combative + combativeUsers
	}
	return ""
}

func (_daocLogs *DaocLogs) calculateTime() string {
	totalMinutes := int(_daocLogs.User.EndTime.Sub(_daocLogs.User.StartTime).Seconds()) / 60
	totalSeconds := int(_daocLogs.User.EndTime.Sub(_daocLogs.User.StartTime).Seconds()) - (60 * totalMinutes)
	return fmt.Sprintf("----- Total Time -----\n%d minutes and %d seconds\n", totalMinutes, totalSeconds)
}

func (_daocLogs *DaocLogs) calculateArmorhits() string {
	head := []int{}
	torso := []int{}
	arm := []int{}
	leg := []int{}
	hand := []int{}
	foot := []int{}
	for _, user := range _daocLogs.Enemy {
		head = append(head, user.ArmorHit.Head...)
		torso = append(torso, user.ArmorHit.Torso...)
		arm = append(arm, user.ArmorHit.Arm...)
		leg = append(leg, user.ArmorHit.Leg...)
		hand = append(hand, user.ArmorHit.Hand...)
		foot = append(foot, user.ArmorHit.Foot...)
	}
	if len(head)+len(torso)+len(arm)+len(leg)+len(hand)+len(foot) > 0 {
		headHitFmt := ""
		torsoHitFmt := ""
		armHitFmt := ""
		legHitFmt := ""
		handHitFmt := ""
		footHitFmt := ""
		ArmorHitFmt := "----- Armor Damaged -----\n"
		if len(head) > 0 {
			headHitFmt = fmt.Sprintf("Head Hit: %d\nHead Damage: %d\n", len(head), sumArr(head))
		}
		if len(torso) > 0 {
			torsoHitFmt = fmt.Sprintf("Torso Hit: %d\nTorso Damage: %d\n", len(torso), sumArr(torso))
		}
		if len(arm) > 0 {
			armHitFmt = fmt.Sprintf("Arm Hit: %d\nArm Damage: %d\n", len(arm), sumArr(arm))
		}
		if len(leg) > 0 {
			legHitFmt = fmt.Sprintf("Leg Hit: %d\nLeg Damage: %d\n", len(leg), sumArr(leg))
		}
		if len(hand) > 0 {
			handHitFmt = fmt.Sprintf("Hand Hit: %d\nHand Damage: %d\n", len(hand), sumArr(hand))
		}
		if len(foot) > 0 {
			footHitFmt = fmt.Sprintf("Foot Hit: %d\nFoot Damage: %d\n", len(foot), sumArr(foot))
		}

		return ArmorHitFmt + headHitFmt + torsoHitFmt + armHitFmt + legHitFmt + handHitFmt + footHitFmt
	}
	return ""
}
