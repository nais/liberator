package nais_io_v1

func (in *Naisjob) GetStatus() *Status {
	return &in.Status
}

func (in *Naisjob) SetStatus(status *Status) {
	in.Status = *status
}
