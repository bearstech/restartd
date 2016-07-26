package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bearstech/restartd/model"
	"github.com/bearstech/restartd/systemd"
)

// struct Handler
// implements Handler interface
type Handler struct {
	// array containing services (names)
	Services []string
	user     string
}

// func Handle
// implemented by Handler interface
func (h *Handler) Handle(m model.Message) (r model.Response) {

	var code model.Response_Codes
	var message string

	// verify if requested unit exists
	ret := systemd.IsUnit(m.GetService(), h.Services)

	// if unit does not exists
	if ret != true {
		// write appropriate message
		code = model.Response_err_missing
		message = fmt.Sprintf("Service %s does not exists",
			m.GetService())
	} else {

		// switch between all supported commands
		switch m.GetCommand() {

		case model.Message_status:
			// Get status for requested unit
			ret, err := systemd.GetStatus(m.GetService())

			// error checking
			// write appropriate messages
			if err != nil {
				message = fmt.Sprintf("Error getting %s service status",
					m.GetService())
				code = model.Response_err_status
				log.Error(message)
			} else {
				message = ret
				code = model.Response_suc_status
			}

			break

		case model.Message_start:
			// start unit
			err := systemd.StartUnit(m.GetService())

			// error checking
			if err != nil {
				message = fmt.Sprintf("Error starting %s service",
					m.GetService)
				code = model.Response_err_start
			} else {
				message = fmt.Sprintf("%s service is started",
					m.GetService())
				code = model.Response_suc_start
			}

			break

		case model.Message_stop:
			// stop unit
			err := systemd.StopUnit(m.GetService())

			// error checking
			if err != nil {
				message = fmt.Sprintf("Error stopping %s service",
					m.GetService)
				code = model.Response_err_stop
			} else {
				message = fmt.Sprintf("%s service is stopped",
					m.GetService())
				code = model.Response_suc_stop
			}

			break

		case model.Message_restart:
			// restart unit
			err := systemd.RestartUnit(m.GetService())

			// error checking
			if err != nil {
				message = fmt.Sprintf("Error restarting %s service",
					m.GetService)
				code = model.Response_err_restart
			} else {
				message = fmt.Sprintf("%s service is restarted",
					m.GetService())
				code = model.Response_suc_restart
			}

			break

		case model.Message_reload:
			// reload unit
			err := systemd.ReloadUnit(m.GetService())

			// error checking
			if err != nil {
				message = fmt.Sprintf("Error reloading %s service",
					m.GetService())
				code = model.Response_err_restart
			} else {
				message = fmt.Sprintf("%s service is reloaded",
					m.GetService())
				code = model.Response_suc_restart
			}

			break

		default:
			code = model.Response_err_cmd
			message = fmt.Sprint("Command %s not supported",
				m.GetCommand)
		}
	}

	// send message to restartctl client
	return model.Response{
		Code:    &code,
		Message: &message,
	}

}
