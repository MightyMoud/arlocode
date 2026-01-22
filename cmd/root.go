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
	"github.com/mightymoud/arlocode/internal/bridge"
	"github.com/mightymoud/arlocode/internal/butler/memory"
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

var systemPrompt = `You are ArloCode, an AI coding agent designed to assist developers with complex coding tasks.
You have access to various tools to help you accomplish your goals.
You should think carefully about which tools to use and when to use them.
You should also provide reasoning for your actions to help the user understand your thought process.

When you need to use a tool, you must specify the tool name and provide the necessary arguments in JSON format.
After executing a tool, you should analyze the results and decide on the next steps.

Always aim to provide clear and concise explanations for your actions and decisions.
Your ultimate goal is to assist the user in completing their coding tasks effectively and efficiently.
`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "arlocode",
	Short: "ArloCode a Coding Agent focused on long running tasks and local models",
	Long:  `AlroCode is an AI coding agent designed to assist developers with complex coding tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		codingAgent := coding_agent.Agent.WithMaxIterations(30).WithMemory([]memory.MemoryEntry{
			{
				Role:    "system",
				Message: systemPrompt,
			},
		})
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

	// Set version template for --version flag
	rootCmd.Version = version
	rootCmd.SetVersionTemplate(fmt.Sprintf("arlocode %s (commit: %s, built: %s)\n", version, commit, date))
}
