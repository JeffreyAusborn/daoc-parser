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
		totalHeals := []int{}
		totalCrit := 0
		for _, ability := range _daocLogs.getUser().Heals {
			totalHeals = append(totalHeals, ability.Output...)
			totalCrit += sumArr(ability.Crit)
		}
		listItems = append(listItems, "\t----- Total Heals -----\n")
		min, max := 0, 0
		if len(totalHeals) > 0 {
			min, max = getMinAndMax(totalHeals)
		}
		listItems = append(listItems, fmt.Sprintf("Total Heals: %d\n", sumArr(totalHeals)))
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
	if _daocLogs.getUser().BlockTotal > 0 {
		listItems = append(listItems, fmt.Sprintf("Block: %d\n", _daocLogs.getUser().BlockTotal))
	}
	if _daocLogs.getUser().ParryTotal > 0 {
		listItems = append(listItems, fmt.Sprintf("Parry: %d\n", _daocLogs.getUser().ParryTotal))
	}
	if _daocLogs.getUser().EvadeTotal > 0 {
		listItems = append(listItems, fmt.Sprintf("Evade: %d\n", _daocLogs.getUser().EvadeTotal))
	}
	if _daocLogs.getUser().TotalStuns > 0 {
		listItems = append(listItems, fmt.Sprintf("Stuns: %d\n", _daocLogs.getUser().TotalStuns))
	}
	if _daocLogs.getUser().TotalStuns > 0 {
		listItems = append(listItems, fmt.Sprintf("You Resisted: %d\n", _daocLogs.getUser().ResistsInTotal))
	}
	if _daocLogs.getUser().BladeturnTotal > 0 {
		listItems = append(listItems, fmt.Sprintf("Bladeturn: %d\n", _daocLogs.getUser().BladeturnTotal))
	}
	return listItems
}

func (_daocLogs *DaocLogs) getCombativeUsers() []string {
	users := []string{}
	if len(_daocLogs.Enemy) > 0 {
		for _, user := range _daocLogs.Enemy {
			userDamaged := 0
			users = append(users, "\t----- "+user.UserName+" -----")
			for _, ability := range _daocLogs.getUser().Spells {
				for _, userHit := range ability.Users {
					if userHit.UserName == user.UserName {
						userDamaged += sumArr(userHit.MovingDamageReceived)
					}
				}
			}
			for _, ability := range _daocLogs.getUser().Styles {
				for _, userHit := range ability.Users {
					if userHit.UserName == user.UserName {
						userDamaged += sumArr(userHit.MovingDamageReceived)
					}
				}
			}
			if userDamaged > 0 {
				users = append(users, fmt.Sprintf("Damage Received: %d", userDamaged))
			}
			if sumArr(user.MovingDamageTotal) > 0 {
				users = append(users, fmt.Sprintf("Damaged You: %d", sumArr(user.MovingDamageTotal)))
			}
			if user.BlockTotal > 0 {
				users = append(users, fmt.Sprintf("Blocks You: %d", user.BlockTotal))
			}
			if user.EvadeTotal > 0 {
				users = append(users, fmt.Sprintf("Evades You: %d", user.EvadeTotal))
			}
			if user.ParryTotal > 0 {
				users = append(users, fmt.Sprintf("Parry You: %d", user.ParryTotal))
			}
			if user.ResistsInTotal > 0 {
				users = append(users, fmt.Sprintf("Resisted You: %d", user.ResistsInTotal))
			}
		}
	}
	return users
}

func (_daocLogs *DaocLogs) calculateTime() []string {
	listItems := []string{}
	totalMinutes := int(_daocLogs.getUser().EndTime.Sub(_daocLogs.getUser().StartTime).Seconds()) / 60
	totalSeconds := int(_daocLogs.getUser().EndTime.Sub(_daocLogs.getUser().StartTime).Seconds()) - (60 * totalMinutes)
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
		if len(head) > 0 {
			listItems = append(listItems, "\t\t----- Head Damaged -----")
			listItems = append(listItems, fmt.Sprintf("Hits: %d", len(head)))
			listItems = append(listItems, fmt.Sprintf("Damage: %d", sumArr(head)))
			min, max := getMinAndMax(head)
			listItems = append(listItems, fmt.Sprintf("Min Damage: %d", min))
			listItems = append(listItems, fmt.Sprintf("Max Damage: %d", max))
			listItems = append(listItems, fmt.Sprintf("Avg Damage: %d", sumArr(head)/len(head)))
		}
		if len(torso) > 0 {
			listItems = append(listItems, "\t\t----- Torso Damaged -----")
			listItems = append(listItems, fmt.Sprintf("Hits: %d", len(torso)))
			listItems = append(listItems, fmt.Sprintf("Damage: %d", sumArr(torso)))
			min, max := getMinAndMax(torso)
			listItems = append(listItems, fmt.Sprintf("Min Damage: %d", min))
			listItems = append(listItems, fmt.Sprintf("Max Damage: %d", max))
			listItems = append(listItems, fmt.Sprintf("Avg Damage: %d", sumArr(torso)/len(torso)))
		}
		if len(arm) > 0 {
			listItems = append(listItems, "\t\t----- Arm Damaged -----")
			listItems = append(listItems, fmt.Sprintf("Hits: %d", len(arm)))
			listItems = append(listItems, fmt.Sprintf("Damage: %d", sumArr(arm)))
			min, max := getMinAndMax(arm)
			listItems = append(listItems, fmt.Sprintf("Min Damage: %d", min))
			listItems = append(listItems, fmt.Sprintf("Max Damage: %d", max))
			listItems = append(listItems, fmt.Sprintf("Avg Damage: %d", sumArr(arm)/len(arm)))
		}
		if len(leg) > 0 {
			listItems = append(listItems, "\t\t----- Leg Damaged -----")
			listItems = append(listItems, fmt.Sprintf("Hits: %d", len(leg)))
			listItems = append(listItems, fmt.Sprintf("Damage: %d", sumArr(leg)))
			min, max := getMinAndMax(leg)
			listItems = append(listItems, fmt.Sprintf("Min Damage: %d", min))
			listItems = append(listItems, fmt.Sprintf("Max Damage: %d", max))
			listItems = append(listItems, fmt.Sprintf("Avg Damage: %d", sumArr(leg)/len(leg)))
		}
		if len(hand) > 0 {
			listItems = append(listItems, "\t\t----- Hand Damaged -----")
			listItems = append(listItems, fmt.Sprintf("Hits: %d", len(hand)))
			listItems = append(listItems, fmt.Sprintf("Damage: %d", sumArr(hand)))
			min, max := getMinAndMax(hand)
			listItems = append(listItems, fmt.Sprintf("Min Damage: %d", min))
			listItems = append(listItems, fmt.Sprintf("Max Damage: %d", max))
			listItems = append(listItems, fmt.Sprintf("Avg Damage: %d", sumArr(hand)/len(hand)))
		}
		if len(foot) > 0 {
			listItems = append(listItems, "\t\t----- Foot Damaged -----")
			listItems = append(listItems, fmt.Sprintf("Hits: %d", len(foot)))
			listItems = append(listItems, fmt.Sprintf("Damage: %d", sumArr(foot)))
			min, max := getMinAndMax(foot)
			listItems = append(listItems, fmt.Sprintf("Min Damage: %d", min))
			listItems = append(listItems, fmt.Sprintf("Max Damage: %d", max))
			listItems = append(listItems, fmt.Sprintf("Avg Damage: %d", sumArr(foot)/len(foot)))
		}
	}
	return listItems
}

// WORKING ON REWORK

func (_daocLogs *DaocLogs) calculateSpells() []string {
	listItems := []string{}
	if len(_daocLogs.getUser().Spells) > 0 {
		totalHits := 0
		totalDamage := []int{}
		totalCrit := []int{}
		totalInterupts := 0
		totalSiphons := 0
		totalResists := 0
		for _, ability := range _daocLogs.getUser().Spells {
			totalHits += len(ability.Output)
			totalDamage = append(totalDamage, ability.Output...)
			totalCrit = append(totalCrit, ability.Crit...)
			totalSiphons += ability.Siphon
			for _, user := range ability.Users {
				totalInterupts += user.Interrupted
			}
			for _, user := range ability.Users {
				totalResists += user.ResistsInTotal
			}
		}

		listItems = append(listItems, fmt.Sprintf("Total Hits: %d", totalHits))
		listItems = append(listItems, fmt.Sprintf("Total Damage: %d\n", sumArr(totalDamage)))
		if len(totalDamage) > 0 {
			minC, maxC := getMinAndMax(totalDamage)
			listItems = append(listItems, fmt.Sprintf("Min Damage: %d\n", minC))
			listItems = append(listItems, fmt.Sprintf("Max Damage: %d\n", maxC))
			listItems = append(listItems, fmt.Sprintf("Average Damage: %d\n", sumArr(totalDamage)/len(totalDamage)))
		}
		listItems = append(listItems, fmt.Sprintf("Total Crit: %d\n", sumArr(totalCrit)))
		if len(totalCrit) > 0 {
			minC, maxC := getMinAndMax(totalCrit)
			listItems = append(listItems, fmt.Sprintf("Min Crit: %d\n", minC))
			listItems = append(listItems, fmt.Sprintf("Max Crit: %d\n", maxC))
			listItems = append(listItems, fmt.Sprintf("Average Crit: %d\n", sumArr(totalCrit)/len(totalCrit)))
		}
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
					minC, maxC := getMinAndMax(ability.Output)
					listItems = append(listItems, fmt.Sprintf("Min Damage: %d\n", minC))
					listItems = append(listItems, fmt.Sprintf("Max Damage: %d\n", maxC))
					listItems = append(listItems, fmt.Sprintf("Average Damage: %d\n", sumArr(ability.Output)/len(ability.Output)))
				}
				if sumArr(ability.Crit) > 0 {
					listItems = append(listItems, fmt.Sprintf("Crit: %d\n", sumArr(ability.Crit)))
					minC, maxC := getMinAndMax(ability.Crit)
					listItems = append(listItems, fmt.Sprintf("Min Crit: %d\n", minC))
					listItems = append(listItems, fmt.Sprintf("Max Crit: %d\n", maxC))
					listItems = append(listItems, fmt.Sprintf("Average Crit: %d\n", sumArr(ability.Crit)/len(ability.Crit)))
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
				}
				for _, user := range ability.Users {
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

		totalHits = 0
		totalDamage = 0
		totalCrit = 0
		totalExtraDamage = 0
		for _, ability := range _daocLogs.getUser().Styles {
			for _, weapon := range ability.Weapons {
				totalHits += len(weapon.Output)
				totalDamage += sumArr(weapon.Output)
				totalCrit += sumArr(weapon.Crit)
				totalExtraDamage += sumArr(weapon.ExtraDamage)
			}
		}

		for _, ability := range _daocLogs.getUser().Styles {
			totalHits = 0
			totalDamage = 0
			totalCrit = 0
			totalExtraDamage = 0
			for _, weapon := range ability.Weapons {
				totalHits += len(weapon.Output)
				totalDamage += sumArr(weapon.Output)
				totalCrit += sumArr(weapon.Crit)
				totalExtraDamage += sumArr(weapon.ExtraDamage)
			}
			listItems = append(listItems, fmt.Sprintf("\t\t----- %s -----", ability.Name))
			listItems = append(listItems, fmt.Sprintf("Total Hits: %d", totalHits))
			listItems = append(listItems, fmt.Sprintf("Total Damage: %d\n", totalDamage))
			if totalCrit > 0 {
				listItems = append(listItems, fmt.Sprintf("Total Crit: %d\n", totalCrit))
			}
			if totalExtraDamage > 0 {
				listItems = append(listItems, fmt.Sprintf("Total Extra Damage: %d\n", totalExtraDamage))
			}
			if ability.Stunned > 0 {
				listItems = append(listItems, fmt.Sprintf("Total Stunned: %d\n", ability.Stunned))
			}
			if ability.Blocked > 0 {
				listItems = append(listItems, fmt.Sprintf("Total Blocked: %d", ability.Blocked))
			}
			if ability.Evaded > 0 {
				listItems = append(listItems, fmt.Sprintf("Total Evaded: %d\n", ability.Evaded))
			}
			if ability.Parried > 0 {
				listItems = append(listItems, fmt.Sprintf("Total Parried: %d\n", ability.Parried))
			}

			if len(ability.GrowthRate) > 0 {
				minG, maxG := getMinAndMax(ability.GrowthRate)
				listItems = append(listItems, fmt.Sprintf("Growth Rate Min: %d\n", minG))
				listItems = append(listItems, fmt.Sprintf("Growth Rate Max: %d\n", maxG))
				listItems = append(listItems, fmt.Sprintf("Growth Rate Average: %d\n", sumArr(ability.GrowthRate)/len(ability.GrowthRate)))
			}

			for _, weapon := range ability.Weapons {
				listItems = append(listItems, fmt.Sprintf("\t----- %s -----", weapon.Name))
				listItems = append(listItems, fmt.Sprintf("Hits: %d", len(weapon.Output)))
				listItems = append(listItems, fmt.Sprintf("Damage: %d\n", sumArr(weapon.Output)))
				if sumArr(weapon.Crit) > 0 {
					listItems = append(listItems, fmt.Sprintf("Crit: %d\n", sumArr(weapon.Crit)))
				}
				if sumArr(weapon.ExtraDamage) > 0 {
					listItems = append(listItems, fmt.Sprintf("Extra Damage: %d\n", sumArr(weapon.ExtraDamage)))
				}
				if len(weapon.GrowthRate) > 0 {
					minG, maxG := getMinAndMax(weapon.GrowthRate)
					listItems = append(listItems, fmt.Sprintf("Growth Rate Min: %d\n", minG))
					listItems = append(listItems, fmt.Sprintf("Growth Rate Max: %d\n", maxG))
					listItems = append(listItems, fmt.Sprintf("Growth Rate Average: %d\n", sumArr(weapon.GrowthRate)/len(weapon.GrowthRate)))
				}
			}
		}
	}
	return listItems
}
