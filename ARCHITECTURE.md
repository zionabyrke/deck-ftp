# DECK
## Deploy Every Commit, Kontinuous

## Scope

| In v1 | Out of v1 |
|---|---|
| FTP adapter | SFTP / other adapters |
| Manifest diffing | Multi-target config |
| `deck.yaml` config | GitHub Actions |
| Local git hook trigger | Webhook daemon |
| Styled CLI output | GUI |

## Components

```mermaid
flowchart TD
    CLI[CLI commands] --> CFG[Config]
    CLI --> DEP[Deployer]
    DEP --> MAN[Manifest]
    DEP --> ADP[Adapter interface]
    ADP --> FTP[FTP Adapter]
    FTP --> HOST[(InfinityFree FTP)]
    CLI --> HOOK[Hook Installer]
    HOOK -.writes.-> GH[.git/hooks/post-commit]
    GH -.runs.-> CLI
```

## Types

```mermaid
classDiagram
    class Adapter {
        <<interface>>
        +Connect() error
        +List(path) []FileInfo
        +Upload(local, remote) error
        +Delete(remote) error
        +Mkdir(remote) error
        +Close() error
    }
    class FTPAdapter
    class Config {
        +Host string
        +Port int
        +User string
        +PassEnv string
        +LocalDir string
        +RemoteDir string
        +Ignore []string
        +BuildCmd string
    }
    class Manifest {
        +Files map~string,string~
        +Load(path) Manifest
        +Save(path) error
        +Diff(other) DiffResult
    }
    class DiffResult {
        +Added []string
        +Changed []string
        +Deleted []string
    }
    class Deployer {
        +Run(dryRun bool) DiffResult
    }

    Adapter <|.. FTPAdapter
    Deployer --> Adapter
    Deployer --> Config
    Deployer --> Manifest
    Manifest --> DiffResult
```

## Flow — `deck push`

```mermaid
sequenceDiagram
    participant U as User
    participant C as CLI
    participant M as Manifest
    participant D as Deployer
    participant A as FTPAdapter
    participant H as InfinityFree

    U->>C: deck push
    C->>M: scan local + load previous
    M-->>C: DiffResult
    C->>D: Run(diff)
    D->>A: Connect()
    A->>H: login
    loop changed files
        D->>A: Upload / Delete
        A->>H: STOR / DELE
    end
    D->>A: Close()
    D->>M: Save
    D-->>U: summary
```

## Flow — auto trigger

```mermaid
sequenceDiagram
    participant U as User
    participant G as Git
    participant Hook as post-commit
    participant C as DECK CLI

    U->>G: commit (main)
    G->>Hook: run
    Hook->>Hook: branch == main ?
    Hook->>C: deck push
    C-->>U: summary
```

## Commands

| Command | Effect |
|---|---|
| `deck init` | write `deck.yaml` |
| `deck diff` | dry run, no network |
| `deck push` | deploy diff |
| `deck status` | manifest vs config check |
| `deck install-hook` | wire `post-commit` on `main` |

## Roadmap (post-v1)

| Version | Adds |
|---|---|
| v2 | GitHub Actions target, secrets via env |
| v3 | SFTP / cPanel API adapters |
| v4 | Webhook daemon, multi-target config |
| v5 | GUI shell over same core |
