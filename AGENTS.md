# cluster-lifecycle-api Agent Instructions

This repository defines Kubernetes API types and client helpers used by MCE and
ACM. It is a Go module at `github.com/stolostron/cluster-lifecycle-api`.

## Repository layout

- `action/v1beta1/`, `view/v1beta1/`, `clusterinfo/v1beta1/`, `imageregistry/v1alpha1/`, and `klusterletconfig/v1alpha1/` contain API types, registration, CRDs, and generated code.
- `helpers/` contains reusable clients and utilities for image registries, TLS profiles, local clusters, klusterlet configuration, and user permissions.
- `hack/` contains code-generation and generated-artifact verification scripts.
- `vendor/` contains the pinned build and Go dependencies used by the Makefile.

## Development commands

Run these from the repository root:

```bash
make build
make test
make verify
make verify-deps
```

`make verify` checks generated deepcopy, Swagger documentation, CRDs, and
codegen output. After changing API definitions or generation markers, run
`make update` and then `make verify`; review generated changes before submitting
them. Do not hand-edit `zz_generated*` files or CRD output.

The GitHub Actions workflow runs `make verify`, `make verify-deps`, `make build`,
and `make test`. The module declares Go 1.25; the workflow currently provisions
Go 1.23, so local validation should use the repository's declared toolchain when
available and CI compatibility should be considered for changes.

## API conventions

- Treat API type definitions and JSON tags as public compatibility surfaces.
- Preserve Kubernetes API conventions: TypeMeta/ObjectMeta, list types, Spec,
  Status, conditions, and explicit versioned packages.
- Keep generated clients, deepcopy code, Swagger documentation, and CRDs in sync
  with source types.
- Use the existing Kubernetes and Open Cluster Management API dependencies and
  annotations rather than introducing parallel representations.
- Keep helper behavior compatible with the API annotations and typed clients;
  add focused unit tests for helper changes.

## Architecture

For system architecture, data flows, and module layout, see
[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## Personal configuration

Read personal config at the start of any task that needs an assignee, email, or project key.
Canonical path: `~/.config/user.local.md` (tool-agnostic, global).
If the file does not exist, fall back to agent memory (`user-config`), then placeholders.
Run `make personalize` to generate or update the file (if this repo uses Fleet Engineering tooling).

## Tool availability

- GitHub operations: GitHub MCP tools are available; the `gh` CLI is not installed in this environment. Use the GitHub MCP server scoped to the repository's organization for pull request operations.
- Jira operations: Jira MCP tools are available; use them for Jira work.
- Separate GitHub organization tokens, when using a CLI in another environment, follow the `GH_TOKEN_<ORG>` convention with the organization lowercased and hyphens replaced by underscores. Never print token values or inspect the environment broadly.

## Fleet Engineering Skills

Fetch and apply the relevant skill when the task matches its domain.

| Skill | When to use |
|---|---|
| [backlog-grooming](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/backlog-grooming/SKILL.md) | Groom Jira backlog readiness and stale work |
| [breaking-changes](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/breaking-changes/SKILL.md) | Detect API, behavior, configuration, or integration breaks |
| [bug-specialist](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/bug-specialist/SKILL.md) | Triage bugs and plan reproductions and fixes |
| [ci-triage](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/ci-triage/SKILL.md) | Diagnose failing pull request checks |
| [coderabbit-sync](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/coderabbit-sync/SKILL.md) | Maintain the Fleet reference CodeRabbit configuration |
| [diagnosing-bugs](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/diagnosing-bugs/SKILL.md) | Debug unclear failures systematically |
| [epic-specialist](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/epic-specialist/SKILL.md) | Plan multi-sprint Jira epics |
| [feature-specialist](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/feature-specialist/SKILL.md) | Define major customer-facing capabilities |
| [finish-work](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/finish-work/SKILL.md) | Commit, push, open a pull request, and update Jira when explicitly requested |
| [init-context-docs](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/init-context-docs/SKILL.md) | Assess and bootstrap repository context documentation |
| [initiative-specialist](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/initiative-specialist/SKILL.md) | Plan internal strategic engineering programs |
| [jira-create](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/jira-create/SKILL.md) | Create Jira issues through the required quality workflow |
| [jira-qe-readiness](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/jira-qe-readiness/SKILL.md) | Check whether Jira work is ready for QE verification |
| [jira-report](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/jira-report/SKILL.md) | Produce Jira portfolio and quality reports |
| [jira-specialist](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/jira-specialist/SKILL.md) | Triage, search, link, and transition Jira work |
| [org-repo-audit](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/org-repo-audit/SKILL.md) | Audit repositories for SDLC readiness |
| [pr-fix](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/pr-fix/SKILL.md) | Fix conflicts, CI failures, and review comments |
| [pr-hygiene](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/pr-hygiene/SKILL.md) | Manage stale pull request lifecycle |
| [pr-review](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/pr-review/SKILL.md) | Review GitHub pull requests with inline findings |
| [pr-review-detailed](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/pr-review-detailed/SKILL.md) | Run layered branch and diff analysis |
| [pr-review-fix](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/pr-review-fix/SKILL.md) | Review and fix local changes before committing |
| [repo-content-audit](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/repo-content-audit/SKILL.md) | Find unlinked or orphaned repository content |
| [repo-setup](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/repo-setup/SKILL.md) | Onboard repositories to the Fleet agentic SDLC |
| [release-notes](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/release-notes/SKILL.md) | Generate release notes from merged pull requests |
| [risk-report](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/risk-report/SKILL.md) | Detect Jira and delivery risk signals |
| [spike-specialist](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/spike-specialist/SKILL.md) | Conduct time-boxed technical research |
| [start-work](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/start-work/SKILL.md) | Create a Jira sub-task for active work |
| [story-specialist](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/story-specialist/SKILL.md) | Define user stories and acceptance criteria |
| [task-specialist](https://raw.githubusercontent.com/OpenShift-Fleet/agentic-sdlc/main/skills/task-specialist/SKILL.md) | Plan internal technical implementation work |

## Context loading

Load deeper context on demand rather than reading the whole repository up front:

- Architecture: `docs/ARCHITECTURE.md`
- API definitions: the relevant versioned package under `action/`, `view/`, `clusterinfo/`, `imageregistry/`, or `klusterletconfig/`
- Generated artifact checks: `hack/` and the Makefile
