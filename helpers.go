package main

func sumArr(arr []int) int {
	total := 0
	for _, item := range arr {
		total += item
	}
	return total
}

func dedupe(usersHit []string) []string {
	temp := make(map[string]int)
	for _, user := range usersHit {
		temp[user] = 1
	}

	tempKeys := []string{}
	for key, _ := range temp {
		if key != "" {
			tempKeys = append(tempKeys, key)
		}
	}
	return tempKeys
}

func getMinAndMax(array []int) (int, int) {
	var max int = array[0]
	var min int = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}
