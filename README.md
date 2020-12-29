# sidenote

sidenote is a command line utility for managing plain text notes per working directory.

## Installation

### Binary releases (linux/amd64 only for now)

See the [releases page](https://github.com/ryot4/sidenote/releases).
Download and extract the archive, and then put the binary into the directory listed in `$PATH`.

### Build from source

Make sure [Go distribution is installed](https://golang.org/doc/install), and then run `go get`.

    $ go get -u github.com/ryot4/sidenote

This installs the binary into `$GOPATH/bin`. Make sure `$GOPATH/bin` is listed in `$PATH`.

### Command-line completions

See `completions` directory. (Currently only Bash is supported. Contributions are welcome)

## Quick start

### Initialization

To prevent the working directory from being cluttered with text files, sidenote uses a dedicated
directory (`.notes`) for notes. First of all, you need to prepare it with `init` subcommand.

    $ sidenote init
    initialized .notes
    $ ls -dl .notes
    drwxr-xr-x 2 ryot4 ryot4 4096 Feb  9 19:03 .notes

If you want to store files outside the current working directory, you can initialize `.notes`
as a symbolic link with `init -l`. The target directory is created if it does not exist.

    $ sidenote init -l ~/Documents/notes
    initialized .notes (-> /home/ryot4/Documents/notes)
    $ ls -l .notes
    lrwxrwxrwx 1 ryot4 ryot4 31 Feb  9 19:03 .notes -> /home/ryot4/Documents/notes

When `.notes` does not exist in the current working directory, sidenote searchs the directory
hierarchy upward. Therefore you only need to run `init` at the top-level working directory.

    $ sidenote path            # Print the relative path to the .notes directory.
    .notes
    $ mkdir subdir; cd subdir
    $ sidenote path
    ../.notes                  # Notes in the parent directory is used.

### Editing notes

You can use your favorite text editor to edit notes. (sidenote refers `$VISUAL` and `$EDITOR`)

    $ sidenote edit todo.txt                # This opens .notes/todo.txt with $EDITOR.
    $ sidenote edit -p issues/issue-123.md  # You can create subdirectories. (-p creates the directory if not exists)

Filenames can be generated based on the current date. (subset of `strftime(3)` format is available)

    $ date +'%Y-%m-%d'
    2020-02-09
    $ sidenote edit -f '%Y-%m-%d.md'             # This opens 2020-02-09.md.
    $ export SIDENOTE_NAME_FORMAT='%Y-%m-%d.md'
    $ sidenote edit                              # Same as above, but no need to specify the format every time.

### Importing existing files

Instead of creating new files, you can also import existing files with `import` subcommand.

    $ sidenote import note.txt
    imported note.txt
    $ sidenote import -d todo.txt                # This removes the original file after import.
    imported todo.txt
    $ sidenote import hello.txt hello-world.txt  # You can specify the name of the imported file.
    imported hello-world.txt

### Displaying notes

To display notes, use `cat` or `show` subcommand:

    $ sidenote cat todo.txt
    $ sidenote show todo.txt  # This opens todo.txt with $PAGER.

### File operations

You can list and remove notes with `ls` and `rm` subcommands, respectively.

    $ sidenote ls
    todo.txt
    $ sidenote ls -l
    Feb  9 21:37 todo.txt
    $ sidenote rm todo.txt
    removed todo.txt

Of course you can use standard command line utilities as well.

    $ find .notes -name todo.txt            # Find notes named todo.txt.
    $ find $(sidenote path) -name todo.txt  # Same as above, but this is convenient when you are not in the top-level directory.

    $ cd $(sidenote path)   # You can operate files as usual after this.
    $ mv todo.txt done.txt

### Searching

Searching can be done with a combination of `path` subcommand and existing searching commands
such as `grep`.

    $ grep -R pattern $(sidenote path)          # Search from all files.
    $ grep -R pattern $(sidenote path 2020/02)  # Search from files in 2020/02/.

### And more

For the full list of subcommands, options and environment varibles, see `sidenote -h` and
`sidenote <command> -h`.

## Tips

### Ignore notes in Git globally

To ignore notes in Git globally, add `.notes` to the file specified by `core.excludesfile`
(by default, this is `~/.config/git/ignore`) in Git config.

    $ mkdir -p ~/.config/git
    $ echo .notes >> ~/.config/git/ignore

### Store files in the directory other than .notes

You can use `-d` option or `SIDENOTE_DIR` environment variable to specify the directory
where files are stored.

    $ sidenote -d .mynotes init           # Use .mynotes for notes.
    $ sidenote -d .mynotes edit todo.txt
    $ export SIDENOTE_DIR=.mynotes
    $ sidenote edit todo.txt              # Same as above.

You can also use absolute paths:

    $ sidenote -d ~/Documents/notes ls  # List notes in ~/Documents/notes.

### Dotfiles are ignored

You cannot use filenames beginning with a dot (`.`), and such files are ignored by `ls` subcommand.

    $ sidenote edit .test
    error: path .test contains dotfile
    $ sidenote edit dir/.test
    error: path dir/.test contains dotfile

    $ git --git-dir=$(sidenote path)/.git init -q  # Initialize the Git repository to put notes under version control.
    $ sidenote ls                                  # This does not list .notes/.git.

### Share the same notes directory from multiple working directories

With `init -l`, you can refer the same directory from multiple working directories:

    $ cd /path/to/project
    $ sidenote init -l ~/Documents/notes/coding
    $ sidenote edit useful-knowledge.adoc

    $ cd /path/to/another-project
    $ sidenote init -l ~/Documents/notes/coding  # Use the same directory.
    $ sidenote ls
    useful-knowledge.adoc
    ...