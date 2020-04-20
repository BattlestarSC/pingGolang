Golang ping program for Cloudflair internship challenge 2020.

Design choices:
- a main function handling user input, metrics, and an interface to the ping go routine
- ping go routine is an interface to all the backend ping functions
- a dialer/net.Conn is created for every ping sent for a few reasons, chief of which is that it allows me to set a precise timeout on the connection (which the docs do not cover well and is unclear if it is reusable or not) and secondly it allows me better debug information (each new conn is more information if it fails). 
