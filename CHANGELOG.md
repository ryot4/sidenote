# ChangeLog

## Unreleased

Miscellaneous changes:

 * Update to Go 1.22.1

## v0.1.11 - 2023-06-25

Improvements and fixes:

 * ls: Fix to print file timestamps when -t is specified with a regular file

Miscellaneous changes:

 * Update to Go 1.19.10

## v0.1.10 - 2022-09-24

Breaking changes:

 * Rename -version option to -V

Improvements and fixes:

 * serve: Fix status code logging when accessing directories

Miscellaneous changes:

 * Update to Go 1.19

## v0.1.9 - 2021-02-28

Breaking changes:

 * Shell completion scripts are now embedded in the binary and can be loaded with `sidenote completion` command.
   See the README file about how to enable command line completion.

Miscellaneous changes:

 * Update to Go 1.16

## v0.1.8 - 2021-01-18

Improvements and fixes:

 * edit: Fix nil pointer dereference on editing new files

Miscellaneous changes:

 * Create the parent directory before archiving binary releases

## v0.1.7 - 2021-01-18

Breaking changes:

 * edit: Remove -x option

New features:

 * Add exec command

Improvements and fixes:

 * import: Preserve file modes when importing existing files
 * rm: Add support for removing multiple notes
 * show: Add support for showing multiple notes

## v0.1.6 - 2020-12-30

Breaking changes:

 * edit: Do not create the parent directory by default (add -p option)
 * ls: Make -t imply -l

New features:

 * Add serve command

Improvements and fixes:

 * Fix notes discovery to handle -d option with relative paths properly

Miscellaneous changes:

 * Update to Go 1.15
 * Convert documents from AsciiDoc to Markdown
 * Add CHANGELOG (this file)

## v0.1.5 - 2020-05-30

Improvements and fixes:

 * Print messages after successful file operations
 * Improve error messages

## v0.1.4 - 2020-05-23

New features:

 * Add import command

## v0.1.3 - 2020-04-10

New features:

 * Add show command

Improvements and fixes:

 * Improve help messages

## v0.1.2 - 2020-04-06

Improvements and fixes:

 * Add -L option to path command
 * Fix path completion in bash

## v0.1.1 - 2020-02-16

 * Initial release
