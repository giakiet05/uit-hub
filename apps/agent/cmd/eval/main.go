package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/config"
	"github.com/giakiet05/uit-hub/apps/agent/internal/eval"
	"github.com/giakiet05/uit-hub/apps/agent/internal/mcpadapter"
	"github.com/giakiet05/uit-hub/apps/agent/internal/logging"
)

// main runs agent eval cases without starting the TUI.
func main() {
	if err := run(context.Background(), os.Args[1:], os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, stderr *os.File) error {
	var casePath string
	var casesDir string
	var outputDir string
	var mcpConfigPath string
	var memoryPath string
	var maxRounds int
	var timeout time.Duration

	flags := flag.NewFlagSet("agent-eval", flag.ContinueOnError)
	flags.SetOutput(stderr)
	flags.StringVar(&casePath, "case", "", "path to one eval case YAML file")
	flags.StringVar(&casesDir, "cases", "", "path to a directory containing eval case YAML files")
	flags.StringVar(&outputDir, "out", "evals/results", "directory for eval result reports")
	flags.StringVar(&mcpConfigPath, "mcp-config", "mcp.yaml", "path to MCP YAML config")
	flags.StringVar(&memoryPath, "memory", "tmp/eval-memory", "memory directory used by eval runs")
	flags.IntVar(&maxRounds, "max-rounds", 24, "maximum agent rounds per eval case")
	flags.DurationVar(&timeout, "timeout", 10*time.Minute, "timeout per eval case")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if casePath == "" && casesDir == "" {
		return fmt.Errorf("provide -case or -cases")
	}

	if err := config.LoadEnv(); err != nil {
		fmt.Fprintf(stderr, "No .env file loaded: %v\n", err)
	}

	if err := os.Setenv("MCP_CONFIG_PATH", mcpConfigPath); err != nil {
		return err
	}
	if err := os.Setenv("MEMORY_PATH", memoryPath); err != nil {
		return err
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	cfg.Agent.MaxRounds = maxRounds

	// Mock SSO Login for eval runs
	ctx = context.WithValue(ctx, mcpadapter.TokenKey, "mock-22520001")

	logger := logging.NewLogger(stderr)
	runner := eval.NewRunner(cfg, logger)

	paths, err := resolveCasePaths(casePath, casesDir)
	if err != nil {
		return err
	}

	var failed int
	for _, path := range paths {
		testCase, err := eval.LoadCase(path)
		if err != nil {
			return fmt.Errorf("load case %q: %w", path, err)
		}

		// Clean up the hardcoded agent workspace to prevent state leakage between cases.
		os.RemoveAll("tmp/agent-files")
		os.MkdirAll("tmp/agent-files", 0755)

		caseCtx, cancel := context.WithTimeout(ctx, timeout)
		result, reportDir, err := runner.RunCase(caseCtx, testCase, outputDir)
		cancel()
		if err != nil {
			return fmt.Errorf("run case %q: %w", testCase.ID, err)
		}

		status := "PASS"
		if !result.Passed {
			status = "FAIL"
			failed++
		}
		logger.Info(
			"Eval case completed",
			"case_id", testCase.ID,
			"status", status,
			"report_dir", reportDir,
			"rounds", result.RoundCount,
			"tool_calls", result.ToolCalls,
			"llm_calls", result.LLMCalls,
		)
	}

	if failed > 0 {
		return fmt.Errorf("%d eval case(s) failed", failed)
	}
	return nil
}

func resolveCasePaths(casePath string, casesDir string) ([]string, error) {
	if casePath != "" {
		return []string{casePath}, nil
	}

	entries, err := os.ReadDir(casesDir)
	if err != nil {
		return nil, err
	}

	paths := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml") {
			paths = append(paths, filepath.Join(casesDir, name))
		}
	}
	sort.Strings(paths)
	return paths, nil
}
