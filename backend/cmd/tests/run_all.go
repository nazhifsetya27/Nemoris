// Runner for all manual tests. Executes tests sequentially and prints summary.
// Run: go run ./cmd/tests/run_all.go
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	base, _ := os.Getwd()
	if _, err := os.Stat("cmd"); err != nil {
		base = filepath.Join(base, "backend")
	}
	_ = os.MkdirAll(filepath.Join(base, "test_results"), 0755)

	tests := []string{
		"test_ai",
		"test_reply_builder",
		"test_reminder_flow",
		"test_repository",
		"test_scheduler",
		"test_health",
	}

	for _, name := range tests {
		fmt.Printf("\n--- Running %s ---\n", name)
		cmd := exec.Command("go", "run", filepath.Join("cmd", "tests", name+".go"))
		cmd.Dir = base
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
	}

	// Aggregate PASS/FAIL from result files
	passTotal, failTotal := 0, 0
	resultNames := []string{"test_ai", "test_reply_builder", "test_reminder_flow", "test_repository", "test_scheduler", "test_health"}
	for _, name := range resultNames {
		p := filepath.Join(base, "test_results", name+"_result.txt")
		content, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "PASS: ") {
				var p, f int
				_, _ = fmt.Sscanf(line, "PASS: %d  FAIL: %d", &p, &f)
				passTotal += p
				failTotal += f
				break
			}
		}
	}

	fmt.Println("\n========== RUN ALL SUMMARY ==========")
	fmt.Printf("PASS: %d\n", passTotal)
	fmt.Printf("FAIL: %d\n", failTotal)
	fmt.Printf("Results saved to %s\n", filepath.Join(base, "test_results"))
}
