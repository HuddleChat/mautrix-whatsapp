# mautrix-whatsapp (Huddle fork)

> **This is a modified version of [mautrix-whatsapp](https://github.com/mautrix/whatsapp).**
>
> Modified by HuddleChat. The `huddle` branch is based on upstream tag `v0.2607.0` and
> carries the changes listed below; `main` tracks upstream unmodified. This fork is
> distributed under the AGPL-3.0, the same licence as upstream.
>
> ### Changes
>
> * **2026-07-25** — Added the `private_chat_name_template` config option, which names DM
>   portal rooms from the receiving login's own contact store. Upstream names DM rooms after
>   the ghost puppet, whose displayname is global across the bridge, so putting address-book
>   fields such as `{{.FullName}}` in `displayname_template` leaks one user's private label
>   for a contact to every other user who has that same contact. The new option renders a
>   per-login name instead, which is private as long as `split_portals` is enabled.
>   Touches `pkg/connector/config.go`, `pkg/connector/chatinfo.go` and
>   `pkg/connector/example-config.yaml`; adds `pkg/connector/privatechatname_test.go`.
> * **2026-07-25** — Re-render names for existing chats when the naming scheme changes.
>   Upstream only resyncs a DM portal when the user opens the chat, so a change to the
>   name templates never reaches rooms that already exist. `resyncContacts` now also
>   refreshes DM portal names, and a `name_scheme_version` field on the user login triggers
>   a one-off resync after connecting when it trails the current scheme. Touches
>   `pkg/connector/userinfo.go`, `pkg/connector/handlewhatsapp.go` and `pkg/waid/dbmeta.go`.
> * **2026-07-25** — Sweep every DM portal, not only those matching a saved contact. The
>   contact-driven sweep missed DMs with people the user never saved, which were exactly the
>   rooms still holding a name derived from another user's address book. Touches
>   `pkg/connector/userinfo.go`.
> * **2026-07-25** — Name the self-chat and LID-only chats as well, and re-render a DM's name
>   when its contact/push name changes. A DM without an explicit name keeps following the
>   global ghost, and `updateDMPortals` reaches DM portals across *all* logins, so any other
>   user's contact sync could rename such a room — including retitling someone's own notes
>   room with a stranger's label for them. Touches `pkg/connector/chatinfo.go`,
>   `pkg/connector/handlewhatsapp.go` and `pkg/connector/userinfo.go`.
> * **2026-07-29** — Fetch each DM's profile picture with the login that owns the room and
>   put it on that login's own portal. WhatsApp decides picture visibility per viewer, but
>   ghost puppets are global, so upstream stores whichever answer the first login to ask was
>   given — and files a refusal as settled, never retried. One user being refused hid a
>   contact from everyone; one being allowed published a contact to people who had hidden it
>   from them. A `dm_avatar_scheme_version` on the user login sweeps existing portals, since
>   `EnqueuePortalResync` skips DMs entirely. Touches `pkg/connector/chatinfo.go`,
>   `pkg/connector/handlewhatsapp.go`, `pkg/connector/userinfo.go` and `pkg/waid/dbmeta.go`.
> * **2026-07-29** — Label backfilled messages so the homeserver can decline to notify for
>   them. Only hungryserv can batch-send; on Synapse mautrix replays history one ordinary
>   timeline event at a time, and push rules cannot tell an imported two-year-old message
>   from a friend messaging right now — so linking an account fires a notification per
>   imported message. Every message out of the history-sync store now carries a
>   `huddle_backfill` content field, which a client-installed override push rule matches with
>   no actions. Live messages never take this path, so nothing the user is waiting for can be
>   silenced by accident. Touches `pkg/connector/backfill.go`.

A Matrix-WhatsApp puppeting bridge based on [whatsmeow](https://github.com/tulir/whatsmeow).

## Documentation
All setup and usage instructions are located on [docs.mau.fi]. Some quick links:

[docs.mau.fi]: https://docs.mau.fi/bridges/go/whatsapp/index.html

* [Bridge setup](https://docs.mau.fi/bridges/go/setup.html?bridge=whatsapp)
  (or [with Docker](https://docs.mau.fi/bridges/general/docker-setup.html?bridge=whatsapp))
* Basic usage: [Authentication](https://docs.mau.fi/bridges/go/whatsapp/authentication.html)

### Features & Roadmap
[ROADMAP.md](ROADMAP.md) contains a general overview of what is supported by the bridge.

## Discussion
Matrix room: [#whatsapp:maunium.net](https://matrix.to/#/#whatsapp:maunium.net)
