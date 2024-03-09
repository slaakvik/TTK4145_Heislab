# TTK4145-Heislab
Heislab i emnet TTK4145 Sanntidsprogrammering



---- [Det jeg har gjort til nå: 8/3] ----
- Jeg har laget to funksjoner: ReceiveConn og TransmitConn
    + 
      + målet med ReceiveConn: dette 





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



 
