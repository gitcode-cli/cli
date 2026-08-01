package contracts

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type workflowDocument struct {
	Jobs map[string]workflowJob `yaml:"jobs"`
}

type workflowJob struct {
	Needs []string       `yaml:"needs"`
	Steps []workflowStep `yaml:"steps"`
}

type workflowStep struct {
	Name string         `yaml:"name"`
	ID   string         `yaml:"id"`
	Uses string         `yaml:"uses"`
	Run  string         `yaml:"run"`
	With map[string]any `yaml:"with"`
}

type artifactContract struct {
	jobID    string
	stepID   string
	nameBase string
	path     string
}

var artifactContracts = []artifactContract{
	{jobID: "test", stepID: "coverage-artifact", nameBase: "coverage-linux-", path: "coverage.out"},
	{jobID: "build", stepID: "binary-artifact", nameBase: "gc-ubuntu-", path: "bin/gc"},
	{jobID: "package", stepID: "wheel-artifact", nameBase: "gc-wheel-", path: "dist/*.whl"},
}

func TestGitCodeWorkflowContract(t *testing.T) {
	doc := readWorkflow(t)
	if problems := validateWorkflow(doc); len(problems) > 0 {
		t.Fatalf("GitCode workflow contract violations:\n- %s", strings.Join(problems, "\n- "))
	}
}

func TestGitCodeWorkflowContractRejectsRegressions(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*workflowDocument)
		want   string
	}{
		{name: "missing package job", mutate: removePackageJob, want: "missing required job"},
		{name: "docker job", mutate: addDockerJob, want: "unexpected job"},
		{name: "static artifact name", mutate: makeCoverageArtifactStatic, want: "upload name"},
		{name: "missing setup-python", mutate: removePythonSetup, want: "setup-python"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := readWorkflow(t)
			tt.mutate(&doc)
			if !containsProblem(validateWorkflow(doc), tt.want) {
				t.Fatalf("validator accepted regression; wanted problem containing %q", tt.want)
			}
		})
	}
}

func readWorkflow(t *testing.T) workflowDocument {
	t.Helper()
	path := filepath.Join("..", "..", ".gitcode", "workflows", "ci.yml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}

	var doc workflowDocument
	if err := yaml.Unmarshal(data, &doc); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	return doc
}

func validateWorkflow(doc workflowDocument) []string {
	var problems []string
	problems = append(problems, validateJobs(doc)...)
	problems = append(problems, validateDependencies(doc)...)
	problems = append(problems, validatePackage(doc)...)
	for _, contract := range artifactContracts {
		problems = append(problems, validateArtifact(doc, contract)...)
	}
	return problems
}

func validateJobs(doc workflowDocument) []string {
	required := []string{"lint", "test", "build", "package"}
	var problems []string
	for _, jobID := range required {
		if _, ok := doc.Jobs[jobID]; !ok {
			problems = append(problems, fmt.Sprintf("missing required job %q", jobID))
		}
	}
	for jobID, job := range doc.Jobs {
		if !contains(required, jobID) {
			problems = append(problems, fmt.Sprintf("unexpected job %q", jobID))
		}
		if jobUsesDocker(job) {
			problems = append(problems, fmt.Sprintf("job %q must not depend on Docker", jobID))
		}
	}
	return problems
}

func validateDependencies(doc workflowDocument) []string {
	var problems []string
	for _, jobID := range []string{"build", "package"} {
		job, ok := doc.Jobs[jobID]
		if ok && !contains(job.Needs, "test") {
			problems = append(problems, fmt.Sprintf("job %q must need test", jobID))
		}
	}
	return problems
}

func validatePackage(doc workflowDocument) []string {
	job, ok := doc.Jobs["package"]
	if !ok {
		return nil
	}
	var hasPythonSetup, hasSmoke bool
	for _, step := range job.Steps {
		if step.Uses == "setup-python" && stringValue(step.With["python-version"]) == "3.12" {
			hasPythonSetup = true
		}
		if strings.Contains(step.Run, "python -m build") && hasEntrypointSmoke(step.Run) {
			hasSmoke = true
		}
	}
	var problems []string
	if !hasPythonSetup {
		problems = append(problems, "package job must use setup-python 3.12")
	}
	if !hasSmoke {
		problems = append(problems, "package job must smoke-test all wheel entrypoints")
	}
	return problems
}

func validateArtifact(doc workflowDocument, contract artifactContract) []string {
	job, ok := doc.Jobs[contract.jobID]
	if !ok {
		return nil
	}
	resolve, ok := findStepByID(job.Steps, contract.stepID)
	if !ok {
		return []string{fmt.Sprintf("job %q missing artifact resolver %q", contract.jobID, contract.stepID)}
	}
	var problems []string
	if !strings.Contains(resolve.Run, contract.nameBase+"${ATOMGIT_RUN_ID:?}") ||
		!strings.Contains(resolve.Run, "$ATOMGIT_OUTPUT") {
		problems = append(problems, fmt.Sprintf("resolver %q must use runtime run ID and output file", contract.stepID))
	}
	expectedName := fmt.Sprintf("${{ steps.%s.outputs.name }}", contract.stepID)
	if !hasArtifactUpload(job.Steps, expectedName, contract.path) {
		problems = append(problems, fmt.Sprintf("job %q upload name/path must be %q and %q", contract.jobID, expectedName, contract.path))
	}
	return problems
}

func findStepByID(steps []workflowStep, id string) (workflowStep, bool) {
	for _, step := range steps {
		if step.ID == id {
			return step, true
		}
	}
	return workflowStep{}, false
}

func hasArtifactUpload(steps []workflowStep, name, path string) bool {
	for _, step := range steps {
		if step.Uses == "upload-artifact" && stringValue(step.With["name"]) == name &&
			stringValue(step.With["path"]) == path {
			return true
		}
	}
	return false
}

func hasEntrypointSmoke(run string) bool {
	for _, command := range []string{"gc version", "gitcode version", "python -m gc_cli version"} {
		if !strings.Contains(run, command) {
			return false
		}
	}
	return true
}

func jobUsesDocker(job workflowJob) bool {
	for _, step := range job.Steps {
		if strings.Contains(strings.ToLower(step.Uses+"\n"+step.Run), "docker") {
			return true
		}
	}
	return false
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func containsProblem(problems []string, want string) bool {
	for _, problem := range problems {
		if strings.Contains(problem, want) {
			return true
		}
	}
	return false
}

func stringValue(value any) string {
	if value == nil {
		return ""
	}
	return fmt.Sprint(value)
}

func removePackageJob(doc *workflowDocument) {
	delete(doc.Jobs, "package")
}

func addDockerJob(doc *workflowDocument) {
	doc.Jobs["docker"] = workflowJob{}
}

func makeCoverageArtifactStatic(doc *workflowDocument) {
	job := doc.Jobs["test"]
	for i := range job.Steps {
		if job.Steps[i].Uses == "upload-artifact" {
			job.Steps[i].With["name"] = "coverage-linux"
		}
	}
	doc.Jobs["test"] = job
}

func removePythonSetup(doc *workflowDocument) {
	job := doc.Jobs["package"]
	filtered := make([]workflowStep, 0, len(job.Steps))
	for _, step := range job.Steps {
		if step.Uses != "setup-python" {
			filtered = append(filtered, step)
		}
	}
	job.Steps = filtered
	doc.Jobs["package"] = job
}
