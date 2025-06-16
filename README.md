# PRD and Tasks

This document contains the prompts to generate the PRD and tasks for generating the gov Go code

- `gov.md` - "spec" to use for the PRD
- `tasks` - contains the mdc files and templates

## Generate PRD

- Use GPT4.1 or Sonnet 3.7
- Use the following prompt
  - Create a PRD for the "gov" feature using ./gov.md for the features
    - use ./tasks/create-prd.mdc to define the steps and structure of the PRD

## Generate Task List from PRD

- Prompt
  - use ./tasks/prd-gov.md to create tasks and subtasks using ./tasks/generate-tasks.mdc

## Generate Code

- Prompt
  - start on task 1.1 and use ./tasks/process-task-list.mdc for control
