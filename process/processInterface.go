package process

type ProcessInfo interface {
	Status(uint32) string
	Name(uint32) (string, error)
	Ppid(uint32) (uint32, error)
	MemInfo(uint32) (bool, error)
}

type ProcessesInfo interface {
	Pids() ([]uint32, error)
}
