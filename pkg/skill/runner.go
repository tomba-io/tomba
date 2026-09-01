package skill

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tomba-io/go/tomba"
	"github.com/tomba-io/tomba/pkg/output"
	"github.com/tomba-io/tomba/pkg/util"
)

// Runner executes a skill's steps using the Tomba API
type Runner struct {
	Tomba       *tomba.Tomba
	Skill       *Skill
	Inputs      map[string]string
	StepResults map[string]map[string]interface{}
}

// NewRunner creates a skill runner
func NewRunner(t *tomba.Tomba, s *Skill, inputs map[string]string) *Runner {
	return &Runner{
		Tomba:       t,
		Skill:       s,
		Inputs:      inputs,
		StepResults: make(map[string]map[string]interface{}),
	}
}

// Run executes all steps and displays the output
func (r *Runner) Run() error {
	fmt.Printf("\n%s Running skill: %s\n", util.SuccessIcon(), util.Bold(r.Skill.Name))
	fmt.Printf("  %s\n\n", util.Gray(r.Skill.Description))

	// Show workflow
	var stepNames []string
	for _, step := range r.Skill.Steps {
		stepNames = append(stepNames, step.Action)
	}
	fmt.Printf("  %s %s\n\n", util.Gray("workflow:"), util.Cyan(strings.Join(stepNames, " -> ")))

	// Execute each step
	totalSteps := len(r.Skill.Steps)
	for i, step := range r.Skill.Steps {
		fmt.Printf("  [%d/%d] %s %s", i+1, totalSteps, util.Bold(step.Name), util.Gray("..."))

		result, err := r.executeStep(step)
		if err != nil {
			fmt.Printf(" %s\n", util.Red("failed"))
			fmt.Printf("         %s\n", util.Red(err.Error()))
			// Continue to next step, some might not depend on this one
			continue
		}

		r.StepResults[step.ID] = result
		fmt.Printf(" %s\n", util.Green("done"))
	}

	// Display output
	fmt.Println()
	r.displayOutput()

	return nil
}

func (r *Runner) executeStep(step SkillStep) (map[string]interface{}, error) {
	// Resolve parameter templates
	params := make(map[string]string)
	for k, v := range step.Params {
		params[k] = ResolveTemplate(v, r.Inputs, r.StepResults)
	}

	switch step.Action {
	case "search":
		p := tomba.Params{"domain": params["domain"]}
		if v, ok := params["limit"]; ok && v != "" {
			p["limit"] = v
		}
		if v, ok := params["page"]; ok && v != "" {
			p["page"] = v
		}
		if v, ok := params["department"]; ok && v != "" {
			p["department"] = v
		}
		result, err := r.Tomba.DomainSearch(p)
		if err != nil {
			return nil, err
		}
		raw, _ := result.Marshal()
		return ParseStepResult(raw)

	case "finder":
		p := tomba.Params{"domain": params["domain"]}
		if v := params["first_name"]; v != "" {
			p["first_name"] = v
		}
		if v := params["last_name"]; v != "" {
			p["last_name"] = v
		}
		if v := params["full_name"]; v != "" {
			p["full_name"] = v
		}
		result, err := r.Tomba.EmailFinder(p)
		if err != nil {
			return nil, err
		}
		raw, _ := result.Marshal()
		return ParseStepResult(raw)

	case "enrich":
		email := params["email"]
		if email == "" || strings.Contains(email, "{{") {
			return nil, fmt.Errorf("no email available from previous step")
		}
		result, err := r.Tomba.Enrichment(tomba.Params{"email": email})
		if err != nil {
			return nil, err
		}
		raw, _ := result.Marshal()
		return ParseStepResult(raw)

	case "verify":
		email := params["email"]
		if email == "" || strings.Contains(email, "{{") {
			return nil, fmt.Errorf("no email available from previous step")
		}
		result, err := r.Tomba.EmailVerifier(tomba.Params{"email": email})
		if err != nil {
			return nil, err
		}
		raw, _ := result.Marshal()
		return ParseStepResult(raw)

	case "count":
		result, err := r.Tomba.Count(params["domain"])
		if err != nil {
			return nil, err
		}
		raw, _ := result.Marshal()
		return ParseStepResult(raw)

	case "status":
		result, err := r.Tomba.Status(params["domain"])
		if err != nil {
			return nil, err
		}
		raw, _ := result.Marshal()
		return ParseStepResult(raw)

	case "sources":
		email := params["email"]
		if email == "" || strings.Contains(email, "{{") {
			return nil, fmt.Errorf("no email available from previous step")
		}
		result, err := r.Tomba.Sources(email)
		if err != nil {
			return nil, err
		}
		raw, _ := result.Marshal()
		return ParseStepResult(raw)

	case "similar":
		result, err := r.Tomba.SimilarDomains(params["domain"])
		if err != nil {
			return nil, err
		}
		raw, _ := result.Marshal()
		return ParseStepResult(raw)

	case "technology":
		result, err := r.Tomba.TechnologyCheck(params["domain"])
		if err != nil {
			return nil, err
		}
		raw, _ := result.Marshal()
		return ParseStepResult(raw)

	case "linkedin":
		p := tomba.Params{"url": params["url"]}
		result, err := r.Tomba.LinkedinFinder(p)
		if err != nil {
			return nil, err
		}
		raw, _ := result.Marshal()
		return ParseStepResult(raw)

	case "author":
		result, err := r.Tomba.AuthorFinder(tomba.Params{"url": params["url"]})
		if err != nil {
			return nil, err
		}
		raw, _ := result.Marshal()
		return ParseStepResult(raw)

	case "phone-finder":
		p := tomba.Params{}
		if v := params["email"]; v != "" && !strings.Contains(v, "{{") {
			p["email"] = v
		}
		if v := params["domain"]; v != "" {
			p["domain"] = v
		}
		result, err := r.Tomba.PhoneFinder(p)
		if err != nil {
			return nil, err
		}
		raw, _ := result.Marshal()
		return ParseStepResult(raw)

	case "phone-validator":
		p := tomba.Params{"phone": params["phone"]}
		if v := params["country_code"]; v != "" {
			p["country_code"] = v
		}
		result, err := r.Tomba.PhoneValidator(p)
		if err != nil {
			return nil, err
		}
		raw, _ := result.Marshal()
		return ParseStepResult(raw)

	default:
		return nil, fmt.Errorf("unknown action: %s", step.Action)
	}
}

func (r *Runner) displayOutput() {
	fmt.Printf("  %s\n", util.Bold(r.Skill.Output.Title))
	fmt.Printf("  %s\n\n", strings.Repeat("─", len(r.Skill.Output.Title)+4))

	// Collect all step results into a unified view
	for stepID, data := range r.StepResults {
		raw, err := json.Marshal(data)
		if err != nil {
			continue
		}

		// Find the step definition to determine the action type
		for _, step := range r.Skill.Steps {
			if step.ID == stepID {
				text := output.DisplayText(string(raw), step.Action)
				if text != "" {
					fmt.Print(text)
				}
				break
			}
		}
	}
}
