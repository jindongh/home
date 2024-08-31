package docker

import (
    "context"
    "log"

    containertypes "github.com/docker/docker/api/types/container"
    apitypes "github.com/docker/docker/api/types"
    "github.com/docker/docker/client"
)

type Docker struct {
    ctx context.Context
    cli *client.Client
}

func Connect() *Docker {
    ctx := context.Background()
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        log.Fatalf("failed to connect to docker", err)
    }
    return &Docker {
        ctx: ctx,
        cli: cli,
    }
}

func (d *Docker) ListContainers() []apitypes.Container{
    containers, err := d.cli.ContainerList(d.ctx, containertypes.ListOptions{All: true})
    if err != nil {
        log.Fatalf("Failed to find container", err)
    } else {
        log.Println("ListContainers found ", len(containers))
    }
    return containers
}

func (d *Docker) FindContainersByName(containers []apitypes.Container, names []string) []apitypes.Container {
    matchContainers := []apitypes.Container{}
    for _, container := range containers {
        for _, containerName := range(container.Names) {
            for _, name := range(names) {
                if containerName == name {
                    log.Println("Found match container", name, container.ID)
                    matchContainers = append(matchContainers, container)
                }
            }
        }
    }
    return matchContainers
}

func (*Docker) AreAllContainersUp(containers []apitypes.Container) bool {
    for _, container := range(containers) {
        if container.State != "running" {
            log.Println("Found container not running", container.Names, container.State)
            return false
        }
    }
    return true
}

func (d *Docker) StartContainers(containers []apitypes.Container) bool {
    allSucceed := true
    log.Println("begin start containers", len(containers))
    for _, container := range(containers) {
        log.Println("begin start container", container.Names)
        err := d.cli.ContainerStart(d.ctx, container.ID, containertypes.StartOptions{})
        log.Println("finish start container", container.Names, err)
        if err != nil {
            allSucceed = false
        }
    }
    return allSucceed
}

func (d *Docker) StopContainers(containers []apitypes.Container) bool {
    allSucceed := true
    log.Println("begin stop containers", len(containers))
    for _, container := range(containers) {
        log.Println("begin stop container", container.Names)
        err := d.cli.ContainerStop(d.ctx, container.ID, containertypes.StopOptions{})
        log.Println("finish stop container", container.Names, err)
        if err != nil {
            allSucceed = false
        }
    }
    return allSucceed
}
