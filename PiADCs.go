package piadcs

import (
	"log"

	adc "github.com/AnnaKnapp/piadcs/ads126x"

	spi "periph.io/x/periph/conn/spi"
)

type Register struct {
	Name    string
	Address byte

	Defaultvalue byte

	//use the function Setregister with values from the constants file or you can also set this manually by looking at the datasheet.
	Setvalue byte
}

//This function sets the register data byte according to the settings specified. The lists of all possible settings for registers are listed in the constants file for that device
func (reg *Register) Setregister(settings []byte) {
	for i := range settings {
		reg.Setvalue = settings[i] | reg.Setvalue
	}
}

//Use this to create a register to then set using the setregister function
func NewRegister(name string, address byte) Register {
	return Register{
		Name:    name,
		Address: address,
	}

}

//This writes data to consecutive registers. You need to specify the starting register and the byte slice of data to write. It will go down the slice and write one byte to each consecutive register starting from the one specified. Please see the datasheet for more information. The WREG opcode is used here and uses the same programming on multiple different TI ADCs
func WriteToConsecutiveRegisters(connection spi.Conn, startingreg byte, datatowrite []byte) {
	towrite := []byte{adc.WREG | startingreg, byte(len(datatowrite) - 1)}
	// for i := range datatowrite {
	// towrite = append(towrite, datatowrite[i])
	// }
	// ... destructures a slice
	towrite = append(towrite, datatowrite...)
	toread := make([]byte, len(towrite))
	if err := connection.Tx(towrite, toread); err != nil {
		log.Fatal(err)
	}
}

//Use this function to check what data is stored at what registers. The starting register and the number of registers to read must be specified.
func ReadFromConsecutiveRegisters(connection spi.Conn, startingreg byte, numbertoread byte) []byte {
	asktoread := []byte{adc.RREG | startingreg, numbertoread - 1}
	blank1 := make([]byte, 2)
	registerdata := make([]byte, int(numbertoread))
	blank2 := make([]byte, int(numbertoread))
	if err := connection.Tx(asktoread, blank1); err != nil {
		log.Fatal(err)
	}
	if err := connection.Tx(blank2, registerdata); err != nil {
		log.Fatal(err)
	}
	return registerdata

}
