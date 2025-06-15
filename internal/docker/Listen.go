package docker

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/ei-sugimoto/wakemae/internal/registry"
	"github.com/pkg/errors"
)

func Listen(reg *registry.Registry) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return errors.Wrap(err, "failed to create docker client")
	}

	ctx := context.Background()

	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return errors.Wrap(err, "failed to list containers")
	}

	for _, c := range containers {
		if err := upsert(ctx, cli, c.ID, reg); err != nil {
			return errors.Wrap(err, "failed to upsert container")
		}
	}

	f := filters.NewArgs()
	f.Add("type", "container")
	f.Add("event", "start")
	f.Add("event", "die")
	f.Add("event", "destory")

	msgs, errs := cli.Events(ctx, events.ListOptions{Filters: f})

	for {
		select {
		case msg := <-msgs:
			// Handle container events
			switch msg.Action {
			case "start":
				if err := upsert(ctx, cli, msg.Actor.ID, reg); err != nil {
					return errors.Wrap(err, "failed to handle container event")
				}
			}
		case err := <-errs:
			if err != nil {
				return errors.Wrap(err, "event stream error")
			}
		}
	}
}

func upsert(ctx context.Context, cli *client.Client, id string, reg *registry.Registry) error {
	inspect, err := cli.ContainerInspect(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to inspect container")
	}

	lbl := inspect.Config.Labels
	name := lbl["dns.name"]
	if name == "" {
		return nil
	}

	domain := lbl["dns.domain"]
	if domain == "" {
		domain = "local"
	}
	fqdn := name + "." + domain

	ip := firstID(inspect)

	if ip != "" {
		reg.AddA(fqdn, ip)
	}

	return nil
}

func firstID(inspect container.InspectResponse) string {
	for _, n := range inspect.NetworkSettings.Networks {
		if n.IPAddress != "" {
			return n.IPAddress
		}
	}
	return ""
}
