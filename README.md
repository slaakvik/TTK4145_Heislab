# TTK4145-Heislab_GuttaHeiser
Heislab i emnet TTK4145 Sanntidsprogrammering
liveshare: 

Må endre alt som er skrevet som pass by reference. Bør istedet la funksjonene våre kun gjøre utregninger, også endre på heisen vår i main hovedsakelig. 
Kanskje kalle på hardware i fsm fortsatt? så lenge vi kun gjør det et sted så gjør det det enklere å debugge. 

Imperative shell functional core.


Tirsdag 5/3:
cost function kjøres for hver nye etasje. må få heisen til å kjøre cost function og lese den for å finne ut av sine ordre. foreløpig kjører ikke heisen cost function. cost funksjonen kjøres, men eneste vi gjør med resultatet er å printe. må få det inn i request listen til heisen. må altså endre litt på onrequestbuttonpress blandt annet. 

Problem med individuell heis: Dersom vi får hall call både opp og ned i samme etasje, og ikke har andre requests, vil heisen først cleare i den retningen den gikk i, deretter begynne å bevege seg i motsatt retning uten å cleare ordren og fortsatt tilsynelatende ha døren åpen. klarer da ikke å ta inn nye ordre og beveger seg opp/ned helt til det ikke går mer. 

---[Gjort tirsdag 5/3]---
Følgende kode slår sammen lokal ip (siste tre siffer) og prosess-id (5 siffer) til en ny string.
Siden det er string kan det kanskje funke å bare ta hele ipv4-stringen + prosess-id.



```Go
import (
	"os"
	"strconv"
)

process_id := os.Getpid()
	fmt.Println("Process Id is", process_id)

	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	var ipv4String string
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			ipv4String = ipv4.String()
			break
		}
	}
	fmt.Println("IPv4: ", ipv4String)

	//uniqueID = ipv4String(-2:)
	uniqueID := ipv4String[len(ipv4String)-3:] + "_" + strconv.Itoa(process_id)
	fmt.Println("Unique ID:", uniqueID)
```

---[Gjort torsdag 29/2]---
Påbegynt funksjon i cost_fns.go som itererer gjennom en array av heisstructs og lager en felles hall call-matrise gjennom or-funksjoner.
Neste steg er å iterere gjennom heisene sine states og generere inputen som trengs til HRAInput, og dermed kan brukes i HRA_funcs-funksjonen.

```Go
func InputToCost(elevatorArray) HRAInput {
	N := len(elevatorArray)
	input := make([][2]bool, N)

	for i, elevator := range elevatorArray {
		// Initialize OR values for the current matrix
		or1 := false
		or2 := false

		// Compute OR operation for corresponding elements
		for _, row := range elevator {
			or1 = or1 || row[0]
			or2 = or2 || row[1]
	}
	input[i] = [2]bool{or1, or2}

	/*

	HallRequests := elevatorArray(1).requests[][:2]
	*/
	input := HRAinput { }
		// Creating a }
	return input
}
```


---[Gjort onsdag 28/2]---
- Hente ut to siste siffer i prosessIDen (variabelnavn eleviId i main) som en integer og kaller variabelen for processID:
	+ processID, err := strconv.Atoi(eleviId[len(eleviId)-2:])
 	+ Atoi returnerer to variabler, derfor må err også være der ^.


---[Plan for onsdag 28/2]---

- Lage flow-chart for systemet, for å sjekke ulike feil eller tilfeller som kan oppstå. Dette skaper viktig oversikt.

- Master-Slave modulene:
	+ Endre på funksjonen DetermineMaster(), slik at den velger det første elementet i p.Peers til å være master. Resterende blir slaver, hvorav den "første" slaven, altså det andre elementet, blir backup for master.
	+ Skal MasterID være slik den er nå? Er det andre måter dette kan gjøres på? 





