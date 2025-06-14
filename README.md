# PRD and Tasks

This document contains the prompts to generate the PRD and tasks for generating the gov Go code

- `tasks` - contains the mdc files and the readme.md which is the "spec"

## Generate PRD

- Add `tasks/create-prd.mdc` and `tasks/readme.md` to the GitHub agent context
- Use GPT4.1
- Use the following prompt
  - Create a PRD and save to the ./tasks/gov-prd.md file. Use ./tasks/create-prd.mdc to define the steps and structure of the PRD. Use ./tasks/README.md for the features to document. Use logical assumptions for any missing information.

- to generate from the prd-template.md
  - Create a PRD and save to the ./tasks/gov-template.md file. Use ./tasks/prd-template.md to define the steps and structure of the PRD. Use ./tasks/README.md for the features to document. Use logical assumptions for any missing information.
