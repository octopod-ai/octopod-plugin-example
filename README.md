# Octopod plugin example

A minimal plugin for [Octopod](https://github.com/octopod-ai/octopod): the local time anywhere, and greetings in the style of a workspace. Copy it to write your own.

```
octopod-plugin.yaml                  the manifest: tools, server, parameters, instructions, personas, cron jobs
INSTRUCTIONS.md                      added to the prompt of the personas using the plugin
skills/*/SKILL.md                    know-how given to those personas
personas/greeter.md                  a persona the plugin brings: its description, model, limits and prompt
personas/greeter/skills/*/SKILL.md   know-how given to that persona alone
server/                              the MCP server (any language), built from its Dockerfile
```

- Every tool is declared in the manifest, with its input schema and its access: `access: {resource, does, id_from}`, the type of resource of Octopod it acts on (`none` for none, as both tools here), `read` or `write`, and the argument holding the id of that resource. Octopod decides it before each call, for the person the work is for and everyone who reads it; no call waits for a person. What the plugin may do is reviewed when it is installed, not discovered at run time.
- Every persona says why it holds each of its tools and core tools (`reasons: {<tool>: why}`). Octopod refuses to install or update a plugin with a tool declaring no access, or a tool held for no reason (`UNDECLARED_ACCESS`), naming each.
- The server answers MCP over HTTP (streamable) on the port and path of the manifest. Each call carries `_meta.octopod`: the workspace, the persona, and the parameters the workspace gave the plugin (`greeting` here).
- The server runs in a hardened container on a network only Octopod reaches; it goes out to the Internet (public addresses only) through Octopod's proxy. Secrets it declares (`secrets: [name]`) come from the owner, never from an agent.
- The plugin brings a persona, `example-greeter`, and a cron job it runs every Monday in each workspace where it works. A persona says only how it thinks (description, model, effort, limits) and which tools it is given, each with its reason: those of its plugin, and tools of Octopod by the names of its catalog (`create_notification`, `list_items`…). Its rights are generated from that, never written by the plugin.

Install it from the chat: « installe le plugin github.com/octopod-ai/octopod-plugin-example », then give it to a persona of a workspace, with its parameters, or have its own persona work in a workspace: « fais travailler example-greeter dans le workspace demo, avec la salutation Bonjour »; its lead then delegates to it.
