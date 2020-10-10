[![Coverage
Status](https://coveralls.io/repos/github/AndreasSko/go-jwlm/badge.svg?branch=master)](https://coveralls.io/github/AndreasSko/go-jwlm?branch=master)

# go-jwlm
A command line tool to easily merge JW Library backups, written in Go.

go-jwlm allows you to merge two .jwlibrary backup files, while giving you
control of the process - your notes are precious, and you shouldn‘t need to
trust a program solving possible merge conflicts for you.

I created this project with the goal of having a tool that is able to work on
multiple operating systems, and even allowing it to be incorporated in other
programs as a library (like an iOS app) in the future. It is - and will be for
quite some time - a work-in-progress project, so I‘m always open for suggestion
and especially reports if you encounter an unexpected behaviour or other bugs. 

The usage is pretty simple: you have one command, you name your backup files -
and press enter. The tool will merge all entries for you. If it encounters a
conflict (like the same note with different content or two markings that
overlap), it will ask you for directions: should it choose the left version or
the right one? After that is finished, you have a nicely merged backup that you
can import into your JW Library App. The first merge process might take some
time because of the number of possible conflicts, depending on how far apart you
backups are. But if you merge them regularly, it should be a matter of seconds
:) 

## Usage
```shell
go-jwlm merge <left-backup> <right-backup> <merged-backup>
```

If a conflict occurs while merging, the tool will ask for directions: should it
choose the left version or the right one. For that, it shows you the actual
entries (I‘m planning to improve that view and add more information, especially
about publications, in the future). If you are not sure what to do, press `?`
for help. 

### Resolve conflicts automatically
Currently, there are three solvers you can use to automatically resolve
conflicts: `chooseLeft`, `chooseRight`, and `chooseNewest` (though the last one
is only usable for Notes). As their names suggest, `chooseLeft` and
`chooseRight` will always choose the same side if a conflict occurs, while
`chooseNewest` always chooses the newest entry. 

You can enable these solvers with the `--bookmarks`, `--markings`, and
`--notes` flags:

```shell
go-jwlm merge <left-backup> <right-backup> <merged-backup> --bookmarks chooseLeft --markings chooseRight --notes chooseNewest
```

The conflict resolvers are helpful for regular merging, when you are 
sure that one sides is always the most up-to-date one. For a first merge, 
it is still recommended to manually solve conflicts, so you don't risk
accidentally overwriting entries.

## Installation 
You can find the compiled binaries for Windows, Linux, and Mac under the
[Release](https://github.com/AndreasSko/go-jwlm/releases) section. 

### Installation using Homebrew (Mac and Linux)
go-jwlm can easily installed using Hombrew:
```shell
brew install andreassko/homebrew-go-jwlm/go-jwlm
```

See the instructions on how to install Homebrew at https://brew.sh

## A word of caution 
It took me a while to trust my own program, but I still keep backups of my
Libraries - and so should you. Go-jwlm is still in alpha-phase, so there is a
possibility that something might get lost because of a yet-to-find bug. So
please keep that in mind and - again - if you found a bug, feel free to open an
issue. 

## Need help?
Something is unclear, you have suggestions for documentation or you found a bug?
Feel free to open an issue. I‘m happy to help, though please be patient if it
takes a while for me to respond :)
