package builtins

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"time"
)

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

// -------------------- internal Master functions ----------------------

func (m *Master) selectTasks(taskList []Task, workerId int) *Task {
	var task *Task
	for i := 0; i < len(taskList); i++ {
		if taskList[i].Status == NotStarted {
			task = &taskList[i]
			task.Status = Executing
			task.WorkerId = workerId
			return task
		}
	}

	return &Task{NoTask, Finished, -1, "", 0}
}

func (m *Master) waitForTask(task *Task) {
	if task.Type != MapTask && task.Type != ReduceTask {
		return
	}

	<-time.After(time.Millisecond * TaskTimeout)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if task.Status == Executing {
		task.Status = NotStarted
		task.WorkerId = -1
	}
}

// --------------------- RPC Handlers of Master ------------------------

// RPC handler that gets the number of reduce tasks
func (m *Master) GetReduceCount(args *GetReduceCountArgs, reply *GetReduceCountReply) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	reply.ReduceCount = len(m.reduceTasks)

	return nil
}

func (m *Master) RequestTask(args *RequestTaskArgs, reply *RequestTaskReply) error {
	m.mutex.Lock()
	var task *Task
	if m.mapTasksNo > 0 {
		task = m.selectTasks(m.mapTasks, args.WorkerId)
	} else if m.reduceTasksNo > 0 {
		task = m.selectTasks(m.reduceTasks, args.WorkerId)
	} else {
		task = &Task{ExitTask, Finished, -1, "", 0}
	}

	reply.TaskFile = task.File
	reply.TaskId = task.Index
	reply.TaskType = task.Type

	m.mutex.Unlock()

	return nil
}

func (m *Master) ReportTaskDone(args *ReportTaskArgs, reply *ReportTaskReply) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var task *Task
	if args.TaskType == MapTask {
		task = &m.mapTasks[args.TaskId]
	} else if args.TaskType == ReduceTask {
		task = &m.reduceTasks[args.TaskId]
	} else {
		fmt.Printf("Incorrect task type to report: %v\n", args.TaskType)
		return nil
	}

	if args.WorkerId == task.WorkerId && task.Status == Executing {
		task.Status = Finished
		if args.TaskType == MapTask && m.mapTasksNo > 0 {
			m.mapTasksNo--
		} else if args.TaskType == ReduceTask && m.reduceTasksNo > 0 {
			m.reduceTasksNo--
		}
	}

	reply.CanExit = (m.mapTasksNo == 0 && m.reduceTasksNo == 0)

	return nil
}

// Thread for listening RPC server

func (m *Master) server() {
	rpc.Register(m)
	rpc.HandleHTTP()
	sockname := masterSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

func (m *Master) Done() bool {
	// Your code here.
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.mapTasksNo == 0 && m.reduceTasksNo == 0
}

// --------------------- Constructor ------------------------

func MakeMaster(files []string, nReduce int) *Master {
	m := Master{}

	// Your code here.
	nMap := len(files)
	m.mapTasksNo = nMap
	m.reduceTasksNo = nReduce
	m.mapTasks = make([]Task, 0, nMap)
	m.reduceTasks = make([]Task, 0, nReduce)

	for i := 0; i < nMap; i++ {
		mTask := Task{MapTask, NotStarted, i, files[i], -1}
		m.mapTasks = append(m.mapTasks, mTask)
	}
	for i := 0; i < nReduce; i++ {
		rTask := Task{ReduceTask, NotStarted, i, "", -1}
		m.reduceTasks = append(m.reduceTasks, rTask)
	}

	m.server()
	return &m
}
