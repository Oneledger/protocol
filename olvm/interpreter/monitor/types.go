/*
	Copyright 2017-2018 OneLedger
*/
package monitor

// TODO: Conflicts with the protocol's status code
type StatusCode int

const (
	STATUS_OK StatusCode = iota
	STATUS_WARNING
	STATUS_ERROR
	STATUS_DEADLOOP
	STATUS_PANIC
	STATUS_ALREADY_RUNNING
)

type RunningMode int

const (
	DEFAULT_MODE RunningMode = iota
	AGGRESIVE_MODE
	CONSERVATIVE_MODE
)

type Monitor struct {
	TickerThreshold int
	RunningMode     RunningMode
	PidFilePath     string
}

type Status struct {
	Details string
	Code    StatusCode
}
