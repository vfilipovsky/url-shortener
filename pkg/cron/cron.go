package cron

type Job interface {
	Run()
	GetName() string
}

type Cron interface {
	Schedule(cronExpression Expression, job Job)
	Start()
}
