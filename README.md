# TTK4145-Heislab
Heislab i emnet TTK4145 Sanntidsprogrammering
liveshare: 


---[Plan for neste gang]---'
- vi diskuterer hva som gjennstår og hva som trengs
- vi må først avklare hva vi skal gjøre i dag, konkret.
- må prøve å bygge TCP forbindelse "automatisk" dersom peersa er på nettet:
    + hente ut IP-adresse og prosess-id kan er funksjoner som trengs og bør lages
    + TCP: må sende ack for hver melding som sendes? I allfall med tanke på "ordre utført"
    + 

   
- ting jeg kommer på:
  + primary-backup logikk og automatikk: master er primary, og en "random" annen heis er backup
        -- tanke: hva skjer med backup'ene når eksempelvis to "individulle heissystemer" merges sammen?
  + master-slave logikk
  + cost-funksjonen må læres å brukes
        -- tanke: er det best at cost-funksjonen brukes av alle heiser uansett? selv enkeltheiser? ja
  + fortsatt litt som gjenstår på TCP, men det er på vei.
  + vi burde gå gjennom main og bestemme oss for hva som trengs å være der til nå, og hva vi evt kan flytte til andre filer slik at koden blir mer      oversiktlig.
  + order og struktur er veldig viktig!




En siste:
- prøv i testClient å skrive elevatorCh<-elev for å skrive elev til kanalen




 
