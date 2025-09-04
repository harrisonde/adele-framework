SHELL := /bin/bash

export PACKAGE_BRANCH = main
export PACKAGE_PUBLICATION_TAG ?=
export PACKAGE_PUBLICATION_TAG_NEXT ?=
export OUT ?=


## builds the adele command line tool
.SILENT:
build\:adele:
	@go build -o ./bin/adele ./cli/adele

# Help command for build commands
.SILENT:
build\:help:
	@echo "═══════════════════════════════════════════════════════════════════"
	@echo "                          BUILD COMMANDS HELP"
	@echo "═══════════════════════════════════════════════════════════════════"
	@echo ""
	@echo "🔨 AVAILABLE BUILD COMMANDS:"
	@echo "  make build:help       - Show this help documentation"
	@echo "  make build:adele      - Build the Adele CLI tool"
	@echo ""
	@echo "📦 BUILD DETAILS:"
	@echo ""
	@echo "  build:adele"
	@echo "  ├── 🎯 Purpose: Compiles the Adele command-line tool"
	@echo "  ├── 📁 Source:  ./cli/adele (Go source code)"
	@echo "  ├── 📤 Output:  ./bin/adele (executable binary)"
	@echo "  └── ⚙️  Action:  go build -o ./bin/adele ./cli/adele"
	@echo ""
	@echo "🔄 TYPICAL WORKFLOW:"
	@echo "  1. Make your changes to CLI code in ./cli/adele/"
	@echo "  2. Build the tool:"
	@echo "     make build:adele"
	@echo "  3. Test your CLI tool:"
	@echo "     ./bin/adele --help"
	@echo "     ./bin/adele [your-command]"
	@echo ""
	@echo "💡 TIPS:"
	@echo "  • The binary is created in ./bin/ directory"
	@echo "  • Add ./bin to your PATH to use 'adele' command globally"
	@echo "  • Run build:adele after any CLI code changes"
	@echo "  • Use 'go run ./cli/adele' for development without building"
	@echo ""
	@echo "🚨 TROUBLESHOOTING:"
	@echo "  • Build errors → Check Go syntax in ./cli/adele/"
	@echo "  • Permission denied → chmod +x ./bin/adele"
	@echo "  • Command not found → Use ./bin/adele or add to PATH"
	@echo ""
	@echo "═══════════════════════════════════════════════════════════════════"

## package tests
.SILENT:
test\:all:
	@go clean -testcache
	make test:cache test:cli test:database test:filesystem test:helpers test:httpservertest:logger test:middleware test:mailer test:middleware test:mux test:session test:render test:rpcserver
test\:cache:
	@go test ./cache/...
test\:cli:
	@go test ./cli/adele/...
test\:database:
	@go test ./database/...
test\:filesystem:
	@go test ./filesystem/...
test\:helpers:
	@go test ./helpers
test\:httpserver:
	@go test ./httpserver
test\:logger:
	@go test ./logger
test\:middleware:
	@go test ./middleware
test\:mailer:
	@go test ./mailer
test\:mux:
	@go test ./mux
test\:session:
	@go test ./session
test\:render:
	@go test ./render
test\:rpcserver:
	@go test ./rpcserver
## coverage: displays test coverage
test\:coverage:
	@go test -cover ./...
test\:coverage\:browser:
	@go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

# Help command for test commands
.SILENT:
test\:help:
	@echo "═══════════════════════════════════════════════════════════════════"
	@echo "                          TESTING COMMANDS HELP"
	@echo "═══════════════════════════════════════════════════════════════════"
	@echo ""
	@echo "🧪 AVAILABLE TEST COMMANDS:"
	@echo "  make test:help                - Show this help documentation"
	@echo "  make test:all                 - Run core package tests (logger, mailer, middleware, mux, session)"
	@echo "  make test:coverage            - Run all tests with coverage summary"
	@echo "  make test:coverage:browser    - Run tests and open detailed coverage in browser"
	@echo ""
	@echo "📦 INDIVIDUAL PACKAGE TESTS:"
	@echo "  make test:cache               - Test caching functionality"
	@echo "  make test:cli                 - Test CLI tool functionality"
	@echo "  make test:database            - Test database operations"
	@echo "  make test:filesystem          - Test filesystem operations"
	@echo "  make test:helpers             - Test helper utilities"
	@echo "  make test:httpserver          - Test HTTP server functionality"
	@echo "  make test:logger              - Test logging system"
	@echo "  make test:middleware          - Test middleware components"
	@echo "  make test:mailer              - Test email functionality"
	@echo "  make test:mux                 - Test HTTP routing"
	@echo "  make test:session             - Test session management"
	@echo "  make test:render              - Test template rendering"
	@echo "  make test:rpcserver           - Test RPC server functionality"
	@echo ""
	@echo "🎯 TEST CATEGORIES:"
	@echo ""
	@echo "  🔄 test:all"
	@echo "  ├── Clears test cache for fresh results"
	@echo "  ├── Runs: logger, mailer, middleware, mux, session tests"
	@echo "  └── Good for: Core functionality validation"
	@echo ""
	@echo "  📊 test:coverage"
	@echo "  ├── Runs all package tests with coverage analysis"
	@echo "  ├── Shows coverage percentage per package"
	@echo "  └── Good for: Quick coverage overview"
	@echo ""
	@echo "  🌐 test:coverage:browser"
	@echo "  ├── Generates detailed HTML coverage report"
	@echo "  ├── Opens coverage.out in your default browser"
	@echo "  └── Good for: Detailed line-by-line coverage analysis"
	@echo ""
	@echo "🔄 TYPICAL WORKFLOWS:"
	@echo ""
	@echo "  Quick validation:"
	@echo "    make test:all"
	@echo ""
	@echo "  Full test suite:"
	@echo "    make test:coverage"
	@echo ""
	@echo "  Detailed analysis:"
	@echo "    make test:coverage:browser"
	@echo ""
	@echo "  Specific package:"
	@echo "    make test:database"
	@echo "    make test:httpserver"
	@echo ""
	@echo "💡 TIPS:"
	@echo "  • Individual package tests run faster than test:all"
	@echo "  • Use test:coverage:browser to find untested code paths"
	@echo "  • Test cache is cleared in test:all for reliable results"
	@echo "  • Coverage reports help identify areas needing more tests"
	@echo ""
	@echo "🚨 TROUBLESHOOTING:"
	@echo "  • Test failures → Check specific package: make test:[package]"
	@echo "  • Cached results → Run: go clean -testcache"
	@echo "  • Coverage not opening → Check if coverage.out exists"
	@echo "  • Slow tests → Run individual packages instead of test:all"
	@echo ""
	@echo "📁 PACKAGE STRUCTURE:"
	@echo "  ./cache/       → Caching and Redis functionality"
	@echo "  ./cli/adele/   → Command-line interface"
	@echo "  ./database/    → Database connections and operations"
	@echo "  ./filesystem/  → File and directory operations"
	@echo "  ./helpers      → Utility functions"
	@echo "  ./httpserver   → HTTP server implementation"
	@echo "  ./logger       → Logging system"
	@echo "  ./middleware   → HTTP middleware components"
	@echo "  ./mailer       → Email sending functionality"
	@echo "  ./mux          → HTTP request routing"
	@echo "  ./session      → Session management"
	@echo "  ./render       → Template rendering engine"
	@echo "  ./rpcserver    → RPC server implementation"
	@echo ""
	@echo "═══════════════════════════════════════════════════════════════════"

## Release workflow to tag and push to the current branch
.SILENT:
release\:verify:
	@if [[ ! $(PACKAGE_PUBLICATION_TAG_NEXT) =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$$ ]]; then \
		echo "Error: Tag '$(PACKAGE_PUBLICATION_TAG_NEXT)' does not follow semantic versioning format"; \
		echo "Expected: vX.Y.Z[-prerelease][+buildmeta]"; \
		echo "Examples: v1.0.0, v10.21.34-alpha.1, v1.2.3+build.123, v2.0.0-beta.1+exp.sha.5114f85"; \
		exit 1; \
	fi

	# Check if tag already exists
	@if git rev-parse $$PACKAGE_PUBLICATION_TAG_NEXT >/dev/null 2>&1; then \
		echo "Error: Tag '$$PACKAGE_PUBLICATION_TAG_NEXT' already exists"; \
		exit 1; \
	fi

	@echo "The next package release will be tagged with $(PACKAGE_PUBLICATION_TAG_NEXT)"

	# Check if working directory is clean
# 	@if ! git diff-index --quiet HEAD --; then \
# 		echo "Error: working directory has uncommitted changes"; \
# 		exit 1; \
# 	fi

# 	# Check if we're on the expected branch
# 	CURRENT_BRANCH=$$(git branch --show-current); \
# 	if [[ ! "$$CURRENT_BRANCH" =~ ^(main|release/.*|ci/.*|hotfix/.*)$$ ]]; then \
# 		echo "Error: Branch '$$CURRENT_BRANCH' not allowed. Must be main, ci/*, release/*, or hotfix/*"; \
# 		exit 1; \
# 	fi

# 	@read -p "Do you wish to proceed with the release? [y/N] " ans && ans=$${ans:-N} ; \
# 	if [ $${ans} = y ] || [ $${ans} = Y ]; then \
# 		echo "Creating and pushing tag: $$PACKAGE_PUBLICATION_TAG_NEXT"; \
# 		git tag $$PACKAGE_PUBLICATION_TAG_NEXT; \
# 		git push origin $$PACKAGE_PUBLICATION_TAG_NEXT; \
# 		echo "✓ Release $$PACKAGE_PUBLICATION_TAG_NEXT created successfully"; \
# 	else \
# 		echo "Release cancelled"; \
# 		exit 1; \
# 	fi

.SILENT:
release\:pull:
	@echo "Checking repository status..."

	# Fetch latest changes
	if ! git fetch origin $(PACKAGE_BRANCH); then \
		echo "Error: Failed to fetch from origin"; \
		exit 1; \
	fi

	# Check if local branch is behind
	LOCAL=$$(git rev-parse HEAD); \
	REMOTE=$$(git rev-parse origin/$(PACKAGE_BRANCH)); \
	if [[ "$$LOCAL" != "$$REMOTE" ]]; then \
		echo "Local branch is behind origin - pulling changes..."; \
		git pull origin $(PACKAGE_BRANCH); \
	else \
		echo "✓ Repository is up to date"; \
	fi

release\:preamble:
	@echo "Please enter a SemVer-compatible version tag for this release."
	@echo ""
	@echo "🏷️  SEMANTIC VERSIONING FORMAT:"
	@echo "  Tags must follow: vMAJOR.MINOR.PATCH[-prerelease][+buildmeta]"
	@echo ""
	@echo "  ✅ Valid examples:"
	@echo "    v1.0.0                    - Basic release"
	@echo "    v10.21.34                 - Multi-digit versions"
	@echo "    v1.0.0-alpha              - Prerelease"
	@echo "    v1.0.0-alpha.1            - Prerelease with number"
	@echo "    v1.0.0-beta.2+build.123   - Prerelease + build metadata"
	@echo "    v2.1.0+exp.sha.5114f85    - Build metadata only"
	@echo ""

#$(eval PACKAGE_PUBLICATION_TAG_NEXT=$(shell read -p "Enter new tag: " tag; echo $$tag))
release\:capture:
	@NEXT_TAG=$$(read -p "Enter new tag: " tag; echo $$tag); \
	 export NEXT_TAG; \
	 echo "Selected tag: $$NEXT_TAG"; \
	 make release:verify PACKAGE_PUBLICATION_TAG_NEXT=$$NEXT_TAG 2>/dev/null || exit 1



.SILENT:
release\:get-current-tag:
	$(eval LATEST_TAG=$(shell git describe --tags --abbrev=0 2>/dev/null || echo "No current tags found"))
	@echo "Current tag: $(LATEST_TAG)"

# Combined release target for convenience
#@make release:pull release:verify
.SILENT:
release:
	@make release:preamble
	@make release:get-current-tag
	@make release:capture
	@echo "✓ Release process completed"


# Help command with release documentation
.SILENT:
release\:help:
	@echo "═══════════════════════════════════════════════════════════════════"
	@echo "                        RELEASE WORKFLOW HELP"
	@echo "═══════════════════════════════════════════════════════════════════"
	@echo ""
	@echo "📋 AVAILABLE COMMANDS:"
	@echo "  make release:help     - Show this help documentation"
	@echo "  make release:pull     - Pull latest changes from origin"
	@echo "  make release:verify   - Verify and create release tag (interactive)"
	@echo "  make release          - Run pull + verify in sequence"
	@echo ""
	@echo "🏷️  SEMANTIC VERSIONING FORMAT:"
	@echo "  Tags must follow: vMAJOR.MINOR.PATCH[-prerelease][+buildmeta]"
	@echo ""
	@echo "  ✅ Valid examples:"
	@echo "    v1.0.0                    - Basic release"
	@echo "    v10.21.34                 - Multi-digit versions"
	@echo "    v1.0.0-alpha              - Prerelease"
	@echo "    v1.0.0-alpha.1            - Prerelease with number"
	@echo "    v1.0.0-beta.2+build.123   - Prerelease + build metadata"
	@echo "    v2.1.0+exp.sha.5114f85    - Build metadata only"
	@echo ""
	@echo "  ❌ Invalid examples:"
	@echo "    1.0.0        - Missing 'v' prefix"
	@echo "    v1.0         - Missing patch version"
	@echo "    v1.0.0-      - Empty prerelease"
	@echo "    v1.0.0+      - Empty build metadata"
	@echo ""
	@echo "🌿 SUPPORTED BRANCHES:"
	@echo "  📦 Production releases: $(PACKAGE_BRANCH)"
	@echo "  🧪 Test releases:       ci/* (any branch starting with 'ci/')"
	@echo ""
	@echo "🔄 TYPICAL WORKFLOW:"
	@echo "  1. Set your target version:"
	@echo "     export PACKAGE_PUBLICATION_TAG_NEXT=v1.2.3"
	@echo ""
	@echo "  2. For PRODUCTION releases:"
	@echo "     git checkout $(PACKAGE_BRANCH)"
	@echo "     make release"
	@echo ""
	@echo "  3. For TEST releases:"
	@echo "     git checkout ci/your-test-branch"
	@echo "     make release"
	@echo ""
	@echo "✅ PRE-RELEASE CHECKS:"
	@echo "  The workflow automatically verifies:"
	@echo "  • Tag follows semantic versioning format"
	@echo "  • Tag doesn't already exist"
	@echo "  • Working directory is clean (no uncommitted changes)"
	@echo "  • You're on the correct branch ($(PACKAGE_BRANCH) or ci/*)"
	@echo "  • Repository is up to date with origin"
	@echo ""
	@echo "🚨 TROUBLESHOOTING:"
	@echo "  • 'Tag already exists' → Check: git tag -l | grep v1.2.3"
	@echo "  • 'Uncommitted changes' → Commit or stash your changes"
	@echo "  • 'Wrong branch' → Switch to $(PACKAGE_BRANCH) or ci/* branch"
	@echo "  • 'Behind origin' → Run: make release:pull"
	@echo ""
	@echo "💡 TIPS:"
	@echo "  • Use ci/ branches to test the release process safely"
	@echo "  • Prerelease versions (v1.0.0-alpha) won't trigger production deploys"
	@echo "  • Build metadata (+build.123) is ignored by version precedence"
	@echo "  • Always test on ci/ branches before releasing to $(PACKAGE_BRANCH)"
	@echo ""
	@echo "📚 MORE INFO:"
	@echo "  • Semantic Versioning: https://semver.org/"
	@echo "  • Git Tagging: https://git-scm.com/book/en/v2/Git-Basics-Tagging"
	@echo ""
	@echo "═══════════════════════════════════════════════════════════════════"
