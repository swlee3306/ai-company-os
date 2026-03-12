package server

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/swlee3306/ai-company-os/internal/audit"
	"github.com/swlee3306/ai-company-os/internal/model"
	"github.com/swlee3306/ai-company-os/internal/store"
)

type runRequest struct {
	TaskID     string `json:"task_id"`
	RunnerType string `json:"runner_type"`
	Pipeline   string `json:"pipeline"`
}

type runDetail struct {
	Run model.Run `json:"run"`
}

func registerRunRoutes(api *gin.RouterGroup, st *store.FileStore, au *audit.FileAudit) {
	api.POST("/tasks/:id/run", func(c *gin.Context) {
		taskID := c.Param("id")
		var body struct {
			Pipeline string `json:"pipeline"`
		}
		_ = c.ShouldBindJSON(&body)
		if body.Pipeline == "" {
			body.Pipeline = "pm_only"
		}

		// read runner settings (best-effort)
		runnerType := "codex_cli"
		runnerBackend := "local_cli"
		pmAgent := ""
		if b, err := st.ReadSettings(); err == nil {
			var s map[string]any
			if json.Unmarshal(b, &s) == nil {
				if rset, ok := s["runner"].(map[string]any); ok {
					if t, ok := rset["type"].(string); ok && t != "" {
						runnerType = t
					}
					if b, ok := rset["backend"].(string); ok && b != "" {
						runnerBackend = b
					}
					if agents, ok := rset["agents"].(map[string]any); ok {
						if s, ok := agents["pm"].(string); ok {
							pmAgent = s
						}
					}
				}
			}
		}

		runID := fmtID("run")
		now := time.Now().UTC().Format(time.RFC3339)
		r := model.Run{ID: runID, TaskID: taskID, RunnerType: runnerType, Pipeline: body.Pipeline, Status: "running", StartedAt: now}

		if err := st.EnsureRunDir(runID); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		dir := st.RunDir(runID)

		// persist request.json
		rq := runRequest{TaskID: taskID, RunnerType: runnerType, Pipeline: body.Pipeline}
		rqb, _ := json.MarshalIndent(rq, "", "  ")
		_ = os.WriteFile(filepath.Join(dir, "request.json"), rqb, 0o644)

		// minimal logs / result placeholders
		_ = os.WriteFile(filepath.Join(dir, "prompt.txt"), []byte("(runner mvp)\n"), 0o644)
		_ = os.WriteFile(filepath.Join(dir, "stdout.log"), []byte("(runner mvp) started\n"), 0o644)
		_ = os.WriteFile(filepath.Join(dir, "stderr.log"), []byte(""), 0o644)

		au.Emit("api", "run.start", map[string]any{"run_id": runID, "task_id": taskID, "runner": runnerType, "pipeline": body.Pipeline})

		if runnerBackend == "openclaw_acp" {
			if pmAgent == "" {
				c.JSON(400, gin.H{"error": "runner.agents.pm required for openclaw_acp"})
				return
			}

			// Read task content to feed into the PM step
			tb, err := st.ReadTasks()
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			var tasks []map[string]any
			_ = json.Unmarshal(tb, &tasks)
			title := ""
			desc := ""
			for _, t := range tasks {
				if t["id"] == taskID {
					title, _ = t["title"].(string)
					desc, _ = t["desc"].(string)
					break
				}
			}

			prompt := "You are PM. Turn this Task into acceptance criteria and execution plan.\n" +
				"Task ID: " + taskID + "\n" +
				"Title: " + title + "\n" +
				"Desc: " + desc + "\n" +
				"Constraints: branch/PR only, QA batch once, audit-first.\n" +
				"Output: a concise plan + next agent steps.\n"
			_ = os.WriteFile(filepath.Join(dir, "prompt.txt"), []byte(prompt), 0o644)

			// Spawn ACP session via OpenClaw Tools Invoke
			args := map[string]any{
				"task":    prompt,
				"runtime": "acp",
				"agentId": pmAgent,
				"thread":  true,
				"mode":    "session",
				"label":   "aicos:" + runID + ":pm",
			}
			res, err := openclawInvoke("sessions_spawn", args)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			_ = os.WriteFile(filepath.Join(dir, "openclaw_spawn.json"), res, 0o644)
			_ = os.WriteFile(filepath.Join(dir, "stdout.log"), []byte("(openclaw_acp) spawned\n"), 0o644)

			// Mark done for MVP (we don't stream logs yet)
			r.Status = "done"
			r.EndedAt = time.Now().UTC().Format(time.RFC3339)
			r.Summary = "openclaw_acp: spawned PM session (log streaming next)"
			resb, _ := json.MarshalIndent(map[string]any{"status": "ok", "summary": r.Summary}, "", "  ")
			_ = os.WriteFile(filepath.Join(dir, "RESULT.json"), resb, 0o644)
			au.Emit("system", "run.done", map[string]any{"run_id": runID, "task_id": taskID, "runner": runnerType, "backend": runnerBackend})
		} else if runnerBackend == "local_cli" {
			// E2.3 local_cli: run a non-interactive command once and capture output.
			cmdStr := ""
			if b, err := st.ReadSettings(); err == nil {
				var s map[string]any
				if json.Unmarshal(b, &s) == nil {
					if rset, ok := s["runner"].(map[string]any); ok {
						if c, ok := rset["command"].(string); ok {
							cmdStr = c
						}
					}
				}
			}
			if cmdStr == "" {
				cmdStr = "codex"
			}

			// Build a simple PM prompt (non-interactive)
			repoPath := ""
			if b, err := st.ReadSettings(); err == nil {
				var s map[string]any
				if json.Unmarshal(b, &s) == nil {
					if ws, ok := s["workspace"].(map[string]any); ok {
						if rp, ok := ws["repo_path"].(string); ok {
							repoPath = rp
						}
					}
				}
			}
			prompt := "You are PM. Provide acceptance criteria and an execution plan.\n" +
				"Task ID: " + taskID + "\n" +
				"Pipeline: " + body.Pipeline + "\n" +
				"Repo: " + repoPath + "\n" +
				"Constraints: branch/PR only, QA batch once, audit-first.\n"
			_ = os.WriteFile(filepath.Join(dir, "prompt.txt"), []byte(prompt), 0o644)

			au.Emit("system", "run.step.start", map[string]any{"run_id": runID, "task_id": taskID, "role": "pm", "backend": runnerBackend})

			// For portability: try `cmd --help` first (exists check), then attempt stdin prompt.
			execCmd := exec.Command(cmdStr, "--help")
			helpOut, helpErr := execCmd.CombinedOutput()
			if helpErr != nil {
				_ = os.WriteFile(filepath.Join(dir, "stderr.log"), helpOut, 0o644)
				r.Status = "failed"
				r.EndedAt = time.Now().UTC().Format(time.RFC3339)
				r.Summary = "local_cli failed: command not runnable: " + helpErr.Error()
				resb, _ := json.MarshalIndent(map[string]any{"status": "fail", "summary": r.Summary}, "", "  ")
				_ = os.WriteFile(filepath.Join(dir, "RESULT.json"), resb, 0o644)
				au.Emit("system", "run.fail", map[string]any{"run_id": runID, "task_id": taskID, "runner": runnerType, "backend": runnerBackend})
				c.JSON(500, gin.H{"error": r.Summary})
				return
			}

			// If runner is codex, use `codex exec` non-interactively.
			if filepath.Base(cmdStr) == "codex" {
				if repoPath == "" {
					c.JSON(400, gin.H{"error": "workspace.repo_path required for codex"})
					return
				}

				// Step: PM
				pmOut, pmErr := exec.Command(cmdStr, "exec", "-C", repoPath, "-s", "workspace-write", prompt).CombinedOutput()
				_ = os.WriteFile(filepath.Join(dir, "pm.stdout.log"), pmOut, 0o644)
				if pmErr != nil {
					_ = os.WriteFile(filepath.Join(dir, "stdout.log"), pmOut, 0o644)
					r.Status = "failed"
					r.EndedAt = time.Now().UTC().Format(time.RFC3339)
					r.Summary = "pm step failed: " + pmErr.Error()
					resb, _ := json.MarshalIndent(map[string]any{"status": "fail", "summary": r.Summary}, "", "  ")
					_ = os.WriteFile(filepath.Join(dir, "RESULT.json"), resb, 0o644)
					au.Emit("system", "run.step.done", map[string]any{"run_id": runID, "task_id": taskID, "role": "pm", "status": "failed"})
					au.Emit("system", "run.fail", map[string]any{"run_id": runID, "task_id": taskID, "runner": runnerType, "backend": runnerBackend})
					c.JSON(500, gin.H{"error": r.Summary})
					return
				}
				_ = os.WriteFile(filepath.Join(dir, "stdout.log"), pmOut, 0o644)
				au.Emit("system", "run.step.done", map[string]any{"run_id": runID, "task_id": taskID, "role": "pm", "status": "ok"})

				if body.Pipeline == "full" {
					// Step: BE (create branch + apply patch)
					branch := "task-" + taskID + "-" + runID
					for _, ch := range []string{":", "/", " ", "\\"} {
						branch = strings.ReplaceAll(branch, ch, "-")
					}
					au.Emit("system", "run.step.start", map[string]any{"run_id": runID, "task_id": taskID, "role": "be"})
					co := exec.Command("git", "checkout", "-B", branch, "main")
					co.Dir = repoPath
					coOut, coErr := co.CombinedOutput()
					_ = os.WriteFile(filepath.Join(dir, "be.git.log"), coOut, 0o644)

					// Ensure clean working tree before applying any patch.
					clean := exec.Command("git", "reset", "--hard")
					clean.Dir = repoPath
					_, _ = clean.CombinedOutput()
					cl2 := exec.Command("git", "clean", "-fd")
					cl2.Dir = repoPath
					_, _ = cl2.CombinedOutput()

					if coErr != nil {
						r.Status = "failed"
						r.EndedAt = time.Now().UTC().Format(time.RFC3339)
						r.Summary = "be step failed: git checkout -b: " + coErr.Error()
						resb, _ := json.MarshalIndent(map[string]any{"status": "fail", "summary": r.Summary}, "", "  ")
						_ = os.WriteFile(filepath.Join(dir, "RESULT.json"), resb, 0o644)
						au.Emit("system", "run.fail", map[string]any{"run_id": runID, "task_id": taskID})
						c.JSON(500, gin.H{"error": r.Summary})
						return
					}

					bePrompt := "You are BE. Make the minimal change to satisfy the Task.\n" +
						"SCOPE LIMIT (MVP): ONLY modify README.md. Do not change any other files.\n" +
						"OUTPUT FORMAT: Output ONLY a unified diff. No commentary, no code fences, no 'file update:' lines.\n" +
						"The diff MUST start with: diff --git a/README.md b/README.md\n" +
						prompt
					beOut, beErr := exec.Command(cmdStr, "exec", "-C", repoPath, "-s", "workspace-write", bePrompt).CombinedOutput()
					_ = os.WriteFile(filepath.Join(dir, "be.raw.log"), beOut, 0o644)
					if beErr != nil {
						r.Status = "failed"
						r.EndedAt = time.Now().UTC().Format(time.RFC3339)
						r.Summary = "be step failed: " + beErr.Error()
						resb, _ := json.MarshalIndent(map[string]any{"status": "fail", "summary": r.Summary}, "", "  ")
						_ = os.WriteFile(filepath.Join(dir, "RESULT.json"), resb, 0o644)
						au.Emit("system", "run.fail", map[string]any{"run_id": runID, "task_id": taskID})
						c.JSON(500, gin.H{"error": r.Summary})
						return
					}
					patch := extractUnifiedDiff(beOut)
					if len(patch) == 0 {
						r.Status = "failed"
						r.EndedAt = time.Now().UTC().Format(time.RFC3339)
						r.Summary = "be step failed: no unified diff found in output"
						resb, _ := json.MarshalIndent(map[string]any{"status": "fail", "summary": r.Summary}, "", "  ")
						_ = os.WriteFile(filepath.Join(dir, "RESULT.json"), resb, 0o644)
						c.JSON(500, gin.H{"error": r.Summary})
						return
					}
					patch = sanitizePatch(patch)
					_ = os.WriteFile(filepath.Join(dir, "be.patch.diff"), patch, 0o644)

					// codex exec in workspace-write mode may have modified files already.
					// Reset to a clean base before applying the extracted patch.
					reset2 := exec.Command("git", "reset", "--hard")
					reset2.Dir = repoPath
					_, _ = reset2.CombinedOutput()

					ap := exec.Command("git", "apply", "--3way", "--whitespace=nowarn", filepath.Join(dir, "be.patch.diff"))
					ap.Dir = repoPath
					apOut, apErr := ap.CombinedOutput()
					_ = os.WriteFile(filepath.Join(dir, "be.apply.log"), apOut, 0o644)
					if apErr != nil {
						r.Status = "failed"
						r.EndedAt = time.Now().UTC().Format(time.RFC3339)
						r.Summary = "be step failed: git apply: " + apErr.Error()
						resb, _ := json.MarshalIndent(map[string]any{"status": "fail", "summary": r.Summary}, "", "  ")
						_ = os.WriteFile(filepath.Join(dir, "RESULT.json"), resb, 0o644)
						c.JSON(500, gin.H{"error": r.Summary})
						return
					}
					add := exec.Command("git", "add", "-A")
					add.Dir = repoPath
					_, _ = add.CombinedOutput()
					// Ensure git has a local identity for commits (avoid 'Author identity unknown').
					cfgName := exec.Command("git", "config", "user.name", "ai-company-os-bot")
					cfgName.Dir = repoPath
					_, _ = cfgName.CombinedOutput()
					cfgEmail := exec.Command("git", "config", "user.email", "ai-company-os-bot@localhost")
					cfgEmail.Dir = repoPath
					_, _ = cfgEmail.CombinedOutput()
					cm := exec.Command("git", "commit", "-m", "feat: "+taskID+" (auto)")
					cm.Dir = repoPath
					cmOut, cmErr := cm.CombinedOutput()
					_ = os.WriteFile(filepath.Join(dir, "be.commit.log"), cmOut, 0o644)
					if cmErr != nil {
						r.Status = "failed"
						r.EndedAt = time.Now().UTC().Format(time.RFC3339)
						r.Summary = "be step failed: git commit: " + cmErr.Error()
						resb, _ := json.MarshalIndent(map[string]any{"status": "fail", "summary": r.Summary}, "", "  ")
						_ = os.WriteFile(filepath.Join(dir, "RESULT.json"), resb, 0o644)
						c.JSON(500, gin.H{"error": r.Summary})
						return
					}
					au.Emit("system", "run.step.done", map[string]any{"run_id": runID, "task_id": taskID, "role": "be", "status": "ok"})

					// Step: PR
					pu := exec.Command("git", "push", "-u", "origin", branch)
					pu.Dir = repoPath
					puOut, puErr := pu.CombinedOutput()
					_ = os.WriteFile(filepath.Join(dir, "pr.push.log"), puOut, 0o644)
					if puErr != nil {
						r.Status = "failed"
						r.EndedAt = time.Now().UTC().Format(time.RFC3339)
						r.Summary = "pr step failed: git push: " + puErr.Error()
						resb, _ := json.MarshalIndent(map[string]any{"status": "fail", "summary": r.Summary}, "", "  ")
						_ = os.WriteFile(filepath.Join(dir, "RESULT.json"), resb, 0o644)
						c.JSON(500, gin.H{"error": r.Summary})
						return
					}
					pr := exec.Command("gh", "pr", "create", "--base", "main", "--head", branch, "--title", "Auto: "+taskID, "--body", "Created by AI Company OS")
					pr.Dir = repoPath
					prOut, prErr := pr.CombinedOutput()
					_ = os.WriteFile(filepath.Join(dir, "pr.stdout.log"), prOut, 0o644)
					if prErr != nil {
						r.Status = "failed"
						r.EndedAt = time.Now().UTC().Format(time.RFC3339)
						r.Summary = "pr step failed: gh pr create: " + prErr.Error()
						resb, _ := json.MarshalIndent(map[string]any{"status": "fail", "summary": r.Summary}, "", "  ")
						_ = os.WriteFile(filepath.Join(dir, "RESULT.json"), resb, 0o644)
						c.JSON(500, gin.H{"error": r.Summary})
						return
					}
					prURL := strings.TrimSpace(string(prOut))
					_ = appendArtifactJSON(st, au, mustJSON(map[string]any{"type": "pr_link", "title": "PR for " + taskID, "uri": prURL, "task_id": taskID, "meta": map[string]any{"run_id": runID}}))

					// Step: QA
					// Ensure web/ dependencies exist (tsc, vite, etc.) before running web build.
					qa := exec.Command("bash", "-lc", "go test ./... && npm -C web ci && npm -C web run build")
					qa.Dir = repoPath
					qaOut, qaErr := qa.CombinedOutput()
					_ = os.WriteFile(filepath.Join(dir, "qa.log"), qaOut, 0o644)
					_ = appendArtifactJSON(st, au, mustJSON(map[string]any{"type": "qa_log", "title": "QA log " + runID, "uri": "file://" + filepath.Join(dir, "qa.log"), "task_id": taskID, "meta": map[string]any{"run_id": runID}}))
					if qaErr != nil {
						r.Status = "failed"
						r.EndedAt = time.Now().UTC().Format(time.RFC3339)
						r.Summary = "qa failed"
						resb, _ := json.MarshalIndent(map[string]any{"status": "fail", "summary": r.Summary, "pr": prURL}, "", "  ")
						_ = os.WriteFile(filepath.Join(dir, "RESULT.json"), resb, 0o644)
						c.JSON(500, gin.H{"error": r.Summary, "pr": prURL})
						return
					}

					// Done
					r.Status = "done"
					r.EndedAt = time.Now().UTC().Format(time.RFC3339)
					r.Summary = "full pipeline complete"
					resb, _ := json.MarshalIndent(map[string]any{"status": "ok", "summary": r.Summary, "pr": prURL}, "", "  ")
					_ = os.WriteFile(filepath.Join(dir, "RESULT.json"), resb, 0o644)
					au.Emit("system", "run.done", map[string]any{"run_id": runID, "task_id": taskID, "backend": runnerBackend})
					c.JSON(201, runDetail{Run: r})
					return
				}

				// pm_only done
				r.Status = "done"
				r.EndedAt = time.Now().UTC().Format(time.RFC3339)
				r.Summary = "pm_only complete"
				resb, _ := json.MarshalIndent(map[string]any{"status": "ok", "summary": r.Summary}, "", "  ")
				_ = os.WriteFile(filepath.Join(dir, "RESULT.json"), resb, 0o644)
				au.Emit("system", "run.done", map[string]any{"run_id": runID, "task_id": taskID, "backend": runnerBackend})
			} else {
				// Attempt to run the CLI once with prompt via stdin (generic fallback)
				cmd := exec.Command(cmdStr)
				cmd.Stdin = bytes.NewReader([]byte(prompt))
				outB, err := cmd.CombinedOutput()
				if err != nil {
					_ = os.WriteFile(filepath.Join(dir, "stdout.log"), outB, 0o644)
					r.Status = "failed"
					r.EndedAt = time.Now().UTC().Format(time.RFC3339)
					r.Summary = "local_cli failed: " + err.Error()
					resb, _ := json.MarshalIndent(map[string]any{"status": "fail", "summary": r.Summary}, "", "  ")
					_ = os.WriteFile(filepath.Join(dir, "RESULT.json"), resb, 0o644)
					au.Emit("system", "run.step.done", map[string]any{"run_id": runID, "task_id": taskID, "role": "pm", "status": "failed"})
					au.Emit("system", "run.fail", map[string]any{"run_id": runID, "task_id": taskID, "runner": runnerType, "backend": runnerBackend})
					c.JSON(500, gin.H{"error": r.Summary})
					return
				}
				_ = os.WriteFile(filepath.Join(dir, "stdout.log"), outB, 0o644)
			}
		} else {
			// Fallback placeholder
			r.Status = "done"
			r.EndedAt = time.Now().UTC().Format(time.RFC3339)
			r.Summary = "runner placeholder"
			resb, _ := json.MarshalIndent(map[string]any{"status": "ok", "summary": r.Summary}, "", "  ")
			_ = os.WriteFile(filepath.Join(dir, "RESULT.json"), resb, 0o644)
			au.Emit("system", "run.done", map[string]any{"run_id": runID, "task_id": taskID, "runner": runnerType, "backend": runnerBackend})
		}

		// Create artifact pointing to stdout.log
		artifactBody := map[string]any{
			"type":    "run_log",
			"title":   "Run log " + runID,
			"uri":     "file://" + filepath.Join(dir, "stdout.log"),
			"task_id": taskID,
			"meta":    map[string]any{"run_id": runID},
		}
		ab, _ := json.Marshal(artifactBody)
		// Reuse artifact create logic via direct append (simple)
		_ = appendArtifactJSON(st, au, ab)

		c.JSON(201, runDetail{Run: r})
	})

	api.GET("/runs", func(c *gin.Context) {
		st.EnsureRunsDir()
		dirs, _ := os.ReadDir(st.RunsDir())
		out := []model.Run{}
		for _, d := range dirs {
			if !d.IsDir() {
				continue
			}
			b, err := os.ReadFile(filepath.Join(st.RunDir(d.Name()), "request.json"))
			if err != nil {
				continue
			}
			var rq runRequest
			if json.Unmarshal(b, &rq) != nil {
				continue
			}
			out = append(out, model.Run{ID: d.Name(), TaskID: rq.TaskID, RunnerType: rq.RunnerType, Pipeline: rq.Pipeline})
		}
		c.JSON(200, out)
	})

	api.GET("/runs/:id", func(c *gin.Context) {
		id := c.Param("id")
		dir := st.RunDir(id)
		b, err := os.ReadFile(filepath.Join(dir, "request.json"))
		if err != nil {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}
		var rq runRequest
		_ = json.Unmarshal(b, &rq)
		r := model.Run{ID: id, TaskID: rq.TaskID, RunnerType: rq.RunnerType, Pipeline: rq.Pipeline}
		c.JSON(200, runDetail{Run: r})
	})
}

// appendArtifactJSON appends an artifact to artifacts.json using the same shape as POST /api/artifacts.
// NOTE: This is an internal helper for the runner MVP.
func appendArtifactJSON(st *store.FileStore, au *audit.FileAudit, body []byte) error {
	var in struct {
		Type      string         `json:"type"`
		Title     string         `json:"title"`
		ProjectID string         `json:"project_id"`
		TaskID    string         `json:"task_id"`
		URI       string         `json:"uri"`
		Meta      map[string]any `json:"meta"`
	}
	if err := json.Unmarshal(body, &in); err != nil {
		return err
	}
	a := model.Artifact{
		ID:        fmtID("A"),
		Type:      in.Type,
		Title:     in.Title,
		ProjectID: in.ProjectID,
		TaskID:    in.TaskID,
		URI:       in.URI,
		CreatedAt: time.Now().UTC(),
		Meta:      in.Meta,
	}
	au.Emit("system", "artifact.create", map[string]any{"artifact_id": a.ID, "project_id": a.ProjectID, "task_id": a.TaskID, "cause": "run"})
	b, err := st.ReadArtifacts()
	if err != nil {
		return err
	}
	var arr []model.Artifact
	_ = json.Unmarshal(b, &arr)
	arr = append(arr, a)
	out, _ := json.MarshalIndent(arr, "", "  ")
	return st.WriteArtifacts(out)
}
