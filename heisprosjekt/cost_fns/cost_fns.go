package cost_fns

import (
	"Heis/driver-go/elevio"
	"Heis/elevator"
	"encoding/json"
	"fmt"
	"os/exec"
)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

type HRAElevState struct {
	Behavior    string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests [][2]bool               `json:"hallRequests"`
	States       map[string]HRAElevState `json:"states"`
}

//	func InputToCost() HRAInput {
//		input := HRAInput{
//			HallRequests: [][2]bool{{false, false}, {true, false}, {false, false}, {false, true}},
//			States: map[string]HRAElevState{
//				"one": HRAElevState{
//					Behavior:    "moving",
//					Floor:       2,
//					Direction:   "up",
//					CabRequests: []bool{false, false, false, true},
//				},
//				"two": HRAElevState{
//					Behavior:    "idle",
//					Floor:       0,
//					Direction:   "stop",
//					CabRequests: []bool{false, false, false, false},
//				},
//			},
//		}
//		return input
//	}
func RunCostFunc(elevMap map[string]elevator.Elevator) map[string]elevator.Elevator {
	commonHallCalls := elevator.OrHallCalls(elevMap)
	tempElevMap := elevMap
	// fmt.Println("[RunCostFunc] kommet inn")
	for k, v := range tempElevMap {
		if v.Failure {
			fmt.Println("CostFunc: Elevator ", k, " has failure")
			delete(tempElevMap, k)
		} else if v.Floor == -1 {
			fmt.Println("CostFunc: Elevator ", k, " has floor -1")
			delete(tempElevMap, k)
		}
	}	
	// fmt.Println("[RunCostFunc] ferdig iterert over tempElevMap")
	input := inputToCost(commonHallCalls, tempElevMap)
	newHRAs := getCostOutput(input)
	// fmt.Println("[RunCostFunc] skal til å itere gjennom newHRAs")
	for k := range newHRAs {
		elevMap[k] = mergeHallAndRequests(elevMap[k], newHRAs[k])
	}
	// fmt.Println("[RunCostFunc] går ut")
	return elevMap
}

func elevToHRAElevState(elev elevator.Elevator) HRAElevState {
	return HRAElevState{
		Behavior:    elevator.EB_ToString(elev.Behaviour),
		Floor:       elev.Floor,
		Direction:   elevator.MD_ToString(elev.Dirn),
		CabRequests: elevator.GetCabCalls(elev),
	}
}

func inputToCost(commonHallCalls [][2]bool, elevMap map[string](elevator.Elevator)) HRAInput {
	// commonHallCalls := orHallCalls(elevMap)

	stateMap := map[string]HRAElevState{}
	for k, v := range elevMap {
		stateMap[k] = elevToHRAElevState(v)
	}

	input := HRAInput{
		HallRequests: commonHallCalls,
		States:       stateMap,
	}
	return input
}

// func InputToCost(elev elevator.Elevator, elevatorId string) HRAInput {
// 	input := HRAInput{
// 		HallRequests: elevator.GetHallCalls(elev),
// 		States: map[string]HRAElevState{
// 			elevatorId: elevToHRAElevState(elev),
// 		},
// 	}
// 	return input
// }

func hra_funcs(input HRAInput, output *map[string][][2]bool, hraExecutable string) {

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		fmt.Println("json.Marshal error: ", err)
		return
	}

	ret, err := exec.Command("./hall_request_assigner/"+hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
	if err != nil {
		fmt.Println("exec.Command error: ", err)
		fmt.Println(string(ret))
		return
	}

	err = json.Unmarshal(ret, &output)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
		return
	}

}
func getCostOutput(input HRAInput) map[string][][2]bool {
	hraExecutable := "hall_request_assigner"
	output := new(map[string][][2]bool)

	hra_funcs(input, output, hraExecutable)

	return *output
}

// Velger å kun direkte sette lik hall calls istedenfor å or'e requests og hall calls.
func mergeHallAndRequests(elev elevator.Elevator, halls [][2]bool) elevator.Elevator {
	for f := 0; f < elevio.NumFloors; f++ {
		for b := 0; b < elevio.NumButtons-1; b++ {
			// if b < 2 {
			elev.Requests[f][b] = halls[f][b]
			// }
		}
	}
	return elev
}

// func main() {

// 	hraExecutable := ""
// 	switch runtime.GOOS {
// 	case "linux":
// 		hraExecutable = "hall_request_assigner"
// 	case "windows":
// 		hraExecutable = "hall_request_assigner.exe"
// 	default:
// 		panic("OS not supported")
// 	}

// 	input := HRAInput{
// 		HallRequests: [][2]bool{{false, false}, {true, false}, {false, false}, {false, true}},
// 		States: map[string]HRAElevState{
// 			"one": HRAElevState{
// 				Behavior:    "moving",
// 				Floor:       2,
// 				Direction:   "up",
// 				CabRequests: []bool{false, false, false, true},
// 			},
// 			"two": HRAElevState{
// 				Behavior:    "idle",
// 				Floor:       0,
// 				Direction:   "stop",
// 				CabRequests: []bool{false, false, false, false},
// 			},
// 		},
// 	}

// 	jsonBytes, err := json.Marshal(input)
// 	if err != nil {
// 		fmt.Println("json.Marshal error: ", err)
// 		return
// 	}

// 	ret, err := exec.Command("./hall_request_assigner/"+hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
// 	if err != nil {
// 		fmt.Println("exec.Command error: ", err)
// 		fmt.Println(string(ret))
// 		return
// 	}

// 	output := new(map[string][][2]bool)
// 	err = json.Unmarshal(ret, &output)
// 	if err != nil {
// 		fmt.Println("json.Unmarshal error: ", err)
// 		return
// 	}

// 	fmt.Printf("output: \n")
// 	for k, v := range *output {
// 		fmt.Printf("%6v :  %+v\n", k, v)
// 	}
// }
