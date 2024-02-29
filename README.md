# TTK4145-Heislab_GuttaHeiser
Heislab i emnet TTK4145 Sanntidsprogrammering
liveshare: 

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





