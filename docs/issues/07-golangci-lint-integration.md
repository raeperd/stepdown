## Parent PRD

#1

## What to build

Submit an upstream PR to golangci-lint to integrate `stepdown` as a supported linter.

## Implementation

Based on recvcheck's integration pattern in golangci-lint, the following files need changes:

1. **`pkg/config/linters_settings.go`** — Add `StepdownSettings` struct:
   ```go
   type StepdownSettings struct {
       Exclusions []string `mapstructure:"exclusions"`
   }
   ```

2. **`pkg/golinters/stepdown/stepdown.go`** — Create wrapper (~15 lines):
   ```go
   func New(settings *config.StepdownSettings) *goanalysis.Linter {
       var cfg stepdown.Settings
       if settings != nil {
           cfg.Exclusions = settings.Exclusions
       }
       return goanalysis.
           NewLinterFromAnalyzer(stepdown.NewAnalyzer(cfg)).
           WithLoadMode(goanalysis.LoadModeTypesInfo)
   }
   ```

3. **`pkg/lint/lintersdb/builder_linter.go`** — Register linter:
   ```go
   linter.NewConfig(stepdown.New(&cfg.Linters.Settings.Stepdown)).
       WithSince("v2.X.0").
       WithLoadForGoAnalysis().
       WithURL("https://github.com/raeperd/stepdown"),
   ```

4. **`pkg/commands/internal/migrate/migrate_linter_names.go`** — Add metadata:
   ```go
   {Name: "stepdown", Presets: []string{"style"}, Slow: true},
   ```

5. **`.golangci.reference.yml`** — Add config example

## Acceptance criteria

- [ ] `stepdown` appears in `golangci-lint help linters`
- [ ] `golangci-lint run` with `stepdown` enabled detects violations
- [ ] `linters-settings.stepdown.exclusions` maps to `Settings.Exclusions`
- [ ] golangci-lint PR is merged upstream

## Blocked by

- Blocked by issues 01 through 06 (all feature slices)

## User stories addressed

- User story 6: enable stepdown in golangci-lint
