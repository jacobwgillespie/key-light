package cmd

import (
	"context"
	"time"

	"github.com/jacobwgillespie/key-light/pkg/lights"
	"github.com/spf13/cobra"
)

var toggleCmd = &cobra.Command{
	Use:   "toggle",
	Short: "Toggle lights on/off",
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

		currentState := -1
		count := 0

		for device := range results {
			lg, err := device.FetchLightGroup(context.Background())
			if err != nil {
				continue
			}

			for _, light := range lg.Lights {
				if currentState == -1 {
					currentState = light.On
				}

				desiredState := 0
				if currentState == 0 {
					desiredState = 1
				}

				light.On = desiredState
			}

			device.UpdateLightGroup(context.Background(), lg)

			count++

			if expected > 0 && count >= expected {
				break
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(toggleCmd)
	toggleCmd.Flags().StringP("timeout", "t", "5s", "Timeout for discovery")
	toggleCmd.Flags().Int("expected", 0, "Expected number of lights to discover")
}
