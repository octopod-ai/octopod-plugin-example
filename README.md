# Octopod plugin example

A minimal plugin for [Octopod](https://github.com/octopod-ai/octopod): the local time anywhere, and greetings in the style of a workspace. Copy it to write your own.

```
octopod-plugin.yaml   the manifest: tools, server, parameters, instructions
INSTRUCTIONS.md       added to the prompt of the personas using the plugin
skills/*/SKILL.md     know-how given to those personas
server/               the MCP server (any language), built from its Dockerfile
```

- Every tool is declared in the manifest, with its input schema and its approval level (`auto`, `agent` for the validator, `human` for the owner): what the plugin may do is reviewed when it is installed, not discovered at run time.
- The server answers MCP over HTTP (streamable) on the port and path of the manifest. Each call carries `_meta.octopod`: the workspace, the persona, and the parameters the workspace gave the plugin (`greeting` here).
- The server runs in a hardened container on a network only Octopod reaches; it goes out to the Internet (public addresses only) through Octopod's proxy. Secrets it declares (`secrets: [name]`) come from the owner, never from an agent.

Install it from the chat: « installe le plugin github.com/octopod-ai/octopod-plugin-example », then give it to a persona of a workspace, with its parameters.
