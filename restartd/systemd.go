package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bearstech/restartd/protocol"
	"github.com/bearstech/restartd/systemd"
)

// struct HandlerSystemd
// implements Handler interface
type HandlerSystemd struct {
	// array containing services (names)
	Services []string
}

// func Handle
// implemented by Handler interface
func (h *HandlerSystemd) Handle(m protocol.Message) (r protocol.Response) {

	var code protocol.Response_Codes
	var message string

	// verify if requested unit exists
	ret := systemd.IsUnit(m.GetService(), h.Services)

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
			ret, err := systemd.GetStatus(m.GetService())

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
			err := systemd.StartUnit(m.GetService())

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
			err := systemd.StopUnit(m.GetService())

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
			err := systemd.RestartUnit(m.GetService())

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
			err := systemd.ReloadUnit(m.GetService())

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
