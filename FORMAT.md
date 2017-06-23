Log format
==========

Logs are [RFC 4180][]-compliant CSV files. All entries describe actions, i.e.
each new message, edit or deletion is a separate entry.

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

    action,type,messageid,userid,emojiid

 - `messageid` (required)
 - `userid` (required)
 - `emojiid` (required)

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

    action,type,id,xaction,name,avatar,game,streaming

 - `xaction` - `leave` or `ban`
    - `leave` - leave or kick
 - `game` - game name
 - `streaming` (boolean)

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
