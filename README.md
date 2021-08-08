# PiADCs
Libraries in Go to interface TI ADCs with Raspberry Pi. Currently only the ADS126x is supported but other ADCs with similar programming could also be added. The ADS1262 precision ADC makes for a cheap but very powerful Data aquisition system when paired with a Raspberry Pi.

## How to set up and install
1. Set up your Raspberry Pi with the Raspberry Pi OS (https://www.raspberrypi.org/software/) and connect it to the ADC. For an example of how to connect them with minimal noise for the highest precision measurments see the hookup example given in the schematics folder.

2. Make sure you have Go installed on your Raspberry Pi (Version 1.16 or higher) - https://golang.org/doc/install

3. Use the go get command
```bash
go get -u github.com/AnnaKnapp/piadcs
```
## How to use
Import the main library and the subfolder with the library for the adc you want to use (at this time on the ads126x is supported)

```go
import(
	"github.com/AnnaKnapp/piadcs"
	adc "github.com/AnnaKnapp/piadcs/ads126x"
)
```

Please check out the examples folder for usage examples. These examples are intended to run on a Raspberry Pi running Raspberry Pi OS and assume the connections shown in the schematics folder. The testing was done using a ProtoCentral ADS126x breakout board and a Raspberry Pi 4 B (other Pi models should also work). 

## Documentation
https://pkg.go.dev/github.com/AnnaKnapp/piadcs