package cli

import (
	"errors"
	"fmt"

	"github.com/aws/amazon-ecs-cli-v2/internal/pkg/term/color"
	"github.com/aws/amazon-ecs-cli-v2/internal/pkg/term/log"
	"github.com/aws/amazon-ecs-cli-v2/internal/pkg/workspace"
)

type Asker struct {
	projectService   projectService
	prompter         prompter
	workspaceService *workspace.Workspace
}

func NewAsker(projectService projectService, prompter prompter) *Asker {
	return &Asker{
		projectService:   projectService,
		prompter:         prompter,
		workspaceService: nil,
	}
}

type SelectAppInput struct {
	Project   string
	Prompt    string
	LocalOnly bool
}

// SelectApp prompts the user to select one of their apps.
// It will fetch all existing apps for a project. If there's
// only one app, it does not prompt them for an app and returns
// the only app.
func (a *Asker) SelectApp(input *SelectAppInput) (string, error) {
	var names []string
	if input.LocalOnly {
		names, err := a.workspaceService.AppNames()
		if err != nil {
			return "", fmt.Errorf("list applications in workspace: %w", err)
		}
		if len(names) == 0 {
			return "", errors.New("no applications found in the workspace")
		}
	} else {
		apps, err := a.projectService.ListApplications(input.Project)
		if err != nil {
			return "", fmt.Errorf("get app names: %w", err)
		}
		for _, app := range apps {
			names = append(names, app.Name)
		}
		if err != nil {
			return "", err
		}
	}

	if len(names) == 0 {
		return "", fmt.Errorf("couldn't find any application in the project")
	}
	if len(names) == 1 {
		log.Infof("Only found one application, defaulting to: %s\n", color.HighlightUserInput(names[0]))
		return names[0], nil
	}
	name, err := a.prompter.SelectOne(input.Prompt, "", names)
	if err != nil {
		return "", fmt.Errorf("select application to delete: %w", err)
	}

	return name, nil
}

type SelectEnvInput struct {
	Project string
	Prompt  string
}

// SelectEnv prompts the user to select an environment associated
// with their project
func (a *Asker) SelectEnv(input *SelectEnvInput) (string, error) {
	envs, err := a.projectService.ListEnvironments(input.Project)
	if err != nil {
		return "", fmt.Errorf("get environments for project %s from metadata store: %w", input.Project, err)
	}
	if len(envs) == 0 {
		log.Infof("Couldn't find any environments associated with project %s, try initializing one: %s\n",
			color.HighlightUserInput(input.Project),
			color.HighlightCode("ecs-preview env init"))
		return "", fmt.Errorf("no environments found in project %s", input.Project)
	}
	if len(envs) == 1 {
		log.Infof("Only found one environment, defaulting to: %s\n", color.HighlightUserInput(envs[0].Name))
		return envs[0].Name, nil
	}

	var names []string
	for _, env := range envs {
		names = append(names, env.Name)
	}

	selectedEnvName, err := a.prompter.SelectOne(input.Prompt, "", names)
	if err != nil {
		return "", fmt.Errorf("select env name: %w", err)
	}

	return selectedEnvName, nil
}

type SelectProjectInput struct {
	HelpPrompt string
	Prompt     string
}

// SelectProject prompts the user to select one of their projects.
func (a *Asker) SelectProject(input *SelectProjectInput) (string, error) {
	projs, err := a.projectService.ListProjects()
	if err != nil {
		return "", err
	}
	var projNames []string
	for _, proj := range projs {
		projNames = append(projNames, proj.Name)
	}
	if len(projNames) == 0 {
		log.Infoln("There are no projects to select.")
		return "", nil
	}
	proj, err := a.prompter.SelectOne(
		input.Prompt,
		input.HelpPrompt,
		projNames,
	)
	return proj, err
}
