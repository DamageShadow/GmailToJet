package main

func difference(slice1 []string, slice2 []string) ([]string, []string){
	diffStr1 := []string{}
	diffStr2 := []string{}
	m :=map [string]int{}

	for _, s1Val := range slice1 {
		m[s1Val] = 1
	}
	for _, s2Val := range slice2 {
		m[s2Val] = m[s2Val] + 2
	}

	for mKey, mVal := range m {
		if mVal==1 {
			diffStr1 = append(diffStr1, mKey)
		}
		if mVal==2 {
			diffStr2 = append(diffStr2, mKey)
		}
	}

	return diffStr1, diffStr2
}


