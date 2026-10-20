# Helwan Penguin

An original 2D platform game built with Go and Ebitengine, targeting Linux and Windows.

**Helwan Penguin is currently under active development.** The game is playable, but some features are still being developed and tested. Feedback, bug reports, and suggestions are welcome and can help improve the game.

## Project Status

- **30 planned levels** across six worlds, with five levels per world.
- **Current playable milestone:** Five playable layouts in World 1, Snowy Beginnings.
- Local save support is implemented as a reusable package; integration with gameplay is planned for a later milestone.
- English-first, designed to work offline at runtime.
- The game is under development, so some mechanics, levels, and features may change as development progresses.

## Requirements

- Go 1.23 or newer.
- Ebitengine dependencies, downloaded during the first build.

## Run

```sh
go mod tidy
go run ./cmd/helwan-penguin
```

## Build on Linux

```sh
go build -o helwan-penguin ./cmd/helwan-penguin
```

## Build on Windows PowerShell

```powershell
go build -o helwan-penguin.exe ./cmd/helwan-penguin
```

## Controls

- Move: Left/Right arrows or A/D.
- Jump: Space, Up, or W.
- Complete a level: Reach the green flag.
- Continue or restart after a result: Enter.
- Quit: Escape.

## Graphics and Assets

The five levels currently implemented in World 1 use artwork drawn directly in code, so external image assets are not required for this milestone.

Other world files and game systems remain in the project's fixed structure for future development.

## Feedback and Bug Reports

Feedback is welcome! If you encounter a bug, find a gameplay issue, or have an idea that could improve Helwan Penguin, please report it through the project's GitHub Issues page.

Constructive feedback helps identify problems, improve gameplay, and make future releases better.

## Development Notice

This is an evolving project, and build and runtime testing should be performed on systems with the required Ebitengine dependencies available. Features and behavior may change as development continues.
