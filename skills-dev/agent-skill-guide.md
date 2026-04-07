# Agent Skill Discovery & Installation Guide

> How an OpenClaw agent discovers AWP worknets, browses their skills, and installs them.

---

## 1. Overview

Each AWP worknet can publish a **skillsURI** — an HTTPS URL pointing to a `SKILL.md` file. This file follows the [OpenClaw Skills format](https://docs.openclaw.ai/tools/skills) and teaches the agent how to interact with the worknet's services (task submission, result retrieval, reward claiming, etc.).

**Flow:**
```
Agent                          AWP API                      Worknet
  │                               │                           │
  ├─ GET /api/worknets?status=Active ─►│                       │
  │◄─ [{worknet_id, skills_uri, ...}]──┤                       │
  │                               │                           │
  ├─ GET /api/worknets/{id}/skills ──►│                        │
  │◄─ {"skillsURI": "https://..."}──┤                        │
  │                               │                           │
  ├─ Fetch SKILL.md from skillsURI ────────────────────────►│
  │◄─ SKILL.md content ───────────────────────────────────────┤
  │                               │                           │
  ├─ Install skill locally        │                           │
  └─ Use skill to interact ───────────────────────────────►│
```

---

## 2. Discover Worknets

### 2.1 List Active Worknets

```javascript
// REST API — list all active worknets with their skills URIs
const res = await fetch(`${API_BASE}/worknets?status=Active&page=1&limit=50`);
const worknets = await res.json();

for (const worknet of worknets) {
  console.log(`[${worknet.worknet_id}] ${worknet.name}`);
  console.log(`  Status: ${worknet.status}`);
  console.log(`  Skills: ${worknet.skills_uri || 'None'}`);
}
```

### 2.2 Get a Single Worknet's Skills URI

```javascript
// REST API — get skills URI for a specific worknet
const res = await fetch(`${API_BASE}/worknets/1/skills`);
const { skillsURI } = await res.json();
// skillsURI = "https://worknet.example.com/SKILL.md"
```

### 2.3 On-Chain Query (optional)

Skills URI is stored on-chain in AWPWorkNet and emitted via the `WorknetRegistered` event. It can be updated by the NFT owner via `awpWorkNet.setSkillsURI()`, which emits `SkillsURIUpdated`.

---

## 3. Fetch & Inspect SKILL.md

The `skillsURI` is an HTTPS URL pointing to a `SKILL.md` file. Fetch it directly:

```bash
curl -s https://worknet.example.com/SKILL.md
```

### 3.1 SKILL.md Format

Each `SKILL.md` follows the OpenClaw skill format — YAML frontmatter + markdown instructions:

```markdown
---
name: my-worknet-tasks
description: Submit and retrieve AI tasks from My Worknet on AWP
user-invocable: true
metadata: {"openclaw":{"requires":{"env":["AWP_WORKNET_API_KEY"]}}}
---

## Instructions

You are connected to My Worknet (AWP Worknet #1).

### Available Actions
- **Submit task**: POST to the coordinator endpoint with task parameters
- **Check result**: GET task status and results by task ID
- **Claim reward**: ...

### Coordinator API
Base URL: `https://coord.myworknet.io`

#### Submit Task
POST /api/tasks
...
```

**Key frontmatter fields:**

| Field | Required | Description |
|-------|----------|-------------|
| `name` | Yes | Unique skill identifier (slug) |
| `description` | Yes | One-line description (shown in skill list) |
| `user-invocable` | No | `true` = exposed as `/my-worknet-tasks` slash command |
| `metadata` | No | JSON object with requirements (env vars, binaries, platform) |

---

## 4. Install Skills

### 4.1 Via ClawHub (if published)

If the worknet developer has published their skill to [ClawHub](https://github.com/openclaw/clawhub), agents can install with one command:

```bash
clawhub install my-worknet-tasks
```

Skills are installed to `./skills/` in the current workspace (or `~/.openclaw/skills/` globally).

### 4.2 Manual Install from URL

Download the `SKILL.md` directly from the worknet's `skillsURI` and place it in the skills directory:

```bash
# Create skill directory
mkdir -p ~/.openclaw/skills/my-worknet-tasks

# Download SKILL.md from the worknet's skillsURI
curl -o ~/.openclaw/skills/my-worknet-tasks/SKILL.md \
  "https://worknet.example.com/SKILL.md"
```

OpenClaw will automatically detect and load the skill on the next session start.

### 4.3 Workspace Install (per-project)

For project-specific skills, place them in the workspace `skills/` directory:

```bash
mkdir -p ./skills/my-worknet-tasks
curl -o ./skills/my-worknet-tasks/SKILL.md \
  "https://worknet.example.com/SKILL.md"
```

### 4.4 Automated Discovery Script

A script that discovers all active worknets and installs their skills:

```bash
#!/usr/bin/env bash
# Install all AWP worknet skills into the local workspace

API_BASE="${API_BASE_URL:-https://<api-host>/api}"
SKILLS_DIR="./skills"

# Fetch active worknets
worknets=$(curl -s "$API_BASE/worknets?status=Active&limit=100")

echo "$worknets" | jq -c '.[]' | while read -r worknet; do
  id=$(echo "$worknet" | jq -r '.worknet_id')
  name=$(echo "$worknet" | jq -r '.name')
  skills_uri=$(echo "$worknet" | jq -r '.skills_uri // empty')

  if [[ -z "$skills_uri" ]]; then
    echo "[$id] $name — no skillsURI, skipping"
    continue
  fi

  # Derive skill slug from worknet name (lowercase, hyphens)
  slug=$(echo "awp-worknet-$id" | tr '[:upper:]' '[:lower:]')
  skill_dir="$SKILLS_DIR/$slug"

  mkdir -p "$skill_dir"
  if curl -sf -o "$skill_dir/SKILL.md" "$skills_uri"; then
    echo "[$id] $name — installed to $skill_dir"
  else
    echo "[$id] $name — failed to fetch $skills_uri"
    rm -rf "$skill_dir"
  fi
done

echo "Done. Restart OpenClaw to load new skills."
```

---

## 5. Skill Loading & Precedence

OpenClaw loads skills from three locations (highest precedence first):

| Priority | Location | Scope |
|----------|----------|-------|
| 1 (highest) | `<workspace>/skills/` | Current project only |
| 2 | `~/.openclaw/skills/` | All projects on this machine |
| 3 (lowest) | Bundled skills | Shipped with OpenClaw |

Additional skill directories can be configured via `skills.load.extraDirs` in `~/.openclaw/openclaw.json`.

Skills are **snapshotted at session start**. If you install a new skill, restart the OpenClaw session to load it.

---

## 6. Using Installed Skills

Once installed, skills are available in two ways:

### 6.1 Slash Command (user-invocable)

If `user-invocable: true` in the frontmatter:
```
/my-worknet-tasks submit a translation task for "Hello World" to Chinese
```

### 6.2 Model-Triggered

OpenClaw injects eligible skill descriptions into the system prompt. The agent can autonomously decide to use a skill when it matches the user's intent:

```
User: "I need to run an image classification task on worknet 3"
Agent: [detects installed awp-worknet-3 skill, follows its instructions to POST to the coordinator]
```

---

## 7. For Worknet Developers: Publishing Skills

### 7.1 Create SKILL.md

Your `SKILL.md` should document how agents interact with your worknet's coordinator API:

```markdown
---
name: awp-worknet-myproject
description: Interact with MyProject worknet — submit AI tasks and collect results
user-invocable: true
---

## You are connected to MyProject Worknet (AWP #42)

### Coordinator API: https://coord.myproject.io

### Submit Task
POST /api/v1/tasks
Content-Type: application/json

{
  "type": "image-classification",
  "input": "https://example.com/image.jpg",
  "callback": "https://..."
}

Response: {"task_id": "abc123", "status": "pending"}

### Check Status
GET /api/v1/tasks/{task_id}

### Get Result
GET /api/v1/tasks/{task_id}/result
```

### 7.2 Host SKILL.md

Host the file at a stable HTTPS URL. Common options:
- GitHub raw URL: `https://raw.githubusercontent.com/org/repo/main/SKILL.md`
- Your own CDN/server: `https://docs.myworknet.io/SKILL.md`

### 7.3 Set skillsURI on-chain

When registering or updating your worknet, include the URL:

```javascript
// At registration
await awpRegistry.registerWorknet({
  name: "MyProject",
  symbol: "MYP",
  worknetManager: "0x0000000000000000000000000000000000000000",
  salt: "0x0000000000000000000000000000000000000000000000000000000000000000",
  minStake: 0n,
  skillsURI: "https://docs.myproject.io/SKILL.md",
});

// Update skillsURI later (NFT owner only, via AWPWorkNet)
await awpWorkNet.setSkillsURI(worknetId, "https://docs.myproject.io/SKILL.md");
```

### 7.4 Publish to ClawHub (optional)

For broader distribution beyond AWP:

```bash
cd my-skill-directory/  # contains SKILL.md
clawhub publish
```

---

## 8. Example: End-to-End

```bash
# 1. Discover worknets with skills
curl -s "${API_BASE}/worknets?status=Active" | \
  jq '.[] | select(.skills_uri != null) | {id: .worknet_id, name: .name, skills: .skills_uri}'

# Output:
# {"id": 1, "name": "ImageNet Classifier", "skills": "https://imagenet-worknet.io/SKILL.md"}
# {"id": 3, "name": "Translation Engine", "skills": "https://translate-worknet.io/SKILL.md"}

# 2. Install a skill
mkdir -p ~/.openclaw/skills/awp-imagenet
curl -o ~/.openclaw/skills/awp-imagenet/SKILL.md https://imagenet-worknet.io/SKILL.md

# 3. Restart OpenClaw and use it
# In OpenClaw:
#   /awp-imagenet classify this image: https://example.com/cat.jpg
```
