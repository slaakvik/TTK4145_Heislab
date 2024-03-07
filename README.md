# TTK4145-Heislab_GuttaHeiser
Heislab i emnet TTK4145 Sanntidsprogrammering
liveshare: 

Imperative shell functional core.

 [___Torsdag 7/3___]
 
 
 Logikk: 
 
 ulike caser i fsm: 
 
 
 button press, newOrdders, updatedElevs + de andre som er der nå.
 
 
 Ved button press: 
 
 
 	hvis slave:  
 		Hvis hall call: lag en kopi av deg selv, legg til den nye requesten inn i request listen til kopien. Send kopien til master (da vil updatedElevs inntreffe).
  		Hvis cab call: legg til den nye requesten til i din request liste. Send deg selv til master. 
 
 
	 Hvis master: 
 		Hvis hall call: lag en kopi av deg selv, legg til den nye requesten inn i request listen til kopien. Kjør cost function på og bruk kopien som deg selv. 
  			Send resultat fra cost function til alle.
  		Hvis cab call: legg til den nye requesten til i din request liste. Kjør cost function.
   			Send resultat fra cost function til alle. 
      		Likt uavhengig call: Send resultat fra costfunksjon til kanalen for newOrders. (ikke via nettet, kun kall på kanalen). 



Ved newOrders: 
NewOrders er en kanal hvor det blir sendt et map som inneholder resultat fra cost function.
	
 	Like oppgaver uavhengig av om du er master eller slave. 
 	når kanalen oppdateres skal slaver lese fra kanalen og utføre ordren som hører til den sin egen ID.
  	Master skal helst ikke lese fra nettet hvilke nye ordre den skal ta? kanskje den kan det? eventuelt ha et eget case for neworders som master sender til, hvor den 	skal utføre ordrene, og slavene ikke skal gjøre noe. 



Ved updatedElevs:
updatedElevs er en kanal hvor det blir sendt en Elevator struct hver gang noe med heisen oppdateres. alle slavene sender heisen sin til denne kanalen. 
	
 	Hvis slave:
		ikke gjør noe, gå ut av casen. 
	Hvis master: 
		finn ut hvilken id den oppdaterte heisen har, og oppdater heisen i arrayen/mapet med tilhørende id.
 		Kjør så cost function på de nye heisene, og send cost function ut til newOrders kanalen. 
		





[__Onsdag 6/3__]:


Master slave logikk: 
sjekk i main om man er master eller slave. f.eks. lag en funksjon som returnerer en bool som er true hvis master og false ellers. 
FSM tar inn et ekstra bool parameter som er "master" som blir true eller false. 
Når heisen er slave, skal den alltid sende ordre til master på button press, og vente på oppdatert cost function. 
Når heisen er master, skal den kjøre cost function når den mottar nye ordre, enten fra seg selv eller ordre fra andre heiser. og deretter sende ut oppdatert cost function, også til seg selv. 

Kanskje en if setning som sjekker if master, og hvis true så kjører den cost function, ellers, sender den ordren med en gang, resten blir vel ganske likt for alle heisene. Kanskje sende ordren til seg selv hvis master? Kan da få problemer dersom heisen ikke er kobblet til nettet. 
Kanskje å la master kjøre cost function og sende oppdatert HRA til alle, men trenger ikke sende til seg selv i tillegg, siden man allerede har infoen. kanskje like greit å bare sende likevel, men ikke nødvendigvis bruke det?

Hver heis gjør seg selv til master dersom den er alene på nettet. (siden den da har lavest prosessID av alle på nettet). Når en ny heis kobler på, skal heisene som er enkeltheiser gjøre seg selv til slave. dette vil gjøre at dersom en heis er master og har en slave, også kommer det en enkeltheis som er master, skal den heisen uten tilhørende slave bli slave selv. 
dersom to enkeltheiser kobler på gjør de begge seg til slave. 
Vi skal derfor alltid ha en sjekk som sjekker om vi har en master, hvis ikke, så skal den med lavest prosessID bli master. 

stop og obstruction funksjonaliteten vår påvirker kun elevio (hardware), vi endrer ikke statesene til heisen vår, hvilket vil bety at vi ikke sender at heisen er obstructet mens den er det. 




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





