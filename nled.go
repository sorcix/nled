package main

import (
	"fmt"
	"os"
)

// Beaglebone pin configuration
const (
	pinRed = "60" // RED: P9_12; GPIO1[28]; 60
	pinGreen = "50" // GREEN: P9_14; GPIO1[18]; 50
	pinBlue = "51" // BLUE: P9_16; GPIO1[19]; 51
)

// Names used as command arguments
const (
	nameRed = "red"
	nameGreen = "green"
	nameBlue = "blue"
)

// SysFS paths and permissions
const (
	path = "/sys/class/gpio"
	pathExport = path+"/export" // Export GPIO pin
	pathLed = path+"/gpio%s"
	pathValue = pathLed+"/value" // Set GPIO pin
	pathDirection = pathLed+"/direction"
	pathActive = pathLed+"/active_low"
	permissions = 0666
)

// Commands
const (
	cmdOn = "on"
	cmdOff = "off"
	cmdInit = "init"
)

// Errors
const (
	errorStatus = "Unable to set GPIO pin %s status to %s!\n"
	errorExport = "Unable to export GPIO pin %s!\n"
)

// Status
const (
	on = "1"
	off = "0"
	high = "high"
	low = "low"
)

func usage() {
	fmt.Printf("Usage:\n\n")
	fmt.Printf(" %s {red|green|blue} {on|off}\n", os.Args[0])
	fmt.Printf(" %s init\n\n", os.Args[0])
}

func write(path string, value []byte) error {

	// Open path for writing
	file, err := os.OpenFile(path, os.O_WRONLY, permissions)
	if err != nil {
		return err
	}

	// Write given value
	file.Write(value)
	file.Close()

	return nil
}

func set(led, status string) bool {

	// Translate status
	switch status {
		case cmdOn:
			status = on
		case cmdOff:
			status = off
		default:
			return false
	}

	// Translate led name to pin number
	switch led {
		case nameRed:
			led = pinRed
		case nameGreen:
			led = pinGreen
		case nameBlue:
			led = pinBlue
		default:
			return false
	}

	// Get path to GPIO file
	file = fmt.Sprintf(pathValue, led)

	// Attempt to write led status to GPIO.
	if err := write(file,[]byte(status)); err != nil {
		die(errorStatus,led,status)
	}

	return true;

}

func exists(path string) bool {

	// Request file info
    info, err := os.Stat(path)

	// Return true only if this path denotes a directory.
	return err != nil && info != nil && info.IsDir()

}

func die(error string, a ...string) {
	fmt.Printf(error, a)
	os.Exit(1)
}

func export() {

	leds := [3]string{pinRed,pinGreen,pinBlue}

	// We need to initialize every led separately.
	for _, led := range leds {

		if exists(fmt.Sprintf(pathLed,led)) {

			// GPIO pin is already initialized
			continue
		}

		if err := write(pathExport,[]byte(led)); err != nil {
			die(errorExport, led)
		}

	}
}

func main() {

	argc := len(os.Args)

	if argc == 3 {

		// Two arguments means switching GPIO value
		if !set(os.Args[1],os.Args[2]) {
			usage()
		}

	} else if argc == 2 && os.Args[1] == cmdInit {

		// Initialize GPIO
		export()

	} else {

		// Show usage
		usage()

	}

}
