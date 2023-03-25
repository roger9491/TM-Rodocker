package cgroups

import (
	"TM-Rodocker/subsystems"

	"github.com/sirupsen/logrus"
)

// cgroup 結構管理
type CgroupManager struct {
	Path	 string	// cgroup 在 hierarchy 中的路徑
	// 資源配置
	Resource *subsystems.ResourceConfig
}


// 建立 CgroupManager 實例
func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}


// 將進程加入到每個subsystem的cgroup中
func (c *CgroupManager) Apply(pid int) error {
	for _, subSysIns := range subsystems.SubsystemsIns {
		subSysIns.Apply(c.Path, pid)
	}
	return nil
}

// 設定每個subsystem的 cgroup資源限制
func (c *CgroupManager) Set(res *subsystems.ResourceConfig) error {
	for _, subSysIns := range subsystems.SubsystemsIns {
		subSysIns.Set(c.Path, res)
	}
	return nil
}

// 移除 各個 subsystem的 cgroup
func (c *CgroupManager) Destroy() error {
	for _, subSysIns := range subsystems.SubsystemsIns {
		if err := subSysIns.Remove(c.Path); err != nil {
			logrus.Warnf("remove cgroup fail %v", err)
		}
	}
	return nil
}