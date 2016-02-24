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

type JobPayload struct {
	Hosts   []*Host
	Script  string
	Command string
	JID     string
	Timeout int
}

type HostResult struct {
	Payload string
	Host    string
	Port    string
	Status  string
}

type HostsResult struct {
	JID         string
	HostsResult []*HostResult
}
