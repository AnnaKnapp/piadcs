//A quick example demonstrating how to set, read from, and reset registers on the ADS126x
package main

import (
	"fmt"
	"log"

	"github.com/AnnaKnapp/piadcs"
	adc "github.com/AnnaKnapp/piadcs/ads126x"

	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/host"
)

func main() {

	//Initilalize periph (see periph documentation)
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	//Set the required GPIO Pins

	//Sets the powerdown pin to the Rasberry Pi pin 27 (see ADS126x datasheet for what this pin does)
	pwdnpin := gpioreg.ByName("27")
	if pwdnpin == nil {
		log.Fatal("Failed to find powerdown pin (27)")
	}
	//Sets the start pin to the Rasberry Pi pin 22 (see ADS126x datasheet for what this pin does)
	startpin := gpioreg.ByName("22")
	if startpin == nil {
		log.Fatal("Failed to find start pin (22)")
	}

	//Set up SPI

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

	//This function resets the ADC - to ensure its ready to write to the registers
	adc.Restart(startpin, pwdnpin)

	// Initialize a register using the built in register object. Register addressed
	// for all registers on the ADS126x are stored in the constants file for easy use
	Power := piadcs.NewRegister("Power", adc.POWER_address)

	//Change the settings for the register. The settings for each register are defined in the constants file. In order for this to work properly either set it with the default value or ensure that you specify ALL possible settings. In this case even though the internal ref is enabled in the default we still need to specify it if we want to enable the vbias otherwise it will be disabled.
	Power.Setregister([]byte{adc.POWER_vbias_enabled, adc.POWER_intref_enabled})

	//This initializes and sets the Interface register - Note how all options are specified regardless of if they are defaults or not
	Interface := piadcs.NewRegister("Interface", adc.INTERFACE_address)
	Interface.Setregister([]byte{adc.INTERFACE_timeout_disabled, adc.INTERFACE_status_enabled, adc.INTERFACE_crc_checksum})

	//This sets Mode0 to the default values for all settings. The default values for programming registers are defined in the ADS126x datasheet and can be found in the constants file
	Mode0 := piadcs.NewRegister("Mode0", adc.MODE0_address)
	Mode0.Setregister([]byte{adc.MODE0_default})

	Mode1 := piadcs.NewRegister("Mode1", adc.MODE1_address)
	Mode1.Setregister([]byte{adc.MODE1_filter_FIR, adc.MODE1_sbADC_ADC1, adc.MODE1_sbmag_none})

	Mode2 := piadcs.NewRegister("Mode2", adc.MODE2_address)
	Mode2.Setregister([]byte{adc.MODE2_bypass_PGAenabled, adc.MODE2_GAIN_32, adc.MODE2_DR_20})

	Inpmux := piadcs.NewRegister("INPMUX", adc.INPMUX_address)
	Inpmux.Setregister([]byte{adc.INPMUX_muxP_AIN9, adc.INPMUX_muxN_AINCOM})

	//This creates a byte slice which we will use to write the values to the registers specified. It is critical that the values listed are in order and no register is skipped. If there is a register that you don't want to change the value of you still need to speficy it if it falls between two others that you want to change. You can use the built in default value for this as shown in this example with Mode0. See the ADS126x datasheet section 9.5.7 for an explanation of why this is the case
	registerdata := []byte{Power.Setvalue, Interface.Setvalue, Mode0.Setvalue, Mode1.Setvalue, Mode2.Setvalue, Inpmux.Setvalue}

	//This actually writes the data to the register. The starting register needs to be specified - in this case it is the POWER register
	piadcs.WriteToConsecutiveRegisters(spi0, Power.Address, registerdata)

	//This reads the data from registers so we can check that the ADC is working and that we correctly wrote the data to the registers
	incomingregdata := piadcs.ReadFromConsecutiveRegisters(spi0, Power.Address, 6)

	//Since the output of ReadFromConsecutiveRegisters is also a byte slice we can compare it to the slice we sent to make sure that there were no communication errors and the registers were written as intended. The registermatch function is intended for this purpose
	if piadcs.RegisterMatch(incomingregdata, registerdata) {
		fmt.Println("registers match")
	} else {
		fmt.Println("registers don't match - Perhaps there was an error in the SPI communication or the ADC isn't powered up")
	}

	//We can also print the register data to check it that way
	fmt.Println(incomingregdata)

	//Changing register values after they are initially set

	//If conversions are running make sure to disable by bringing start pin low or sending the stop command them before writing to registers - uncomment one of the following to do this
	//startpin.Out(gpio.Low)
	//adc.Stopcommand(spi, startpin)

	//In this example we want to change only MODE2 and INPMUX registers so that we can check the internal temperature of the ADC chip (see thermocouple measurement example for a real use case of this)

	//It is important to reset the setvalue of the register that you want to change by setting it to zero
	Mode2.Setvalue = 0

	Mode2.Setregister([]byte{adc.MODE2_GAIN_1, adc.MODE2_DR_20})

	Inpmux.Setvalue = 0
	//This allows us to read the internal temperature sensor (see the datasheet section 9.3.4)
	Inpmux.Setregister([]byte{adc.INPMUX_muxP_tempSensorP, adc.INPMUX_muxN_tempSensorN})

	registerdata2 := []byte{Mode2.Setvalue, Inpmux.Setvalue}

	//Write new values to the registers
	piadcs.WriteToConsecutiveRegisters(spi0, Mode2.Address, registerdata2)

	incomingregdata2 := piadcs.ReadFromConsecutiveRegisters(spi0, Power.Address, 2)

	if piadcs.RegisterMatch(incomingregdata2, registerdata2) {
		fmt.Println("registers match")
	} else {
		fmt.Println("registers don't match - Perhaps there was an error in the SPI communication or the ADC isn't powered up")
	}

	fmt.Println(incomingregdata2)

}
