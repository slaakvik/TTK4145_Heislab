package fsm

import (
    "Heis/cost_fns"
    "Heis/driver-go/elevio"
    "Heis/elevator"
    "Heis/requests"
    "fmt"
    "sync"
)

// Mutex to synchronize access to shared map and channels
var (
    mapMutex              sync.Mutex
    sendMapToSlavesCh     chan<- map[string]elevator.Elevator
    /* getElevFromSlave      chan<- elevator.Elevator
    receiveMapFromMasterCh <-chan map[string]elevator.Elevator */
    newOrderCh            chan<- map[string]elevator.Elevator
)

// Initialize shared channels
func InitializeChannels(sendCh chan<- map[string]elevator.Elevator, getCh chan<- elevator.Elevator, receiveCh <-chan map[string]elevator.Elevator, orderCh chan<- map[string]elevator.Elevator) {
    sendMapToSlavesCh = sendCh
   /*  getElevFromSlave = getCh
    receiveMapFromMasterCh = receiveCh */
    newOrderCh = orderCh
}

// Finite state machine

func ButtonsAndRequests(masterPort string, elevatorID string, isMaster bool, elevUpdateRealtimeCh <-chan elevator.Elevator,
    drv_buttons chan elevio.ButtonEvent, sendMyselfToMaster chan elevator.Elevator, doorTimerChForBtnFSM chan bool) {
    fmt.Println("[ButtonsAndRequests] akkurat kommet inni")
    elev := elevator.InitElev()
    elev.ElevID = elevatorID

    // Initialize mapOfElevs
    mapMutex.Lock()
    mapOfElevs := make(map[string]elevator.Elevator)
    mapOfElevs[elev.ElevID] = elev
    mapMutex.Unlock()

    fmt.Println("______________")

    for {
        select {
        case a := <-elevUpdateRealtimeCh:
            fmt.Printf("[ButtonsAndRequests] mottok en heis på elevUpdateRealtimeCh: %v\n", a)
            mapMutex.Lock()
            elev = a
            mapOfElevs[elev.ElevID] = elev
            mapMutex.Unlock()
            fmt.Printf("[ButtonsAndRequests] mapOfElevs ser slik ut nå: %v\n", mapOfElevs)
        case a := <-drv_buttons:
            btn_floor := a.Floor
            btn_type := a.Button
            fmt.Printf("[ButtonsAndRequests] Button: %+v\n", a)
            if elev.Behaviour != elevator.EB_Moving && requests.ShouldClearImmediately(elev, btn_floor, btn_type) {
                fmt.Println("[ButtonsAndRequests] kom meg forbi if-state i drv_buttons")
                elev.Behaviour = elevator.EB_DoorOpen
                mapMutex.Lock()
                mapOfElevs[elev.ElevID] = elev
                mapMutex.Unlock()
                go SendElevToMaster(isMaster, elev, sendMyselfToMaster)
                elevio.SetDoorOpenLamp(true)
                elevio.SetButtonLamp(btn_type, btn_floor, true)
                doorTimerChForBtnFSM <- true

            } else {
                if isMaster {
                    mapMutex.Lock()
                    elev.Requests[btn_floor][btn_type] = true
                    mapOfElevs[elev.ElevID] = elev
                    mapOfElevs := cost_fns.RunCostFunc(mapOfElevs)
                    mapMutex.Unlock()
                    sendMapToSlavesCh <- mapOfElevs
                    fmt.Println("[ButtonsAndRequests] Sendte til alle")
                    newOrderCh <- mapOfElevs

                } else {
                    mapMutex.Lock()
                    elev.Requests[btn_floor][btn_type] = true
                    mapOfElevs[elev.ElevID] = elev
                    mapMutex.Unlock()
                    go SendElevToMaster(isMaster, elev, sendMyselfToMaster)
                }
            }
        }
    }
}

func FloorObstrStop(masterPort string, isMaster bool, elevatorId string, elevUpdateRealtimeCh chan<- elevator.Elevator, drv_floors chan int, doorTimerCh chan bool, timedOut chan int, lightsCh chan<- int, sendMyselfToMaster chan elevator.Elevator) {
    elev := elevator.InitElev()
    elev.ElevID = elevatorId
    fmt.Println("[FloorObstrStop] kommet meg inni")
    if elevio.GetFloor() == -1 {
        elev = elevator.OnInitBetweenFloors(elev)
        elevUpdateRealtimeCh <- elev
        elevator.Elevator_print(elev)
    }
    fmt.Println("[FloorObstrStop] kommet til for-loopen")
    for {
        select {
        case a := <-drv_floors:
            fmt.Printf("[FloorObstrStop] Floor: %+v\n", a)
            elev.Floor = a
            elevUpdateRealtimeCh <- elev
            elevio.SetFloorIndicator(elev.Floor)
            go SendElevToMaster(isMaster, elev, sendMyselfToMaster)
            switch elev.Behaviour {
            case elevator.EB_Moving:
                if requests.ShouldStop(elev) {
                    elevio.SetMotorDirection(elevio.MD_Stop)
                    elev.Behaviour = elevator.EB_DoorOpen
                    elevUpdateRealtimeCh <- elev
                    go SendElevToMaster(isMaster, elev, sendMyselfToMaster)
                    lightsCh <- 1
                    doorTimerCh <- true
                }

            default:

            }

        case <-timedOut:
            fmt.Println("[FloorObstrStop] gikk inn i timedout")
            elev = requests.OnDoorTimeout(elev, doorTimerCh, lightsCh, elevUpdateRealtimeCh)
            fmt.Println("gikk inn i timedout")
            elevUpdateRealtimeCh <- elev

            SendElevToMaster(isMaster, elev, sendMyselfToMaster)
            fmt.Println("[FloorObstrStop] Managed to send floor update")
        }
    }
}

func SendElevToMaster(isMaster bool, elev elevator.Elevator, sendMyselfToMaster chan elevator.Elevator) {
    if !isMaster {
        fmt.Println("[SendElevToMaster] kommet meg inni, og jeg er ikke master ")
        fmt.Println("[FloorObstrStop] sender heisen til sendMyselfToMasterCh")
        sendMyselfToMaster <- elev
    }
    fmt.Println("[SendElevToMaster] går ut")
}
