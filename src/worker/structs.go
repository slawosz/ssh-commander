package worker

type Worker interface {
	Work(*WorkerPayload) *HostResult
}

type IScheduler interface {
	Start()
	Push(*JobPayload)
	ResultsChan() chan *HostsResult
}

type Host struct {
	Host     string
	Port     string
	User     string
	Password string
}

type JobPayload struct {
	Hosts   []*Host
	Command string
	JID     string
	Timeout int
}

type WorkerPayload struct {
	*Host
	Command string
}

type HostResult struct {
	Payload string
	Host    string
	Port    string
}

type HostsResult struct {
	JID         string
	HostsResult []*HostResult
}
