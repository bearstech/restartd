//Package systemd use to link restartd to go-systemd using dbus package
package systemd

import (
	"errors"
	"fmt"
	"github.com/coreos/go-systemd/dbus"
	"strings"
)

const DONE = "done"

// IsUnit verify that requested unit is declared in a config file
func IsUnit(u string, s []string) bool {

	for _, v := range s {
		if v == u {
			return true
		}
	}

	return false

}

// GetStatus fetch status for a requested unit
func GetStatus(unitName string) (string, error) {

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
	unitsStatus, err := conn.ListUnits()
	if err != nil {
		message := fmt.Sprintf("Error getting %s service status",
			serviceName)
		return message, err
	}

	// for each units, find the requested one
	for _, v := range unitsStatus {
		if strings.Contains(v.Name, serviceName) {
			// create basic response message
			message := LoadedStatusMessage(v)

			// return message to restartctl client
			return message, err
		}
	}

	// go deeper
	unitsFiles, err := conn.ListUnitFiles()
	if err != nil {
		message := fmt.Sprintf("Error getting %s service file",
			serviceName)
		return message, err
	}

	// search for unit
	for _, v := range unitsFiles {
		if strings.Contains(v.Path, serviceName) {
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

// CreateServiceName creates service name
func CreateServiceName(unitName string) string {
	return fmt.Sprintf("%s.service", unitName)
}

// LoadedStatusMessage creates a basic status message (used with loaded units)
func LoadedStatusMessage(unit dbus.UnitStatus) string {

	return fmt.Sprintf("Name: %s\n\tDescription: %s\n\tLoad: "+
		"%s\n\tActive: %s\n\tState: %s\n", unit.Name, unit.Description,
		unit.LoadState, unit.ActiveState, unit.SubState)

}

// UnloadedStatusMessage creates a basic status message (used with unloaded units)
func UnloadedStatusMessage(unitFile dbus.UnitFile) string {

	return fmt.Sprintf("Name: %s\n\tStatus: %s\n", unitFile.Path,
		unitFile.Type)

}

func dbusConn(unitName string, closure func(serviceName string, conn *dbus.Conn, ch chan string) error) error {
	// concatenante uinitName + .service in a serviceName string
	serviceName := CreateServiceName(unitName)

	ch := make(chan string)
	// create systemd-dbus conn
	conn, err := dbus.New()
	// ensure conn is closed
	if err != nil {
		return err
	}
	defer conn.Close()

	err = closure(serviceName, conn, ch)
	if err != nil {
		return err
	}
	// wait for done signal
	msg := <-ch
	if msg != DONE {
		return errors.New("Systemd error :" + msg)
	}

	return nil
}

// StartUnit starts unit
func StartUnit(unitName string) error {
	return dbusConn(unitName, func(serviceName string, conn *dbus.Conn, ch chan string) error {
		_, err := conn.StartUnit(serviceName, "replace", ch)
		return err
	})
}

// StopUnit stops unit
func StopUnit(unitName string) error {
	return dbusConn(unitName, func(serviceName string, conn *dbus.Conn, ch chan string) error {
		_, err := conn.StopUnit(serviceName, "replace", ch)
		return err
	})
}

// RestartUnit restarts unit
func RestartUnit(unitName string) error {
	return dbusConn(unitName, func(serviceName string, conn *dbus.Conn, ch chan string) error {
		_, err := conn.RestartUnit(serviceName, "replace", ch)
		return err
	})
}

// ReloadUnit reloads unit
func ReloadUnit(unitName string) error {
	return dbusConn(unitName, func(serviceName string, conn *dbus.Conn, ch chan string) error {
		_, err := conn.ReloadUnit(serviceName, "replace", ch)
		return err
	})
}
