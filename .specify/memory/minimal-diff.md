# Skill: Minimal Diff & Surgical Changes

When modifying existing files, you MUST adhere to the Principle of Minimal Diff to keep code changes focused, safe, and easy to review.

## Minimal Diff Principles
- **Surgical Edits**: Touch ONLY the exact lines of code required to implement the spec or fix the issue.
- **No Unrelated Formatting**: Do NOT reformat untouched lines, reorder imports, or alter whitespace outside the immediate scope of your task.
- **Preserve Local Style**: Match the formatting, naming conventions, and idioms of the surrounding code, even if it differs from your default preference.
- **Avoid Unrequested Refactoring**: Do NOT fix adjacent code smells, rename untouched variables, or "clean up" unrelated functions unless explicitly requested in the task.
- **Focus Git Diffs**: Ensure every line present in `git diff` maps directly to an active task in `specs/`.