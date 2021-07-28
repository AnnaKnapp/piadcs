package main

import (
	"fmt"
	"log"
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

func main() {

	exiter := make(chan os.Signal) //This is the channel that will receive the exit signal (ctrl-c)

	signal.Notify(exiter, os.Interrupt, syscall.SIGTERM) //this is triggered when ctrl-c is pressed

	datafile, err := os.Create("test.txt") //this makes the datafile. If you don't change it it will make it right in the examples folder.
	if err != nil {
		log.Fatal(err)
	}
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
	startpin := gpioreg.ByName("22")
	if startpin == nil {
		log.Fatal("Failed to find start pin (22)")
	}

	//Sets the data ready pin to the Rasberry Pi pin 12 (see ADS126x datasheet for what this pin does)
	drdypin := gpioreg.ByName("4")
	if drdypin == nil {
		log.Fatal("Failed to find data ready pin (4)")
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
	Power.Setregister([]byte{adc.POWER_vbias_enabled})

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
	Inpmux.Setregister([]byte{adc.INPMUX_muxP_AIN9, adc.INPMUX_muxN_AINCOM})

	//This creates a byte slice which we will use to write the values to the registers specified. It is critical that the values listed are in order and no register is skipped. If there is a register that you don't want to change the value of you still need to speficy it if it falls between two others that you want to change. You can use the built in default value for this as shown in this example with Mode0. See the ADS126x datasheet section 9.5.7 for an explanation of why this is the case
	registerdata := []byte{Power.Setvalue, Interface.Setvalue, Mode0.Setvalue, Mode1.Setvalue, Mode2.Setvalue, Inpmux.Setvalue}

	//This actually writes the data to the register
	piadcs.WriteToConsecutiveRegisters(spi0, Power.Address, registerdata)

	//This reads the data from registers so we can check that the ADC is working and that we correctly wrote the data to the registers
	incomingregdata := piadcs.ReadFromConsecutiveRegisters(spi0, Power.Address, 6)

	fmt.Println(incomingregdata)

	//In order to have timestamps for incoming data we need a starting point
	beginning := time.Now()

	//these are used to keep track of how well the data is being transferred
	var successes int
	var failures int

	//this go function will run concurrently with the rest of the code. It waits for you to press ctrl-c and then runs the closer function to exit gracefully
	go closer(exiter, datafile)

	//this go function isn't needed but I keep it to keep track of how successful the communication between the ADC and the Pi is. It will print the error rate every 5 seconds
	go func() {
		for {
			time.Sleep(time.Second * 5)
			errorrate := float32(failures) / float32(successes)
			fmt.Println(errorrate)
		}
	}()

	//This starts the conversions on the ADS126x. It is critical that this is here otherwise there won't be any data coming in when we try to read
	if err := startpin.Out(gpio.High); err != nil {
		log.Fatal(err)
	}

	//this for loop will take data continuously, convert it. and write it to the file until ctrl-c is pressed
	for {
		adcdata, err := adc.ContinuousReadCHK(spi0, drdypin)
		if err == nil {
			converteddata := adc.ConvertData(adcdata)
			timestamp := time.Since(beginning).Milliseconds()
			outputstring := strconv.FormatInt(int64(timestamp), 10) + "," + strconv.FormatFloat(float64(converteddata), 'f', -1, 64) + "\n"
			// this writes the converted data to the file with the format "time, data"
			datafile.WriteString(outputstring)
			successes = successes + 1
		} else {
			//if you are getting a high error rate you can print the error to see whats going wrong. I have found that its impossible to get an error rate of 0 and there will always be some instances of the SPI communication failing. I belevie this is because of the raspberry pi operating system not being real time and the CPU taking a break to go do something else. Errors are not recorded in the datafile.
			failures = failures + 1
		}
	}

}
