package process

import (
	"time"

	"exemple.com/psutil/cpu"
)

const MAX_PATH = 260

type Processes struct {
	pids []uint32
}

type ProcessEntryMy struct {
	Size            uint32
	Usage           uint32
	ProcessID       uint32
	DefaultHeapID   uintptr
	ModuleID        uint32
	Threads         uint32
	ParentProcessID uint32
	PriClassBase    int32
	Flags           uint32
	ExeFile         [MAX_PATH]uint16
}

type Process struct {
	pid          uint32
	name         string
	status       string
	ppid         uint32
	memInfo      bool
	sigInfo      *SignalInfoStat
	createTime   int64
	lastCPUTimes *cpu.TimesStat
	lastCPUTime  time.Time
	tgid         int32
}

type MemoryInfoStat struct {
	RSS    uint64 `json:"rss"`
	VMS    uint64 `json:"vms"`
	HWM    uint64 `json:"hwm"`
	Data   uint64 `json:"data"`
	Stack  uint64 `json:"stack"`
	Locked uint64 `json:"locked"`
	Swap   uint64 `json:"swap"`
}

type SignalInfoStat struct {
	PendingProcess uint64 `json:"pending_process"`
	PendingThread  uint64 `json:"pending_thread"`
	Blocked        uint64 `json:"blocked"`
	Ignored        uint64 `json:"ignored"`
	Caught         uint64 `json:"caught"`
}

type ProcessMmoryCounters struct {
	Cd                         byte
	PageFaultCount             uint32
	PeakWorkingSetSize         byte
	WorkingSetSize             byte
	QuotaPeakPagedPoolUsage    byte
	QuotaPagedPoolUsage        byte
	QuotaPeakNonPagedPoolUsage byte
	QuotaNonPagedPoolUsage     byte
	PagefileUsage              byte
	PeakPagefileUsage          byte
}
