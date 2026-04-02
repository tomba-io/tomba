package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/tomba-io/tomba/pkg/skill"
	"github.com/tomba-io/tomba/pkg/start"
	"github.com/tomba-io/tomba/pkg/util"
)

// skillCmd represents the skill parent command
var skillCmd = &cobra.Command{
	Use:     "skill",
	Aliases: []string{"sk"},
	Short:   "ClawHub — manage and run reusable Tomba workflows.",

	Run:     skillRunDefault,
	Example: skillExample,
}

var skillListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available skills (built-in + installed).",
	Run:   skillListRun,
}

var skillRunCmd = &cobra.Command{
	Use:   "run [skill-name]",
	Short: "Run a skill by name.",
	Args:  cobra.MinimumNArgs(1),
	Run:   skillRunRun,
}

var skillInfoCmd = &cobra.Command{
	Use:   "info [skill-name]",
	Short: "Show detailed information about a skill.",
	Args:  cobra.ExactArgs(1),
	Run:   skillInfoRun,
}

var skillCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new custom skill interactively.",
	Run:   skillCreateRun,
}

var skillRemoveCmd = &cobra.Command{
	Use:   "remove [skill-name]",
	Short: "Remove an installed custom skill.",
	Args:  cobra.ExactArgs(1),
	Run:   skillRemoveRun,
}

var skillExportCmd = &cobra.Command{
	Use:   "export [skill-name]",
	Short: "Export a skill as YAML to stdout.",
	Args:  cobra.ExactArgs(1),
	Run:   skillExportRun,
}

var skillImportCmd = &cobra.Command{
	Use:   "import [file.yaml]",
	Short: "Import a skill from a YAML file.",
	Args:  cobra.ExactArgs(1),
	Run:   skillImportRun,
}

func init() {
	skillCmd.AddCommand(skillListCmd, skillRunCmd, skillInfoCmd, skillCreateCmd, skillRemoveCmd, skillExportCmd, skillImportCmd)
}

// skillRunDefault shows the skill list when no subcommand is given
func skillRunDefault(cmd *cobra.Command, args []string) {
	fmt.Println()
	printClawHubBanner()
	fmt.Println()
	skillListRun(cmd, args)
}

func printClawHubBanner() {
	fmt.Printf("  %s\n", util.Bold(util.Cyan("ClawHub")))
	fmt.Printf("  %s\n", util.Gray("Reusable Tomba API workflows — chain multiple API calls into one command"))
}

// skillListRun lists all available skills
func skillListRun(cmd *cobra.Command, args []string) {
	skills, err := skill.ListSkills()
	if err != nil {
		fmt.Printf("%s %s\n", util.ErrorIcon(), util.Red(err.Error()))
		return
	}

	// Group: built-in vs user
	var builtins, custom []*skill.Skill
	builtinNames := make(map[string]bool)
	for _, s := range skill.BuiltinSkills() {
		builtinNames[s.Name] = true
	}
	for _, s := range skills {
		if builtinNames[s.Name] {
			builtins = append(builtins, s)
		} else {
			custom = append(custom, s)
		}
	}

	fmt.Printf("\n  %s Built-in Skills (%d)\n\n", util.Bold(""), len(builtins))
	fmt.Printf("  %-20s %-55s %s\n", "Name", "Description", "Steps")
	fmt.Printf("  %s\n", strings.Repeat("─", 90))
	for _, s := range builtins {
		var actions []string
		for _, step := range s.Steps {
			actions = append(actions, step.Action)
		}
		workflow := strings.Join(actions, " -> ")
		fmt.Printf("  %-20s %-55s %s\n", util.Green(s.Name), s.Description, util.Gray(workflow))
	}

	if len(custom) > 0 {
		fmt.Printf("\n  %s Custom Skills (%d)\n\n", util.Bold(""), len(custom))
		fmt.Printf("  %-20s %-55s %s\n", "Name", "Description", "Steps")
		fmt.Printf("  %s\n", strings.Repeat("─", 90))
		for _, s := range custom {
			var actions []string
			for _, step := range s.Steps {
				actions = append(actions, step.Action)
			}
			workflow := strings.Join(actions, " -> ")
			fmt.Printf("  %-20s %-55s %s\n", util.Cyan(s.Name), s.Description, util.Gray(workflow))
		}
	}

	fmt.Printf("\n  Run a skill:   %s\n", util.Bold("tomba skill run <name> --target <value>"))
	fmt.Printf("  Skill details: %s\n", util.Bold("tomba skill info <name>"))
	fmt.Printf("  Create custom: %s\n\n", util.Bold("tomba skill create"))
}

// skillRunRun executes a skill
func skillRunRun(cmd *cobra.Command, args []string) {
	init := start.New(conn)
	skillName := args[0]

	s, err := skill.FindSkill(skillName)
	if err != nil {
		fmt.Printf("%s %s\n", util.ErrorIcon(), util.Red(err.Error()))
		return
	}

	// Collect inputs
	inputs := make(map[string]string)

	// Check if --target provides the primary input
	if init.Target != "" && len(s.Inputs) > 0 {
		inputs[s.Inputs[0].Name] = init.Target
	}

	// Also parse extra args as key=value
	for _, arg := range args[1:] {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 2 {
			inputs[parts[0]] = parts[1]
		}
	}

	// Prompt for any missing required inputs
	for _, input := range s.Inputs {
		if _, ok := inputs[input.Name]; ok {
			continue
		}
		if input.Default != "" && !input.Required {
			inputs[input.Name] = input.Default
			continue
		}
		if !input.Required && input.Default == "" {
			continue
		}

		prompt := promptui.Prompt{
			Label:   input.Description,
			Default: input.Default,
		}
		val, err := prompt.Run()
		if err != nil {
			fmt.Printf("%s Cancelled\n", util.WarningIcon())
			return
		}
		inputs[input.Name] = val
	}

	// Run the skill
	runner := skill.NewRunner(init.Tomba, s, inputs)
	if err := runner.Run(); err != nil {
		fmt.Printf("%s %s\n", util.ErrorIcon(), util.Red(err.Error()))
	}
}

// skillInfoRun shows detailed info about a skill
func skillInfoRun(cmd *cobra.Command, args []string) {
	s, err := skill.FindSkill(args[0])
	if err != nil {
		fmt.Printf("%s %s\n", util.ErrorIcon(), util.Red(err.Error()))
		return
	}

	fmt.Println()
	fmt.Printf("  %s %s\n", util.Bold("Skill:"), util.Green(s.Name))
	fmt.Printf("  %s %s\n", util.Bold("Description:"), s.Description)
	fmt.Printf("  %s %s\n", util.Bold("Author:"), s.Author)
	fmt.Printf("  %s %s\n", util.Bold("Version:"), s.Version)

	if len(s.Tags) > 0 {
		fmt.Printf("  %s %s\n", util.Bold("Tags:"), util.Gray(strings.Join(s.Tags, ", ")))
	}

	// Inputs
	fmt.Printf("\n  %s\n\n", util.Bold("Inputs:"))
	for _, input := range s.Inputs {
		req := util.Gray("optional")
		if input.Required {
			req = util.Yellow("required")
		}
		defaultStr := ""
		if input.Default != "" {
			defaultStr = fmt.Sprintf(" (default: %s)", util.Gray(input.Default))
		}
		fmt.Printf("    %-15s %-10s %s %s%s\n", util.Cyan(input.Name), fmt.Sprintf("[%s]", input.Type), req, input.Description, defaultStr)
	}

	// Steps (workflow)
	fmt.Printf("\n  %s\n\n", util.Bold("Workflow:"))
	for i, step := range s.Steps {
		connector := "├──"
		if i == len(s.Steps)-1 {
			connector = "└──"
		}
		fmt.Printf("    %s %s %s\n", connector, util.Bold(step.Name), util.Gray(fmt.Sprintf("(%s)", step.Action)))

		for k, v := range step.Params {
			pipe := "│  "
			if i == len(s.Steps)-1 {
				pipe = "   "
			}
			fmt.Printf("    %s   %s: %s\n", pipe, util.Gray(k), util.Gray(v))
		}
	}

	// Output
	fmt.Printf("\n  %s %s (%s)\n", util.Bold("Output:"), s.Output.Title, s.Output.Format)
	if len(s.Output.Fields) > 0 {
		fmt.Printf("  %s %s\n", util.Bold("Fields:"), strings.Join(s.Output.Fields, ", "))
	}

	// Usage example
	fmt.Printf("\n  %s\n", util.Bold("Usage:"))
	if len(s.Inputs) > 0 {
		var example string
		switch s.Inputs[0].Type {
		case "domain":
			example = "stripe.com"
		case "email":
			example = "b.mohamed@tomba.io"
		case "url":
			example = "https://example.com/article"
		case "linkedin":
			example = "https://www.linkedin.com/in/someone"
		default:
			example = "value"
		}
		fmt.Printf("    tomba skill run %s --target %s\n", s.Name, example)
		if len(s.Inputs) > 1 {
			var extras []string
			for _, inp := range s.Inputs[1:] {
				extras = append(extras, fmt.Sprintf("%s=value", inp.Name))
			}
			fmt.Printf("    tomba skill run %s --target %s %s\n", s.Name, example, strings.Join(extras, " "))
		}
	}
	fmt.Println()
}

// skillCreateRun creates a new custom skill interactively
func skillCreateRun(cmd *cobra.Command, args []string) {
	fmt.Println()
	printClawHubBanner()
	fmt.Printf("\n  %s\n\n", util.Bold("Create a new skill"))

	// Name
	namePrompt := promptui.Prompt{Label: "Skill name"}
	name, err := namePrompt.Run()
	if err != nil {
		return
	}

	// Description
	descPrompt := promptui.Prompt{Label: "Description"}
	desc, err := descPrompt.Run()
	if err != nil {
		return
	}

	s := &skill.Skill{
		Name:        strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), " ", "-")),
		Description: desc,
		Author:      "user",
		Version:     "1.0.0",
		Tags:        []string{},
	}

	// Inputs
	fmt.Printf("\n  %s (press Enter with empty name to finish)\n\n", util.Bold("Define inputs"))
	for {
		inNamePrompt := promptui.Prompt{Label: "Input name (or empty to finish)"}
		inName, err := inNamePrompt.Run()
		if err != nil || inName == "" {
			break
		}

		typeSelect := promptui.Select{
			Label: "Input type",
			Items: []string{"string", "email", "domain", "url", "linkedin"},
		}
		_, inType, err := typeSelect.Run()
		if err != nil {
			break
		}

		inDescPrompt := promptui.Prompt{Label: "Input description"}
		inDesc, _ := inDescPrompt.Run()

		s.Inputs = append(s.Inputs, skill.SkillInput{
			Name:        inName,
			Description: inDesc,
			Type:        inType,
			Required:    true,
		})
	}

	// Steps
	fmt.Printf("\n  %s (press Enter with empty action to finish)\n\n", util.Bold("Define workflow steps"))
	actions := []string{"search", "finder", "enrich", "verify", "count", "status", "sources", "similar", "technology", "linkedin", "author", "phone-finder", "phone-validator", "reveal"}

	stepNum := 1
	for {
		actionSelect := promptui.Select{
			Label: fmt.Sprintf("Step %d action (Ctrl+C to finish)", stepNum),
			Items: actions,
		}
		_, action, err := actionSelect.Run()
		if err != nil {
			break
		}

		idPrompt := promptui.Prompt{
			Label:   "Step ID",
			Default: fmt.Sprintf("step%d", stepNum),
		}
		id, _ := idPrompt.Run()

		stepNamePrompt := promptui.Prompt{
			Label:   "Step display name",
			Default: action,
		}
		stepName, _ := stepNamePrompt.Run()

		// Collect params
		params := make(map[string]string)
		fmt.Printf("    %s (use {{input.name}} or {{steps.id.field}} for templates)\n", util.Gray("Add parameters:"))
		for {
			pkPrompt := promptui.Prompt{Label: "  Param key (or empty to finish)"}
			pk, err := pkPrompt.Run()
			if err != nil || pk == "" {
				break
			}
			pvPrompt := promptui.Prompt{Label: "  Param value"}
			pv, _ := pvPrompt.Run()
			params[pk] = pv
		}

		s.Steps = append(s.Steps, skill.SkillStep{
			ID:     id,
			Name:   stepName,
			Action: action,
			Params: params,
		})
		stepNum++
	}

	if len(s.Steps) == 0 {
		fmt.Printf("%s Skill must have at least one step\n", util.ErrorIcon())
		return
	}

	// Output
	formatSelect := promptui.Select{
		Label: "Output format",
		Items: []string{"summary", "table", "list"},
	}
	_, outFormat, _ := formatSelect.Run()

	titlePrompt := promptui.Prompt{Label: "Output title", Default: s.Name + " results"}
	outTitle, _ := titlePrompt.Run()

	s.Output = skill.SkillOutput{
		Format: outFormat,
		Title:  outTitle,
	}

	// Save
	if err := skill.SaveSkill(s); err != nil {
		fmt.Printf("%s %s\n", util.ErrorIcon(), util.Red(err.Error()))
		return
	}

	fmt.Printf("\n%s Skill '%s' created successfully!\n", util.SuccessIcon(), util.Green(s.Name))
	fmt.Printf("  Run it: %s\n\n", util.Bold(fmt.Sprintf("tomba skill run %s --target <value>", s.Name)))
}

// skillRemoveRun removes a custom skill
func skillRemoveRun(cmd *cobra.Command, args []string) {
	if err := skill.RemoveSkill(args[0]); err != nil {
		fmt.Printf("%s %s\n", util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	fmt.Printf("%s Skill '%s' removed\n", util.SuccessIcon(), args[0])
}

// skillExportRun exports a skill as YAML
func skillExportRun(cmd *cobra.Command, args []string) {
	s, err := skill.FindSkill(args[0])
	if err != nil {
		fmt.Printf("%s %s\n", util.ErrorIcon(), util.Red(err.Error()))
		return
	}

	data, err := yaml.Marshal(s)
	if err != nil {
		fmt.Printf("%s %s\n", util.ErrorIcon(), util.Red(err.Error()))
		return
	}

	fmt.Println(string(data))
}

// skillImportRun imports a skill from a YAML file
func skillImportRun(cmd *cobra.Command, args []string) {
	path := args[0]

	s, err := skill.LoadSkill(path)
	if err != nil {
		fmt.Printf("%s %s\n", util.ErrorIcon(), util.Red(err.Error()))
		return
	}

	if err := skill.SaveSkill(s); err != nil {
		fmt.Printf("%s %s\n", util.ErrorIcon(), util.Red(err.Error()))
		return
	}

	fmt.Printf("%s Skill '%s' imported successfully!\n", util.SuccessIcon(), util.Green(s.Name))
	fmt.Printf("  Run it: %s\n", util.Bold(fmt.Sprintf("tomba skill run %s --target <value>", s.Name)))
}

// skillInitCmd is used internally to create sample skill files for new users
func SkillInit() error {
	dir := skill.SkillDir()
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return nil // already initialized
	}
	return skill.EnsureSkillDir()
}
