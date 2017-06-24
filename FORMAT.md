Log format
==========

Logs are mostly [RFC 4180][]-compliant CSV files. Unlike RFC 4180 CSV, records
are delimited by line feeds and fields aren't allowed to contain record
delimiters. There are two escape sequences: `\n`, which replaces U+000A LINE
FEED, and `\\` for U+005C REVERSE SOLIDUS.

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

[RFC 4180]: https://www.ietf.org/rfc/rfc4180.txt

## Channel entry types

### `message`

    action,type,id,authorid,timestamp,tts,content

 - `authorid` (required)
 - `timestamp` (required)
 - `tts` (boolean)

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

    action,type,id,color,pos,perms

 - `color` (required)
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

TODO: permission overwrites

### `emoji`

    action,type,id,name,colons

 - `name` (required)
 - `colons` (boolean)
