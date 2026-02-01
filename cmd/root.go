/*
Copyright © 2026 Mahmoud Mousa <m.mousa@hey.com>

Licensed under the GNU GPL License, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
https://www.gnu.org/licenses/gpl-3.0.en.html

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mightymoud/arlocode/internal/bridge"
	"github.com/mightymoud/arlocode/internal/coding_agent"
	"github.com/mightymoud/arlocode/internal/tui/app"
	"github.com/spf13/cobra"
)

// Version information - set via ldflags at build time
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var (
	instructions string
	outputPath   string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "arlocode",
	Short: "ArloCode a Coding Agent focused on long running tasks and local models",
	Long:  `ArloCode is an AI coding agent designed to assist developers with complex coding tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		if instructions != "" {
			if err := runHeadless(instructions, outputPath); err != nil {
				os.Exit(1)
			}
			return
		}

		codingAgent := coding_agent.Agent
		agentBridge := bridge.NewDirectBridge(codingAgent)

		// Create the app model using the new constructor
		m := app.NewAppModel(agentBridge)

		p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

		go func() {
			for event := range agentBridge.Events() {
				switch e := event.(type) {
				case bridge.TextChunkEvent:
					p.Send(app.AgentTextChunkMsg(e.Text))
				case bridge.TextCompleteEvent:
					p.Send(app.AgentTextCompleteMsg(""))
				case bridge.ThinkingChunkEvent:
					p.Send(app.AgentThinkingChunkMsg(e.Text))
				case bridge.ThinkingCompleteEvent:
					p.Send(app.AgentThinkingCompleteMsg(""))
				case bridge.ToolCallEvent:
					p.Send(app.ToolCallMsg(e.ToolCall))
				case bridge.ErrorEvent:
					// Handle error event if needed
				case bridge.TurnCompleteEvent:
					// Handle turn complete event if needed
				}
			}
		}()
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error: %v", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.Flags().StringVarP(&instructions, "instructions", "i", "", "Run in headless mode with the provided instruction text")
	rootCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Path to write the ATIF trajectory JSON")

	// Set version template for --version flag
	rootCmd.Version = version
	rootCmd.SetVersionTemplate(fmt.Sprintf("arlocode %s (commit: %s, built: %s)\n", version, commit, date))
}

func runHeadless(prompt string, path string) error {
	codingAgent := coding_agent.Agent.WithMaxIterations(500)
	agentBridge := bridge.NewDirectBridge(codingAgent)
	defer agentBridge.Close()

	if err := agentBridge.Run(context.Background(), prompt); err != nil {
		return err
	}

	var runErr error
	for {
		event, ok := <-agentBridge.Events()
		if !ok {
			break
		}
		switch e := event.(type) {
		case bridge.ErrorEvent:
			if e.Err != nil {
				runErr = e.Err
			}
		case bridge.TurnCompleteEvent:
			_, err := agentBridge.ExportATIF(path)
			if err != nil {
				return err
			}
			if runErr != nil {
				return runErr
			}
			return nil
		}
	}

	return runErr
}
