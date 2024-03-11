# TTK4145-Heislab
Heislab i emnet TTK4145 Sanntidsprogrammering




---- [Master-slave: 11/3] -----
- Slik det er nå vil det muligens oppstå problemer, i.f.t heartbeat, hvis flere slaves sier de er master. Dette på grunn av den ekstra p.Master == "" i peers.Receiver-funksjonen
- Ide: man sjekker hele tiden hvor mange som er på nett. Dersom lengden av p.Peers > 1, velger man den som har minst prosesside til å være master. 
      Når man blir frakoblet nett, næææ glem det




 
