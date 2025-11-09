package runtime

import (
    "fmt"
    "os"
    "path/filepath"
    "syscall"

    "github.com/opencontainers/runc/libcontainer"
    "github.com/opencontainers/runc/libcontainer/configs"
)

// Runtime is a wrapper around "libcontainer.Factory"
// It's responsible for creating, running, stopping and deleting containers

// Declaring a runtime type

type Runtime struct {
	factory 	libcontainer.Factory
	baseDir 	string
}

func NewRuntime(baseDir string) (*Runtime, error){
	// Let's ensure runtime state directory exisits

	err := os.Mkdir(baseDir, 0755);
	if err != nil {
		return nil, err
	}

	// Creating the libcontainer factory
	factory, err := libcontainer.New(baseDir, libcontainer.Cgroups, libcontainer.InitArgs(os.Args[0], "init"));

	if err != nil{
		return nil, fmt.Errorf("Failed to create factory: %v", err);
	}


	return &Runtime{
		factory: factory,
		baseDir: baseDir,
	}, nil

}

func (r *Runtime) CreateContainer(cfg *ContainerConfig) (libcontainer.Container, error){

	libCfg, err := cfg.ToLibcontainerConfig();

	if err != nil{
		return nil, err;
	}

	container, err := r.factory.Create(cfg.ID, libCfg);

	if err != nil {
		return nil, fmt.Errorf("Failed to Create container: %v", err);
	}

	return container, nil;
}

func (r* Runtime) StartContainer(container libcontainer.Container, cfg ContainerConfig) error {
	process := &libcontainer.Process{
		Args: cfg.Cmd,
		Env: []string{"PATH=/bin:/sbin:/usr/bin:/usr/sbin"},
		Cwd: "/",
		Stdin: os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	fmt.Printf("[coconutainer] starting container %s ...", cfg.ID);


	err := container.Run(process);

	if err != nil {
		return fmt.Errorf("Failed to run process: %v", err);
	}

	_, err := process.Wait()
	return err;
	
}


func (r *Runtime) StopContainer(container libcontainer.Container) error {

	procs, err := container.Processes()

	if err != nil {
		return err
	}

	for _, pid := range procs {
		syscall.kill(pid, syscall.SIGTERM);
	}

	return nil;
}

func (r *Runtime) DestroyContainer(container libcontainer.Container) error {
	
	err := container.Destroy();

	if err != nil {
		return fmt.Errorf("Failed to Destroy Container: %v", err);
	}
	
	return nil;
}