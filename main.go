package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gosuri/uilive"
)

type DaocLogs struct {
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
	stylesPerformed       int
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

	enemyEvade       int
	enemyParry       int
	enemyBlock       int
	enemyDamageTotal []int
}

func main() {
	writer := uilive.New()
	writer.Start()
	defer writer.Stop()
	logPath := flag.String("file", "", "Path to chat.log")
	streamLogs := flag.Bool("stream", false, "")
	flag.Parse()
	if *logPath != "" {
		openLogFile(*logPath, writer)
		if *streamLogs {
			for range time.Tick(time.Second * 3) {
				openLogFile(*logPath, writer)
			}
		}
	} else {
		fmt.Println("Requied argument missing: --file path/to/chat.log")
	}
}

func openLogFile(logPath string, writer *uilive.Writer) {
	f, err := os.OpenFile(logPath, os.O_RDONLY|os.O_EXCL, 0666)
	defer f.Close()
	if err == nil {
		iterateLogFile(f, writer)
	}
}

func iterateLogFile(f *os.File, writer *uilive.Writer) {
	var daocLogs DaocLogs
	reader := bufio.NewReader(f)
	style := false
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break

		}

		/*
			Player offensive
		*/
		match, _ := regexp.MatchString("You attack.*with your.*and hit for.*damage", line)
		if match {
			damage := strings.Split(line, "and hit for ")[1]
			damage = strings.Split(damage, " damage")[0]
			damage = strings.Split(damage, " ")[0]
			damageInt, _ := strconv.Atoi(damage)
			daocLogs.movingDamageTotal = append(daocLogs.movingDamageTotal, damageInt)
			if style {
				style = false
				daocLogs.movingDamageStyles = append(daocLogs.movingDamageStyles, damageInt)
			} else {
				daocLogs.movingDamageBaseMelee = append(daocLogs.movingDamageBaseMelee, damageInt)
			}

			user := strings.Split(line, "You attack ")[1]
			user = strings.Split(user, " with your")[0]
			daocLogs.usersHit = append(daocLogs.usersHit, user)
		}
		match, _ = regexp.MatchString("You prepare to perform", line)
		if match {
			daocLogs.stylesPerformed += 1
			style = true
		}
		match, _ = regexp.MatchString("You cast a ", line)
		if match {
			daocLogs.spellsPerformed += 1
		}
		match, _ = regexp.MatchString("You begin casting a ", line)
		if match {
			daocLogs.castedSpellsPerformed += 1
		}
		match, _ = regexp.MatchString("You hit.*for.*damage!", line)
		if match {
			damage := strings.Split(line, " for ")[1]
			damage = strings.Split(damage, " damage")[0]
			damage = strings.Split(damage, " ")[0]
			damageInt, _ := strconv.Atoi(damage)
			match, _ = regexp.MatchString("extra damage", line)
			if match {
				daocLogs.movingExtraDamage = append(daocLogs.movingExtraDamage, damageInt)
			} else {
				daocLogs.movingDamageSpells = append(daocLogs.movingDamageSpells, damageInt)
			}
			daocLogs.movingDamageTotal = append(daocLogs.movingDamageTotal, damageInt)

			user := strings.Split(line, "You hit ")[1]
			user = strings.Split(user, " for")[0]
			daocLogs.usersHit = append(daocLogs.usersHit, user)
		}
		match, _ = regexp.MatchString("You miss!.*", line)
		if match {
			daocLogs.missesTotal += 1
		}
		match, _ = regexp.MatchString("resists the effect.*", line)
		if match {
			daocLogs.resistsTotal += 1
		}
		match, _ = regexp.MatchString("You gather energy from your surroundings.*", line)
		if match {
			daocLogs.siphonTotal += 1
		}
		match, _ = regexp.MatchString("fully healed", line)
		if match {
			daocLogs.overHeals += 1
		}
		match, _ = regexp.MatchString("You critically hit", line)
		if match {
			damage := strings.Split(line, "for an additional ")[1]
			damage = strings.Split(damage, " damage")[0]
			damageInt, _ := strconv.Atoi(damage)
			match, _ = regexp.MatchString("extra damage", line)
			daocLogs.movingDamageTotal = append(daocLogs.movingDamageTotal, damageInt)
			daocLogs.movingCritDamage = append(daocLogs.movingCritDamage, damageInt)
		}
		/*
			Player Defensive
		*/
		match, _ = regexp.MatchString("you block", line)
		if match {
			daocLogs.blockTotal += 1
		}
		match, _ = regexp.MatchString("you evade", line)
		if match {
			daocLogs.evadeTotal += 1
		}
		match, _ = regexp.MatchString("you parry", line)
		if match {
			daocLogs.parryTotal += 1
		}
		match, _ = regexp.MatchString("Your ablative absorbs", line)
		if match {
			ablative := strings.Split(line, "ablative absorbs ")[1]
			ablative = strings.Split(ablative, " damage")[0]
			ablativeInt, _ := strconv.Atoi(ablative)
			daocLogs.totalAblative += ablativeInt
		}
		/*
			Support
		*/
		match, _ = regexp.MatchString("You heal yourself for", line)
		if match {
			healing := strings.Split(line, "You heal yourself for ")[1]
			healing = strings.Split(healing, " hit points")[0]
			healingInt, _ := strconv.Atoi(healing)
			daocLogs.totalSelfHeal += healingInt
		}
		match, _ = regexp.MatchString("is stunned and cannot move", line)
		if match {
			daocLogs.totalStuns += 1
		}
		/*
			Pets
		*/
		match, _ = regexp.MatchString("The.*casts a spell!", line)
		if match {
		}
		match, _ = regexp.MatchString("Your.*attacks.*and hits for.*damage!", line)
		if match {
		}
		match, _ = regexp.MatchString("The.*attacks.*and misses!", line)
		if match {
		}
		/*
			Misc
		*/
		match, _ = regexp.MatchString("You gain a total of.*experience points", line)
		if match {
			exp := strings.Split(line, "You gain a total of ")[1]
			exp = strings.Split(exp, " experience")[0]
			exp = strings.ReplaceAll(exp, ",", "")
			expInt, _ := strconv.Atoi(exp)
			daocLogs.experienceGained += expInt
		}
		match, _ = regexp.MatchString("dies", line)
		if match {
			daocLogs.totalKills += 1
		}
		/*
			Enemy
		*/
		match, _ = regexp.MatchString("parries your attack", line)
		if match {
			daocLogs.enemyParry += 1
		}
		match, _ = regexp.MatchString("evades your attack", line)
		if match {
			daocLogs.enemyEvade += 1
		}
		match, _ = regexp.MatchString("blocks your attack", line)
		if match {
			daocLogs.enemyBlock += 1
		}
		match, _ = regexp.MatchString("hits your.*for ", line)
		if match {
			damage := strings.Split(line, "for ")[1]
			damage = strings.Split(damage, " damage")[0]
			damage = strings.Split(damage, " ")[0]
			damageInt, _ := strconv.Atoi(damage)
			daocLogs.enemyDamageTotal = append(daocLogs.enemyDamageTotal, damageInt)
		}

		/*
			Start and end
		*/
		if daocLogs.startTime.IsZero() {
			match, _ = regexp.MatchString("Chat Log Opened", line)
			if match {
				timeObj, err := time.Parse("Mon Jan 02 15:04:05 2006", strings.TrimSuffix(strings.Split(line, ": ")[1], "\r\n"))
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
				daocLogs.startTime = timeObj
			}
		}
		match, _ = regexp.MatchString("Chat Log Closed", line)
		if match {
			timeObj, err := time.Parse("Mon Jan 02 15:04:05 2006", strings.TrimSuffix(strings.Split(line, ": ")[1], "\r\n"))
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			daocLogs.endTime = timeObj
		}

	}
	fmt.Fprintf(writer, "Dark Age of Camelot - Chat Parser\nWritten by: Theorist\nIf you have any feedback, feel free to DM in Discord.\n\n")
	fmt.Fprintln(writer, "********** Melee **********")
	fmt.Fprintf(writer, "Styles Performed: %d\n", len(daocLogs.movingDamageStyles))
	fmt.Fprintf(writer, "Base Hits Performed: %d\n", len(daocLogs.movingDamageBaseMelee))
	fmt.Fprintf(writer, "Misses: %d\n", daocLogs.missesTotal)
	fmt.Fprintf(writer, "Total Style Damage: %d\n", sumArr(daocLogs.movingDamageStyles))
	fmt.Fprintf(writer, "Total Base Hit Damage: %d\n", sumArr(daocLogs.movingDamageBaseMelee))
	fmt.Fprintf(writer, "Total Melee Damage: %d\n", sumArr(daocLogs.movingDamageStyles)+sumArr(daocLogs.movingDamageBaseMelee))
	fmt.Fprintln(writer, "********** Spells **********")
	fmt.Fprintf(writer, "Total Spells Performed: %d\n", daocLogs.spellsPerformed)
	fmt.Fprintf(writer, "Casted Spells Performed: %d\n", daocLogs.castedSpellsPerformed)
	fmt.Fprintf(writer, "Insta Spells Performed: %d\n", daocLogs.spellsPerformed-daocLogs.castedSpellsPerformed)
	fmt.Fprintf(writer, "Spells with Damage: %d\n", len(daocLogs.movingDamageSpells))
	fmt.Fprintf(writer, "Total Resists: %d\n", daocLogs.resistsTotal)
	fmt.Fprintf(writer, "Total Siphons: %d\n", daocLogs.siphonTotal)
	fmt.Fprintf(writer, "Spell Damage: %d\n", sumArr(daocLogs.movingDamageSpells))
	fmt.Fprintf(writer, "Spell Extra Damage: %d\n", sumArr(daocLogs.movingExtraDamage))
	fmt.Fprintln(writer, "********** Criticals **********")
	fmt.Fprintf(writer, "Total Crits: %d\n", len(daocLogs.movingCritDamage))
	fmt.Fprintf(writer, "Total Crit Damage: %d\n", sumArr(daocLogs.movingCritDamage))
	fmt.Fprintln(writer, "********** Defensives **********")
	fmt.Fprintf(writer, "Total Blocks: %d\n", daocLogs.blockTotal)
	fmt.Fprintf(writer, "Total Parrys: %d\n", daocLogs.parryTotal)
	fmt.Fprintf(writer, "Total Evades: %d\n", daocLogs.evadeTotal)
	fmt.Fprintf(writer, "Total Stuns: %d\n", daocLogs.totalStuns)
	fmt.Fprintf(writer, "Total Self Heals: %d\n", daocLogs.totalSelfHeal)
	fmt.Fprintf(writer, "Total Ablative Absorbs: %d\n", daocLogs.totalAblative)
	fmt.Fprintln(writer, "********** Enemy **********")
	fmt.Fprintf(writer, "Enemy Total Hits: %d\n", len(daocLogs.enemyDamageTotal))
	fmt.Fprintf(writer, "Enemy Total Parrys: %d\n", daocLogs.enemyParry)
	fmt.Fprintf(writer, "Enemy Total Evades: %d\n", daocLogs.enemyEvade)
	fmt.Fprintf(writer, "Enemy Total Blocks: %d\n", daocLogs.enemyBlock)
	fmt.Fprintf(writer, "Total Damage Taken From Enemy: %d\n", sumArr(daocLogs.enemyDamageTotal))
	fmt.Fprintln(writer, "********** Total **********")
	fmt.Fprintf(writer, "Total Damage: %d\n", sumArr(daocLogs.movingDamageTotal))
	fmt.Fprintf(writer, "Total Experience Gained: %d\n", daocLogs.experienceGained)
	fmt.Fprintf(writer, "Total Killed: %d\n", daocLogs.totalKills)
	fmt.Fprintf(writer, "Users Hit: %s\n", strings.Join(dedupe(daocLogs.usersHit), ","))
	totalMinutes := int(daocLogs.endTime.Sub(daocLogs.startTime).Seconds()) / 60
	totalSeconds := int(daocLogs.endTime.Sub(daocLogs.startTime).Seconds()) - (60 * totalMinutes)
	fmt.Fprintf(writer, "Total Time: %d minutes and %d seconds\n", totalMinutes, totalSeconds)
}

func sumArr(arr []int) int {
	total := 0
	for _, item := range arr {
		total += item
	}
	return total
}

func dedupe(usersHit []string) []string {
	temp := make(map[string]int)
	// Add global users to the temp map
	for _, user := range usersHit {
		temp[user] = 1
	}

	tempKeys := []string{}
	// Create a slice of the temp keys to use with passpass
	for key, _ := range temp {
		if key != "" {
			tempKeys = append(tempKeys, key)
		}
	}
	return tempKeys
}
