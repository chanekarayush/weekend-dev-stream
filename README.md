# Local Git Heatmap

## Based on [this article](https://flaviocopes.com/go-git-contributions/)

Local Git Heatmap is a [Github](https://github.com) inspired heatmap for a
given users commits.

It renders a rich heatmap of the commit history across multiple repostries that
have been added by the user.

All your data stays local, never leaves your PC.

### Installation

1. Clone this repository to you device.
2. Open terminal in the clone directory
3. Build the go project using the following command:

```bash
go build .
```

4. Install the project using the command:

```bash
go install .
```

5. Run the following command to ensure the cli tool works

```bash
local-git-heatmap -h
```

### Usage

1. Once installed you can understand the flags and their usage by using the `-h`
   flag

```bash
local-git-heatmap -h

```

2. You can add only 1 directory at once. and all the added directories will be
   stored at `~/.local-git-heatmap`

```bash
local-git-heatmap --add-dir=<dir-name> -email=<email@example.com>
```
