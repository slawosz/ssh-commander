package worker

import (
	"sync"
)

type JobScheduler interface {
	Start()
	PushJob(*SchedulerPayload)
}

type Scheduler struct {
	chanin chan *SchedulerPayload
	worker Worker
	pool   int
}

func NewScheduler(w Worker, pool int) *Scheduler {
	in := make(chan *SchedulerPayload)
	return &Scheduler{
		chanin: in,
		worker: w,
		pool:   pool,
	}
}

func (s *Scheduler) Start() {
	limiter := make(chan bool, s.pool)

	for job := range s.chanin {
		go func(j *SchedulerPayload) {
			// waitGroup is to wait until all hosts will do their job
			wg := &sync.WaitGroup{}
			workerResults := make(chan *HostResult, len(j.Hosts))

			for _, hostPayload := range j.Hosts {
				wg.Add(1)
				limiter <- true // blocks when limiter chan buffer is full

				go func(hostPayload *Host) {
					res := s.worker.Work(hostPayload)
					workerResults <- res

					<-limiter //release place in channel
					wg.Done()
				}(hostPayload)
			}
			wg.Wait()

			close(workerResults) // TODO: test if without it reduce blocks
			hostsResult := make([]*HostResult, 0)

			for r := range workerResults {
				hostsResult = append(hostsResult, r)
			}

			j.ResultsChan <- hostsResult
		}(job)
	}
}

func (s *Scheduler) PushJob(j *SchedulerPayload) {
	s.chanin <- j
}
