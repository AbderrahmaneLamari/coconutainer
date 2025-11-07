package runtime

import (
	"os"
	"path/filepath"
	"github.com/opencontainers/runc/libcontainer/configs"
)

// ContainerConfig Defines a high-level configuration for a container.
// The engine will understand this (Coconutainer) specific configuration and translate it
// into a libcontainer.Config.

type ContainerConfig struct {
	ID 			string		// Unique container id
	Rootfs  	string 		// Path to container FS
	Cmd    		[]string 	// Command to run in the container
	Hostname 	string 		// Container hostname
	Network  	bool 		// Whether to enable networking
	Mounts   	[]string 	// List of mounts: Extra mounts
	CpuLimit 	int 		// CPU limit (optional)
	MemoryLimit int64 		// Memory limit (optional)
}


// ToLibcontainerConfig Translates the high-level ContainerConfig into a libcontainer.Config
func (c *ContainerConfig) ToLibcontainerConfig() (*configs.Config, error ){

	rootfsPath, err := filepath.Abs(c.Rootfs);
	if err != nil{
		return nil, err
	}

	namespaces := []configs.Namespace{
		{Type: configs.NEWNS},     // Mount
		{Type: configs.NEWUTS},	   // Hostname
		{Type: configs.NEWIPC},	   // Interprocess
		{Type: configs.NEWPID},	   // Process ID isolation
	}

	if c.Network{
		namespaces = append(namespaces, configs.Namespace{Type: configs.NEWNET}) // Network Namespace
	}

	// Mount the basic filesystem directories
	mounts := []*configs.Mount{
		{
			Source: "proc",
			Destination: "/proc",
			Device: "proc",
		},
		{
			Srouce: "tmpfs",
			Destination: "/dev",
			Device: "tmpfs",
			Flags: syscall.MS_NOSUID | syscall.MS_STRICTATIME,
			Data: "mode=755",
		}
	}

	// Let's add user-defined mounts to our contianers!

	for _, mnt := range c.Mounts {
		parts := strings.SplitN(mnt, ":", 2);

		if len(parts) == 2 {
			mounts = append(mounts, &configs.Mount{
				Source: parts[0],
				Destination: parts[1],
				Device: "bind",
				Flags: syscall.MS_BIND | syscall.MS_REC,
			});

		}
	}

	cfg := &configs.Config{
		Rootfs: rootfsPath,
		Readonlyfs: false,
		Hostname: c.Hostname,
		Namespaces: namespaces,
		Mounts: mounts,
		Capabilities: &configs.Capabilites{},
		Cgroups: &configs.Cgroup{
			Name: c.ID,
			Parent: "coconutainer",
			Resources: &configs.Resources{
				Memory: cMemoryLimit,
				CpuQuota: int64(c.CpuLimit * 1000),
			}
		}
	}

	return cfg, nil;
}
