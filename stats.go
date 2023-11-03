package main

import (
	"fmt"
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
			meleeDamage = fmt.Sprintf("\tMelee Hit: %d\n\tMelee Damage: %d\n", len(_daocLogs.User.MovingDamageStyles)+len(_daocLogs.User.MovingDamageBaseMelee), sumArr(_daocLogs.User.MovingDamageStyles)+sumArr(_daocLogs.User.MovingDamageBaseMelee))
		}
		if len(_daocLogs.User.MovingDamageSpells)+len(_daocLogs.User.MovingExtraDamage) > 0 {
			spellDamage = fmt.Sprintf("\tSpell Hit: %d\n\tSpell Damage: %d\n", len(_daocLogs.User.MovingDamageSpells)+len(_daocLogs.User.MovingExtraDamage), sumArr(_daocLogs.User.MovingDamageSpells)+sumArr(_daocLogs.User.MovingExtraDamage))
		}
		if len(_daocLogs.User.MovingCritDamage) > 0 {
			critDamage = fmt.Sprintf("\tCrit Hit: %d\n\tCrit Damage: %d\n", len(_daocLogs.User.MovingCritDamage), sumArr(_daocLogs.User.MovingCritDamage))
		}

		if _daocLogs.User.ResistsTotal > 0 {
			spellsResists = fmt.Sprintf("\tResits: %d\n", _daocLogs.User.ResistsTotal)
		}
		if _daocLogs.User.MissesTotal > 0 {
			meleeMiss = fmt.Sprintf("\tMisses: %d\n", _daocLogs.User.MissesTotal)
		}
		if _daocLogs.User.SiphonTotal > 0 {
			siphons = fmt.Sprintf("\tSiphons: %d\n", _daocLogs.User.SiphonTotal)
		}
		if _daocLogs.User.TotalKills > 0 {
			kills = fmt.Sprintf("\tKills: %d\n", _daocLogs.User.TotalKills)
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
	if len(_daocLogs.User.TotalSelfHeal)+len(_daocLogs.User.TotalAbsorbed)+len(_daocLogs.User.TotalHeals) > 0 {
		healAndAbsorb := "----- Healing & Absorb -----\n"
		selfHeal := ""
		selfAbsorb := ""
		allHeal := ""
		overHeal := ""
		if len(_daocLogs.User.TotalSelfHeal) > 0 {
			selfHeal = fmt.Sprintf("\tSelf Heals: %d\n", sumArr(_daocLogs.User.TotalSelfHeal))
		}
		if len(_daocLogs.User.TotalHeals) > 0 {
			allHeal = fmt.Sprintf("\tAll Heals: %d\n", sumArr(_daocLogs.User.TotalHeals))
		}
		if len(_daocLogs.User.TotalAbsorbed) > 0 {
			selfAbsorb = fmt.Sprintf("\tAbsorbed: %d\n", sumArr(_daocLogs.User.TotalAbsorbed))
		}
		if _daocLogs.User.OverHeals > 0 {
			overHeal = fmt.Sprintf("\tOverHeal Count: %d\n", _daocLogs.User.OverHeals)
		}
		return healAndAbsorb + selfHeal + allHeal + overHeal + selfAbsorb
	}
	return ""
}

func (_daocLogs *DaocLogs) calculateDensives() string {
	if len(_daocLogs.User.TotalSelfHeal)+len(_daocLogs.User.TotalAbsorbed)+len(_daocLogs.User.TotalHeals) > 0 {
		defensives := "----- Defensives -----\n"
		block := ""
		parry := ""
		evade := ""
		stuns := ""
		if _daocLogs.User.BlockTotal > 0 {
			block = fmt.Sprintf("\tBlock: %d\n", _daocLogs.User.BlockTotal)
		}
		if _daocLogs.User.ParryTotal > 0 {
			parry = fmt.Sprintf("\tParry: %d\n", _daocLogs.User.ParryTotal)
		}
		if _daocLogs.User.EvadeTotal > 0 {
			evade = fmt.Sprintf("\tEvade: %d\n", _daocLogs.User.EvadeTotal)
		}
		if _daocLogs.User.TotalStuns > 0 {
			stuns = fmt.Sprintf("\tStuns: %d\n", _daocLogs.User.TotalStuns)
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
	totalMinutes := int(_daocLogs.User.EndTime.Sub(_daocLogs.User.StartTime).Seconds()) / 60
	totalSeconds := int(_daocLogs.User.EndTime.Sub(_daocLogs.User.StartTime).Seconds()) - (60 * totalMinutes)
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

		return ArmorHitFmt + headHitFmt + torsoHitFmt + armHitFmt + legHitFmt + handHitFmt + footHitFmt
	}
	return ""
}
