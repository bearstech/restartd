// systemd package
// use to link restartd to go-systemd using dbus package
package systemd

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bearstech/restartd/protocol"
	"github.com/coreos/go-systemd/dbus"
	"strings"
)

// struct HandlerSystemd
// implements Handler interface
type HandlerSystemd struct {
	// array containing services (names)
	Services []string
}

// func isUnit()
// verify that requested unit is declared in a config file
func isUnit(u string, s []string) bool {

	for _, v := range s {
		if v == u {
			return true
		}
	}

	return false

}

// func Handle
// implemented by Handler interface
func (h *HandlerSystemd) Handle(m protocol.Message) (r protocol.Response) {

	var code protocol.Response_Codes
	var message string

	// verify if requested unit exists
	ret := isUnit(m.GetService(), h.Services)

	// if unit does not exists
	if ret != true {
		// write appropriate message
		code = protocol.Response_err_missing
		message = fmt.Sprintf("Service %s does not exists",
			m.GetService())
	} else {

		// switch between all supported commands
		switch m.GetCommand() {

		case protocol.Message_status:
			// Get status for requested unit
			ret, err := getStatus(m.GetService())

			// error checking
			// write appropriate messages
			if err != nil {
				message = fmt.Sprintf("Error getting %s service status",
					m.GetService())
				code = protocol.Response_err_status
				log.Error(message)
			} else {
				message = ret
				code = protocol.Response_suc_status
			}

			break

		case protocol.Message_start:
			// start unit
			err := startUnit(m.GetService())

			// error checking
			if err != nil {
				message = fmt.Sprintf("Error starting %s service",
					m.GetService)
				code = protocol.Response_err_start
			} else {
				message = fmt.Sprintf("%s service is started",
					m.GetService())
				code = protocol.Response_suc_start
			}

			break

		case protocol.Message_stop:
			// stop unit
			err := stopUnit(m.GetService())

			// error checking
			if err != nil {
				message = fmt.Sprintf("Error stopping %s service",
					m.GetService)
				code = protocol.Response_err_stop
			} else {
				message = fmt.Sprintf("%s service is stopped",
					m.GetService())
				code = protocol.Response_suc_stop
			}

			break

		case protocol.Message_restart:
			// restart unit
			err := restartUnit(m.GetService())

			// error checking
			if err != nil {
				message = fmt.Sprintf("Error restarting %s service",
					m.GetService)
				code = protocol.Response_err_restart
			} else {
				message = fmt.Sprintf("%s service is restarted",
					m.GetService())
				code = protocol.Response_suc_restart
			}

			break

		case protocol.Message_reload:
			// reload unit
			err := reloadUnit(m.GetService())

			// error checking
			if err != nil {
				message = fmt.Sprintf("Error reloading %s service",
					m.GetService())
				code = protocol.Response_err_restart
			} else {
				message = fmt.Sprintf("%s service is reloaded",
					m.GetService())
				code = protocol.Response_suc_restart
			}

			break

		default:
			code = protocol.Response_err_cmd
			message = fmt.Sprint("Command %s not supported",
				m.GetCommand)
		}
	}

	// send message to restartctl client
	return protocol.Response{
		Code:    &code,
		Message: &message,
	}

}

// getStatus
// Fetch status for a requested unit
func getStatus(unitName string) (string, error) {

	var found bool

	// concatenante uinitName + .service in a serviceName string
	serviceName := createServiceName(unitName)

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
			message := loadedStatusMessage(v)

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
			message := unloadedStatusMessage(v)

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
func createServiceName(unitName string) string {
	return fmt.Sprintf("%s.service", unitName)
}

// loadedStatusMessage
// create a basic status message (used with loaded units)
func loadedStatusMessage(unit dbus.UnitStatus) string {

	return fmt.Sprintf("Name: %s\n\tDescription: %s\n\tLoad: "+
		"%s\n\tActive: %s\n\tState: %s\n", unit.Name, unit.Description,
		unit.LoadState, unit.ActiveState, unit.SubState)

}

// unloadedStatusMessage
// ceate a basic status message (used with unloaded units)
func unloadedStatusMessage(unitFile dbus.UnitFile) string {

	return fmt.Sprintf("Name: %s\n\tStatus: %s\n", unitFile.Path,
		unitFile.Type)

}

// startUnit
func startUnit(unitName string) error {

	// concatenante uinitName + .service in a serviceName string
	serviceName := createServiceName(unitName)

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
func stopUnit(unitName string) error {

	// concatenante uinitName + .service in a serviceName string
	serviceName := createServiceName(unitName)

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

func restartUnit(unitName string) error {

	// concatenante uinitName + .service in a serviceName string
	serviceName := createServiceName(unitName)

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

func reloadUnit(unitName string) error {

	// concatenante uinitName + .service in a serviceName string
	serviceName := createServiceName(unitName)

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
