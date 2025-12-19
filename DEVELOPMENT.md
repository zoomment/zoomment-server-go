# Development Guide

## Quick Start

### Option 1: Auto-reload with Air (Recommended for Development)

Air automatically rebuilds and restarts your server when you change code:

```bash
# Install Air (one-time setup)
go install github.com/air-verse/air@latest

# Run with auto-reload
make dev-air
# OR
air
```

**Benefits:**
- ✅ No manual rebuilds needed
- ✅ Automatic restart on code changes
- ✅ Faster development workflow

### Option 2: Manual Development

If you prefer manual control or don't want to install Air:

```bash
# Run directly (rebuilds on each run)
make dev
# OR
go run cmd/server/main.go
```

**Note:** You'll need to manually stop and restart when you change code.

### Option 3: Build and Run Binary

For production-like testing:

```bash
# Build once
make build

# Run the binary
make run
# OR
./bin/server
```

**Note:** You need to rebuild (`make build`) every time you change code.

---

## When Do You Need to Rebuild?

### ✅ **Automatic (No Rebuild Needed)**
- When using `make dev-air` or `air` - Air handles everything automatically
- When using `make dev` or `go run` - Go rebuilds automatically on each run

### 🔄 **Manual Rebuild Required**
- When using `make build` + `make run` - You must rebuild after code changes
- When deploying to production - Always rebuild for production

---

## Development Workflow Recommendations

### For Active Development:
```bash
# Terminal 1: Run server with auto-reload
make dev-air

# Terminal 2: Run tests
make test

# Terminal 3: Run linter
make lint
```

### For Quick Testing:
```bash
# Simple run (auto-rebuilds each time)
make dev
```

### For Production Build:
```bash
# Build optimized binary
make build

# Test the binary
make run
```

---

## Troubleshooting

### Air not found?
```bash
go install github.com/air-verse/air@latest
```

### Air not detecting changes?
- Check `.air.toml` configuration
- Ensure files are in `include_ext` list
- Check `exclude_dir` doesn't exclude your files

### Build errors?
- Check `build-errors.log` file (created by Air)
- Run `go build ./cmd/server` manually to see errors

---

## File Watching

Air watches for changes in:
- `.go` files
- Template files (`.tpl`, `.tmpl`, `.html`)

Air ignores:
- `bin/`, `tmp/`, `vendor/`, `testdata/` directories
- `*_test.go` files
- Files matching exclude patterns

---

## Performance Tips

1. **Use Air for development** - Fastest iteration cycle
2. **Use `make dev` for quick tests** - No setup needed
3. **Use `make build` for production** - Optimized binary
