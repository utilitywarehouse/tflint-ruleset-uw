package rules

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func ReferenceLink(ruleName string) string {
	return fmt.Sprintf("https://github.com/utilitywarehouse/tflint-ruleset-uw/blob/main/rules/%s.md", toSnakeCase(ruleName))
}

func getVar(name string, target any, runner tflint.Runner) error {
	variableSchema := &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "variable",
				LabelNames: []string{"variable_name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "default"},
					},
				},
			},
		},
	}
	content, err := runner.GetModuleContent(variableSchema, nil)
	if err != nil {
		return fmt.Errorf("failed to get module content: %w", err)
	}

	for _, block := range content.Blocks {
		if block.Labels[0] == name {
			attr, ok := block.Body.Attributes["default"]
			if !ok {
				return fmt.Errorf("variable %s has no default value", name)
			}

			err = runner.EvaluateExpr(attr.Expr, target, nil)
			if err != nil {
				return fmt.Errorf("failed to evaluate expression: %w", err)
			}
			return nil
		}
	}
	return fmt.Errorf("variable %q not found", name)
}

// The runner have the functions GetOriginalwd() and GetModulePath(), but they don't work correctly
func modulePath(runner tflint.Runner) string {
	files, _ := runner.GetFiles()
	for fullPath := range files {
		return path.Dir(fullPath)
	}
	return "<module root not found>"
}

func mustGetStringListVar(name string, runner tflint.Runner) ([]string, error) {
	var result []string
	err := getVar(name, &result, runner)
	if err != nil {
		return nil, fmt.Errorf("plugin error at %s: failed to get required value: %w", modulePath(runner), err)
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("plugin error at %s: required value is empty: %q", modulePath(runner), name)
	}
	return result, nil
}

func mustGetStringVar(name string, runner tflint.Runner) (string, error) {
	var result string
	err := getVar(name, &result, runner)
	if err != nil {
		return "", fmt.Errorf("plugin error at %s: failed to get required value: %w", modulePath(runner), err)
	}
	if len(result) == 0 {
		return "", fmt.Errorf("plugin error at %s: required value is empty: %q", modulePath(runner), name)
	}
	return result, nil
}

// Ideally uses the variable value, but in dev and prod it can also auto-detect
// it. This helps with local submodules, where the environment default value is
// usually left empty, and received as a module attribute.
func GetEnv(runner tflint.Runner) (string, error) {
	env, err := mustGetStringVar("environment", runner)
	if err != nil {
		path := modulePath(runner)
		if strings.Contains(path, "aws/prod") {
			env, err = "prod", nil
		}
		if strings.Contains(path, "aws/dev") {
			env, err = "dev", nil
		}
	}
	return env, err
}

func GetTeams(runner tflint.Runner) ([]string, error) {
	env, err := GetEnv(runner)
	if err != nil {
		return nil, err
	}
	return mustGetStringListVar("teams_"+env, runner)
}

func GetTeamPrefixes(runner tflint.Runner) ([]string, error) {
	env, err := GetEnv(runner)
	if err != nil {
		return nil, err
	}
	return mustGetStringListVar("team_prefixes_"+env, runner)
}

func GetBucketNameExceptions(runner tflint.Runner) []string {
	env, err := GetEnv(runner)
	if err != nil {
		return []string{}
	}

	exceptions, err := mustGetStringListVar("bucket_name_exceptions_"+env, runner)
	if err != nil {
		return []string{}
	}

	return exceptions
}
