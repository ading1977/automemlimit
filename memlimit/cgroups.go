//go:build linux
// +build linux

package memlimit

import (
	"github.com/containerd/cgroups"
	v2 "github.com/containerd/cgroups/v2"
)

const (
	cgroupPath       = "/proc/self/cgroup"
	cgroupMountPoint = "/sys/fs/cgroup"
)

func FromCgroup() (uint64, error) {
	switch cgroups.Mode() {
	case cgroups.Legacy:
		return FromCgroupV1()
	case cgroups.Hybrid, cgroups.Unified:
		return FromCgroupV2()
	}
	return 0, ErrNoCgroup
}

func FromCgroupV1() (uint64, error) {
	cg, err := cgroups.Load(cgroups.SingleSubsystem(cgroups.V1, cgroups.Memory), cgroups.RootPath)
	if err != nil {
		return 0, err
	}

	metrics, err := cg.Stat()
	if err != nil {
		return 0, err
	} else if metrics.Memory == nil {
		return 0, ErrNoLimit
	}

	return metrics.Memory.HierarchicalMemoryLimit, nil
}

func FromCgroupV2() (uint64, error) {
	m, err := v2.LoadManager(cgroupMountPoint, cgroupPath)
	if err != nil {
		return 0, err
	}

	stats, err := m.Stat()
	if err != nil {
		return 0, err
	} else if stats.Memory == nil {
		return 0, ErrNoLimit
	}

	return stats.Memory.UsageLimit, nil
}