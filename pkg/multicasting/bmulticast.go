package multicasting

import "b1multicasting/pkg/basic"

//IMPLEMENTING THE BMULTICAST ALGORITHM
//To B-multicast(g, m): for each process p of the group g , send(p, m)
func (c *Conns) XMulticast(p string, m basic.Message) error {
	for i := 0; i < len(c.conns); i++ {
		err := c.conns[i].Send(p, m)
		if err != nil {
			return err
		}
	}
	return nil
}

//On receive(m) at p: B-deliver(m) at p.
func XDeliver(m basic.Message) error {

	return nil
}
