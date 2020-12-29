# Shell command-line completions

## Bash

Copy `sidenote.bash` to `/etc/bash_completion.d/`. The directory may vary depending on the distribution or OS.

    # install -o root -g root -m 0644 sidenote.bash /etc/bash_completion.d/sidenote

Instead of installing system-wide, you can source the file from your bashrc:

    . /path/to/sidenote.bash