package docker

import (
	"context"
	"log"

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
	f.Add("event", "destroy")

	msgs, errs := cli.Events(ctx, events.ListOptions{Filters: f})

	for {
		select {
		case msg := <-msgs:
			log.Println(msg)
			// Handle container events
			switch msg.Action {
			case "start":
				if err := upsert(ctx, cli, msg.Actor.ID, reg); err != nil {
					return errors.Wrap(err, "failed to handle container event")
				}
			case "die", "destroy":
				if err := remove(ctx, cli, msg.Actor.ID, reg); err != nil {
					log.Printf("failed to remove container from registry: %v", err)
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
	fqdn := lbl["wakemae.domain"]
	if fqdn == "" {
		return nil
	}

	ip := firstID(inspect)

	if ip != "" {
		reg.AddA(fqdn, ip)
		reg.AddContainer(id, fqdn)
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

func remove(ctx context.Context, cli *client.Client, id string, reg *registry.Registry) error {
	fqdn, exists := reg.RemoveContainer(id)
	if !exists {
		return nil
	}

	// Try to get IP from inspect, but if it fails, we'll remove all records for this fqdn
	inspect, err := cli.ContainerInspect(ctx, id)
	if err != nil {
		// Container might already be deleted, remove all A records for this fqdn
		reg.Del(fqdn)
		return nil
	}

	ip := firstID(inspect)
	if ip != "" {
		reg.RemoveA(fqdn, ip)
	} else {
		reg.Del(fqdn)
	}

	return nil
}
