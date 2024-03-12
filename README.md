# TTK4145-Heislab
Heislab i emnet TTK4145 Sanntidsprogrammering
liveshare: 

Imperative shell functional core.


[___Tirsdag 12/3___]
før den har fått connection fra master. Vi må gjøre slik at master tar kun med connectionsa som er i mappet i betrakning (kanskje det allerede er slik?)
- noe er galt med når slaven skal transmitte. Så master og slave har bygd tcp-connection, så tar master slave i betrakning? Mens slaven kjører fsm må den sende sin heis til master, og connection som slaven skal bruke er fortsatt nil, altså connection er ennå ikke oppdatert for slaven? vi burde kanskje heller sende conn i en channel hos slaven, så vi er sikre på at når det forsøkes å sendes så har den faktisk en connection. Vi må sjekke litt mer nøye på alle kanalene vi har og hva vi egt trenger. Tegne opp trådene og hvor de sender. Foreløpig er det litt uklart hvilke tråder som kjøres og hvilke kanaler som mottar/ ikke mottar data til hvilken tid. Vi må bestemme oss for om vi skal ha isMaster sjekk inni alle trådene eller at Master og Slave har sine respektive tråder som de kjører. Dersom sistnevnte er det viktig at det ofte sjekkes om man er master ikke, da dette kan oppdateres og endres fortløpende.
- Vi må ta en titt på den for-loopen Jonathan snakka om. Her må master også sjekkes på noe vis.

Vi må gå gjennom flyten i trådene og kanalene mer nøye og tegne det opp, fordi det er her det går galt. En case krever en data som kun oppdateres i en annen osv. 


Hva jeg sjekka:
- NotifyMaster er det trøbbel hos. Det er her slave prøver å skrive til en connection som er nil. For mye baseres på rett timing, og det er ikke bra. De ulike delene av systemet må vente på hverandre, hvis nødvendig. Og de må flyten må også være rett
