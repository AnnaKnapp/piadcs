package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/AnnaKnapp/piadcs"
	adc "github.com/AnnaKnapp/piadcs/ads126x"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/host"
)

//This function allows the program to be closed by pressing ctrl-c and exit gracefully
func closer(exiter chan os.Signal, dataFile *os.File) {
	<-exiter
	dataFile.Sync()
	dataFile.Close()
	fmt.Printf("Exited gracefully")
	os.Exit(0)
}

//K-type thermocouple conversion polynomial constants
//constants for converting Temp to Emf
var c0 float64 = -0.176004136860e-1
var c1 float64 = 0.389212049750e-1
var c2 float64 = 0.185587700320e-4
var c3 float64 = -0.994575928740e-7
var c4 float64 = 0.318409457190e-9
var c5 float64 = -0.560728448890e-12
var c6 float64 = 0.560750590590e-15
var c7 float64 = -0.320207200030e-18
var c8 float64 = 0.971511471520e-22
var c9 float64 = -0.121047212750e-25

var a0 float64 = 0.118597600000
var a1 float64 = -0.118343200000e-3
var a2 float64 = 0.126968600000e3

//Based on NIST's data
func calculateEmfFromTemp(temp float64) float64 {
	emf := c0 + c1*temp + c2*math.Pow(temp, 2) + c3*math.Pow(temp, 3) + c4*math.Pow(temp, 4) + c5*math.Pow(temp, 5) + c6*math.Pow(temp, 6) + c7*math.Pow(temp, 7) + c8*math.Pow(temp, 8) + c9*math.Pow(temp, 9) + a0*math.Pow(math.E, (a1*math.Pow(temp-a2, 2)))
	return emf
}

//K-type thermocouple conversion polynomial constants
//constants for converting EMF to Temp
var d0 float64 = 0
var d1 float64 = 2.5083551e1
var d2 float64 = 7.860106e-2
var d3 float64 = -2.503131e-1
var d4 float64 = 8.315270e-2
var d5 float64 = -1.228034e-2
var d6 float64 = 9.804036e-4
var d7 float64 = -4.413030e-5
var d8 float64 = 1.057734e-6
var d9 float64 = -1.052755e-8

//Based on NIST's data
func calculateTempFromEmf(emf float64) float64 {
	temp := d0 + d1*emf + d2*math.Pow(emf, 2) + d3*math.Pow(emf, 3) + d4*math.Pow(emf, 4) + d5*math.Pow(emf, 5) + d6*math.Pow(emf, 6) + d7*math.Pow(emf, 7) + d8*math.Pow(emf, 8) + d9*math.Pow(emf, 9)
	return temp
}

func main() {

	exiter := make(chan os.Signal) //This is the channel that will receive the exit signal (ctrl-c)

	signal.Notify(exiter, os.Interrupt, syscall.SIGTERM) //this is triggered when ctrl-c is pressed

	datafile, err := os.Create("test.txt") //this makes the datafile. If you don't change it it will make it right here in the examples folder.
	if err != nil {
		log.Fatal(err)
	}

	//this go function will run concurrently with the rest of the code. It waits for you to press ctrl-c and then runs the closer function to exit gracefully
	go closer(exiter, datafile)

	//see periph documentation
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
	//Sets the powerdown pin to the Rasberry Pi pin 27 (see ADS126x datasheet for what this pin does)
	pwdnpin := gpioreg.ByName("27")
	if pwdnpin == nil {
		log.Fatal("Failed to find powerdown pin (27)")
	}
	//Sets the start pin to the Rasberry Pi pin 26 (see ADS126x datasheet for what this pin does)
	startpin := gpioreg.ByName("26")
	if startpin == nil {
		log.Fatal("Failed to find start pin (26)")
	}

	//Sets the data ready pin to the Rasberry Pi pin 12 (see ADS126x datasheet for what this pin does)
	drdypin := gpioreg.ByName("12")
	if drdypin == nil {
		log.Fatal("Failed to find data ready pin (12)")
	}

	//Since we will be reading the data ready pin we need to configure it as an input pin and specify which edge we will be looking for to know that new conversion data is ready to be read
	if err := drdypin.In(gpio.PullUp, gpio.FallingEdge); err != nil {
		log.Fatal(err)
	}

	//This function resets the ADC
	adc.InitSetup(startpin, pwdnpin)

	//See periph documentation
	port, err := spireg.Open("")
	if err != nil {
		log.Fatal(err)
	}

	// In this example we are using the raspberry Pi's SPI0. Other SPI buses would work but this corresponds the example schematic provided. The speed specified here is 6MHz. Slower will also work as long as the data rate isn't too high. I have not had success with higher speeds but it may be possible with some extra steps
	spi0, err := port.Connect(physic.KiloHertz*600, spi.Mode1, 8)
	if err != nil {
		log.Fatal(err)
	}

	//The following sets the values that we want to write to the register. All register options are defined in the constants.go file. In this example we set the interface register such that the status byte is enabled and the checksum byte is enabled in checksum mode.
	Power := piadcs.NewRegister("Power", adc.POWER_address)
	Power.Setregister([]byte{adc.POWER_vbias_enabled, adc.POWER_intref_enabled})

	Interface := piadcs.NewRegister("Interface", adc.INTERFACE_address)
	Interface.Setregister([]byte{adc.INTERFACE_status_enabled, adc.INTERFACE_crc_checksum})

	//This sets Mode0 to the default values. The default values for programming registers are defined in the ADS126x datasheet
	Mode0 := piadcs.NewRegister("Mode0", adc.MODE0_address)
	Mode0.Setregister([]byte{adc.MODE0_default})

	//This sets the ADC filter to sync4 mode (defined in the datasheet)
	Mode1 := piadcs.NewRegister("Mode1", adc.MODE1_address)
	Mode1.Setregister([]byte{adc.MODE1_filter_FIR})

	//This sets the gain to 1v/v and the data rate to 400 samples per second
	Mode2 := piadcs.NewRegister("Mode2", adc.MODE2_address)
	Mode2.Setregister([]byte{adc.MODE2_GAIN_1, adc.MODE2_DR_20})

	Inpmux := piadcs.NewRegister("INPMUX", adc.INPMUX_address)
	Inpmux.Setregister([]byte{adc.INPMUX_muxP_tempSensorP, adc.INPMUX_muxN_tempSensorN})

	//This creates a byte slice which we will use to write the values to the registers specified. It is critical that the values listed are in order and no register is skipped. If there is a register that you don't want to change the value of you still need to speficy it if it falls between two others that you want to change. You can use the built in default value for this as shown in this example with Mode0. See the ADS126x datasheet section 9.5.7 for an explanation of why this is the case
	registerdata := []byte{Power.Setvalue, Interface.Setvalue, Mode0.Setvalue, Mode1.Setvalue, Mode2.Setvalue, Inpmux.Setvalue}

	//This actually writes the data to the register
	piadcs.WriteToConsecutiveRegisters(spi0, Power.Address, registerdata)

	beginning := time.Now()

	var temperatures []float64
	var n int
	var sum float64

	for {

		//Start reading the onboard ADC temp (as an estimate of cold junction)
		if err := startpin.Out(gpio.High); err != nil {
			log.Fatal(err)
		}

		//Take 5 temperature readings
		for n < 5 {
			n = n + 1
			adcdata, err := adc.ContinuousReadCHK(spi0, drdypin)
			if err == nil {
				converteddata := adc.ConvertData(adcdata)
				temperature := (((converteddata * 1000000) - 122400) / 420) + 25
				temperatures = append(temperatures, temperature)
				//fmt.Println(temperature)
			}
		}

		n = 0

		//Average the temperature readings
		for _, v := range temperatures {
			sum += v
		}
		tempaverage := sum / float64(len(temperatures))

		//reset variables for next cycle
		sum = 0
		temperatures = []float64{}

		//convert ambient temp to a voltage that we can use for cold junction compensation. The 0.7 degrees subtracted is done due to self-heating of the ADS1262/3
		ambient := calculateEmfFromTemp(tempaverage - 0.7)

		//Stop conversions so we can change registers to read from the thermocouple instead of the internal temp sensor
		if err := startpin.Out(gpio.Low); err != nil {
			log.Fatal(err)
		}

		//Set PGA gain to 32 and data rate to 20
		Mode2.Setvalue = 0
		Mode2.Setregister([]byte{adc.MODE2_GAIN_32, adc.MODE2_DR_20})

		//set input pins to AIN9 for positive and AINCOM for negative
		Inpmux.Setvalue = 0
		Inpmux.Setregister([]byte{adc.INPMUX_muxP_AIN9, adc.INPMUX_muxN_AINCOM})

		//Write the new values to the ADC registers
		piadcs.WriteToConsecutiveRegisters(spi0, Mode2.Address, []byte{Mode2.Setvalue, Inpmux.Setvalue})

		//Restart conversions on the ADC
		if err := startpin.Out(gpio.High); err != nil {
			log.Fatal(err)
		}

		//Take thermocouple voltage reading
		adcdata, err := adc.ContinuousReadCHK(spi0, drdypin)
		if err == nil {
			converteddata := adc.ConvertData(adcdata)
			tctemp := (calculateTempFromEmf(ambient + ((converteddata / 32) * 1000) - 0.2))
			fmt.Println(tctemp)
			timestamp := time.Since(beginning).Milliseconds()
			outputstring := strconv.FormatInt(int64(timestamp), 10) + "," + strconv.FormatFloat(float64(converteddata), 'f', -1, 64) + "," + strconv.FormatFloat(float64(tctemp), 'f', -1, 64) + "\n"
			// this writes the converted data to the file with the format "time, data"
			datafile.WriteString(outputstring)
		} else {
			//if you are getting a high error rate you can print the error to see whats going wrong. I have found that its impossible to get an error rate of 0 and there will always be some instances of the SPI communication failing. I belevie this is because of the raspberry pi operating system not being real time and the CPU taking a break to go do something else. Errors are not recorded in the datafile.
			fmt.Println("checksum fail")
		}

		//Stop conversions
		if err := startpin.Out(gpio.Low); err != nil {
			log.Fatal(err)
		}

		//Switch back to measuring chip temperature to estimate cold junction
		//Set the PGA gain to 1 and the data rate to 20
		Mode2.Setvalue = 0
		Mode2.Setregister([]byte{adc.MODE2_GAIN_1, adc.MODE2_DR_20})

		Inpmux.Setvalue = 0
		Inpmux.Setregister([]byte{adc.INPMUX_muxP_tempSensorP, adc.INPMUX_muxN_tempSensorN})

		//Write new values to
		piadcs.WriteToConsecutiveRegisters(spi0, Mode2.Address, []byte{Mode2.Setvalue, Inpmux.Setvalue})

	}

}
