package process

import (
	"fmt"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"github.com/tinycedar/lily/common"
	"golang.org/x/sys/windows"
)

func (p *Process) Ppid(pid uint32) (uint32, error) {
	return parentProcessID(pid)
}

func parentProcessID(pid uint32) (uint32, error) {
	var p *Process
	pe, err := getProcesses(pid)
	if err != nil {
		return 0, fmt.Errorf("could not get Name: %s", err)
	}

	ppid := pe.ParentProcessID
	p = &Process{
		ppid: ppid,
	}
	return p.ppid, nil
}

func (p *Process) MemInfo(pid uint32) (bool, error) {
	return processMemInfo(pid)
}

func processMemInfo(pid uint32) (bool, error) {
	mkernel32 := syscall.MustLoadDLL("kernel32.dll")
	procGetProcessMemoryInfo1 := mkernel32.MustFindProc("GetProcessMemoryInfo")
	handle := OpenProcessHandle(int(pid))
	var memory ProcessMmoryCounters
	const infoSize = unsafe.Sizeof(memory)
	result, _, _ := procGetProcessMemoryInfo1.Call(uintptr(handle), uintptr(unsafe.Pointer(&memory)), infoSize)
	fmt.Println(result)
	if result != 0 {
		return true, nil

	}
	return false, nil
}

type Processs uintptr

const (
	PROCESS_ALL_ACCESS        = 0x1F0FFF
	PROCESS_QUERY_INFORMATION = 0x0400
	PROCESS_VM_READ           = 0x0010
	PROCESS_SET_INFORMATION   = 0x0200
)

const (
	ProcessMemoryPriority = iota
	ProcessMemoryExhaustionInfo
	ProcessAppMemoryInfo
	ProcessInPrivateInfo
	ProcessPowerThrottling
	ProcessReservedValue1
	ProcessTelemetryCoverageInfo
	ProcessProtectionLevelInfo
	ProcessLeapSecondInfo
	ProcessInformationClassMax
)

// Status return status process:
// R: running, Z: zombie, W: wait,
// L: lock, S: sleep.
func (p *Process) Status(pid uint32) string {
	return statusProcess(pid)
}

func statusProcess(pid uint32) string {

	handle := OpenProcessHandle(int(pid))

	return string(handle)
}
func OpenProcessHandle(processId int) Processs {
	var (
		mkernel32                  = syscall.MustLoadDLL("kernel32.dll")
		procOpenProcess1           = mkernel32.MustFindProc("OpenProcess")
		procGetProcessInformation1 = mkernel32.MustFindProc("GetProcessInformation")
	)
	handle, _, _ := procOpenProcess1.Call(ptr(PROCESS_SET_INFORMATION), ptr(true), ptr(processId))
	q, _, _ := procGetProcessInformation1.Call(handle, ProcessInformationClassMax)
	_ = q
	return Processs(q)
}

func ptr(val interface{}) uintptr {
	switch val.(type) {
	case string:
		return uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(val.(string))))
	case int:
		return uintptr(val.(int))
	default:
		return uintptr(0)
	}
}

func (p *Process) Name(pid uint32) (string, error) {
	return processName(pid)
}

func processName(pid uint32) (string, error) {
	var p *Process
	pe, err := getProcesses(pid)
	if err != nil {
		return "", fmt.Errorf("could not get Name: %s", err)
	}

	name := parseProcessName(pe.ExeFile)
	p = &Process{
		name: name,
	}
	return p.name, nil
}

func getProcesses(pid uint32) (*ProcessEntryMy, error) {

	var procEntry ProcessEntryMy

	var (
		mkernel32                     = syscall.MustLoadDLL("kernel32.dll")
		procCreateToolhelp32Snapshot1 = mkernel32.MustFindProc("CreateToolhelp32Snapshot")
		procProcess32FirstW1          = mkernel32.MustFindProc("Process32FirstW")
		procProcess32NextW1           = mkernel32.MustFindProc("Process32NextW")
		procCloseHandle1              = mkernel32.MustFindProc("CloseHandle")
	)

	snapshot, _, _ := procCreateToolhelp32Snapshot1.Call(syscall.TH32CS_SNAPPROCESS, 0)
	// if err != nil {
	// 	common.Error("Fail to syscall CreateToolhelp32Snapshot: %v", err)
	// 	return nil, nil
	// }
	defer procCloseHandle1.Call(snapshot)
	procEntry.Size = uint32(unsafe.Sizeof(procEntry))
	r, _, err := procProcess32FirstW1.Call(snapshot, uintptr(unsafe.Pointer(&procEntry)))
	if r == 0 {
		common.Error("Fail to syscall Process32First: %v", err)
		return nil, nil
	}
	for {
		if procEntry.ProcessID == pid {
			return &procEntry, nil
		}
		r, _, _ := procProcess32NextW1.Call(snapshot, uintptr(unsafe.Pointer(&procEntry)))
		if r == 0 {
			break
		}
	}
	return nil, nil
}

func parseProcessName(exeFile [syscall.MAX_PATH]uint16) string {
	for i, v := range exeFile {
		if v <= 0 {
			return string(utf16.Decode(exeFile[:i]))
		}
	}
	return ""
}

func (p *Processes) Pids() ([]uint32, error) {
	return pidsWithContext()
}

func getPids() (uint32, error) {
	return uint32(syscall.Getpid()), nil
}

func pidsWithContext() ([]uint32, error) {
	var p *Processes
	var ret []uint32
	var read uint32 = 0
	var psSize uint32 = 1024
	const dwordSize uint32 = 4

	for {
		ps := make([]uint32, psSize)
		if err := windows.EnumProcesses(ps, &read); err != nil {
			return nil, err
		}
		if uint32(len(ps)) == read { // ps buffer was too small to host every results, retry with a bigger one
			psSize += 1024
			continue
		}
		for _, pid := range ps[:read/dwordSize] {
			ret = append(ret, uint32(pid))
		}
		break
	}
	p = &Processes{
		pids: ret,
	}
	return p.pids, nil
}
