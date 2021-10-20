# B1_multicasting
###***Martina Salvati***
##MULTICAST APPLICATION
Three different implementation of multicast
###Building
Run ```make``` to build the application's container image.
###Algorithm's type
Adding -t=TYPE_OF_ALGORITHM flag will set the type of algorithm with which the application can exchange message.
The possible type of algorithm : 
- TOC : totally ordered centralized with sequencer
- TOD : totally ordered distributed with scalar clock
- CO  : casually ordered with vector clock 
###Delayed Messages
Adding the -delay=MAX_DELAY flag will turn on delayed message mode. The maximum delay can be set in MAX_DELAY. In this mode, messages are occasionally delayed by a random amount of time. The nondeterminism provided by this randomization helps stress the underlying algorithm.
###Verbose Mode
Adding the -v (--verbose) flag will turn on verbose mode, which will print logging information to standard error. This information includes details about all messages sent and received, as well as round timeout information.
###Testing
The project comes with a automated test suite located under the test directory. Each test configures a set of processes with different numbers of messages to send.
To run all tests, run the command ```make test```.


Each test creates a subprocess for all processes in the system and then reads in the messages that the subprocesses deliver. 

- In case of TOC/TOD :
The tests then verify that the correct number of messages were delivered, and that all processes delivered the same messages in the same order.
- In case of CO : The tests then verify that the correct number of messages were delivered, and that all processes delivered the same message respecting the casual order.
