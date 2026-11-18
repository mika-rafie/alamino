<h1 align="center"><code>git push alamino</code></h1>
<p align="center">
  <a href="https://github.com/mika-rafie/alamino/actions/workflows/release.yml"
    ><img
      alt="Release"
      src="https://img.shields.io/github/actions/workflow/status/mika-rafie/alamino/release.yml?style=flat-square&label=release"
  /></a>
  <a href="https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-blue?style=flat-square"
    ><img
      alt="Platform"
      src="https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-blue?style=flat-square"
  /></a>
</p>

<h3 align="center">干掉所有 slop，开出干净的 PR。</h3>

<p align="center"><a href="README.md">English</a> · <strong>简体中文</strong></p>

<p align="center">
  <img src="https://raw.githubusercontent.com/mika-rafie/alamino/main/demo.gif" alt="alamino demo" width="800" />
</p>

`alamino` 在你真实的远端前面放了一个本地 git 代理。
把分支推给 `alamino` 而不是 `origin`，它会拉起一个用完即弃的 worktree，跑一条 AI 驱动的校验流水线，**只有每一项检查都通过后**才把分支转发到配置的推送目标，并自动开出一个干净的 PR。

- **不阻塞** —— 流水线在隔离的 worktree 里跑，不打断你手头的工作。
- **不挑 agent** —— 支持 `claude`、`codex`、`grok`、`rovodev`、`opencode`、`pi`、`copilot`、`antigravity`，也可通过 `acpx` 使用 `cursor` / `devin` / `acp:<target>`，并支持有序 fallback。
- **agent 原生** —— `/alamino` 既能让编码 agent 完成一个任务再过网关，也能直接为已提交的工作过网关：它跑完流水线、让流水线应用安全的修复，剩下的升级给你。
- **人始终说了算** —— 自动修复，还是逐条审查 findings，你决定。
- **默认就是干净 PR** —— 推送、开 PR、盯 CI、自动修复失败，一气呵成。

完整文档：<https://mika-rafie.github.io/alamino/>

## 工作原理

```
        你的分支
            │  git push alamino
            ▼
   ┌───────────────────────────────────────────────┐
   │  用完即弃的 worktree —— 你的工作原地不动          │
   │  review → test → docs → lint → push → PR → CI │
   └───────────────────────────────────────────────┘
            │  每项检查变绿
            ▼
        干净的 PR，已替你开好
```

每一步要么自己通过，要么停下来给你一条 **finding** 让你处理。
安全、机械性的修复会自动应用；任何牵涉到你**意图**的，都会升级给你来 **approve（批准）**、**fix（修复）** 或 **skip（跳过）**。
在每项检查都变绿之前，没有任何东西会到达配置的推送目标。

## 安装

```sh
curl -fsSL https://raw.githubusercontent.com/mika-rafie/alamino/main/docs/install.sh | sh
```

Windows、Go install、Nix 以及从源码构建的说明，见[安装指南](https://mika-rafie.github.io/alamino/start-here/installation/)。

## 快速上手

```sh
$ alamino init
  ✓ Gate initialized

    repo  /Users/you/src/my-repo
    gate  alamino → /Users/you/.alamino/repos/abc123def456.git
  remote  git@github.com:you/my-repo.git
   skill  /alamino installed for agents at user level

  Push through the gate with:
  git push alamino <branch>

$ git checkout my-branch

# 在分支里干点活……

$ git push alamino
  * Pipeline started

  Run alamino to review.

$ alamino
# 打开当前运行的 TUI
```

如果是 GitHub fork 贡献，让 `origin` 指向父仓库，并用 `alamino init --fork-url <your-fork-url>` 初始化。

在 TUI 里你逐条处理 **finding**：**auto-fix** 类自动替你应用（或由你 approve 放行），**ask-user** 类需要你判断，由你 approve、fix 或 skip。
每项检查变绿后，网关会把你的分支转发到配置的推送目标并替你开好 PR，不用手动 `git push origin`，也不用手写 PR 正文。
想让编码 agent 无人值守地走完同一套流程？
用 `/alamino`（见下文）。

## 触发网关的三种方式

每一处改动都走同一条流水线。改动就绪时，挑一个最贴合你当下工作方式的入口：

- **`git push alamino`** —— 显式的 Git 路径。把已提交的分支推给网关 remote，而不是 `origin`。
- **`alamino`** —— TUI。改完之后运行它（无需先提交），向导会带你建分支、提交、推过网关，然后挂到这次运行上。`alamino -y` 会把这一切自动做完。
- **`/alamino`** —— agent skill。用 `/alamino <task>` 让编码 agent 完成一个任务再过网关，或用裸 `/alamino` 为已提交的工作过网关。它跑完流水线、让流水线应用安全的修复，并在任何需要人来拍板的地方停下来问你。

`alamino init` 会为 Claude Code 及其他 agent 安装 `/alamino` skill。底层上这个 skill 驱动的是 `alamino axi` —— 同一套审批流程的非交互式 TOON 接口。

完整的首次运行走查见[快速上手](https://mika-rafie.github.io/alamino/start-here/quick-start/)。

## 开发

```sh
make build   # 构建 bin/alamino（带版本信息）
make test    # 运行 go test -race ./...（不含 e2e 套件）
make e2e     # 运行打了标签的端到端 agent 旅程套件
make e2e-record # agent 线格式变化时，重新录制 e2e fixtures
make lint    # 检查生成的 skill 是否漂移，并跑 go vet ./...
make skill   # 重新生成已提交的 alamino skill 文件
make fmt     # 运行 gofmt -w .
make demo    # 重新生成 demo.gif 和 demo.mp4（需要 vhs 和 ffmpeg）
make docs    # 在 docs/dist 构建 Astro 文档站
```

完整 target 列表见 `Makefile`。

`make e2e-record` 会用真实的 `claude`、`codex`、`opencode` 和 `antigravity` CLI 覆盖 `internal/e2e/fixtures/`，会消耗真实 API 额度，提交前应当审查。


## 致谢

alamino 是 Kun Chen 的 [no-mistakes](https://github.com/kunchenguid/no-mistakes) 的分支（fork）。
本项目依据本仓库中的 MIT 许可证发布。
原始版权声明见 `LICENSE`。
