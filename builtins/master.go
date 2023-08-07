package builtins

import "sync"

type TaskStatus int
type TaskType int
type JobStage int

const (
	MapTask TaskType = iota
	ReduceTask
	NoTask
	ExitTask
)

const (
	NotStarted TaskStatus = iota
	Executing
	Finished
)

type Task struct {
	Type     TaskType
	Status   TaskStatus
	Index    int
	File     string
	WorkerId int
}

type Master struct {
	mutex         sync.Mutex
	mapTasks      []Task
	mapTasksNo    int
	reduceTasks   []Task
	reduceTasksNo int
}
