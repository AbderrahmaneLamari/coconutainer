package runtime

import (
    "fmt"
    "os"
    "os/exec"
    "syscall"

    "github.com/opencontainers/runc/libcontainer"
    "github.com/opencontainers/runc/libcontainer/console"
    "github.com/opencontainers/runc/libcontainer/utils"
    "github.com/opencontainers/runc/libcontainer/configs"

    "coconutainer/types"
)

func RunContainer(cfg *types.ContainerConfig) error {
    // 1. Convert high-level config → libcontainer config
    lcCfg, err := cfg.ToLibcontainerConfig()
    if err != nil {
        return fmt.Errorf("invalid config: %w", err)
    }

    // 2. Create a factory storing state in /run/coco
    factory, err := libcontainer.New(
        "/run/coco",
        libcontainer.Cgroupfs,
    )
    if err != nil {
        return fmt.Errorf("factory error: %w", err)
    }

    // 3. Create container instance
    container, err := factory.Create(cfg.ID, lcCfg)
    if err != nil {
        return fmt.Errorf("create error: %w", err)
    }

    // Prepare arguments (command executed inside container)
    process := &libcontainer.Process{
        Args:   cfg.Cmd,
        Env:    append(os.Environ(), "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/bin"),
        User:   "0:0",
        Cwd:    "/",
        Stdin:  os.Stdin,
        Stdout: os.Stdout,
        Stderr: os.Stderr,
    }

    // 4. Setup console (TTY)
    console, err := console.NewPty()
    if err != nil {
        return fmt.Errorf("pty error: %w", err)
    }
    process.ConsoleSocket = console

    // 5. Start the container process
    if err := container.Start(process); err != nil {
        return fmt.Errorf("start error: %w", err)
    }

    // 6. Wait for container process to exit
    status, err := process.Wait()
    if err != nil {
        return fmt.Errorf("wait error: %w", err)
    }

    fmt.Printf("Container exited with: %v\n", status)

    // 7. Destroy container (delete cgroups, state dir)
    if err := container.Destroy(); err != nil {
        return fmt.Errorf("destroy error: %w", err)
    }

    return nil
}
