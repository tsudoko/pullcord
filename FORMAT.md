Log format
==========

Logs are [TSV][] files. Records are delimited by line feeds.

There are three escape sequences:

    \t U+0009 CHARACTER TABULATION
    \n U+000A LINE FEED
    \\ U+005C REVERSE SOLIDUS

All entries describe actions, i.e. each new message, edit or deletion is a
separate entry.

Fields common for all entries:

    action,type,id

All common fields are required. Unless specified otherwise, other fields can be
empty. Boolean fields contain their name if the value is true or are empty if
it's false.

`Action` can be one of `add`, `edit`, `del`. `Type` corresponds to the Discord
type, such as `message` or `reaction`. `Id` is the ID of the object or the
subject of the action, depending on the entry type.

[TSV]: https://en.wikipedia.org/wiki/Tab-separated_values

## Channel entry types

### `message`

    action,type,id,authorid,editedtime,tts,content

 - `authorid` (required)
 - `editedtime` - ISO 8601 timestamp (µs) of last edit or deletion
 - `tts` (boolean)

Sample timestamp: `2017-06-24T13:06:38.555000+00:00`

### `attachment`

    action,type,messageid,id

 - `messageid` (required)
 - `id` (required)

### `reaction`

    action,type,messageid,userid,emoji

 - `messageid` (required)
 - `userid` - can be empty if there are more than 100 reactions with a given emoji
 - `emoji` (required) - character or `<emojiname>:<emojiid>`

### `embed`

    action,type,messageid,json

 - `messageid` (required)
 - `json` (required) - JSON-encoded [embed contents](https://discordapp.com/developers/docs/resources/channel#embed-object)

### `pin`

    action,type,messageid

 - `messageid` (required)

## Server entry types

### `guild`

    action,type,id,name,icon,splash,ownerid,afkchanid,afktimeout,embeddable,embedchanid,mfalevel

 - `name` (required)
 - `ownerid` (required)
 - `embeddable` (boolean)

### `user`

    action,type,id,name,avatar

### `ban`

    action,type,userid,reason

 - `userid` (required)
 - `reason` - ban reason, can't be seen by regular members by default

### `role`

    action,type,id,color,hoist,pos,perms

 - `color` (required)
 - `hoist` (boolean) - if this role is pinned in the user listing
 - `pos` (required)
 - `perms` (required)

### `roleassign`

    action,type,userid,roleid

 - `userid` (required)
 - `roleid` (required)

### `channel`

    action,type,id,chantype,pos,name,topic

 - `chantype` (required) - `text` or `voice`
 - `pos` (required)
 - `name` (required)

### `permoverwrite`

    action,type,id,overwritetype,allow,deny

 - `overwritetype` (required)
 - `allow` (required)
 - `deny` (required)

### `emoji`

    action,type,id,name,colons

 - `name` (required)
 - `colons` (boolean)
