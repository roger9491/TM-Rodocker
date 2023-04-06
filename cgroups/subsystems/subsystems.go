package subsystems



// 資源限制配置的結構體
type ResourceConfig struct {
	MemoryLimit string	// 記憶體限制
	CpuShare    string	// CPU時間片權重
	CpuSet      string	// CPU核心數
}


// Subsystem interface
type Subsystem interface {
	Name() string	// subsystem name
	// 設定某個cgroup 在這個subsystem中的資源限制
	// path: cgroup在hierarchy中的名稱
	Set(path string, res *ResourceConfig) error
	// 進程添加進某個cgroup
	Apply(path string, pid int) error
	// 移除某個cgroup
	Remove(path string) error
}


var (
	SubsystemsIns = []Subsystem{
		&MemorySubSystem{},
		&CpuSubSystem{},
		&CpusetSubSystem{},
	}
)