# TTK4145-Heislab
Heislab i emnet TTK4145 Sanntidsprogrammering



---- [Det jeg har gjort til nå: 8/3] ----
- Jeg har laget to funksjoner: ReceiveConn og TransmitConn
    + 
      + målet med ReceiveConn: dette 


```Go
p.Master = ""
if id != "" {
    if _, idExists := lastSeen[id]; !idExists {
        p.Master = id
        updated = true
        // Inform other elevators about the master
        for k := range lastSeen {
            if k != id {
                conn.WriteTo([]byte("MASTER:"+id), addr) // Notify others about the master
            }
        }
    }

    lastSeen[id] = time.Now()
}



---- [Master-slave: 11/3] -----
- Slik det er nå vil det muligens oppstå problemer, i.f.t heartbeat, hvis flere slaves sier de er master. Dette på grunn av den ekstra p.Master == "" i peers.Receiver-funksjonen




 
