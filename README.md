# PRD and Tasks

This document contains the prompts to generate the PRD and tasks for generating the gov Go code

- `gov.md` - "spec" to use for the PRD
- `tasks` - contains the mdc files and templates

## Generate PRD

- Use GPT4.1 or Sonnet 3.7
- Use the following prompt
  - Create a PRD and save to ./tasks/gov-prd.md using ./tasks/create-prd.mdc and ./tasks/gov-template.md to define the steps and structure of the PRD. Use ./gov.md for the features to document. Use logical assumptions for any missing information.

## Generate Task List from PRD

- Prompt
  - use ./tasks/gov-prd.md and create tasks and subtasks using ./tasks/generate-tasks.mdc

## Generate Code

- Prompt
  - start on task 1.1 and use ./tasks/process-task-list.mdc for control
