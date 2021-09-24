package multicasting

import (
	"log"
)

//TOTAL ORDERED CENTRALIZED MULTICAST
//Ogni processo invia il proprio messaggio di update al sequencer
//Il sequencer assegna ad ogni messaggio di update un numero di sequenza univoco e poi invia in multicast (TOCMulticast) il messaggio
//a tutti i processi, che eseguono gli aggiornamenti in ordine in base al numero di sequenza
func (seq *Sequencer) TOCMulticast(g string, m SeqMessage, c *Conns) error {
	ch := make(chan bool, len(c.conns))
	//da aggiungere la garanzia che i messaggi siano arrivati a destinazione tramite un pool di threads e utilizzo canali
	for i := 0; i < len(c.conns); i++ {
		i := i
		go func() {
			err := c.conns[i].SendSeq(g, m, &ch) //one to one send operation
			if err != nil {
				return
			}
		}()
		//inserisco la risposta nel canale gestito dalla pool
	}
	//check if the message correctly arrived to the nodes
	for i := 0; i < len(c.conns); i++ {
		r := <-ch //lettura del canale
		if r != true {
			log.Println("Message not arrived to nodes ", c.conns[i].GetTarget())
			//prova a rinviarlo
		} else {
			log.Println("Message correctly sent to ", c.conns[i].GetTarget())
		}
	}
	return nil
}
