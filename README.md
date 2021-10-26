# B1_multicasting
***Martina Salvati***

## MULTICAST APPLICATION

Three different implementation of multicast
### Building
Run ```docker-compose up``` to build the application's containers image.
### Algorithm's type
- B   : basic multicast
- TOC : totally ordered centralized with sequencer
- TOD : totally ordered distributed with scalar clock
- CO  : casually ordered with vector clock

### Testing
The project comes with a test suite located under the test directory (pkg/multicasting/test). Each test configures a set of processes with different numbers of messages to send.

Each test creates a subprocess for all processes in the system and then reads in the messages that the subprocesses deliver.

- In case of TOC/TOD :
  The tests then verify that the correct number of messages were delivered, and that all processes delivered the same messages in the same order.
- In case of CO : The tests then verify that the correct number of messages were delivered, and that all processes delivered the same message respecting the casual order.
