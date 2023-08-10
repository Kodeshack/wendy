# Changelog

<a name="unreleased"></a>
## [Unreleased]

### Fixes
- Fix returning error on existing dir when created twice in the same run


<a name="v0.6.2"></a>
## [v0.6.2] - 2023-08-09
### Fixes
- Fix modified files triggering ErrorOnExistingFile


<a name="v0.6.1"></a>
## [v0.6.1] - 2023-08-09
### Features
- Add atomicity to file writing


<a name="v0.6.0"></a>
## [v0.6.0] - 2023-08-09
### Refactoring
- Refactor internals to create an in memory file system first


<a name="v0.5.0"></a>
## [v0.5.0] - 2023-08-09
### Features
- Add support for modifiying files
- Add lint failure reporting to step summary
- Add coverage report to step summary


<a name="v0.4.0"></a>
## [v0.4.0] - 2023-07-12
### Features
- Add ErrorOnExistingFile option
- Add go report card
- Add codecove coverage report

### Fixes
- Fix release commit/tagging order


<a name="v0.3.0"></a>
## [v0.3.0] - 2023-05-12
### Features
- Add test-watch task using gotestsum
- Add license info to README
- Add bades to README

### Fixes
- Fix not ignoring non existing dir when CleanDir is set to true

### Refactoring
- Refactor RootDir to OutputDir to be more clear


<a name="v0.2.4"></a>
## [v0.2.4] - 2023-05-12
### Features
- Add dependabot config
- Add missing license file

### Fixes
- Fix error swallowing
- Fix release task


<a name="v0.2.3"></a>
## [v0.2.3] - 2023-05-09
### Fixes
- Fix release task to generate changelog after creating tag


<a name="v0.2.2"></a>
## [v0.2.2] - 2023-05-09
### Fixes
- Fix changelog template to not show unrelease header if there are none


<a name="v0.2.1"></a>
## [v0.2.1] - 2023-05-09
### Features
- Add workflow to automatically create a GitHub release on pushing a version tag
- Add release task with automatic changelog generation


<a name="v0.2.0"></a>
## [v0.2.0] - 2023-05-09
### Features
- Add option to clean output directory before generation
- Add option to error when trying to create an existing directory


<a name="v0.1.1"></a>
## v0.1.1 - 2023-05-06
### Features
- Add FileFromTemplate


[Unreleased]: https://github.com/Kodeshack/wendy/compare/v0.6.2...HEAD
[v0.6.2]: https://github.com/Kodeshack/wendy/compare/v0.6.1...v0.6.2
[v0.6.1]: https://github.com/Kodeshack/wendy/compare/v0.6.0...v0.6.1
[v0.6.0]: https://github.com/Kodeshack/wendy/compare/v0.5.0...v0.6.0
[v0.5.0]: https://github.com/Kodeshack/wendy/compare/v0.4.0...v0.5.0
[v0.4.0]: https://github.com/Kodeshack/wendy/compare/v0.3.0...v0.4.0
[v0.3.0]: https://github.com/Kodeshack/wendy/compare/v0.2.4...v0.3.0
[v0.2.4]: https://github.com/Kodeshack/wendy/compare/v0.2.3...v0.2.4
[v0.2.3]: https://github.com/Kodeshack/wendy/compare/v0.2.2...v0.2.3
[v0.2.2]: https://github.com/Kodeshack/wendy/compare/v0.2.1...v0.2.2
[v0.2.1]: https://github.com/Kodeshack/wendy/compare/v0.2.0...v0.2.1
[v0.2.0]: https://github.com/Kodeshack/wendy/compare/v0.1.1...v0.2.0
