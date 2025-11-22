package types

import (
    "github.com/opencontainers/runc/libcontainer/configs"
)

// ContainerConfig represents a high-level configuration for a container.
// This structure is engine-friendly and can be converted into libcontainer.Config later.
type ContainerConfig struct {

    // Basic container identity
    ID          string            `json:"id"`
    Hostname    string            `json:"hostname"`

    // Root filesystem path
    Rootfs      string            `json:"rootfs"`

    // Process to run inside container
    Command     []string          `json:"command"`
    Env         []string          `json:"env"`

    // Namespaces to enable (pid, net, uts, ipc, mount...)
    Namespaces  []configs.NamespaceType `json:"namespaces"`

    // Linux capabilities (bounding, permitted, effective...)
    Capabilities *configs.Capabilities  `json:"capabilities"`

    // Cgroup limits
    Resources   *configs.Resources      `json:"resources"`

    // Mounts = volumes inside the container
    Mounts      []configs.Mount         `json:"mounts"`

    // Network configuration (expand later)
    Network     *NetworkConfig          `json:"network"`

    // Runtime flags
    ReadonlyFS  bool                    `json:"readonly_fs"`
    NoNewPrivs  bool                    `json:"no_new_privs"`

    // Annotations for extra metadata
    Labels      map[string]string       `json:"labels"`
}

// Future networking config (simple version)
type NetworkConfig struct {
    Enabled   bool     `json:"enabled"`
    Interface string   `json:"interface"` // e.g., eth0, or palm0
    IPAddress string   `json:"ip"`
    Gateway   string   `json:"gateway"`
    DNS       []string `json:"dns"`
}
