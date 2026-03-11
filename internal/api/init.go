package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/swlee3306/ai-company-os/internal/audit"
	"github.com/swlee3306/ai-company-os/internal/store"
)

type initSettings struct {
	Driver struct {
		Selected string `json:"selected"`
	} `json:"driver"`
	Approval struct {
		PolicyText string `json:"policy_text"`
	} `json:"approval"`
	Runner struct {
		Backend string            `json:"backend"`
		Type    string            `json:"type"`
		Command string            `json:"command"`
		Workdir string            `json:"workdir"`
		Agents  map[string]string `json:"agents"`
	} `json:"runner"`
}

func init() {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize settings and secrets for running the OS (wizard)",
		RunE: func(cmd *cobra.Command, args []string) error {
			st := store.NewFileStore(defaultDataDir())
			au := audit.NewFileAudit(st)

			in := bufio.NewReader(os.Stdin)
			fmt.Println("AI Company OS init")
			fmt.Println("- Creates settings.json and secrets.json under:")
			fmt.Println("  ", defaultDataDir())

			backend := prompt(in, "Runner backend (local_cli/openclaw_acp)", "local_cli")
			if backend != "local_cli" && backend != "openclaw_acp" {
				return fmt.Errorf("invalid backend")
			}

			s := initSettings{}
			s.Driver.Selected = "k3d"
			s.Approval.PolicyText = "HIGH: prod deploy / secrets / cluster admin requires approval."
			s.Runner.Backend = backend
			s.Runner.Type = "codex_cli"
			s.Runner.Command = prompt(in, "Runner command (e.g., codex)", "codex")
			s.Runner.Workdir = ""
			s.Runner.Agents = map[string]string{}

			secrets := map[string]string{}

			if backend == "openclaw_acp" {
				s.Runner.Agents["pm"] = prompt(in, "PM agentId (ACP)", "")
				secrets["OPENCLAW_GATEWAY_URL"] = prompt(in, "OPENCLAW_GATEWAY_URL", "http://127.0.0.1:18789")
				secrets["OPENCLAW_GATEWAY_TOKEN"] = prompt(in, "OPENCLAW_GATEWAY_TOKEN", "")
			} else {
				// local_cli secrets for providers (optional)
				secrets["OPENAI_API_KEY"] = prompt(in, "OPENAI_API_KEY (optional)", "")
			}

			// write settings
			b, _ := json.MarshalIndent(s, "", "  ")
			if err := st.WriteSettings(b); err != nil {
				return err
			}

			// write secrets (0600)
			secOut := map[string]string{}
			for k, v := range secrets {
				if strings.TrimSpace(v) == "" {
					continue
				}
				secOut[k] = v
			}
			secB, _ := json.MarshalIndent(secOut, "", "  ")
			if err := st.WriteSecrets(secB); err != nil {
				return err
			}

			au.Emit("cli", "init", map[string]any{"runner_backend": backend})
			fmt.Println("Init complete.")
			fmt.Println("- settings.json written")
			fmt.Println("- secrets.json written (0600)")
			return nil
		},
	}
	rootCmd.AddCommand(cmd)
}

func prompt(in *bufio.Reader, label string, def string) string {
	if def != "" {
		fmt.Printf("%s [%s]: ", label, def)
	} else {
		fmt.Printf("%s: ", label)
	}
	line, _ := in.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return def
	}
	return line
}
