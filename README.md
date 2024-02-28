# TTK4145-Heislab_GuttaHeiser
Heislab i emnet TTK4145 Sanntidsprogrammering
liveshare: 


---[Gjort onsdag 28/2]---
- Hente ut to siste siffer i prosessIDen (variabelnavn eleviId i main) som en integer og kaller variabelen for processID:
	+ processID, err := strconv.Atoi(eleviId[len(eleviId)-2:])
 	+ Atoi returnerer to variabler, derfor må err også være der ^.


---[Plan for onsdag 28/2]---

- Lage flow-chart for systemet, for å sjekke ulike feil eller tilfeller som kan oppstå. Dette skaper viktig oversikt.

- Master-Slave modulene:
	+ Endre på funksjonen DetermineMaster(), slik at den velger det første elementet i p.Peers til å være master. Resterende blir slaver, hvorav den "første" slaven, altså det andre elementet, blir backup for master.
	+ Skal MasterID være slik den er nå? Er det andre måter dette kan gjøres på? 





