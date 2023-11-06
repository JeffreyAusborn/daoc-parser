package main

import (
	"fmt"
)

func (_daocLogs *DaocLogs) calculateDamageOut() []string {
	listItems := []string{}
	if len(_daocLogs.User.MovingDamageTotal) > 0 {
		listItems = append(listItems, "\t\t----- Damage Out -----")
		if len(_daocLogs.User.MovingDamageStyles)+len(_daocLogs.User.MovingDamageBaseMelee) > 0 {
			listItems = append(listItems, fmt.Sprintf("Melee Hit: %d\n", len(_daocLogs.User.MovingDamageStyles)+len(_daocLogs.User.MovingDamageBaseMelee)))
			listItems = append(listItems, fmt.Sprintf("Melee Damage: %d\n", sumArr(_daocLogs.User.MovingDamageStyles)+sumArr(_daocLogs.User.MovingDamageBaseMelee)))
		}
		if len(_daocLogs.User.MovingDamageSpells)+len(_daocLogs.User.MovingExtraDamage) > 0 {
			listItems = append(listItems, fmt.Sprintf("Spell Hit: %d", len(_daocLogs.User.MovingDamageSpells)+len(_daocLogs.User.MovingExtraDamage)))
			listItems = append(listItems, fmt.Sprintf("Spell Damage: %d\n", sumArr(_daocLogs.User.MovingDamageSpells)+sumArr(_daocLogs.User.MovingExtraDamage)))
		}
		if len(_daocLogs.User.MovingCritDamage) > 0 {
			listItems = append(listItems, fmt.Sprintf("Crit Hit: %d", len(_daocLogs.User.MovingCritDamage)))
			listItems = append(listItems, fmt.Sprintf("Crit Damage: %d\n", sumArr(_daocLogs.User.MovingCritDamage)))
		}
		if _daocLogs.User.ResistsOutTotal > 0 {
			listItems = append(listItems, fmt.Sprintf("Resits: %d\n", _daocLogs.User.ResistsOutTotal))
		}
		if _daocLogs.User.MissesTotal > 0 {
			listItems = append(listItems, fmt.Sprintf("Misses: %d\n", _daocLogs.User.MissesTotal))
		}
		if _daocLogs.User.SiphonTotal > 0 {
			listItems = append(listItems, fmt.Sprintf("Siphons: %d\n", _daocLogs.User.SiphonTotal))
		}
		if _daocLogs.User.TotalKills > 0 {
			listItems = append(listItems, fmt.Sprintf("Kills: %d\n", _daocLogs.User.TotalKills))
		}

		if len(_daocLogs.getUser().Styles) > 0 {
			listItems = append(listItems, "\t\t----- Styles -----")
			for _, style := range _daocLogs.getUser().Styles {
				minG, maxG := getMinAndMax(style.GrowtRate)
				minD, maxD := getMinAndMax(style.Output)
				listItems = append(listItems, fmt.Sprintf("\t----- %s -----", style.Name))
				listItems = append(listItems, fmt.Sprintf("Damage: %d\n", sumArr(style.Output)))
				listItems = append(listItems, fmt.Sprintf("Damage Min: %d\n", minD))
				listItems = append(listItems, fmt.Sprintf("Damage Max: %d\n", maxD))
				listItems = append(listItems, fmt.Sprintf("Damage Average: %d\n", sumArr(style.Output)/len(style.Output)))
				listItems = append(listItems, fmt.Sprintf("Growth Rate Min: %d\n", minG))
				listItems = append(listItems, fmt.Sprintf("Growth Rate Max: %d\n", maxG))
				listItems = append(listItems, fmt.Sprintf("Growth Rate Average: %d\n", sumArr(style.GrowtRate)/len(style.GrowtRate)))
			}
		}

		if len(_daocLogs.getUser().Spells) > 0 {
			listItems = append(listItems, "\t\t----- Pets and Dots -----")
			for _, spell := range _daocLogs.getUser().DotsNPets {
				listItems = append(listItems, fmt.Sprintf("\t----- %s -----", spell.Name))
				listItems = append(listItems, fmt.Sprintf("Damage: %d\n", sumArr(spell.Output)))
			}
		}
	}
	return listItems
}

func (_daocLogs *DaocLogs) calculateDamageIn() []string {
	listItems := []string{}
	totalAllDamage := 0

	for _, user := range _daocLogs.Enemy {
		totalAllDamage += sumArr(user.MovingDamageTotal)
	}
	if totalAllDamage > 0 {
		listItems = append(listItems, "\t\t----- Damage In -----\n")
		listItems = append(listItems, fmt.Sprintf("Total Damage: %d\n", totalAllDamage))
		for _, user := range _daocLogs.Enemy {
			listItems = append(listItems, fmt.Sprintf("\t----- %s -----", user.UserName))
			listItems = append(listItems, fmt.Sprintf("Hits: %d", len(user.MovingDamageTotal)))
			listItems = append(listItems, fmt.Sprintf("Damage: %d\n", sumArr(user.MovingDamageTotal)))
		}
	}

	if _daocLogs.getUser().TotalDeaths > 0 {
		listItems = append(listItems, fmt.Sprintf("Deaths: %d\n", daocLogs.getUser().TotalDeaths))
	}
	return listItems
}

func (_daocLogs *DaocLogs) calculateHeal() []string {
	listItems := []string{}
	if len(_daocLogs.User.TotalSelfHeal)+len(_daocLogs.User.TotalAbsorbed)+len(_daocLogs.User.TotalHeals) > 0 {
		listItems = append(listItems, "\t\t----- Healing & Self Absorb -----\n")
		if len(_daocLogs.User.TotalHeals) > 0 {
			listItems = append(listItems, "\t----- Total Heals -----\n")
			min, max := getMinAndMax(_daocLogs.User.TotalHeals)
			listItems = append(listItems, fmt.Sprintf("All Heals: %d\n", sumArr(_daocLogs.User.TotalHeals)))
			listItems = append(listItems, fmt.Sprintf("Minimum Heals: %d\n", min))
			listItems = append(listItems, fmt.Sprintf("Maximum Heals: %d\n", max))
			listItems = append(listItems, fmt.Sprintf("Average Heals: %d\n", sumArr(_daocLogs.User.TotalHeals)/len(_daocLogs.User.TotalHeals)))
		}
		if _daocLogs.User.OverHeals > 0 {
			listItems = append(listItems, fmt.Sprintf("OverHeal Count: %d\n", _daocLogs.User.OverHeals))
		}
		if len(_daocLogs.User.TotalSelfHeal) > 0 {
			min, max := getMinAndMax(_daocLogs.User.TotalSelfHeal)
			listItems = append(listItems, "\t----- Self Heals -----\n")
			listItems = append(listItems, fmt.Sprintf("Total Self Heals: %d\n", sumArr(_daocLogs.User.TotalSelfHeal)))
			listItems = append(listItems, fmt.Sprintf("Minimum Self Heals: %d\n", min))
			listItems = append(listItems, fmt.Sprintf("Maximum Self Heals: %d\n", max))
			listItems = append(listItems, fmt.Sprintf("Average Self Heals: %d\n", sumArr(_daocLogs.User.TotalSelfHeal)/len(_daocLogs.User.TotalSelfHeal)))
		}
		if len(_daocLogs.User.TotalHealsCrits) > 0 {
			listItems = append(listItems, "\t----- Crit Heals -----\n")
			min, max := getMinAndMax(_daocLogs.User.TotalHealsCrits)
			listItems = append(listItems, fmt.Sprintf("Heal Crits: %d\n", sumArr(_daocLogs.User.TotalHealsCrits)))
			listItems = append(listItems, fmt.Sprintf("Minimum Crit Heals: %d\n", min))
			listItems = append(listItems, fmt.Sprintf("Maximum Crit Heals: %d\n", max))
			listItems = append(listItems, fmt.Sprintf("Average Crit Heals: %d\n", sumArr(_daocLogs.User.TotalHealsCrits)/len(_daocLogs.User.TotalHealsCrits)))
		}
		if len(_daocLogs.User.TotalAbsorbed) > 0 {
			listItems = append(listItems, "\t----- Absorbs -----\n")
			min, max := getMinAndMax(_daocLogs.User.TotalAbsorbed)
			listItems = append(listItems, fmt.Sprintf("Absorbed: %d\n", sumArr(_daocLogs.User.TotalAbsorbed)))
			listItems = append(listItems, fmt.Sprintf("Minimum Absorbed: %d\n", min))
			listItems = append(listItems, fmt.Sprintf("Maximum Absorbed: %d\n", max))
			listItems = append(listItems, fmt.Sprintf("Average Absorbed: %d\n", sumArr(_daocLogs.User.TotalAbsorbed)/len(_daocLogs.User.TotalAbsorbed)))
		}
		for _, user := range _daocLogs.Friendly {
			listItems = append(listItems, fmt.Sprintf("\t----- %s -----", user.UserName))
			listItems = append(listItems, fmt.Sprintf("Healed: %d\n", sumArr(user.TotalHeals)))
		}

	}
	return listItems
}

func (_daocLogs *DaocLogs) calculateDensives() []string {
	listItems := []string{}
	if _daocLogs.User.BlockTotal+_daocLogs.User.ParryTotal+_daocLogs.User.EvadeTotal+_daocLogs.User.TotalStuns+_daocLogs.User.ResistsInTotal > 0 {
		listItems = append(listItems, "\t\t----- Defensives -----\n")
		if _daocLogs.User.BlockTotal > 0 {
			listItems = append(listItems, fmt.Sprintf("Block: %d\n", _daocLogs.User.BlockTotal))
		}
		if _daocLogs.User.ParryTotal > 0 {
			listItems = append(listItems, fmt.Sprintf("Parry: %d\n", _daocLogs.User.ParryTotal))
		}
		if _daocLogs.User.EvadeTotal > 0 {
			listItems = append(listItems, fmt.Sprintf("Evade: %d\n", _daocLogs.User.EvadeTotal))
		}
		if _daocLogs.User.TotalStuns > 0 {
			listItems = append(listItems, fmt.Sprintf("Stuns: %d\n", _daocLogs.User.TotalStuns))
		}
		if _daocLogs.User.TotalStuns > 0 {
			listItems = append(listItems, fmt.Sprintf("You Resisted: %d\n", _daocLogs.User.ResistsInTotal))
		}
	}
	return listItems
}

func (_daocLogs *DaocLogs) getCombativeUsers() []string {
	users := []string{}
	totalAllDamage := []int{}

	for _, user := range _daocLogs.Enemy {
		totalAllDamage = append(totalAllDamage, user.MovingDamageTotal...)
	}
	users = append(users, "\t\t----- Damage In -----")
	users = append(users, fmt.Sprintf("Total Damage: %d\n", sumArr(totalAllDamage)))
	users = append(users, fmt.Sprintf("Total Hits: %d", len(totalAllDamage)))
	users = append(users, "\t\t----- Combatives -----")
	for _, user := range _daocLogs.Enemy {
		users = append(users, "\t----- "+user.UserName+" -----")
		users = append(users, fmt.Sprintf("Damage In: %d", sumArr(user.MovingDamageReceived)))
		users = append(users, fmt.Sprintf("Damage Out: %d", sumArr(user.MovingDamageTotal)))
		users = append(users, fmt.Sprintf("Blocks: %d", user.BlockTotal))
		users = append(users, fmt.Sprintf("Evades: %d", user.EvadeTotal))
		users = append(users, fmt.Sprintf("Parry: %d", user.ParryTotal))
		users = append(users, fmt.Sprintf("Resisted You: %d", user.ResistsInTotal))
	}
	return users
}

func (_daocLogs *DaocLogs) calculateTime() []string {
	listItems := []string{}
	totalMinutes := int(_daocLogs.User.EndTime.Sub(_daocLogs.User.StartTime).Seconds()) / 60
	totalSeconds := int(_daocLogs.User.EndTime.Sub(_daocLogs.User.StartTime).Seconds()) - (60 * totalMinutes)
	listItems = append(listItems, "\t\t----- Total Time -----")
	listItems = append(listItems, fmt.Sprintf("%d minutes and %d seconds", totalMinutes, totalSeconds))
	return listItems
}

func (_daocLogs *DaocLogs) getKills() []string {
	listItems := []string{}
	for _, user := range _daocLogs.Enemy {
		if user.TotalDeaths > 0 {
			listItems = append(listItems, fmt.Sprintf("%s: %d", user.UserName, user.TotalDeaths))
		}
	}
	return listItems
}

func (_daocLogs *DaocLogs) calculateArmorhits() []string {
	listItems := []string{}
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
		listItems = append(listItems, "\t\t----- Armor Damaged -----")
		if len(head) > 0 {
			listItems = append(listItems, fmt.Sprintf("Head Hit: %d", len(head)))
			listItems = append(listItems, fmt.Sprintf("Head Damage: %d", sumArr(head)))
		}
		if len(torso) > 0 {
			listItems = append(listItems, fmt.Sprintf("Torso Hit: %d", len(torso)))
			listItems = append(listItems, fmt.Sprintf("Torso Damage: %d", sumArr(torso)))
		}
		if len(arm) > 0 {
			listItems = append(listItems, fmt.Sprintf("Arm Hit: %d", len(arm)))
			listItems = append(listItems, fmt.Sprintf("Arm Damage: %d", sumArr(arm)))
		}
		if len(leg) > 0 {
			listItems = append(listItems, fmt.Sprintf("Leg Hit: %d", len(leg)))
			listItems = append(listItems, fmt.Sprintf("Leg Damage: %d", sumArr(leg)))
		}
		if len(hand) > 0 {
			listItems = append(listItems, fmt.Sprintf("Hand Hit: %d", len(hand)))
			listItems = append(listItems, fmt.Sprintf("Hand Damage: %d", sumArr(hand)))
		}
		if len(foot) > 0 {
			listItems = append(listItems, fmt.Sprintf("Foot Hit: %d", len(foot)))
			listItems = append(listItems, fmt.Sprintf("Foot Damage: %d", sumArr(foot)))
		}
	}
	return listItems
}
