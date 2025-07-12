# fsner

`fsner` is a command-line file system scanner written in Go. It recursively searches file paths from a given root directory and filters results using regular expressions. It also supports bash autocompletion.

## Features

- Regex-based file name search
- Recursive directory scanning
- Option to limit number of results
- Optional file size display
- Bash autocompletion for flags and folder suggestions

## Installation

```bash
git clone https://github.com/Ege-Guler/fsner.git
cd fsner/cmd/fsner
go build -o fsner
sudo mv fsner /usr/local/bin
```
### Bash Autocompletion

To enable autocompletion for `fsner` source `.fsner-completion.bash` in your `.bashrc`
```
echo 'source /full/path/to/.fsner_completion.bash' >> ~/.bashrc
source ~/.bashrc

```

# Usage

`fsner --pattern "<regex>" --root <directory> [options]`

## Required Flags

- `--pattern`, `-p <value>`: Regex pattern to match filenames
- `--root`, `-r <value>`: Root directory to begin scanning

## Optional Flags

- `--max`, `-m <value>`: Maximum number of results to return (default: unlimited)
- `--file-size`, `-s`: Display file sizes in MiB for larger files, KiB for smaller ones
- `--verbose`, `-v`: Enable verbose output
- `--version`, `-V`: Print version info
- `--help`, `-h`: Show help message

## Example

```bash
fsner --pattern ".*\.go$" --root ~/projects --max 10 --file-size
```