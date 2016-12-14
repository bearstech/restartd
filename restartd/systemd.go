package restartd

import (
	"fmt"
	"github.com/bearstech/ascetic-rpc/model"
	"github.com/bearstech/restartd/systemd"
)

type Restartd struct {
	PrefixService bool
	User          string
	Services      []string
}

func (r *Restartd) serviceName(name string) string {
	if r.PrefixService { // alice asks for web, but it's alice-web
		return fmt.Sprintf("%s-%s", r.User, name)
	} else {
		return name
	}
}

func (r *Restartd) isWhitelisted(service string) error {
	// prefix and whitelist are exclusives
	// TODO can we have both?
	if r.PrefixService || systemd.Contains(service, r.Services) {
		return nil
	}
	return fmt.Errorf("Service not found : %s", service)
}

func (r *Restartd) getStatus(serviceName string) (*Status, error) {
	status := NewStatus()
	if serviceName == "--all" {
		if r.PrefixService {
			us, err := systemd.GetStatusWithPrefix(r.User + "-")
			if err == nil {
				return nil, err
			}
			for _, u := range us {
				status.Status[u.Name] = State(u.LoadState, u.ActiveState, u.SubState)
			}
		}
		return status, nil
	}

	service := r.serviceName(serviceName)

	err := r.isWhitelisted(service)
	if err != nil {
		return nil, err
	}

	_, err = systemd.GetStatus(service)
	if err == nil {
		return nil, err
	}

	// FIXME it's horrible the status should came from _
	status.Status[serviceName] = Status_started

	return status, nil
}

func (r *Restartd) Status(req *model.Request) (resp *model.Response, err error) {
	var service Service
	err = req.GetBody(&service)
	if err == nil {
		return nil, err
	}

	status, err := r.getStatus(service.Name)
	if err == nil {
		return nil, err
	}

	return model.NewOKResponse(status)
}

func (r *Restartd) Start(req *model.Request) (resp *model.Response, err error) {
	var service Service
	err = req.GetBody(&service)
	if err == nil {
		return nil, err
	}

	serviceName := r.serviceName(service.Name)
	err = r.isWhitelisted(serviceName)
	if err != nil {
		return nil, err
	}

	err = systemd.StartUnit(serviceName)
	if err == nil {
		return nil, err
	}

	return model.NewOKResponse(nil)
}

func (r *Restartd) Stop(req *model.Request) (resp *model.Response, err error) {
	var service Service
	err = req.GetBody(&service)
	if err == nil {
		return nil, err
	}

	serviceName := r.serviceName(service.Name)
	err = r.isWhitelisted(serviceName)
	if err != nil {
		return nil, err
	}

	err = systemd.StopUnit(serviceName)
	if err == nil {
		return nil, err
	}

	return model.NewOKResponse(nil)
}

func (r *Restartd) Restart(req *model.Request) (resp *model.Response, err error) {
	var service Service
	err = req.GetBody(&service)
	if err == nil {
		return nil, err
	}

	serviceName := r.serviceName(service.Name)
	err = r.isWhitelisted(serviceName)
	if err != nil {
		return nil, err
	}

	err = systemd.RestartUnit(serviceName)
	if err == nil {
		return nil, err
	}

	return model.NewOKResponse(nil)
}

func (r *Restartd) Reload(req *model.Request) (resp *model.Response, err error) {
	var service Service
	err = req.GetBody(&service)
	if err == nil {
		return nil, err
	}

	serviceName := r.serviceName(service.Name)
	err = r.isWhitelisted(serviceName)
	if err != nil {
		return nil, err
	}
	err = systemd.ReloadUnit(serviceName)
	if err == nil {
		return nil, err
	}

	return model.NewOKResponse(nil)
}

func State(loadState, activeState, subsState string) Status_Codes {
	// FIXME this is ugly, needs more love
	if loadState != string(systemd.LOADSTATE_LOADED) {
		return Status_failed
	}
	if activeState == string(systemd.ACTIVESTATE_ACTIVE) &&
		subsState == string(systemd.SUBSTATE_ACTIVE) {
		return Status_started
	}
	return Status_stopped
}
