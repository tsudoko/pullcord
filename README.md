pullcord
========

`Pullcord` is a Discord archiver. It downloads channel logs, server logs,
attachments, avatars, server icons, server splashes and emoji.

Install
-------
```
go get github.com/tsudoko/pullcord/cmd/pullcord
```

This installs `pullcord` to `$GOPATH/bin`, by default `~/go/bin`.

How to use
----------

Right now the only implemented mode is `history`. It is used for downloading
channel history, server history and all related files. A `realtime` mode is
planned.

By default `pullcord` downloads data from every channel and server the account
is connected to, with exception of DMs. To fine-tune this behavior, filtering
options such as `-c`, `-s`, `-C`, `-S` and `-dm` can be used. All files are
downloaded to the current working directory; creating a new empty directory is
recommended.

`Pullcord` exits as soon as it encounters any error.

Basic usage:

    pullcord -<mode> [auth_options] [filter_options]

All options can be seen by running `pullcord -h`.

Log format
----------

See [FORMAT.md](FORMAT.md).
