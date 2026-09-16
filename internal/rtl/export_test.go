package rtl

// ClassListForTest names the walker's result element so the registry sweep in
// package rtl_test can range over it.
type ClassListForTest = classList

func ClassListsForTest(filename string, src []byte) []ClassListForTest {
	lists, err := classLists(filename, src)
	if err != nil {
		panic(err)
	}
	return lists
}

func SplitForTest(class string) (variant, value string) {
	variant, value, _ = split(class)
	return variant, value
}

func IsSideKeyedForTest(variant string) bool { return isSideKeyed(variant) }
