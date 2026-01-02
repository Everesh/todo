# Contributing to Todo TUI

Thanks for your interest in contributing to Todo TUI! We welcome contributions of all kinds — bug reports, feature requests, documentation improvements, tests, and code. This document explains the recommended workflow and project expectations to make contributions easier for everyone.

If you're unsure where to start, check the issue tracker for items labeled `good first issue` or `help wanted`.

---

## Table of contents

- [Reporting bugs & requesting features](#reporting-bugs--requesting-features)
- [Getting the code](#getting-the-code)
- [Branching & commits](#branching--commits)
- [Pull requests](#pull-requests)
- [Code style & testing](#code-style--testing)
- [Documentation](#documentation)
- [Issue and PR templates](#issue-and-pr-templates)
- [Code of Conduct](#code-of-conduct)
- [Quick contributing checklist](#quick-contributing-checklist)

---

## Reporting bugs & requesting features

- Search existing issues before opening a new one to avoid duplicates.
- When opening an issue, include:
  - A clear and descriptive title.
  - Steps to reproduce (for bugs).
  - Expected vs actual behavior.
  - Relevant logs, error output, or screenshots/GIFs (for UI behavior).
  - Your OS and terminal details if relevant (Windows/macOS/Linux, terminal emulator).
- For feature requests, describe the motivation, proposed behavior, and any UX ideas or mockups.

---

## Getting the code

1. Fork the repository on GitHub.
2. Clone your fork:
   ```bash
   git clone https://github.com/<your-username>/todo.git
   cd todo
   ```
3. Add the upstream remote (optional but recommended):
   ```bash
   git remote add upstream https://github.com/nirabyte/todo.git
   git fetch upstream
   ```

---

## Branching & commits

- Create a descriptive branch for your work:
  - Example: `fix/timer-parse` or `feat/theme-selection`.
  ```bash
  git checkout -b feat/your-feature
  ```
- Keep changes focused: one logical change per branch/PR.
- Write clear commit messages in the imperative tense:
  - `Add timer parsing for hhmm format`
  - `Fix crash when todos.json is empty`
- Squash small, related commits before merging or let maintainers squash during merge if requested.

---

## Pull requests

- Push your branch to your fork and open a Pull Request (PR) against the `main` branch of the upstream repository.
- In the PR description:
  - Explain the motivation and what changed.
  - Reference related issues with `#issue_number`.
  - Include screenshots/GIFs for UI changes when applicable.
  - List steps to manually test the change.
- Be responsive to review feedback; address requested changes promptly.
- If your PR is a work in progress, mark it as a draft PR.

---

## Code style & testing

- Follow existing code patterns in the repository.
- For Go code:
  - Run `gofmt` / `go fmt` and `go vet` where appropriate:
    ```bash
    go fmt ./...
    go vet ./...
    ```
  - Ensure the project builds:
    ```bash
    go build ./...
    ```
- Add tests for new functionality or bug fixes when practical.
- When adding new dependencies, prefer well-maintained and widely-used packages.

---

## Documentation

- Keep the `README.md` and other docs up to date with user-facing changes.
- Add usage examples or GIFs for UI changes.
- If you create new configuration options or commands, document them in the README or a dedicated docs file.

---

## Issue and PR templates

- Use the provided issue and PR templates where applicable to make reporting and reviewing easier.
- If you add a new template, place it under `.github/ISSUE_TEMPLATE/` or `.github/PULL_REQUEST_TEMPLATE.md`.

---

## Code of Conduct

- Be respectful, constructive, and inclusive.
- When participating (issues, PRs, reviews), follow the project's code of conduct.
- If you witness or experience unacceptable behavior, contact the maintainers or open an issue flagged for moderation.

---

## Quick contributing checklist

Use this checklist before submitting a PR:

- [ ] I searched existing issues and PRs to avoid duplication.
- [ ] My changes are limited to a single purpose.
- [ ] The project builds locally: `go build ./...`
- [ ] I ran `go fmt ./...` and fixed formatting issues.
- [ ] I added tests where applicable, or described manual testing steps.
- [ ] I updated documentation (README or other docs) if needed.
- [ ] My PR description explains the motivation and testing steps.

---

## Questions or help

If you need help getting started:
- Open an issue and add the `help wanted` label and a short description of where you'd like to contribute.
- Tag maintainers in discussion threads only when necessary.

Thank you for helping improve Todo TUI — contributions make the project better for everyone!
