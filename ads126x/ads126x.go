//Package supports use of Texas Insturments ADS1262 and ADS1263 32-bit analog to digital converters
package ads126x

import (
	"errors"
	"log"
	"math"
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/spi"
)

//funcs to write - read data, convert data, read pulse, startup, data to file

//function to restart the ADS126x based on fig 159 from the datasheet
func InitSetup(start, pwdn gpio.PinIO) {
	if err := pwdn.Out(gpio.Low); err != nil {
		log.Fatal(err)
	}

	time.Sleep(500 * time.Millisecond)

	pwdn.Out(gpio.High)

	if err := start.Out(gpio.Low); err != nil {
		log.Fatal(err)
	}

	time.Sleep(2 * time.Second)

}

var conversionbytes []byte = make([]byte, 6)
var empty []byte = make([]byte, 6)
var checksumfail bool = false

//This function requires that the checksum be enabled in checksum mode and the status byte enabled. It reads the data in continuous mode - meaning that it waits for the data ready signal on the DRDY pin and then begins reading. The output is an unconverted 32 bit integer. If the checksum fails, SPI fails, or DRDY pin times out it will output an error and a value of zero.
func ContinuousReadCHK(connection spi.Conn, drdy gpio.PinIO) (int32, error) {

	if drdy.WaitForEdge(-1) {

		if err := connection.Tx(empty, conversionbytes); err != nil {
			return 0, errors.New("SPI connection failed")
		} else if conversionbytes[5] != (conversionbytes[1]+conversionbytes[2]+conversionbytes[3]+conversionbytes[4]+0x9B)&255 {
			return 0, errors.New("Checksum Failed - data transmission error occurred")
		} else {
			rawdata := int(conversionbytes[1])<<24 | int(conversionbytes[2])<<16 | int(conversionbytes[3])<<8 | int(conversionbytes[4])
			tobeconverted := int32(rawdata)
			return tobeconverted, nil
		}
	}
	return 0, errors.New("Pin timeout")
}

//This converts the output of the read function to a voltage between -2.5 and 2.5. It does not account for gain.
func ConvertData(data int32) float64 {
	converteddata := float64(data) * float64(2.5/math.Pow(2, 31))
	return converteddata
}
