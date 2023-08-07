package builtins

import (
	"os"
	"path/filepath"
	"strconv"
)

const TempDir = "/var/temp"

type GetReduceCountArgs struct {
}

type GetReduceCountReply struct {
	ReduceCount int
}

type RequestTaskArgs struct {
	WorkerId int
}

type RequestTaskReply struct {
	TaskType TaskType
	TaskId   int
	TaskFile string
}

type ReportTaskArgs struct {
	WorkerId int
	TaskType TaskType
	TaskId   int
}

type ReportTaskReply struct {
	CanExit bool
}

func masterSock() string {
	s := filepath.Join(TempDir, "mapreduce_Master-")
	s += strconv.Itoa(os.Getuid())
	return s

}
