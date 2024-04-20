package Django

type Criteria struct {
	Input  InputCriteria
	Output OutputCriteria
}

// InputCriteria is used to verify the input of the test coverage model.
type InputCriteria struct {
	Types      map[string]int
	Parameters map[string]int
}

// OutputCriteria is used to verify the output of the test coverage model.
type OutputCriteria struct {
	Types       map[string]int
	CodeClasses map[int]int
	Codes       map[int]int
	Properties  map[string]int
}

// func (x *HsuanFuzz) getCoverageLevels(mapInfos map[uint32][]*ResponseInfo) Coverage {
func getCoverageLevels(mapInfos map[int][]*ResponseInfo) Coverage {
	cov := Coverage{}
	// cov.Levels = append(cov.Levels, tmpLevels[path][method])
	cov.Levels = append(cov.Levels, 1)



	goals := map[string]map[string]Criteria{}
	seeds := map[string]map[string]Criteria{}
	tmpLevels := map[string]map[string]int{}


	for path := range x.openAPI.Paths {
		goals[path] = map[string]Criteria{}

	}




	return cov
}

func isIndividualIncrease(a []int, b []int, print bool) (bool, Coverage) {

	flag := false
	for i := 0; i < len(a); i++ {
		if a[i] > b[i] {
			flag = true
		} else {
			a[i] = b[i]
		}
	}
	return flag, Coverage{Levels: a}
}
