// systemd package
// use to link restartd to go-systemd using dbus package
package systemd

import (
	"fmt"
	"github.com/coreos/go-systemd/dbus"
	"strings"
)

// func isUnit()
// verify that requested unit is declared in a config file
func IsUnit(u string, s []string) bool {

	for _, v := range s {
		if v == u {
			return true
		}
	}

	return false

}

// getStatus
// Fetch status for a requested unit
func GetStatus(unitName string) (string, error) {

	var found bool

	// concatenante uinitName + .service in a serviceName string
	serviceName := CreateServiceName(unitName)

	// create systemd-dbus conn
	conn, err := dbus.New()
	// ensure conn is closed
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	// call to systemd-dbus
	// step 1, get all **loaded** units
	UnitsStatus, err := conn.ListUnits()
	if err != nil {
		message := fmt.Sprintf("Error getting %s service status",
			serviceName)
		return message, err
	}

	// for each units, find the requested one
	for _, v := range UnitsStatus {
		if strings.Contains(v.Name, serviceName) == true {
			// create basic response message
			message := LoadedStatusMessage(v)

			// return message to restartctl client
			return message, err
		}
	}

	// go deeper
	UnitsFiles, err := conn.ListUnitFiles()
	if err != nil {
		message := fmt.Sprintf("Error getting %s service file",
			serviceName)
		return message, err
	}

	// search for unit
	for _, v := range UnitsFiles {

		found = strings.Contains(v.Path, serviceName)

		if found == true {
			// create basic response message
			message := UnloadedStatusMessage(v)

			// return message to stopctl client
			return message, err
		}
	}

	// if unit not found
	message := fmt.Sprintf("Error %s service or unit not found", serviceName)

	// return an error
	return message, err

}

// createServiceName
func CreateServiceName(unitName string) string {
	return fmt.Sprintf("%s.service", unitName)
}

// loadedStatusMessage
// create a basic status message (used with loaded units)
func LoadedStatusMessage(unit dbus.UnitStatus) string {

	return fmt.Sprintf("Name: %s\n\tDescription: %s\n\tLoad: "+
		"%s\n\tActive: %s\n\tState: %s\n", unit.Name, unit.Description,
		unit.LoadState, unit.ActiveState, unit.SubState)

}

// unloadedStatusMessage
// ceate a basic status message (used with unloaded units)
func UnloadedStatusMessage(unitFile dbus.UnitFile) string {

	return fmt.Sprintf("Name: %s\n\tStatus: %s\n", unitFile.Path,
		unitFile.Type)

}

// startUnit
func StartUnit(unitName string) error {

	// concatenante uinitName + .service in a serviceName string
	serviceName := CreateServiceName(unitName)

	ch := make(chan string)

	// create systemd-dbus conn
	conn, err := dbus.New()
	// ensure conn is closed
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	_, err = conn.StartUnit(serviceName, "replace", ch)
	if err != nil {
		return err
	}

	// wait for done signal
	<-ch

	return err

}

// stopUnit
func StopUnit(unitName string) error {

	// concatenante uinitName + .service in a serviceName string
	serviceName := CreateServiceName(unitName)

	ch := make(chan string)

	// create systemd-dbus conn
	conn, err := dbus.New()
	// ensure conn is closed
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	_, err = conn.StopUnit(serviceName, "replace", ch)
	if err != nil {
		return err
	}

	// wait for done signal
	<-ch

	return err

}

func RestartUnit(unitName string) error {

	// concatenante uinitName + .service in a serviceName string
	serviceName := CreateServiceName(unitName)

	ch := make(chan string)

	// create systemd-dbus conn
	conn, err := dbus.New()
	// ensure conn is closed
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	_, err = conn.RestartUnit(serviceName, "replace", ch)
	if err != nil {
		return err
	}

	// wait for done signal
	<-ch

	return err

}

func ReloadUnit(unitName string) error {

	// concatenante uinitName + .service in a serviceName string
	serviceName := CreateServiceName(unitName)

	ch := make(chan string)

	// create systemd-dbus conn
	conn, err := dbus.New()
	// ensure conn is closed
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	_, err = conn.ReloadUnit(serviceName, "replace", ch)
	if err != nil {
		return err
	}

	// wait for done signal
	<-ch

	return err

}
