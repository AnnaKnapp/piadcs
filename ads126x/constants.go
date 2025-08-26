package ads126x

//Opcodes - Commands are used to access the configuration and data registers and also to control the ADC. Many of the ADC commands are stand-alone (that is, single-byte). The register write and register read commands, however, are multibyte, consisting of two opcode bytes plus the register data byte or bytes. (see section 9.5 of the datasheet)
//Write registers (opcode)
const WREG byte = 0x40

//Read registers (opcode)
const RREG byte = 0x20

//ADC1 self offset calibration (opcode)
const SFOCAL1 byte = 0x19

//ADC1 system gain calibration (opcode)
const SYGCAL1 byte = 0x17

//ADC1 system offset calibration (opcode)
const SYOCAL1 byte = 0x16

//Read ADC1 data (opcode)
const RDATA1 byte = 0x12

//Stop ADC1 conversions (opcode)
const STOP1 byte = 0x0A

//Start ADC1 conversions (opcode)
const START1 byte = 0x08

//Reset the ADC (opcode)
const RESET byte = 0x06

//First register is the Device identification register which is read only. Reading this register gives device version and revision.
const (
	ID_address byte = 0x00
)

//Second register is the Power register.
const (
	POWER_address byte = 0x01
	POWER_default byte = 0x11

	//reset indicator (Indicates ADC reset has occurred. Clear this bit to detect the next device reset.)
	POWER_reset_no  byte = 0b00000000
	POWER_reset_yes byte = 0b00010000

	//level shift voltage enable (Enables the internal level shift voltage to the AINCOM pin. VBIAS = (VAVDD + VAVSS)/2)
	POWER_vbias_disabled byte = 0b00000000
	POWER_vbias_enabled  byte = 0b00000010

	//Internal reference ebable (Enables the 2.5 V internal voltage reference. Note the IDAC and temperature sensor require the internal voltage reference.) This should always be enabled unless you are using an external reference
	POWER_intref_disabled byte = 0b00000000
	POWER_intref_enabled  byte = 0b00000001
)

//Third register is the interface register. It sets serial interface timeout, status byte, and checksum byte
const (
	INTERFACE_address byte = 0x02
	INTERFACE_default byte = 0x05

	//Serial Interface Time-Out Enable (Enables the serial interface automatic time-out mode)
	INTERFACE_timeout_disabled byte = 0b00000000
	INTERFACE_timeout_enabled  byte = 0b00001000

	//Status Byte Enable (Enables the inclusion of the status byte during conversion data read-back)
	INTERFACE_status_disabled byte = 0b00000000
	INTERFACE_status_enabled  byte = 0b00000100

	//Checksum Byte Enable (Enables the inclusion of the checksum byte during conversion data read-back)
	INTERFACE_crc_disabled byte = 0b00000000 //checksum disabled
	INTERFACE_crc_checksum byte = 0b00000001 //checksum mode
	INTERFACE_crc_crc      byte = 0b00000010 //Cyclic redundancy check
)

//Fourth register is the Mode0 - it sets reference mux polarity, ADC conversion mode, chop mode and conversion delay.
const (
	MODE0_address byte = 0x03
	MODE0_default byte = 0x00

	//Reference Mux Polarity Reversal (Reverses the ADC1 reference multiplexer output polarity)
	MODE0_refrev_normalpolarity  byte = 0b00000000
	MODE0_refrev_reversepolarity byte = 0b10000000

	//ADC Conversion Run Mode (Selects the ADC conversion (run) mode)
	MODE0_runmode_continuous byte = 0b00000000 // continuous conversion
	MODE0_runmode_pulse      byte = 0b01000000 //one shot conversion

	//Chop Mode Enable (Enables the ADC chop and IDAC rotation options)
	MODE0_chop_disabled              byte = 0b00000000 //Input chop and IDAC rotation disabled
	MODE0_chop_chopenabled           byte = 0b00010000 //input chop enabled
	MODE0_chop_IDACrotation          byte = 0b00100000 //IDAC rotation enabled
	MODE0_chop_chop_and_IDACrotation byte = 0b00110000

	//Conversion Delay (Provides additional delay from conversion start to the beginning of the actual conversion)
	MODE0_delay_none  byte = 0b00000000
	MODE0_delay_8µ7s  byte = 0b00000001 //8.7 µs
	MODE0_delay_17µs  byte = 0b00000010 //17 µs
	MODE0_delay_35µs  byte = 0b00000011 //35 µs
	MODE0_delay_69µs  byte = 0b00000100 //69 µs
	MODE0_delay_139µs byte = 0b00000101 //139 µs
	MODE0_delay_278µs byte = 0b00000110 //278 µs
	MODE0_delay_555µs byte = 0b00000111 //555 µs
	MODE0_delay_1m1s  byte = 0b00001000 //1.1ms
	MODE0_delay_2m2s  byte = 0b00001001 //2.2ms
	MODE0_delay_4m4s  byte = 0b00001010 //4.4ms
	MODE0_delay_8m8s  byte = 0b00001011 //8.8ms
)

//The fifth register is Mode 1 - It sets the digital filter and Sensor Bias
const (
	MODE1_address byte = 0x04
	MODE1_default byte = 0x80

	// Digital Filter (Configures the ADC digital filter)
	MODE1_filter_sinc1 byte = 0b00000000
	MODE1_filter_sinc2 byte = 0b00100000
	MODE1_filter_sinc3 byte = 0b01000000
	MODE1_filter_sinc4 byte = 0b01100000
	MODE1_filter_FIR   byte = 0b10000000 //default

	//Sensor Bias ADC Connection (Selects the ADC to connect the sensor bias)
	MODE1_sbADC_ADC1 byte = 0b00000000 //Sensor bias connected to ADC1 mux out (default)
	MODE1_sbADC_ADC2 byte = 0b00010000 //Sensor bias connected to ADC2 mux out

	//Sensor Bias Polarity Selects the sensor bias for pull-up or pull-down
	MODE1_sbpol_pullUp   byte = 0b00000000 //Sensor bias pull-up mode (AINP pulled high, AINN pulled low) (default)
	MODE1_sbpol_pullDown byte = 0b00001000 //Sensor bias pull-down mode (AINP pulled low, AINN pulled high)

	// Sensor Bias Magnitude (Selects the sensor bias current magnitude or the bias resistor)
	MODE1_sbmag_none  byte = 0b00000000 //No sensor bias current or resistor (default)
	MODE1_sbmag_500nA byte = 0b00000001 //0.5-µA sensor bias current
	MODE1_sbmag_2µA   byte = 0b00000010 //2-µA sensor bias current
	MODE1_sbmag_10µA  byte = 0b00000011 //10-µA sensor bias current
	MODE1_sbmag_50µA  byte = 0b00000100 //50-µA sensor bias current
	MODE1_sbmag_200µA byte = 0b00000101 //200-µA sensor bias current
	MODE1_sbmag_10MΩ  byte = 0b00000110 //10-MΩ resistor
)

// The sixth register is Mode 2. It sets the PGA (programable gain amplifier) and the data rate
const (
	MODE2_address byte = 0x05
	MODE2_default byte = 0x04

	//PGA Bypass Mode Selects PGA bypass mode
	MODE2_bypass_PGAenabled  byte = 0b00000000 //default
	MODE2_bypass_PGAdisabled byte = 0b10000000

	//PGA Gain - selects PGA gain
	MODE2_GAIN_1  byte = 0b00000000 //1 V/V default
	MODE2_GAIN_2  byte = 0b00010000 //1 V/V
	MODE2_GAIN_4  byte = 0b00100000 //1 V/V
	MODE2_GAIN_8  byte = 0b00110000 //1 V/V
	MODE2_GAIN_16 byte = 0b01000000 //1 V/V
	MODE2_GAIN_32 byte = 0b01010000 //1 V/V

	//Data rate (Selects the ADC data rate. In FIR filter mode, the available data rates are limited to 2.5, 5, 10 and 20 SPS.)
	MODE2_DR_2_5   byte = 0b00000000 //2.5 samples per second
	MODE2_DR_5     byte = 0b00000001
	MODE2_DR_10    byte = 0b00000010
	MODE2_DR_16_6  byte = 0b00000011 //16.6~ samples per second
	MODE2_DR_20    byte = 0b00000100
	MODE2_DR_50    byte = 0b00000101
	MODE2_DR_60    byte = 0b00000110
	MODE2_DR_100   byte = 0b00000111
	MODE2_DR_400   byte = 0b00001000
	MODE2_DR_1200  byte = 0b00001001
	MODE2_DR_2400  byte = 0b00001010
	MODE2_DR_4800  byte = 0b00001011
	MODE2_DR_7200  byte = 0b00001100
	MODE2_DR_14400 byte = 0b00001101
	MODE2_DR_19200 byte = 0b00001110
	MODE2_DR_38400 byte = 0b00001111
)

//The input multiplexer register sets the input channels
const (
	INPMUX_address byte = 0x06
	INPMUX_default byte = 0x01

	//Positive Input Multiplexer (Selects the positive input multiplexer.)
	INPMUX_muxP_AIN0           byte = 0b00000000 //positive input is AIN0 (default)
	INPMUX_muxP_AIN1           byte = 0b00010000 //positive input is AIN1
	INPMUX_muxP_AIN2           byte = 0b00100000 //positive input is AIN1
	INPMUX_muxP_AIN3           byte = 0b00110000 //positive input is AIN3
	INPMUX_muxP_AIN4           byte = 0b01000000 //positive input is AIN4
	INPMUX_muxP_AIN5           byte = 0b01010000 //positive input is AIN5
	INPMUX_muxP_AIN6           byte = 0b01100000 //positive input is AIN6
	INPMUX_muxP_AIN7           byte = 0b01110000 //positive input is AIN7
	INPMUX_muxP_AIN8           byte = 0b10000000 //positive input is AIN8
	INPMUX_muxP_AIN9           byte = 0b10010000 //positive input is AIN9
	INPMUX_muxP_AINCOM         byte = 0b10100000 //positive input is AINCOM
	INPMUX_muxP_tempSensorP    byte = 0b10110000 //Temperature sensor monitor positive
	INPMUX_muxP_analogSupplyP  byte = 0b11000000 //Analog power supply monitor positive
	INPMUX_muxP_digitalSupplyP byte = 0b11010000 //Digital power supply monitor positive
	INPMUX_muxP_TDACP          byte = 0b11100000 //TDAC test signal positive
	INPMUX_muxP_float          byte = 0b11110000 //Float (open connection)

	INPMUX_muxN_AIN0           byte = 0b00000000 //negative input is AIN0
	INPMUX_muxN_AIN1           byte = 0b00000001 //negative input is AIN1 (default)
	INPMUX_muxN_AIN2           byte = 0b00000010 //negative input is AIN2
	INPMUX_muxN_AIN3           byte = 0b00000011 //negative input is AIN3
	INPMUX_muxN_AIN4           byte = 0b00000100 //negative input is AIN4
	INPMUX_muxN_AIN5           byte = 0b00000101 //negative input is AIN5
	INPMUX_muxN_AIN6           byte = 0b00000110 //negative input is AIN6
	INPMUX_muxN_AIN7           byte = 0b00000111 //negative input is AIN7
	INPMUX_muxN_AIN8           byte = 0b00001000 //negative input is AIN8
	INPMUX_muxN_AIN9           byte = 0b00001001 //negative input is AIN9
	INPMUX_muxN_AINCOM         byte = 0b00001010 //negative input is AINCOM
	INPMUX_muxN_tempSensorN    byte = 0b00001011 //Temperature sensor monitor negative
	INPMUX_muxN_analogSupplyN  byte = 0b00001100 //Analog power supply monitor negative
	INPMUX_muxN_digitalSupplyN byte = 0b00001101 //Digital power supply monitor negative
	INPMUX_muxN_TDACP          byte = 0b00001110 //TDAC test signal negative
	INPMUX_muxN_float          byte = 0b00001111 //Float (open connection)
)

//Offset calibration register - Use calibration to correct internal ADC errors or overall system errors. the value of the offset calibration register is subtracted from the filter output and then multiplied by the full-scale register value divided by 400000h. The data are then clipped to a 32-bit value to provide the final output. (see section 9.4.9 of the datasheet for more information)
const (
	OFCAL0_address byte = 0x07
	OFCAL1_address byte = 0x08
	OFCAL2_address byte = 0x09
)

//Full-scale calibration register - Use calibration to correct internal ADC errors or overall system errors. the value of the offset calibration register is subtracted from the filter output and then multiplied by the full-scale register value divided by 400000h. The data are then clipped to a 32-bit value to provide the final output. (see section 9.4.9 of the datasheet for more information)
const (
	FSCAL0_address byte = 0x0A
	FSCAL1_address byte = 0x0B
	FSCAL2_address byte = 0x0C
)

// IDAC Multiplexer register - selects analog input pins to connect to IDAC1 and IDAC2 (IDAC refers to current sources)
const (
	IDACMUX_address byte = 0x0D
	IDACMUX_default byte = 0xBB

	//IDAC2 Output Multiplexer Selects the analog input pin to connect IDAC2
	IDACMUX_mux2_AIN0   byte = 0b00000000 //IDAC2 output is AIN0
	IDACMUX_mux2_AIN1   byte = 0b00010000 //IDAC2 output is AIN1
	IDACMUX_mux2_AIN2   byte = 0b00100000 //IDAC2 output is AIN1
	IDACMUX_mux2_AIN3   byte = 0b00110000 //IDAC2 output is AIN3
	IDACMUX_mux2_AIN4   byte = 0b01000000 //IDAC2 output is AIN4
	IDACMUX_mux2_AIN5   byte = 0b01010000 //IDAC2 output is AIN5
	IDACMUX_mux2_AIN6   byte = 0b01100000 //IDAC2 output is AIN6
	IDACMUX_mux2_AIN7   byte = 0b01110000 //IDAC2 output is AIN7
	IDACMUX_mux2_AIN8   byte = 0b10000000 //IDAC2 output is AIN8
	IDACMUX_mux2_AIN9   byte = 0b10010000 //IDAC2 output is AIN9
	IDACMUX_mux2_AINCOM byte = 0b10100000 //IDAC2 output is AINCOM
	IDACMUX_mux2_none   byte = 0b10110000 //No connection (default)

	//IDAC1 Output Multiplexer Selects the analog input pin to connect IDAC1
	IDACMUX_mux1_AIN0   byte = 0b00000000 //IDAC1 output is AIN0
	IDACMUX_mux1_AIN1   byte = 0b00000001 //IDAC1 output is AIN1
	IDACMUX_mux1_AIN2   byte = 0b00000010 //IDAC1 output is AIN2
	IDACMUX_mux1_AIN3   byte = 0b00000011 //IDAC1 output is AIN3
	IDACMUX_mux1_AIN4   byte = 0b00000100 //IDAC1 output is AIN4
	IDACMUX_mux1_AIN5   byte = 0b00000101 //IDAC1 output is AIN5
	IDACMUX_mux1_AIN6   byte = 0b00000110 //IDAC1 output is AIN6
	IDACMUX_mux1_AIN7   byte = 0b00000111 //IDAC1 output is AIN7
	IDACMUX_mux1_AIN8   byte = 0b00001000 //IDAC1 output is AIN8
	IDACMUX_mux1_AIN9   byte = 0b00001001 //IDAC1 output is AIN9
	IDACMUX_mux1_AINCOM byte = 0b00001010 //IDAC1 output is AINCOM
	IDACMUX_mux1_none   byte = 0b00001011 //No connection (default)

)

//IDAC magnitude register (IDAC refers to current sources)
const (
	IDACMAG_address byte = 0x0E
	IDACMAG_default byte = 0x00
	//IDAC2 Output Multiplexer Selects the analog input pin to connect IDAC2
	IDACMAG_mag2_off  byte = 0b00000000 //IDAC2 is off (default)
	IDACMAG_mag2_50   byte = 0b00010000 //IDAC2 output is 50 µA
	IDACMAG_mag2_100  byte = 0b00100000 //IDAC2 output is 100 µA
	IDACMAG_mag2_250  byte = 0b00110000 //IDAC2 output is 250 µA
	IDACMAG_mag2_500  byte = 0b01000000
	IDACMAG_mag2_750  byte = 0b01010000
	IDACMAG_mag2_1000 byte = 0b01100000
	IDACMAG_mag2_1500 byte = 0b01110000
	IDACMAG_mag2_2000 byte = 0b10000000
	IDACMAG_mag2_3000 byte = 0b10010000

	//IDAC1 Output Multiplexer Selects the analog input pin to connect IDAC1
	IDACMAG_mag1_off  byte = 0b00000000 //IDAC1 is off
	IDACMAG_mag1_50   byte = 0b00000001 //IDAC1 output is 50 µA
	IDACMAG_mag1_100  byte = 0b00000010 //IDAC1 output is 100 µA
	IDACMAG_mag1_250  byte = 0b00000011 //IDAC1 output is 250 µA
	IDACMAG_mag1_500  byte = 0b00000100
	IDACMAG_mag1_750  byte = 0b00000101
	IDACMAG_mag1_1000 byte = 0b00000110
	IDACMAG_mag1_1500 byte = 0b00000111
	IDACMAG_mag1_2000 byte = 0b00001000
	IDACMAG_mag1_2500 byte = 0b00001001
	IDACMAG_mag1_3000 byte = 0b00001010
)

//Reference multiplexer register - selects reference inputs
const (
	REFMUX_address byte = 0x0F
	REFMUX_default byte = 0x00
	//Reference Positive Input (Selects the positive reference input)
	REFMUX_rmuxP_internalRef   byte = 0b00000000 //Internal 2.5 V reference - P (default)
	REFMUX_rmuxP_AIN0          byte = 0b00001000 //External AIN0
	REFMUX_rmuxP_AIN2          byte = 0b00010000
	REFMUX_rmuxP_AIn4          byte = 0b00011000
	REFMUX_rmuxP_internalVavdd byte = 0b00100000 // Internal analog supply (VAVDD )
	//Reference Negative Input (Selects the negative reference input)
	REFMUX_rmuxN_internalRef   byte = 0b00000000 //Internal 2.5 V reference - N (default)
	REFMUX_rmuxN_AIN1          byte = 0b00000001 //External AIN1
	REFMUX_rmuxN_AIN3          byte = 0b00000010
	REFMUX_rmuxN_AIN5          byte = 0b00000011
	REFMUX_rmuxP_internalVavss byte = 0b00000100 //Internal analog supply (VAVSS)
)

//TDACP control register - Test DAC (positive)
const (
	TDACP_address byte = 0x10
	TDACP_default byte = 0x00
	// TDACP Output Connection (Connects TDACP output to pin AIN6)
	TDACP_outP_none byte = 0b00000000
	TDACP_outP_AIN6 byte = 0b10000000
	// MAGP Output Magnitude Select the TDACP output magnitude. (The TDAC output voltages are ideal and are with respect to VAVSS)
	TDACP_magP_4_5       byte = 0b00001001 // 4.5 V
	TDACP_magP_3_5       byte = 0b00001000 // 3.5 V
	TDACP_magP_3         byte = 0b00000111 // 3 V
	TDACP_magP_2_75      byte = 0b00000110 // 2.75 V
	TDACP_magP_2_625     byte = 0b00000101 // 2.625 v
	TDACP_magP_2_5626    byte = 0b00000100 // 2.5625 V
	TDACP_magP_2_53125   byte = 0b00000011 // 2.53125 V
	TDACP_magP_2_515625  byte = 0b00000010 // 2.515625 V
	TDACP_magP_2_5078125 byte = 0b00000001 // 2.5078125 V
	TDACP_magP_2_5       byte = 0b00000000 // 2.5 V
	TDACP_magP_2_4921875 byte = 0b00010001 // 2.4921875 V
	TDACP_magP_2_484375  byte = 0b00010010 // 2.484375 V
	TDACP_magP_2_46875   byte = 0b00010011 // 2.46875 V
	TDACP_magP_2_4375    byte = 0b00010100 // 2.4375 V
	TDACP_magP_2_375     byte = 0b00010101 // 2.375 V
	TDACP_magP_2_25      byte = 0b00010110 // 2.25 V
	TDACP_magP_2         byte = 0b00010111 // 2 V
	TDACP_magP_1_5       byte = 0b00011000 // 1.5 V
	TDACP_magP_0_5       byte = 0b00011001 // 0.5 V
)

//TDACN control register - Test DAC (negative)
const (
	TDACN_address byte = 0x10
	TDACN_default byte = 0x00
	// TDACN Output Connection (Connects TDACN output to pin AIN6)
	TDACN_outN_none byte = 0b00000000
	TDACN_outN_AIN7 byte = 0b10000000
	// MAGN Output Magnitude Select the TDACN output magnitude. (The TDAC output voltages are ideal and are with respect to VAVSS)
	TDACN_magN_4_5       byte = 0b00001001 // 4.5 V
	TDACN_magN_3_5       byte = 0b00001000 // 3.5 V
	TDACN_magN_3         byte = 0b00000111 // 3 V
	TDACN_magN_2_75      byte = 0b00000110 // 2.75 V
	TDACN_magN_2_625     byte = 0b00000101 // 2.625 v
	TDACN_magN_2_5626    byte = 0b00000100 // 2.5625 V
	TDACN_magN_2_53125   byte = 0b00000011 // 2.53125 V
	TDACN_magN_2_515625  byte = 0b00000010 // 2.515625 V
	TDACN_magN_2_5078125 byte = 0b00000001 // 2.5078125 V
	TDACN_magN_2_5       byte = 0b00000000 // 2.5 V
	TDACN_magN_2_4921875 byte = 0b00010001 // 2.4921875 V
	TDACN_magN_2_484375  byte = 0b00010010 // 2.484375 V
	TDACN_magN_2_46875   byte = 0b00010011 // 2.46875 V
	TDACN_magN_2_4375    byte = 0b00010100 // 2.4375 V
	TDACN_magN_2_375     byte = 0b00010101 // 2.375 V
	TDACN_magN_2_25      byte = 0b00010110 // 2.25 V
	TDACN_magN_2         byte = 0b00010111 // 2 V
	TDACN_magN_1_5       byte = 0b00011000 // 1.5 V
	TDACN_magN_0_5       byte = 0b00011001 // 0.5 V
)
