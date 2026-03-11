package server

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
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
			prompt := "You are PM. Provide acceptance criteria and an execution plan.\n" +
				"Task ID: " + taskID + "\n" +
				"Pipeline: " + body.Pipeline + "\n"
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

			r.Status = "done"
			r.EndedAt = time.Now().UTC().Format(time.RFC3339)
			r.Summary = "local_cli ok"
			resb, _ := json.MarshalIndent(map[string]any{"status": "ok", "summary": r.Summary}, "", "  ")
			_ = os.WriteFile(filepath.Join(dir, "RESULT.json"), resb, 0o644)
			au.Emit("system", "run.step.done", map[string]any{"run_id": runID, "task_id": taskID, "role": "pm", "status": "ok"})
			au.Emit("system", "run.done", map[string]any{"run_id": runID, "task_id": taskID, "runner": runnerType, "backend": runnerBackend})
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
