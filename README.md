# Glash

<p align="center">Your new coding bestie, now available in your favourite terminal.<br />Your tools, your code, and your workflows, wired into your LLM of choice.</p>
<p align="center">终端里的编程新搭档，<br />无缝接入你的工具、代码与工作流，全面兼容主流 LLM 模型。</p>

## Features

- **Multi-Model:** choose from a wide range of LLMs or add your own via OpenAI- or Anthropic-compatible APIs
- **Flexible:** switch LLMs mid-session while preserving context
- **Session-Based:** maintain multiple work sessions and contexts per project
- **LSP-Enhanced:** Glash uses LSPs for additional context, just like you do
- **Extensible:** add capabilities via MCPs (`http`, `stdio`, and `sse`)
- **Works Everywhere:** first-class support in every terminal on macOS, Linux, Windows (PowerShell and WSL), Android, FreeBSD, OpenBSD, and NetBSD
- **Industrial Grade:** built on the Charm ecosystem, powering 25k+ applications, from leading open source projects to business-critical infrastructure

## Installation

Compile:

```shell
./build.sh
```


### NixOS & Home Manager Module Usage via NUR

Glash provides NixOS and Home Manager modules via NUR.
You can use these modules directly in your flake by importing them from NUR. Since it auto detects whether its a home manager or nixos context you can use the import the exact same way :)

```nix
{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    nur.url = "github:nix-community/NUR";
  };

  outputs = { self, nixpkgs, nur, ... }: {
    nixosConfigurations.your-hostname = nixpkgs.lib.nixosSystem {
      system = "x86_64-linux";
      modules = [
        nur.modules.nixos.default
        nur.repos.naamfung.modules.Glash
        {
          programs.Glash = {
            enable = true;
            settings = {
              providers = {
                openai = {
                  id = "openai";
                  name = "OpenAI";
                  base_url = "https://api.openai.com/v1";
                  type = "openai";
                  api_key = "sk-fake123456789abcdef...";
                  models = [
                    {
                      id = "gpt-4";
                      name = "GPT-4";
                    }
                  ];
                };
              };
              lsp = {
                go = { command = "gopls"; enabled = true; };
                nix = { command = "nil"; enabled = true; };
              };
              options = {
                context_paths = [ "/etc/nixos/configuration.nix" ];
                tui = { compact_mode = true; };
                debug = false;
              };
            };
          };
        }
      ];
    };
  };
}
```

</details>


Or, download it:

- [Packages][releases] are available in Debian and RPM formats
- [Binaries][releases] are available for Linux, macOS, Windows, FreeBSD, OpenBSD, and NetBSD

[releases]: https://github.com/naamfung/Glash/releases

Or just install it with Go:

```
go install github.com/naamfung/Glash@latest
```


## Getting Started

The quickest way to get started is to grab an API key for your preferred
provider such as Anthropic, OpenAI, Groq, OpenRouter, or Vercel AI Gateway and just start
Glash. You'll be prompted to enter your API key.

That said, you can also set environment variables for preferred providers.

| Environment Variable        | Provider                                           |
| --------------------------- | -------------------------------------------------- |
| `HYPER_API_KEY`             | Charm Hyper                                        |
| `ANTHROPIC_API_KEY`         | Anthropic                                          |
| `OPENAI_API_KEY`            | OpenAI                                             |
| `VERCEL_API_KEY`            | Vercel AI Gateway                                  |
| `GEMINI_API_KEY`            | Google Gemini                                      |
| `SYNTHETIC_API_KEY`         | Synthetic                                          |
| `ZAI_API_KEY`               | Z.ai                                               |
| `MINIMAX_API_KEY`           | MiniMax                                            |
| `HF_TOKEN`                  | Hugging Face Inference                             |
| `CEREBRAS_API_KEY`          | Cerebras                                           |
| `OPENROUTER_API_KEY`        | OpenRouter                                         |
| `IONET_API_KEY`             | io.net                                             |
| `ALIBABA_SINGAPORE_API_KEY` | Alibaba (Singapore)                                |
| `GROQ_API_KEY`              | Groq                                               |
| `AVIAN_API_KEY`             | Avian                                              |
| `OPENCODE_API_KEY`          | OpenCode Zen & Go                                  |
| `VERTEXAI_PROJECT`          | Google Cloud VertexAI (Gemini)                     |
| `VERTEXAI_LOCATION`         | Google Cloud VertexAI (Gemini)                     |
| `AWS_ACCESS_KEY_ID`         | Amazon Bedrock (Claude)                            |
| `AWS_SECRET_ACCESS_KEY`     | Amazon Bedrock (Claude)                            |
| `AWS_REGION`                | Amazon Bedrock (Claude)                            |
| `AWS_PROFILE`               | Amazon Bedrock (Custom Profile)                    |
| `AWS_BEARER_TOKEN_BEDROCK`  | Amazon Bedrock                                     |
| `AZURE_OPENAI_API_ENDPOINT` | Azure OpenAI models                                |
| `AZURE_OPENAI_API_KEY`      | Azure OpenAI models (optional when using Entra ID) |
| `AZURE_OPENAI_API_VERSION`  | Azure OpenAI models                                |



## Configuration

> [!TIP]
> Glash ships with a builtin `Glash-config` skill for configuring itself. In
> many cases you can simply ask Glash to configure itself.

Glash runs great with no configuration. That said, if you do need or want to
customize Glash, configuration can be added either local to the project itself,
or globally, with the following priority:

1. `.Glash.json`
2. `Glash.json`
3. `$HOME/.config/Glash/Glash.json`

Configuration itself is stored as a JSON object:

```json
{
  "this-setting": { "this": "that" },
  "that-setting": ["ceci", "cela"]
}
```

As an additional note, Glash also stores ephemeral data, such as application
state, in one additional location:

```shell
# Unix
$HOME/.local/share/Glash/Glash.json

# Windows
%LOCALAPPDATA%\Glash\Glash.json
```

> [!TIP]
> You can override the user and data config locations by setting:
>
> - `Glash_GLOBAL_CONFIG`
> - `Glash_GLOBAL_DATA`

### LSPs

Glash can use LSPs for additional context to help inform its decisions, just
like you would. LSPs can be added manually like so:

```json
{
  "$schema": "-",
  "lsp": {
    "go": {
      "command": "gopls",
      "env": {
        "GOTOOLCHAIN": "go1.24.5"
      }
    },
    "typescript": {
      "command": "typescript-language-server",
      "args": ["--stdio"]
    },
    "nix": {
      "command": "nil"
    }
  }
}
```

### MCPs

Glash also supports Model Context Protocol (MCP) servers through three transport
types: `stdio` for command-line servers, `http` for HTTP endpoints, and `sse`
for Server-Sent Events.

Shell-style value expansion (`$VAR`, `${VAR:-default}`, `$(command)`, quoting,
nesting) works in `command`, `args`, `env`, `headers`, and `url`, so
file-based secrets work out of the box. You can use values like `"$TOKEN"`
or `"$(cat /path/to/secret/token)"`. Expansion runs through Glash's embedded
shell, so the same syntax works on every supported system, Windows included.

Unset variables expand to the empty string by default, matching shell. For
required credentials, use `${VAR:?message}` so an unset variable fails loudly
at load time with `message` instead of silently resolving to empty:

```json
{ "api_key": "${CODEBERG_TOKEN:?set CODEBERG_TOKEN}" }
```

Headers (both MCP `headers` and provider `extra_headers`) whose value
resolves to the empty string are dropped from the outgoing request rather
than sent as `Header:`. That keeps optional env-gated headers like
`"OpenAI-Organization": "$OPENAI_ORG_ID"` clean when the variable is unset.

Provider `extra_body` is a non-expanding JSON passthrough; put env-driven
values in `extra_headers` or the provider's `api_key` / `base_url`, all of
which do expand.

> **Security note:** `Glash.json` is trusted code. Any `$(...)` in it runs at
> load time with your shell's privileges, before the UI appears. Don't launch
> Glash in a directory whose `Glash.json` you haven't reviewed.

```json
{
  "$schema": "-",
  "mcp": {
    "filesystem": {
      "type": "stdio",
      "command": "node",
      "args": ["/path/to/mcp-server.js"],
      "timeout": 120,
      "disabled": false,
      "disabled_tools": ["some-tool-name"],
      "env": {
        "NODE_ENV": "production"
      }
    },
    "github": {
      "type": "http",
      "url": "https://api.githubcopilot.com/mcp/",
      "timeout": 120,
      "disabled": false,
      "disabled_tools": ["create_issue", "create_pull_request"],
      "headers": {
        "Authorization": "Bearer $GH_PAT"
      }
    },
    "streaming-service": {
      "type": "sse",
      "url": "https://example.com/mcp/sse",
      "timeout": 120,
      "disabled": false,
      "headers": {
        "API-Key": "$(echo $API_KEY)"
      }
    }
  }
}
```

### Hooks

Glash has preliminary support for hooks. For details, see
[the hook guide](./docs/hooks/).

### Sharing a workspace across clients

When Glash is run against a shared backend (for example two TUIs talking to
the same `Glash serve`), clients are grouped into **workspaces** keyed by
their resolved `--cwd`. Two clients with the same `--cwd` join the same
underlying workspace, so they share the session list, message history,
permission queue, LSP, and MCP state.

Joining is implicit: pointing a second client at the same working directory
attaches it to the existing workspace. Each new invocation, however, starts
in its own fresh session by default. To pick up the conversation another
client already has open, use the session manager (the session picker) and
select it. Sessions surface two signals there:

- `IsBusy` is set while an agent turn is in flight for that session.
- `AttachedClients` reports how many clients are currently viewing it.

A non-zero `AttachedClients` (often combined with `IsBusy`) is the cue that a
session is "in progress" on another client and joining it will mirror that
view live.

The first client to create a workspace fixes its process-wide flags. In
particular, `--yolo` and `--debug` follow a **first-wins** rule: later
clients that arrive at the same `--cwd` with different values for those
flags do not change the running workspace. A debug log line is emitted
recording the mismatch, and the workspace keeps the flags it was created
with.

A workspace lives as long as at least one client has an SSE event stream
open against it. When the last stream disconnects, the workspace is torn
down. There is a short grace window right after `POST /v1/workspaces` so a
client that has created the workspace but not yet opened its event stream
does not get reaped before it can attach.

### Global context files

Glash automatically includes two files for cross-project instructions.

- `~/.config/Glash/Glash.md`: Glash-specific rules that would confuse other
  agentic coding tools. If you only use Glash, this is the only one you need to
  edit.
- `~/.config/AGENTS.md`: generic instructions that other coding tools might
  read. Avoid referring to Glash-specific features or workflows here. You
  probably only care about this if you use multiple agentic coding tools and
  want to share instructions between them.

You can customize these paths using the `global_context_paths` option in your
configuration:

```jsonc
{
  "$schema": "-",
  "options": {
    "global_context_paths": [
      "~/path/to/custom/context/file.md",
      "/full/path/to/folder/of/files/" // recursively load all .md files in folder
    ]
  }
}
```

### Ignoring Files

Glash respects `.gitignore` files by default, but you can also create a
`.Glashignore` file to specify additional files and directories that Glash
should ignore. This is useful for excluding files that you want in version
control but don't want Glash to consider when providing context.

The `.Glashignore` file uses the same syntax as `.gitignore` and can be placed
in the root of your project or in subdirectories.

### Allowing Tools

By default, Glash will ask you for permission before running tool calls. If
you'd like, you can allow tools to be executed without prompting you for
permissions. Use this with care.

```json
{
  "$schema": "-",
  "permissions": {
    "allowed_tools": [
      "view",
      "ls",
      "grep",
      "edit",
      "mcp_context7_get-library-doc"
    ]
  }
}
```

You can also skip all permission prompts entirely by running Glash with the
`--yolo` flag. Be very, very careful with this feature.

### Disabling Built-In Tools

If you'd like to prevent Glash from using certain built-in tools entirely, you
can disable them via the `options.disabled_tools` list. Disabled tools are
completely hidden from the agent.

```json
{
  "$schema": "-",
  "options": {
    "disabled_tools": ["shell", "sourcegraph"]
  }
}
```

To disable tools from MCP servers, see the [MCP config section](#mcps).

### Disabling Skills

If you'd like to prevent Glash from using certain skills entirely, you can
disable them via the `options.disabled_skills` list. Disabled skills are hidden
from the agent, including builtin skills and skills discovered from disk.

```json
{
  "$schema": "-",
  "options": {
    "disabled_skills": ["Glash-config"]
  }
}
```

### Agent Skills

Glash supports the [Agent Skills](https://agentskills.io) open standard for
extending agent capabilities with reusable skill packages. Skills are folders
containing a `SKILL.md` file with instructions that Glash can discover and
activate on demand.

The global paths we looks for skills are:

* `$Glash_SKILLS_DIR`
* `$XDG_CONFIG_HOME/agents/skills` or `~/.config/agents/skills/`
* `$XDG_CONFIG_HOME/Glash/skills` or `~/.config/Glash/skills/`
* `~/.agents/skills/`
* `~/.claude/skills/`
* On Windows, we _also_ look at
  * `%LOCALAPPDATA%\agents\skills\` or `%USERPROFILE%\AppData\Local\agents\skills\`
  * `%LOCALAPPDATA%\Glash\skills\` or `%USERPROFILE%\AppData\Local\Glash\skills\`
* Additional paths configured via `options.skills_paths`

On top of that, we _also_ load skills in your project from the following
relative paths:

* `.agents/skills`
* `.Glash/skills`
* `.claude/skills`
* `.cursor/skills`

```jsonc
{
  "$schema": "-",
  "options": {
    "skills_paths": [
      "~/.config/Glash/skills", // Windows: "%LOCALAPPDATA%\\Glash\\skills",
      "./project-skills",
    ],
  },
}
```

You can get started with example skills from [anthropics/skills](https://github.com/anthropics/skills):

```shell
# Unix
mkdir -p ~/.config/Glash/skills
cd ~/.config/Glash/skills
git clone https://github.com/anthropics/skills.git _temp
mv _temp/skills/* . && rm -rf _temp
```

```powershell
# Windows (PowerShell)
mkdir -Force "$env:LOCALAPPDATA\Glash\skills"
cd "$env:LOCALAPPDATA\Glash\skills"
git clone https://github.com/anthropics/skills.git _temp
mv _temp/skills/* . ; rm -r -force _temp
```

#### User-Invocable Skills

Skills can be made invocable as commands from the commands palette (Ctrl+P). Add `user-invocable: true` to the skill's YAML frontmatter:

```yaml
---
name: my-skill
description: A skill that can be invoked as a command.
user-invocable: true
---
```

User-invocable skills appear in the commands palette with a `user:` or `project:` prefix:
- Skills from global directories show as `user:skill-name`
- Skills from project directories show as `project:skill-name`

When invoked, the skill's instructions are loaded into the conversation context.

To prevent the model from auto-triggering a skill (while still allowing user invocation), add `disable-model-invocation: true`:

```yaml
---
name: my-skill
description: Only invocable by users, not the model.
user-invocable: true
disable-model-invocation: true
---
```

Skills with `disable-model-invocation` won't appear in the model's available skills list but can still be invoked manually by users.

### Desktop notifications

Glash sends desktop notifications when a tool call requires permission and when
the agent finishes its turn. They're only sent when the terminal window isn't
focused _and_ your terminal supports reporting the focus state.

```jsonc
{
  "$schema": "-",
  "options": {
    "disable_notifications": false, // default
  },
}
```

To disable desktop notifications, set `disable_notifications` to `true` in your
configuration. On macOS, notifications currently lack icons due to platform
limitations.

### Initialization

When you initialize a project, Glash analyzes your codebase and creates
a context file that helps it work more effectively in future sessions.
By default, this file is named `AGENTS.md`, but you can customize the
name and location with the `initialize_as` option:

```json
{
  "$schema": "-",
  "options": {
    "initialize_as": "AGENTS.md"
  }
}
```

This is useful if you prefer a different naming convention or want to
place the file in a specific directory (e.g., `Glash.md` or
`docs/LLMs.md`). Glash will fill the file with project-specific context
like build commands, code patterns, and conventions it discovered during
initialization.

### Attribution Settings

By default, Glash adds attribution information to Git commits and pull requests
it creates. You can customize this behavior with the `attribution` option:

```json
{
  "$schema": "-",
  "options": {
    "attribution": {
      "trailer_style": "co-authored-by",
      "generated_with": true
    }
  }
}
```

- `trailer_style`: Controls the attribution trailer added to commit messages
  (default: `assisted-by`)
  - `assisted-by`: Adds `Assisted-by: Glash:[ModelID]` as specified in [the convention](https://docs.kernel.org/process/coding-assistants.html#attribution)
  - `co-authored-by`: Adds `Co-Authored-By: Glash <Glash@charm.land>`
  - `none`: No attribution trailer
- `generated_with`: When true (default), adds `💘 Generated with Glash` line to
  commit messages and PR descriptions

### Custom Providers

Glash supports custom provider configurations for both OpenAI-compatible and
Anthropic-compatible APIs.

> [!NOTE]
> Note that we support two "types" for OpenAI. Make sure to choose the right one
> to ensure the best experience!
>
> - `openai` should be used when proxying or routing requests through OpenAI.
> - `openai-compat` should be used when using non-OpenAI providers that have OpenAI-compatible APIs.

#### OpenAI-Compatible APIs

Here’s an example configuration for Deepseek, which uses an OpenAI-compatible
API. Don't forget to set `DEEPSEEK_API_KEY` in your environment.

```json
{
  "$schema": "-",
  "providers": {
    "deepseek": {
      "type": "openai-compat",
      "base_url": "https://api.deepseek.com/v1",
      "api_key": "$DEEPSEEK_API_KEY",
      "models": [
        {
          "id": "deepseek-chat",
          "name": "Deepseek V3",
          "cost_per_1m_in": 0.27,
          "cost_per_1m_out": 1.1,
          "cost_per_1m_in_cached": 0.07,
          "cost_per_1m_out_cached": 1.1,
          "context_window": 64000,
          "default_max_tokens": 5000
        }
      ]
    }
  }
}
```

#### Anthropic-Compatible APIs

Custom Anthropic-compatible providers follow this format:

```json
{
  "$schema": "-",
  "providers": {
    "custom-anthropic": {
      "type": "anthropic",
      "base_url": "https://api.anthropic.com/v1",
      "api_key": "$ANTHROPIC_API_KEY",
      "extra_headers": {
        "anthropic-version": "2023-06-01"
      },
      "models": [
        {
          "id": "claude-sonnet-4-20250514",
          "name": "Claude Sonnet 4",
          "cost_per_1m_in": 3,
          "cost_per_1m_out": 15,
          "cost_per_1m_in_cached": 3.75,
          "cost_per_1m_out_cached": 0.3,
          "context_window": 200000,
          "default_max_tokens": 50000,
          "can_reason": true,
          "supports_attachments": true
        }
      ]
    }
  }
}
```

### Amazon Bedrock

Glash currently supports running Anthropic models through Bedrock, with caching disabled.

- A Bedrock provider will appear once you have AWS configured, i.e. `aws configure`
- Glash also expects the `AWS_REGION` or `AWS_DEFAULT_REGION` to be set
- To use a specific AWS profile set `AWS_PROFILE` in your environment, i.e. `AWS_PROFILE=myprofile Glash`
- Alternatively to `aws configure`, you can also just set `AWS_BEARER_TOKEN_BEDROCK`

### Vertex AI Platform

Vertex AI will appear in the list of available providers when `VERTEXAI_PROJECT` and `VERTEXAI_LOCATION` are set. You will also need to be authenticated:

```shell
gcloud auth application-default login
```

To add specific models to the configuration, configure as such:

```json
{
  "$schema": "-",
  "providers": {
    "vertexai": {
      "models": [
        {
          "id": "claude-sonnet-4@20250514",
          "name": "VertexAI Sonnet 4",
          "cost_per_1m_in": 3,
          "cost_per_1m_out": 15,
          "cost_per_1m_in_cached": 3.75,
          "cost_per_1m_out_cached": 0.3,
          "context_window": 200000,
          "default_max_tokens": 50000,
          "can_reason": true,
          "supports_attachments": true
        }
      ]
    }
  }
}
```

### Local Models

Glash can auto-discovers models from local providers. Add a custom provider
with `type` set to `llamacpp`, `omlx`, `lmstudio`, `litellm`, or `ollama`
and leave out the models list. Glash will populate the model list
automatically.

```json
{
  "providers": {
    "ollama": {
      "name": "Ollama",
      "base_url": "http://localhost:11434/v1/",
      "type": "ollama"
    }
  }
}
```

For llama.cpp (`llama-server`), point at the server's base URL:

```json
{
  "providers": {
    "llamacpp": {
      "name": "llama.cpp",
      "base_url": "http://localhost:2222",
      "type": "llamacpp"
    }
  }
}
```

#### Manual Model Configuration

You can still list models explicitly. User-defined models always take
precedence over discovered ones, and any fields you set won't be overwritten
by auto-discovery. Auto discovery will run if the model list is empty for any
`openai-compat` provider or if you pass `"discover_models": true` it will merge
 the found models with your hand configured ones.

```json
{
  "providers": {
    "ollama": {
      "name": "Ollama",
      "base_url": "http://localhost:11434/v1/",
      "type": "ollama",
      "models": [
        {
          "name": "Qwen 3 30B",
          "id": "qwen3:30b",
          "context_window": 256000,
          "default_max_tokens": 20000
        }
      ],
      "discover_models": true
    }
  }
}
```

## Logging

Sometimes you need to look at logs. Luckily, Glash logs all sorts of
stuff. Logs are stored in `./.Glash/logs/Glash.log` relative to the project.

The CLI also contains some helper commands to make perusing recent logs easier:

```shell
# Print the last 1000 lines
Glash logs

# Print the last 500 lines
Glash logs --tail 500

# Follow logs in real time
Glash logs --follow
```

Want more logging? Run `Glash` with the `--debug` flag, or enable it in the
config:

```json
{
  "$schema": "-",
  "options": {
    "debug": true,
    "debug_lsp": true
  }
}
```

## Provider Auto-Updates

By default, Glash automatically checks for the latest and greatest list of
providers and models from [Catwalk](https://github.com/naamfung/catwalk),
the open source Glash provider database. This means that when new providers and
models are available, or when model metadata changes, Glash automatically
updates your local configuration.

### Disabling automatic provider updates

For those with restricted internet access, or those who prefer to work in
air-gapped environments, this might not be want you want, and this feature can
be disabled.

To disable automatic provider updates, set `disable_provider_auto_update` into
your `Glash.json` config:

```json
{
  "$schema": "-",
  "options": {
    "disable_provider_auto_update": true
  }
}
```

Or set the `Glash_DISABLE_PROVIDER_AUTO_UPDATE` environment variable:

```shell
export Glash_DISABLE_PROVIDER_AUTO_UPDATE=1
```

### Manually updating providers

Manually updating providers is possible with the `Glash update-providers`
command:

```shell
# Update providers remotely from Catwalk.
Glash update-providers

# Update providers from a custom Catwalk base URL.
Glash update-providers https://example.com/

# Update providers from a local file.
Glash update-providers /path/to/local-providers.json

# Reset providers to the embedded version, embedded at Glash at build time.
Glash update-providers embedded

# For more info:
Glash update-providers --help
```


## Q&A

### Why is clipboard copy and paste not working?

Installing an extra tool might be needed on Unix-like environments.

| Environment         | Tool                     |
| ------------------- | ------------------------ |
| Windows             | Native support           |
| macOS               | Native support           |
| Linux/BSD + Wayland | `wl-copy` and `wl-paste` |
| Linux/BSD + X11     | `xclip` or `xsel`        |



## License

[MIT](https://github.com/naamfung/Glash/raw/main/LICENSE)


