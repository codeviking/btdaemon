package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var ErrNoSSID = errors.New("no SSID")

// TODO: Make these paths configurable
const airport = "/System/Library/PrivateFrameworks/Apple80211.framework/Resources/airport"
const blueutil = "/opt/homebrew/bin/blueutil"

func ssid() (string, error) {
	bin := exec.Command(airport, "-I")
	stdout, err := bin.Output()
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(string(stdout), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}
		if strings.TrimSpace(parts[0]) == "SSID" {
			return strings.TrimSpace(parts[1]), nil
		}
	}

	return "", ErrNoSSID
}

func toggleBluetooth(on bool) error {
	var bit int
	if on {
		bit = 1
	}
	return exec.Command(blueutil, "--power", strconv.Itoa(bit)).Run()
}

func main() {
	cmd := cobra.Command{
		Use:   "btdaemon <ssid...>",
		Short: "daemon that safely enables & disables bluetooth based on your wireless SSID",
		RunE: func(cmd *cobra.Command, args []string) error {
			lp := "/var/log/net.codeviking.btdaemon/stdout.log"
			if err := os.MkdirAll(filepath.Dir(lp), 0755); err != nil {
				return err
			}
			lf, err := os.OpenFile(lp, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
			if err != nil {
				return err
			}
			logger := log.New(lf, "", log.Ldate|log.Ltime|log.LUTC)

			ssids := args
			if len(ssids) == 0 {
				err := errors.New("no trusted SSIDs, nothing to do")
				logger.Print(err)
				return nil
			}

			logger.Printf("bluetooth will be enabled when connected to SSIDs: %s\n", strings.Join(ssids, ", "))

			t := time.NewTicker(15 * time.Second)
		NEXTTICK:
			for {
				select {
				case <-t.C:
					logger.Println("querying ssid...")
					current, err := ssid()
					if err != nil {
						logger.Printf("error querying ssid: %s\n", err.Error())
						continue NEXTTICK
					}

					for _, trusted := range ssids {
						if trusted == current {
							logger.Printf("%s is trusted, enabling bluetooth...\n", current)
							if err := toggleBluetooth(true); err != nil {
								logger.Printf("error enabling bluetooth: %s\n", err.Error())
							}
							continue NEXTTICK
						}
					}

					logger.Printf("%s is not trusted, disabling bluetooth...\n", current)
					if err := toggleBluetooth(false); err != nil {
						logger.Printf("error disabling bluetooth: %s\n", err.Error())
					}
				}
			}
		},
	}

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
