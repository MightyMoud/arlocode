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
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mightymoud/arlocode/internal/coding_agent"
	state "github.com/mightymoud/arlocode/internal/tui"
	"github.com/mightymoud/arlocode/internal/tui/app"
	"github.com/spf13/cobra"
)

// Version information - set via ldflags at build time
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "arlocode",
	Short: "ArloCode a Coding Agent focused on long running tasks and local models",
	Long:  `AlroCode is an AI coding agent designed to assist developers with complex coding tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		appState := state.Get()

		// Create the app model using the new constructor
		m := app.NewAppModel()

		codingAgent := coding_agent.Agent.WithMaxIterations(10).
			WithOnThinkingChunk(func(s string) {
				appState.Program().Send(app.AgentThinkingChunkMsg(s))
			}).
			WithOnThinkingComplete(func() {
				appState.Program().Send(app.AgentThinkingCompleteMsg(""))
			}).
			WithOnTextChunk(func(s string) {
				appState.Program().Send(app.AgentTextChunkMsg(s))
			}).
			WithOnStreamComplete(func() {
				appState.Program().Send(app.AgentTextCompleteMsg(""))
			})

		appState.SetAgent(codingAgent)

		p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
		appState.SetProgram(p)
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

	// Set version template for --version flag
	rootCmd.Version = version
	rootCmd.SetVersionTemplate(fmt.Sprintf("arlocode %s (commit: %s, built: %s)\n", version, commit, date))
}
