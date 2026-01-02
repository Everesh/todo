# Pull Request Template

<!--
Thank you for contributing to Todo TUI! Please fill out the sections below to help maintainers review your change faster.
You can remove sections that don't apply.
-->

## Summary

<!--
Provide a concise summary of the change and the motivation behind it.
Example: "Add HH:MM timer parsing and tests" or "Fix crash when todos.json is empty".
-->
Type: (fix/feat/docs/refactor/chore)  
Scope: (e.g. timer, ui, data, build)  

What changed and why:

## Related Issues

<!--
Link any related issues or pull requests.
Use "Closes #<issue_number>" to auto-close an issue when the PR is merged.
-->
Fixes: #<!-- issue number -->

## How to Test

<!--
Provide step-by-step instructions to manually validate the change.
If applicable include sample commands, test data, or screenshots/GIFs showing the UI behavior.
-->
1. Build/run steps:
   ```bash
   # example
   go build ./...
   ./build/todo
   ```
2. Steps to reproduce expected behavior:
   - ...
3. Edge cases tested:
   - ...

## Screenshots / GIFs (if applicable)

<!--
Paste or link to screenshots/GIFs demonstrating UI changes.
-->

## Checklist

<!-- Replace [ ] with [x] when completed -->
- [ ] I read the contribution guidelines in `CONTRIBUTING.md`
- [ ] My changes follow the project's coding style and patterns
- [ ] I ran `go fmt ./...` and `go vet ./...`
- [ ] The project builds locally: `go build ./...`
- [ ] I added or updated tests for any new behavior
- [ ] All tests pass locally
- [ ] I updated documentation (README or relevant docs) if needed
- [ ] I updated `CHANGELOG.md` or included release notes if applicable

## Types of Changes

<!--
Indicate which type of change this PR introduces. Check all that apply.
-->
- [ ] Bugfix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to change)
- [ ] Documentation update

## Implementation Notes (optional)

<!--
Any notes for reviewers: design decisions, trade-offs, known limitations, backward compatibility concerns, dependency upgrades, etc.
-->

## Security & Privacy

<!--
If your change touches data storage, networking, or IPC, mention privacy/security considerations.
-->

## Reviewer Suggestions (optional)

<!--
If you want specific maintainers to review, list them here.
-->

---

Thank you for your contribution! We appreciate clear, small PRs that are easy to review — it helps us move faster.
