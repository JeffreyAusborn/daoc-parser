package main

import (
	"fmt"
)

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
	if len(_daocLogs.getUser().Heals) > 0 {
		totalHeals := 0
		totalCrit := 0
		for _, ability := range _daocLogs.getUser().Heals {
			totalHeals += sumArr(ability.Output)
			totalCrit += sumArr(ability.Crit)
		}
		listItems = append(listItems, "\t----- Total Heals -----\n")
		min, max := 0, 0
		if totalHeals > 0 {
			min, max = getMinAndMax(_daocLogs.getUser().TotalHeals)
		}
		listItems = append(listItems, fmt.Sprintf("Total Heals: %d\n", totalHeals))
		listItems = append(listItems, fmt.Sprintf("Total Crit Heals: %d\n", totalCrit))
		listItems = append(listItems, fmt.Sprintf("Minimum Heal: %d\n", min))
		listItems = append(listItems, fmt.Sprintf("Maximum Heal: %d\n", max))
	}
	if len(_daocLogs.getUser().TotalAbsorbed) > 0 {
		listItems = append(listItems, "\t----- Absorbs -----\n")
		min, max := getMinAndMax(_daocLogs.getUser().TotalAbsorbed)
		listItems = append(listItems, fmt.Sprintf("Absorbed: %d\n", sumArr(_daocLogs.getUser().TotalAbsorbed)))
		listItems = append(listItems, fmt.Sprintf("Minimum Absorbed: %d\n", min))
		listItems = append(listItems, fmt.Sprintf("Maximum Absorbed: %d\n", max))
		listItems = append(listItems, fmt.Sprintf("Average Absorbed: %d\n", sumArr(_daocLogs.getUser().TotalAbsorbed)/len(_daocLogs.getUser().TotalAbsorbed)))
	}
	if len(_daocLogs.getUser().Heals) > 0 {
		for _, ability := range _daocLogs.getUser().Heals {
			listItems = append(listItems, fmt.Sprintf("\t----- %s -----", ability.Name))
			listItems = append(listItems, fmt.Sprintf("Heal Count: %d", len(ability.Output)))
			listItems = append(listItems, fmt.Sprintf("Heals: %d\n", sumArr(ability.Output)))
			listItems = append(listItems, fmt.Sprintf("Crit: %d\n", sumArr(ability.Crit)))

			for _, user := range ability.Users {
				if user.Interrupted > 0 {
					listItems = append(listItems, fmt.Sprintf("%s interrupted you %d times", user.UserName, user.Interrupted))
				}
			}

			for _, user := range ability.Users {
				if sumArr(user.TotalHeals) > 0 {
					listItems = append(listItems, fmt.Sprintf("Healed %s: %d\n", user.UserName, sumArr(user.TotalHeals)))
				}
			}
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

// WORKING ON REWORK

func (_daocLogs *DaocLogs) calculateSpells() []string {
	listItems := []string{}
	if len(_daocLogs.getUser().Spells) > 0 {
		totalHits := 0
		totalDamage := 0
		totalCrit := 0
		totalInterupts := 0
		totalSiphons := 0
		totalResists := 0
		for _, ability := range _daocLogs.getUser().Spells {
			totalHits += len(ability.Output)
			totalDamage += sumArr(ability.Output)
			totalCrit += sumArr(ability.Crit)
			totalSiphons += ability.Siphon
			for _, user := range ability.Users {
				totalInterupts += user.Interrupted
			}
			for _, user := range ability.Users {
				totalResists += user.ResistsInTotal
			}
		}

		listItems = append(listItems, fmt.Sprintf("Total Hits: %d", totalHits))
		listItems = append(listItems, fmt.Sprintf("Total Damage: %d\n", totalDamage))
		listItems = append(listItems, fmt.Sprintf("Total Crit: %d\n", totalCrit))
		listItems = append(listItems, fmt.Sprintf("Total Siphons: %d\n", totalSiphons))
		listItems = append(listItems, fmt.Sprintf("Total Interupts Received: %d\n", totalInterupts))
		listItems = append(listItems, fmt.Sprintf("Total Resists: %d\n", totalResists))

		for _, ability := range _daocLogs.getUser().Spells {
			totalInterupts = 0
			totalResists = 0
			for _, user := range ability.Users {
				totalInterupts += user.Interrupted
				totalResists += user.ResistsInTotal
			}
			if len(ability.Output) > 0 || (len(ability.Output) == 0 && (totalInterupts > 0 || totalResists > 0 || ability.Siphon > 0)) {
				listItems = append(listItems, fmt.Sprintf("\t----- %s -----", ability.Name))
				if len(ability.Output) > 0 {
					listItems = append(listItems, fmt.Sprintf("Hits: %d", len(ability.Output)))
				}
				if sumArr(ability.Output) > 0 {
					listItems = append(listItems, fmt.Sprintf("Damage: %d\n", sumArr(ability.Output)))
				}
				if sumArr(ability.Crit) > 0 {
					listItems = append(listItems, fmt.Sprintf("Crit: %d\n", sumArr(ability.Crit)))
				}
				if totalResists > 0 {
					listItems = append(listItems, fmt.Sprintf("Resists: %d\n", totalResists))
				}
				if ability.Siphon > 0 {
					listItems = append(listItems, fmt.Sprintf("Siphon: %d\n", ability.Siphon))
				}
				if ability.Stunned > 0 {
					listItems = append(listItems, fmt.Sprintf("Stunned: %d\n", ability.Stunned))
				}
				for _, user := range ability.Users {
					if user.Interrupted > 0 {
						listItems = append(listItems, fmt.Sprintf("%s interrupted you %d times", user.UserName, user.Interrupted))
					}
					if sumArr(user.MovingDamageReceived) > 0 {
						listItems = append(listItems, fmt.Sprintf("You hit %s for: %d", user.UserName, sumArr(user.MovingDamageReceived)))
					}
				}
			}
		}
	}
	return listItems
}

func (_daocLogs *DaocLogs) calculateStyles() []string {
	listItems := []string{}
	if len(_daocLogs.getUser().Styles) > 0 {
		totalHits := 0
		totalDamage := 0
		totalCrit := 0
		totalExtraDamage := 0
		for _, ability := range _daocLogs.getUser().Styles {
			totalHits += len(ability.Output)
			totalDamage += sumArr(ability.Output)
			totalCrit += sumArr(ability.Crit)
			totalExtraDamage += sumArr(ability.ExtraDamage)
		}

		listItems = append(listItems, fmt.Sprintf("Total Hits: %d", totalHits))
		listItems = append(listItems, fmt.Sprintf("Total Damage: %d\n", totalDamage))
		listItems = append(listItems, fmt.Sprintf("Total Crit: %d\n", totalCrit))
		listItems = append(listItems, fmt.Sprintf("Total Extra Damage: %d\n", totalExtraDamage))

		for _, ability := range _daocLogs.getUser().Styles {
			listItems = append(listItems, fmt.Sprintf("\t----- %s -----", ability.Name))
			listItems = append(listItems, fmt.Sprintf("Hits: %d", len(ability.Output)))
			listItems = append(listItems, fmt.Sprintf("Damage: %d\n", sumArr(ability.Output)))
			if sumArr(ability.ExtraDamage) > 0 {
				listItems = append(listItems, fmt.Sprintf("Extra Damage: %d\n", sumArr(ability.ExtraDamage)))
			}
			if sumArr(ability.Crit) > 0 {
				listItems = append(listItems, fmt.Sprintf("Crit: %d\n", sumArr(ability.Crit)))
			}
			if ability.Stunned > 0 {
				listItems = append(listItems, fmt.Sprintf("Stunned: %d\n", ability.Stunned))
			}
			if ability.Blocked > 0 {
				listItems = append(listItems, fmt.Sprintf("Blocked: %d", ability.Blocked))
			}
			if ability.Evaded > 0 {
				listItems = append(listItems, fmt.Sprintf("Evaded: %d\n", ability.Evaded))
			}
			if ability.Parried > 0 {
				listItems = append(listItems, fmt.Sprintf("Parried: %d\n", ability.Parried))
			}
			if len(ability.GrowthRate) > 0 {
				minG, maxG := getMinAndMax(ability.GrowthRate)
				listItems = append(listItems, fmt.Sprintf("Growth Rate Min: %d\n", minG))
				listItems = append(listItems, fmt.Sprintf("Growth Rate Max: %d\n", maxG))
				listItems = append(listItems, fmt.Sprintf("Growth Rate Average: %d\n", sumArr(ability.GrowthRate)/len(ability.GrowthRate)))
			}
		}
	}
	return listItems
}
