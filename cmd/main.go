package main

import (
	"github.com/opencontainers/cgroups"
	"github.com/opencontainers/runc/libcontainer"
	"github.com/opencontainers/runc/libcontainer/configs"
	"syscall"
	"os"
	"fmt"
)

func main() {

	// devices := configs.DefaultDevices();

	caps := &configs.Capabilities{
		Bounding: []string{
			"CAP_CHOWN", "CAP_DAC_OVERRIDE", "CAP_FSETID",
			"CAP_FOWNER", "CAP_MKNOD", "CAP_NET_RAW",
			"CAP_SETGID", "CAP_SETUID", "CAP_SETFCAP",
			"CAP_SETPCAP", "CAP_NET_BIND_SERVICE",
			"CAP_SYS_CHROOT", "CAP_KILL", "CAP_AUDIT_WRITE",
		},
		Effective: []string{
			"CAP_CHOWN", "CAP_DAC_OVERRIDE", "CAP_FSETID",
			"CAP_FOWNER", "CAP_MKNOD", "CAP_NET_RAW",
			"CAP_SETGID", "CAP_SETUID", "CAP_SETFCAP",
			"CAP_SETPCAP", "CAP_NET_BIND_SERVICE",
			"CAP_SYS_CHROOT", "CAP_KILL", "CAP_AUDIT_WRITE",
		},
		Permitted: []string{
			"CAP_CHOWN", "CAP_DAC_OVERRIDE", "CAP_FSETID",
			"CAP_FOWNER", "CAP_MKNOD", "CAP_NET_RAW",
			"CAP_SETGID", "CAP_SETUID", "CAP_SETFCAP",
			"CAP_SETPCAP", "CAP_NET_BIND_SERVICE",
			"CAP_SYS_CHROOT", "CAP_KILL", "CAP_AUDIT_WRITE",
		},
		Inheritable: []string{},
		Ambient:     []string{},
	}

	mounts := []*configs.Mount{
		{Source: "proc", Destination: "/proc", Device: "proc", Flags: 0, Data: ""},
		{Source: "tmpfs", Destination: "/tmp", Device: "tmpfs", Flags: syscall.MS_NOSUID | syscall.MS_NODEV | syscall.MS_NOEXEC, Data: "mode=1777,size=65536k"},
		{Source: "/dev/null", Destination: "/dev/null", Device: "bind", Flags: syscall.MS_BIND, Data: ""},
		{Source: "/dev/zero", Destination: "/dev/zero", Device: "bind", Flags: syscall.MS_BIND, Data: ""},
		{Source: "/dev/urandom", Destination: "/dev/urandom", Device: "bind", Flags: syscall.MS_BIND, Data: ""},
	}
	namespaces := configs.Namespaces([]configs.Namespace{
		{Type: configs.NEWNS},
		{Type: configs.NEWPID},
		{Type: configs.NEWUTS},
		{Type: configs.NEWIPC},
		{Type: configs.NEWNET},
		{Type: configs.NEWUSER},
	})

	cgroup := &cgroups.Cgroup{
		Name:   "jaouad-container",
		Parent: "system.slice",
		Resources: &cgroups.Resources{
			Memory:      128 * 1024 * 1024, // 128 MB
			CpuShares:   512,               // relative CPU weight
			CpuQuota:    -1,                // no quota
			CpusetCpus:  "",
			BlkioWeight: 500,
			// max 64 processes
		},
	}

	configurations := &configs.Config{
		Rootfs:          "./rootfs",
		Capabilities:    caps,
		Mounts:          mounts,
		Namespaces:      namespaces,
		Cgroups:         cgroup,
		NoNewPrivileges: true,
		Readonlyfs:      false,
		MaskPaths:       []string{"/proc/kcore"},
		ReadonlyPaths:   []string{},
	}
	container, err := libcontainer.Create("./rootfs", "jaouad-contaienr", configurations)

	if err != nil {
		panic(err)
	}

	process := &libcontainer.Process{
        Args:   []string{"/bin/sh"},
        Env:    []string{"PATH=/bin"},
        Stdin:  os.Stdin,
        Stdout: os.Stdout,
        Stderr: os.Stderr,
    }
	// 5. Run the process
    if err := container.Run(process); err != nil {
        panic(err)
    }
	// 6. Wait for it to exit
    status, err := process.Wait()
    if err != nil {
        panic(err)
    }

	fmt.Println("Container exited with status:", status)

    // 7. Destroy container (cleanup cgroups + state dir)
    if err := container.Destroy(); err != nil {
        panic(err)
    }
}
