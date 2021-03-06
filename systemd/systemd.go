//Package systemd use to link restartd to go-systemd using dbus package
package systemd

import (
	"errors"
	"fmt"
	"github.com/coreos/go-systemd/dbus"
	"strings"
	"time"
)

const (
	DONE = "done"
)

type SubState string

const (
	SUBSTATE_ACTIVE    SubState = "active"
	SUBSTATE_DEAD      SubState = "dead"
	SUBSTATE_EXITED    SubState = "exited"
	SUBSTATE_FAILED    SubState = "failed"
	SUBSTATE_LISTENING SubState = "listening"
	SUBSTATE_MOUNTED   SubState = "mounted"
	SUBSTATE_PLUGGED   SubState = "plugged"
	SUBSTATE_RUNNING   SubState = "running"
	SUBSTATE_WAITING   SubState = "waiting"
)

type ActiveState string

const (
	ACTIVESTATE_ACTIVE   ActiveState = "active"
	ACTIVESTATE_FAILED   ActiveState = "failed"
	ACTIVESTATE_INACTIVE ActiveState = "inactive"
)

type LoadState string

const (
	LOADSTATE_LOADED    LoadState = "loaded"
	LOADSTATE_MASKED    LoadState = "masked"
	LOADSTATE_NOT_FOUND LoadState = "not-found"
)

type State struct {
	Active ActiveState
	Load   LoadState
	Sub    SubState
	//
	Since time.Time
}

type Unit struct {
	Name        string
	Path        string
	Type        string
	Description string
	State       *State
}

// Contains verify that requested unit is declared in a config file
func Contains(needle string, haystack []string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}

// GetStatusWithPrefix return status of unit matching the prefix
func GetStatusWithPrefix(prefix string) ([]*Unit, error) {
	statuz := []*Unit{}
	// create systemd-dbus conn
	conn, err := dbus.New()
	// ensure conn is closed
	defer conn.Close()
	if err != nil {
		return statuz, err
	}

	unitsStatus, err := conn.ListUnits()
	if err != nil {
		return statuz, err
	}
	for _, us := range unitsStatus {
		if strings.HasPrefix(us.Name, prefix) {
			statuz = append(statuz, LoadedStatusMessage(us))
		}
	}
	return statuz, nil
}

// GetStatus fetch status for a requested unit
func GetStatus(unitName string) (*Unit, error) {

	// concatenante uinitName + .service in a serviceName string
	serviceName := CreateServiceName(unitName)

	// create systemd-dbus conn
	conn, err := dbus.New()
	// ensure conn is closed
	defer conn.Close()
	if err != nil {
		return nil, err
	}

	// call to systemd-dbus
	// step 1, get all **loaded** units
	unitsStatus, err := conn.ListUnits()
	if err != nil {
		return nil, err
	}

	// for each units, find the requested one
	for _, v := range unitsStatus {
		if strings.Contains(v.Name, serviceName) {
			// create basic response message
			message := LoadedStatusMessage(v)

			// return message to restartctl client
			return message, nil
		}
	}

	// go deeper
	unitsFiles, err := conn.ListUnitFiles()
	if err != nil {
		return nil, err
	}

	// search for unit
	for _, v := range unitsFiles {
		if strings.Contains(v.Path, serviceName) {
			// create basic response message
			message := UnloadedStatusMessage(v)
			// return message to stopctl client
			return message, nil
		}
	}

	// if unit not found
	return nil, fmt.Errorf("Error %s service or unit not found", serviceName)
}

// CreateServiceName creates service name
func CreateServiceName(unitName string) string {
	return fmt.Sprintf("%s.service", unitName)
}

// LoadedStatusMessage creates a basic status message (used with loaded units)
func LoadedStatusMessage(unit dbus.UnitStatus) *Unit {
	return &Unit{
		Name:        unit.Name,
		Description: unit.Description,
		State: &State{
			Load:   LoadState(unit.LoadState),
			Active: ActiveState(unit.ActiveState),
			Sub:    SubState(unit.SubState),
		},
	}
}

// UnloadedStatusMessage creates a basic status message (used with unloaded units)
func UnloadedStatusMessage(unitFile dbus.UnitFile) *Unit {
	return &Unit{
		Path: unitFile.Path,
		Type: unitFile.Type,
	}
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
