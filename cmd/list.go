package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/jacobwgillespie/key-light/pkg/lights"
	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List available key lights on the network",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		expected, _ := cmd.Flags().GetInt("expected")

		timeout := cmd.Flag("timeout").Value.String()
		timeoutDuration, err := time.ParseDuration(timeout)
		if err != nil {
			return err
		}

		discoveryCtx, cancelFn := context.WithTimeout(context.Background(), timeoutDuration)
		defer cancelFn()

		results, err := lights.DiscoverLights(discoveryCtx, timeoutDuration)
		if err != nil {
			return err
		}

		count := 0
		for device := range results {

			info, err := device.FetchDeviceInfo(context.Background())
			if err != nil {
				continue
			}
			group, err := device.FetchLightGroup(context.Background())
			if err != nil {
				continue
			}

			fmt.Printf("\n%s\n", aurora.Magenta(aurora.Bold(info.DisplayName)))
			fmt.Printf("%s %s\n", aurora.Faint("Kind       "), info.ProductName)
			fmt.Printf("%s %s\n", aurora.Faint("Serial     "), info.SerialNumber)
			fmt.Printf("%s %s\n", aurora.Faint("Firmware   "), info.FirmwareVersion)

			for _, light := range group.Lights {
				on := aurora.Green("ON")
				if light.On == 0 {
					on = aurora.Red("OFF")
				}

				fmt.Printf("%s %s\n", aurora.Faint("State      "), aurora.Bold(on))
				fmt.Printf("%s %d%%\n", aurora.Faint("Brightness "), light.Brightness)
				fmt.Printf("%s %d\n", aurora.Faint("Temp       "), light.Temperature)
			}

			count++

			if expected > 0 && count >= expected {
				break
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("timeout", "t", "5s", "Timeout for discovery")
	toggleCmd.Flags().Int("expected", 0, "Expected number of lights to discover")
}
