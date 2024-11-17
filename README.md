# timebox

Golang library to interface with a Divoom Timebox Mini.

The provided 'bluetooth' library can be used to establish a Bluetooth RFCOMM socket on Linux-based systems. This can be supplied to the main library to communicate with the Timebox to set the active screen, update the date/temperature/weather, and display the score.



## References

I have found a few good sources for the underlying protocol, however none of them are comprehensive. https://github.com/MarcG046/timebox/blob/master/doc/protocol.md has a good subset of information which is 100% accurate to the Timebox Mini. https://github.com/RomRider/node-divoom-timebox-evo was written for the Timebox Evo, and there are some differences between the Mini and the Evo (such as which screens are supported), however there is also some overlap in the protocols and is a good source to augment the first protocol link.

There are some other projects, primarily in Python, which interface with the Timebox Mini which provided additional information, such as https://github.com/jbfuzier/timeboxmini/ and https://github.com/mathoudebine/homeassistant-timebox-mini. 
