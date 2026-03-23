package metrics

type Snapshot struct {
	ProposalsTotal uint64
	CommitLatency  float64
	Throughput     float64
	ElectionMillis float64
}

type Collector struct {
	last Snapshot
}

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) Update(snapshot Snapshot) {
	c.last = snapshot
}

func (c *Collector) Snapshot() Snapshot {
	return c.last
}
