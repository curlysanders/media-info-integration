# Git workflow

`main` is the only long-lived branch. Changes reach it through focused pull requests from short-lived `feature/`, `fix/`, `chore/`, `refactor/`, or `docs/` branches.

## Standard workflow

```text
main
 │
 │ git switch -c docs/example
 ▼
docs/example
 │
 │ signed commits
 │ git push
 ▼
GitHub pull request
 │
 ├── CI passes
 ├── branch is current
 └── conversations are resolved
 │
 ▼
Squash merge to main
```

Start from an up-to-date `main`:

```bash
git switch main
git pull
git switch -c docs/example
```

`main` accepts only pull requests with current CI and resolved conversations. Use squash merge and delete the short-lived branch afterward. Commits on non-`main` branches must be verified signed; the squash commit created on `main` is not required to be signed.
