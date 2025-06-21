/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ei-sugimoto/wakemae/internal/dns"
	"github.com/ei-sugimoto/wakemae/internal/docker"
	"github.com/ei-sugimoto/wakemae/internal/registry"
	"github.com/ei-sugimoto/wakemae/internal/yaml"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func run() {
	log.Println("Starting wakemae...")

	cfg, err := yaml.ParseFile("/etc/wakemae/config.yml")
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	rg := registry.NewRegistry()
	go func() {
		if err := docker.Listen(rg); err != nil {
			log.Printf("Failed to start docker listener: %v", err)
		}
	}()

	go func() {
		timeout, err := time.ParseDuration(cfg.DNS.Timeout)
		if err != nil {
			log.Printf("Invalid timeout format, using default (5s): %v", err)
			timeout = 5 * time.Second
		}
		if err := dns.Serve(cfg.DNS.BindAddress, rg, cfg.DNS.Upstream, timeout); err != nil {
			log.Printf("Failed to start DNS server on port 53: %v", err)
		}
	}()

	log.Println("wakemae is running...")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
