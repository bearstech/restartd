package restartd

func NewStatus() *Status {
	return &Status{
		Status: make(map[string]Status_Codes),
	}
}
