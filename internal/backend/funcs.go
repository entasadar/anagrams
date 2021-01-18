package backend

import "sort"

type sortRunesInterface []rune

func (s sortRunesInterface) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s sortRunesInterface) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortRunesInterface) Len() int {
	return len(s)
}

func SortString(inputStr string) string {
	runes := []rune(inputStr)
	sort.Sort(sortRunesInterface(runes))
	return string(runes)
}
