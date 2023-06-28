package commands

import "github.com/todopeer/cli/api"

// the shared task variables

// for TaskList
var (
	statusForQuery            []string
	mapStatusShort2TaskStatus = map[string]api.TaskStatus{
		"n": api.TaskStatusNotStarted,
		"i": api.TaskStatusDoing,
		"d": api.TaskStatusDone,
		"p": api.TaskStatusPaused,

		"not_started": api.TaskStatusNotStarted,
		"doing": api.TaskStatusDoing,
		"done": api.TaskStatusDone,
		"paused": api.TaskStatusPaused,
	}
)

// for TaskUpdate
var (
	varName         string
	varDescription  string
	varDueDate      string
	varTriggerPause bool
)

// for EventUpdate
var (
	varStartAtStr string
	varEndAtStr   string

	varDayoffsetStr string
	varNewTaskIDStr string
	
	varDurationStr string
)

// for Me
var (
	flagSimpleOutput bool
)

// for TaskStart
var (
	varDurationOffset string
)