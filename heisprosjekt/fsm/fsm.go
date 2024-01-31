package fsm

import (
	"Heis/driver-go/elevio"
)

// Finite state machine

func Fsm_onInitBetweenFloors() {
	elevio.SetMotorDirection(elevio.MD_Down)
}


func Fsm_onRequestButtonPress(btn_floor int, btn_type elevio.ButtonType) {
	// TODO: Implement FSM logic for handling request button press
	

	select {
		case elevio.get


	case }
}
