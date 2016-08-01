package restartd

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bearstech/restartd/model"
	"github.com/bearstech/restartd/systemd"
)

type Error struct {
	message string
	code    model.Response_Codes
}

func (e *Error) Error() string { // implements error
	return e.message
}

// struct Handler
// implements Handler interface
type Handler struct {
	// array containing services (names)
	Services      []string
	User          string
	PrefixService bool
}

// func Handle
// implemented by Handler interface
func (h *Handler) Handle(m model.Message) (r model.Response) {

	var code model.Response_Codes
	var message string

	if m.GetService() == "--all" {
		if m.GetCommand() == model.Message_status {
			for service := range h.Services {
				//FIXME find all service status
				fmt.Println(service)
			}
		}
	}

	service := m.GetService()
	if h.PrefixService { // alice asks for web, but it's alice-web
		service = fmt.Sprintf("%s-%s", h.User, service)
	}

	// verify if requested unit exists
	ret := systemd.Contains(service, h.Services)

	// if unit does not exists
	if ret != true {
		// write appropriate message
		code = model.Response_err_missing
		message = fmt.Sprintf("Service %s does not exists",
			service)
	} else {

		// switch between all supported commands
		switch m.GetCommand() {

		case model.Message_status:
			// Get status for requested unit
			ret, err := systemd.GetStatus(service)

			// error checking
			// write appropriate messages
			if err != nil {
				message = fmt.Sprintf("Error getting %s service status",
					service)
				code = model.Response_err_status
				log.Error(message)
			} else {
				message = ret
				code = model.Response_suc_status
			}

			break

		case model.Message_start:
			// start unit
			err := systemd.StartUnit(service)

			// error checking
			if err != nil {
				message = fmt.Sprintf("Error starting %s service",
					m.GetService)
				code = model.Response_err_start
			} else {
				message = fmt.Sprintf("%s service is started",
					service)
				code = model.Response_suc_start
			}

			break

		case model.Message_stop:
			// stop unit
			err := systemd.StopUnit(service)

			// error checking
			if err != nil {
				message = fmt.Sprintf("Error stopping %s service",
					m.GetService)
				code = model.Response_err_stop
			} else {
				message = fmt.Sprintf("%s service is stopped",
					service)
				code = model.Response_suc_stop
			}

			break

		case model.Message_restart:
			// restart unit
			err := systemd.RestartUnit(service)

			// error checking
			if err != nil {
				message = fmt.Sprintf("Error restarting %s service",
					m.GetService)
				code = model.Response_err_restart
			} else {
				message = fmt.Sprintf("%s service is restarted",
					service)
				code = model.Response_suc_restart
			}

			break

		case model.Message_reload:
			// reload unit
			err := systemd.ReloadUnit(service)

			// error checking
			if err != nil {
				message = fmt.Sprintf("Error reloading %s service",
					service)
				code = model.Response_err_restart
			} else {
				message = fmt.Sprintf("%s service is reloaded",
					service)
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
