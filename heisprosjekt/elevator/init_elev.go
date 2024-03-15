package elevator

import (
	"Heis/driver-go/elevio"
	"fmt"
	"os"
	"strings"
)

func InitElev() Elevator { //elevator_unitialized?

	return Elevator{
		Floor: -1,
		Dirn:  elevio.MD_Stop,
		Requests: initializeRequests(),
		Behaviour: EB_Idle,
		ElevID:    "",
		Failure:   false,
		//DoorOpenDuration_s: 3.0,
		//DoorOpen:           false,
	}
}

func OnInitBetweenFloors(elev Elevator) Elevator{

	elevio.SetMotorDirection(elevio.MD_Down)
	elev.Behaviour = EB_Moving
	elev.Dirn = elevio.MD_Down
	return elev
}



// Function to read cab calls from localBackup.txt, if the file exists. If not it returns an all false list of length
// Place this function in init_elev.go
// Intended to only be run _once_ inside the initElev() function
func checkAndLoadCabCalls() []bool {
	fileInfo, err := os.Stat("localBackup.txt")
	if os.IsNotExist(err) ||(fileInfo != nil && fileInfo.Size() == 0) {
		cabCalls := make([]bool, elevio.NumFloors)
		for i := range cabCalls {
			cabCalls[i] = false
		}
		return cabCalls
	} else {
		file, err := os.ReadFile("localBackup.txt")
		if err != nil {
			fmt.Println(err)
		}
		var cabCalls []bool
		fileContent := strings.TrimSpace(string(file))
		for _, char := range fileContent {
			if char == 't' {
				cabCalls = append(cabCalls, true)
			} else if char == 'f' {
				cabCalls = append(cabCalls, false)
			}
		}
		return cabCalls
	}
}

// Function to merge requests and initializeCabCalls
// Place this function in init_elev.go
// Put this function in init_elev() by setting Requests: initializeRequests
func initializeRequests() [elevio.NumFloors][elevio.NumButtons]bool {
	initializeCabCalls := checkAndLoadCabCalls()
	requests := [elevio.NumFloors][elevio.NumButtons]bool{}

	for floor := range requests {
		requests[floor][0] = false
		requests[floor][1] = false
		requests[floor][2] = initializeCabCalls[floor]
	}
	return requests
}
