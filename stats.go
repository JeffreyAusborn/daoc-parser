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
	}
	return listItems
}

func (_daocLogs *DaocLogs) calculateDamageIn() []string {
	listItems := []string{}
	totalAllDamage := []int{}

	for _, user := range _daocLogs.Enemy {
		totalAllDamage = append(totalAllDamage, user.MovingDamageTotal...)
	}

	totalDamage := sumArr(totalAllDamage)

	if totalDamage > 0 {
		listItems = append(listItems, "\t\t----- Damage In -----\n")
		if totalDamage > 0 {
			listItems = append(listItems, fmt.Sprintf("Hits: %d", len(totalAllDamage)))
			listItems = append(listItems, fmt.Sprintf("Damage: %d\n", totalDamage))
		}
		if _daocLogs.getUser().TotalDeaths > 0 {
			listItems = append(listItems, fmt.Sprintf("Deaths: %d\n", daocLogs.getUser().TotalDeaths))
		}
	}
	return listItems
}

func (_daocLogs *DaocLogs) calculateHeal() []string {
	listItems := []string{}
	if len(_daocLogs.User.TotalSelfHeal)+len(_daocLogs.User.TotalAbsorbed)+len(_daocLogs.User.TotalHeals) > 0 {
		listItems = append(listItems, "\t\t----- Healing & Absorb -----\n")
		if len(_daocLogs.User.TotalSelfHeal) > 0 {
			listItems = append(listItems, fmt.Sprintf("Self Heals: %d\n", sumArr(_daocLogs.User.TotalSelfHeal)))
		}
		if len(_daocLogs.User.TotalHeals) > 0 {
			listItems = append(listItems, fmt.Sprintf("All Heals: %d\n", sumArr(_daocLogs.User.TotalHeals)))
		}
		if len(_daocLogs.User.TotalHealsCrits) > 0 {
			listItems = append(listItems, fmt.Sprintf("Heal Crits: %d\n", sumArr(_daocLogs.User.TotalHealsCrits)))
		}
		if len(_daocLogs.User.TotalAbsorbed) > 0 {
			listItems = append(listItems, fmt.Sprintf("Absorbed: %d\n", sumArr(_daocLogs.User.TotalAbsorbed)))
		}
		if _daocLogs.User.OverHeals > 0 {
			listItems = append(listItems, fmt.Sprintf("OverHeal Count: %d\n", _daocLogs.User.OverHeals))
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
	users = append(users, "\t\t----- Combatives -----")
	for _, user := range _daocLogs.Enemy {
		users = append(users, "\t----- "+user.UserName+" -----")
		users = append(users, fmt.Sprintf("Damage Taken: %d", sumArr(user.MovingDamageReceived)))
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
