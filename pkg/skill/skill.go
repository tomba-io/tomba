package skill

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/tomba-io/tomba/pkg/config"
)

// Skill defines a reusable multi-step Tomba API workflow
type Skill struct {
	Name        string            `yaml:"name" json:"name"`
	Description string            `yaml:"description" json:"description"`
	Author      string            `yaml:"author" json:"author"`
	Version     string            `yaml:"version" json:"version"`
	Tags        []string          `yaml:"tags" json:"tags"`
	Inputs      []SkillInput      `yaml:"inputs" json:"inputs"`
	Steps       []SkillStep       `yaml:"steps" json:"steps"`
	Output      SkillOutput       `yaml:"output" json:"output"`
	Variables   map[string]string `yaml:"variables,omitempty" json:"variables,omitempty"`
}

// SkillInput defines an input parameter for the skill
type SkillInput struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
	Type        string `yaml:"type" json:"type"` // string, email, domain, url, linkedin
	Required    bool   `yaml:"required" json:"required"`
	Default     string `yaml:"default,omitempty" json:"default,omitempty"`
}

// SkillStep defines one API call in the workflow
type SkillStep struct {
	ID        string            `yaml:"id" json:"id"`
	Name      string            `yaml:"name" json:"name"`
	Action    string            `yaml:"action" json:"action"` // search, finder, enrich, verify, count, status, sources, similar, technology, linkedin, author, phone-finder, phone-validator, reveal
	Params    map[string]string `yaml:"params" json:"params"` // supports {{input.name}} and {{steps.id.field}} templates
	SaveAs    string            `yaml:"save_as,omitempty" json:"save_as,omitempty"`
	Condition string            `yaml:"condition,omitempty" json:"condition,omitempty"` // skip step if condition is false
}

// SkillOutput defines how to display the final result
type SkillOutput struct {
	Format string   `yaml:"format" json:"format"` // table, list, summary
	Title  string   `yaml:"title" json:"title"`
	Fields []string `yaml:"fields" json:"fields"`
}

// SkillDir returns the directory where skills are stored
func SkillDir() string {
	return filepath.Join(config.Home(), ".tomba", "skills")
}

// EnsureSkillDir creates the skills directory if it doesn't exist
func EnsureSkillDir() error {
	dir := SkillDir()
	return os.MkdirAll(dir, 0755)
}

// LoadSkill reads a skill from a YAML file
func LoadSkill(path string) (*Skill, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read skill file: %w", err)
	}

	var s Skill
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("invalid skill format: %w", err)
	}

	if s.Name == "" {
		return nil, fmt.Errorf("skill must have a name")
	}
	if len(s.Steps) == 0 {
		return nil, fmt.Errorf("skill must have at least one step")
	}

	return &s, nil
}

// SaveSkill writes a skill to a YAML file
func SaveSkill(s *Skill) error {
	if err := EnsureSkillDir(); err != nil {
		return err
	}

	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Errorf("cannot serialize skill: %w", err)
	}

	filename := strings.ReplaceAll(strings.ToLower(s.Name), " ", "-") + ".yaml"
	path := filepath.Join(SkillDir(), filename)

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("cannot write skill file: %w", err)
	}

	return nil
}

// ListSkills returns all installed skills (built-in + user)
func ListSkills() ([]*Skill, error) {
	var skills []*Skill

	// Built-in skills
	for _, s := range BuiltinSkills() {
		skills = append(skills, s)
	}

	// User skills from disk
	dir := SkillDir()
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return skills, nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return skills, nil
	}

	for _, entry := range entries {
		if entry.IsDir() || (!strings.HasSuffix(entry.Name(), ".yaml") && !strings.HasSuffix(entry.Name(), ".yml")) {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		s, err := LoadSkill(path)
		if err != nil {
			continue
		}
		skills = append(skills, s)
	}

	return skills, nil
}

// FindSkill looks up a skill by name
func FindSkill(name string) (*Skill, error) {
	lower := strings.ToLower(name)

	// Check built-in first
	for _, s := range BuiltinSkills() {
		if strings.ToLower(s.Name) == lower || strings.ReplaceAll(strings.ToLower(s.Name), " ", "-") == lower {
			return s, nil
		}
	}

	// Check user skills
	dir := SkillDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("skill '%s' not found", name)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		s, err := LoadSkill(path)
		if err != nil {
			continue
		}
		if strings.ToLower(s.Name) == lower || strings.ReplaceAll(strings.ToLower(s.Name), " ", "-") == lower {
			return s, nil
		}
	}

	return nil, fmt.Errorf("skill '%s' not found. Run 'tomba skill list' to see available skills", name)
}

// RemoveSkill deletes a user skill file
func RemoveSkill(name string) error {
	lower := strings.ToLower(name)

	// Cannot remove built-in
	for _, s := range BuiltinSkills() {
		if strings.ToLower(s.Name) == lower || strings.ReplaceAll(strings.ToLower(s.Name), " ", "-") == lower {
			return fmt.Errorf("cannot remove built-in skill '%s'", name)
		}
	}

	dir := SkillDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("skill '%s' not found", name)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		s, err := LoadSkill(path)
		if err != nil {
			continue
		}
		if strings.ToLower(s.Name) == lower || strings.ReplaceAll(strings.ToLower(s.Name), " ", "-") == lower {
			return os.Remove(path)
		}
	}

	return fmt.Errorf("skill '%s' not found", name)
}

// ResolveTemplate replaces {{input.x}} and {{steps.id.field}} placeholders
func ResolveTemplate(tmpl string, inputs map[string]string, stepResults map[string]map[string]interface{}) string {
	result := tmpl

	// Replace {{input.name}}
	for k, v := range inputs {
		result = strings.ReplaceAll(result, "{{input."+k+"}}", v)
	}

	// Replace {{steps.id.field}} — supports nested with dot notation
	for stepID, data := range stepResults {
		for k, v := range flattenMap(data, "") {
			placeholder := "{{steps." + stepID + "." + k + "}}"
			result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", v))
		}
	}

	return result
}

// flattenMap flattens a nested map into dot-notation keys
func flattenMap(m map[string]interface{}, prefix string) map[string]interface{} {
	flat := make(map[string]interface{})
	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch val := v.(type) {
		case map[string]interface{}:
			for fk, fv := range flattenMap(val, key) {
				flat[fk] = fv
			}
		case []interface{}:
			flat[key] = fmt.Sprintf("[%d items]", len(val))
			// Also store first element fields for convenience
			if len(val) > 0 {
				if first, ok := val[0].(map[string]interface{}); ok {
					for fk, fv := range flattenMap(first, key+".0") {
						flat[fk] = fv
					}
				}
			}
		default:
			flat[key] = v
		}
	}
	return flat
}

// ParseStepResult parses raw JSON bytes into a map
func ParseStepResult(raw []byte) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}
	return data, nil
}
