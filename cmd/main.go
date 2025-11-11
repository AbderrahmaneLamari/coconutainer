package main

import (
    "log"
    "coconutainer/cmd/coconutainer/internal/runtime"
)

func main() {
	r, err := runtime.NewRuntime("/var/lib/coconutainer/containers")
	if err != nil {
		log.Fatal(err)
	}

	cfg := &runtime.ContainerConfig{
		ID:       "demo1",
		Rootfs:   "./rootfs", // must contain /bin/sh, etc.
		Cmd:      []string{"/bin/sh"},
		Hostname: "demo1",
	}

	container, err := r.CreateContainer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := r.StartContainer(container, cfg); err != nil {
		log.Fatal(err)
	}

	r.DestroyContainer(container)
}
