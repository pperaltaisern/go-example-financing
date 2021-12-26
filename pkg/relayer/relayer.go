package relayer

type Relayer struct{}

func (r *Relayer) Start() {
	r.UnpublishedEvents()
}

func (r *Relayer) UnpublishedEvents() {}
