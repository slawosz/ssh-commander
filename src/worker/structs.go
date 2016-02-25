package worker

type Worker interface {
	Work(*Host) *HostResult
}

type Host struct {
	Host     string
	Port     string
	User     string
	Password string
	Commands []string
	Prompt   string
	Exit     string
}

type SchedulerPayload struct {
	Hosts       []*Host
	ResultsChan chan []*HostResult
}

type HostResult struct {
	Payload string
	Host    string
	Port    string
	Status  string
}
