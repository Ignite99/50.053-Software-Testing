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


	// iterate through all the respective paths
	for path := range x.Paths {
		goals[path] = map[string]Criteria{}

	// 	// for each path, get respective method and operation
	// 	for method, operation := range x.openAPI.Paths[path].Operations() {

	// 		// Request
	// 		ic := InputCriteria{
	// 			Types: map[string]int{}, 
	// 			Parameters: map[string]int{},
	// 		}

	// 		if operation.RequestBody == nil {

	// 			// Request Parameters
	// 			for _, parameter := range append(x.openAPI.Paths[path].Parameters, operation.Parameters...) {

	// 				ic.Parameters[parameter.Value.Name]++

	// 			}

	// 		} else {
	// 			for mediaType, content := range operation.RequestBody.Value.Content {

	// 				// Request Types
	// 				// JSON only
	// 				if strings.Contains(strings.ToLower(mediaType), "json") {
	// 					ic.Types["json"]++
	// 				}

	// 				// Request Parameters
	// 				ex, err := example.GetBodyExample(example.ModeRequest, content)
	// 				if err != nil {
	// 					if x.strictMode {
	// 						panic(err)
	// 					}
	// 				}

	// 				for _, key := range getJSONKeys(ex) {

	// 					ic.Parameters[key]++

	// 				}

	// 			}

	// 		}

	// }




	return cov
}

func isIndividualIncrease(a []int, b []int) (bool, Coverage) {

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
