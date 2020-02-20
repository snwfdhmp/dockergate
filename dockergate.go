package dockergate

import (
	"errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type Gate struct {
	Endpoints    map[string]Endpoint
	DockerClient *client.Client
}

func NewGate() *Gate {
	cli, err := client.NewEnvClient()
	if err != nil {
		logrus.Fatalf("cannot create docker client: %s", err)
	}

	endpoints := make(map[string]Endpoint)
	return &Gate{
		Endpoints:    endpoints,
		DockerClient: cli,
	}
}

func (g *Gate) Install(imageName, endpoint, port string) error {
	if _, ok := g.Endpoints[endpoint]; ok {
		return errors.New("endpoint already exists")
	}

	ctx := context.Background()
	logrus.Infof("%s: endpoint created for image %s", endpoint, imageName)

	resp, err := g.DockerClient.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, "")
	if err != nil {
		return err
	}

	if err := g.DockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	logrus.Infof("%s: started container %s (exposing %d)", endpoint, resp.ID, port)

	g.Endpoints[endpoint] = Endpoint{
		Name:        endpoint,
		Port:        port,
		ContainerID: resp.ID,
	}

	logrus.Infof("%s: endpoint successfully added", endpoint)
}

type Endpoint struct {
	Name        string
	ContainerID string
	Port        string
}
