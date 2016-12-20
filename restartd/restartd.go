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

func (r *Restartd) getAllStatus() ([]*Status_State, error) {
	states := []*Status_State{}
	if r.PrefixService {
		units, err := systemd.GetStatusWithPrefix(r.User + "-")
		if err != nil {
			return nil, err
		}
		for _, unit := range units {
			states = append(states, &Status_State{
				Name:  unit.Name,
				State: statusState(unit.State),
			})
		}
	} else {
		for _, service := range r.Services {
			s, err := r.getStatus(service)
			if err != nil {
				return nil, err
			}
			states = append(states, s)
		}
	}
	return states, nil
}

func (r *Restartd) getStatus(serviceName string) (*Status_State, error) {
	service := r.serviceName(serviceName)

	err := r.isWhitelisted(service)
	if err != nil {
		return nil, err
	}

	st, err := systemd.GetStatus(service)
	if err != nil {
		return nil, err
	}

	return &Status_State{
		Name:  serviceName,
		State: statusState(st.State),
	}, nil
}

func statusState(s *systemd.State) Status_States {
	if s.Active == systemd.ACTIVESTATE_ACTIVE {
		return Status_started
	}
	if s.Active == systemd.ACTIVESTATE_INACTIVE {
		return Status_stopped
	}
	return Status_failed
}

/*

RPC

*/

func (r *Restartd) StatusAll(req *model.Request) (resp *model.Response, err error) {
	status, err := r.getAllStatus()
	if err == nil {
		return nil, err
	}

	return model.NewOKResponse(&Status{Status: status})
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
