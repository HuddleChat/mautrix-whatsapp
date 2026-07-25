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
