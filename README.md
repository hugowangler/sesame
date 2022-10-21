# sesame

`sesame` is a CLI tool that helps you effortlessly navigate to your favourite `git` projects from the terminal. The main idea it to have a small CLI tool that you can build upon with aliases to make navigation even easier.

The CLI is built using [Cobra](https://github.com/spf13/cobra) together with [Viper](https://github.com/spf13/viper) to handle configuration.

## Installation

`sesame` can be installed by downloading the binary from the [releases page](https://github.com/hugowangler/sesame/releases)

### Building from source

To build from source, you need to have Go installed. After cloning the project you can simply run the following command

```bash
go install
```

## Usage

To see the available commands, run `sesame --help`:

```
Usage:
    sesame [command]

Available Commands:
  add         Adds any found repositories starting from PATH
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  open        Opens a repository that is stored in your config in your browser

Flags:
    --config string   config file (default is $HOME/.config/.sesame.yaml)
    -h, --help            help for sesame
```

### Adding git projects to sesame
You do not have to manually add your git projects to sesame in order to make it possible to open them. You can simply run the following command to add your projects to the config file:

```bash
sesame add PATH
```

The above command will perform a walk in the file tree rooted in the given `PATH` and whenever a `.git` directory is found, the config file is parsed. From the config file the remote origin is Regex matched to construct the URL of the repository. The URL is then stored in the config file together with the repository name.

As a result of this you only have to specify your outermost git project directory and `sesame` will find all the git projects in the file tree rooted in that directory. For example:
    
```bash
sesame add ~/FolderWithAllMyProjects
```

### Opening a git project in your browser
To open a git project in your browser, you can simply run the following command:

```bash
sesame open REPO_NAME
```

### Aliasing
To make navigation event easier, I recommend at least aliasing the `sesame open` command to something shorter. I personally use the following alias:

```bash
alias so="sesame open"
```

To take it even further, you can create aliases to open your most used repositories directly.

### Config file
The config file is stored in the following location by default:

```bash
$HOME/.config/.sesame.yaml
```

It is not created until you add your first repositories.