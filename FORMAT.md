Log format
==========

Logs are [TSV][] files. Records are delimited by line feeds. Since field
contents are type-dependent, the first record does not contain field names.
There are no quoted fields, the following escape sequences are used instead:

    \t U+0009 CHARACTER TABULATION
    \n U+000A LINE FEED
    \\ U+005C REVERSE SOLIDUS

All entries describe actions, i.e. each new message, edit or deletion is a
separate entry.

Keep in mind some records had fields added to them after the initial version
of the format was developed, so you shouldn't assume all fields listed here
are always present. If a new field is introduced, it's always added after the
last field of the current record, the order of existing ones does not change.

Fields common for all entries:

    time,fetchtype,action,type,id

All common fields are required. Unless specified otherwise, other fields can be
empty. Boolean fields contain their name if the value is true or are empty if
it's false.

`Time` is the entry creation time. `Fetchtype` describes if the entry has been
retrieved from the history (`history`) or written as soon as the event it
describes had happened (`realtime`). `Action` is `add` or `del`. `Type`
corresponds to the Discord type, such as `message` or `reaction`. `Id` is the
ID of the object of the action.

If an object has been `add`ed previously and a new `add` entry with the same ID
appears, it's considered an update to the existing object. `Del` for an object
that hasn't been `add`ed at any time means the object had existed in the past,
but it wasn't present at the time of fetching.

[TSV]: https://en.wikipedia.org/wiki/Tab-separated_values

## Channel entry types

### `message`

    time,fetchtype,action,type,id,authorid,editedtime,tts,content,webhook,usernameoverride,avataroverride,msgtype

 - `authorid` (required)
 - `editedtime` - ISO 8601 timestamp (µs) of last edit
 - `tts` (boolean)
 - `webhook` (boolean)
 - `usernameoverride` - username shown if the author is a webhook
 - `avataroverride` - avatar shown if the autor is a webhook
 - `msgtype` - one of `` (empty string), `recipient_add`, `recipient_remove`, `call`, `channel_name_change`, `channel_icon_change`, `channel_pinned_message`, `guild_member_join`, `reply`, `application_command`, `unknown-[id]`

Sample timestamp: `2017-06-24T13:06:38.555000+00:00`

### `attachment`

    time,fetchtype,action,type,id,messageid,filename

 - `id` (required)
 - `messageid` (required)
 - `filename` (optional)

### `reaction`

    action,type,userid,messageid,emoji,count

 - `userid` - can be empty if there are more than 100 reactions with a given emoji
 - `messageid` (required)
 - `emoji` (required) - character or `<emojiname>:<emojiid>`
 - `count` (required) - number of unlisted users if there's no user ID present or `1` otherwise

### `embed`

    action,type,messageid,json

 - `messageid` (required)
 - `json` (required) - JSON-encoded [embed contents](https://discordapp.com/developers/docs/resources/channel#embed-object)

### `pin`

    action,type,messageid

 - `messageid` (required)

## Server entry types

### `guild`

    time,fetchtype,action,type,id,name,icon,splash,ownerid,afkchanid,afktimeout,embeddable,embedchanid

 - `name` (required)
 - `ownerid` (required)
 - `embeddable` (boolean)

### `member`

    time,fetchtype,action,type,userid,username,discriminator,avatar,nick,roles

 - `username` (required) - global name
 - `nick` - server nickname
 - `roles` - comma-separated role IDs

### `ban`

    action,type,userid,reason

 - `userid` (required)
 - `reason` - ban reason, can't be seen by regular members by default

### `role`

    time,fetchtype,action,type,id,name,color,pos,perms,hoist

 - `name` (required)
 - `color` (required)
 - `hoist` (boolean) - if this role is pinned in the user listing
 - `pos` (required)
 - `perms` (required)

### `channel`

    time,fetchtype,action,type,id,chantype,pos,name,topic,nsfw,category,recipients,icon

 - `chantype` (required) - `text`, `voice`, `category`, `dm`, `groupdm`, `news` or `store`
 - `pos` (required)
 - `name` (required if not `dm` or `groupdm`)
 - `nsfw` (boolean)
 - `category` - ID of the parent category for a channel
 - `recipients` - comma-separated list of user IDs, only relevant to DM channels
 - `icon` - only relevant to DM channels

### `permoverwrite`

    time,fetchtype,action,type,id,overwritetype,allow,deny

 - `overwritetype` (required)
 - `allow` (required)
 - `deny` (required)

### `emoji`

    time,fetchtype,action,type,id,name,nocolons

 - `name` (required)
 - `nocolons` (boolean)
