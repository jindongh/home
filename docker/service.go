package docker

import (
    "os"
    "strings"
)

const (
    MEDIA = "media"
    PHOTO = "photo"
    BOOK = "book"
    DOWNLOAD = "download"
    MONITOR = "monitor"
)
var SERVICE_LIST = []string {
    MEDIA,
    PHOTO,
    BOOK,
    DOWNLOAD,
    MONITOR,
}
var containerNamesByService = map[string][]string{
    DOWNLOAD: []string{"/aria2-pro", "/ariang"},
    MONITOR: []string{"/monitoring-grafana"},
}

type Service struct {
    client *Docker
}
type ServiceStatus struct {
    Name string
    Url string
    IsUp bool
}


func NewService() *Service {
    return &Service{
        client: Connect(),
    }
}

func (s *Service) GetServiceStatus() []ServiceStatus {
    serviceStatus := []ServiceStatus{}
    containers := s.client.ListContainers()
    serviceConfigs := GetServiceConfigs()
    urlByName := map[string]string{}
    for _, config := range(serviceConfigs) {
        urlByName[config.Name] = config.Url
    }
    for _, service := range(SERVICE_LIST) {
        containerNames := containerNamesByService[service]
        if len(containerNames) == 0 {
            continue
        }
        serviceContainers := s.client.FindContainersByName(containers, containerNames)
        isUp := s.client.AreAllContainersUp(serviceContainers)
        serviceStatus = append(serviceStatus, ServiceStatus{
            Name: service,
            Url: urlByName[service],
            IsUp: isUp,
        })
    }
    return serviceStatus
}

func (s *Service) StartService(name string) bool {
    containerNames := containerNamesByService[name]
    containers := s.client.ListContainers()
    serviceContainers := s.client.FindContainersByName(containers, containerNames)
    return s.client.StartContainers(serviceContainers)
}

func (s *Service) StopService(name string) bool {
    containerNames := containerNamesByService[name]
    containers := s.client.ListContainers()
    serviceContainers := s.client.FindContainersByName(containers, containerNames)
    return s.client.StopContainers(serviceContainers)

}
type ServiceConfig struct {
    Name string
    Url string
}
func GetServiceConfigs() []ServiceConfig {
    serviceConfigs := []ServiceConfig{}
    for _, service := range(SERVICE_LIST) {
      serviceConfigs = append(serviceConfigs, ServiceConfig{
          Name: service,
          Url: os.Getenv("URL_" + strings.ToUpper(service)),
      })
    }
    return serviceConfigs
}
