package utils

import (
	"fmt"
	"github.com/opencontainers/runc/libcontainer/configs"
    "coconutainer/types"
)


func (c *ContainerConfig) ToLibcontainerConfig() (*config.Config, error) {

	if c.Rootfs == "" {
		return nil, fmt.Errorf("RootFS must be provided");
	}

	ns := []configs.Namespaces{

		{Type: configs.NEWPID},
        {Type: configs.NEWIPC},
        {Type: configs.NEWUTS},
        {Type: configs.NEWNS},
        {Type: configs.NEWUSER},
        {Type: configs.NEWNET},
        {Type: configs.NEWMOUNT},
	} 

	// Override with user-provided namespaces
    if len(c.Namespaces) > 0 {
        ns = []configs.Namespace{}
        for _, n := range c.Namespaces {
            ns = append(ns, configs.Namespace{Type: n})
        }
    }

	cfg := &configs.Config{
        Rootfs:     c.Rootfs,
        Labels:     c.Labels,
        NoNewKeyring: true,
        Readonlyfs: c.ReadonlyFS,

        Namespaces: ns,

        Capabilities: c.Capabilities,

        Devices: configs.DefaultDevices,

        Mounts: c.Mounts,

        Cgroups: &configs.Cgroup{
            Name:   c.ID,
            Parent: "coconutainer",
            Resources: c.Resources,
        },

        NoNewPrivileges: c.NoNewPrivs,
    }

    return cfg, nil

}
