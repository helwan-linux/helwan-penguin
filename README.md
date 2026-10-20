# Helwan Penguin

An original 2D platform game built with Go and Ebitengine, targeting Linux and Windows.

## Fixed project plan

- 30 levels across six worlds, five levels per world.
- Current playable milestone: five layouts in World 1, Snowy Beginnings.
- Local save support is implemented as a reusable package; gameplay save integration is a later milestone.
- English-first, offline at runtime.

## Requirements

- Go 1.23 or newer
- Ebitengine dependencies downloaded during the first build

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

- Move: Left/Right arrows or A/D
- Jump: Space, Up, or W
- Complete a level: reach the green flag
- Continue/restart after a result: Enter
- Quit: Escape

The five first-world levels use code-drawn artwork, so no external images are required for this milestone. The other world files and systems remain in the fixed tree for later implementation. Build and runtime testing still need to be performed on a machine with the Ebitengine dependencies available.
