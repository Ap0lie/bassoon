# bassoon 项目工作交接：Phase 1

## 项目与代码状态

- 仓库：`Ap0lie/bassoon`（GitHub 私有仓库）
- 工作分支：`main`
- 已推送提交：`67ebe8a7 feat(phase-1): 迁移 Kimi K2.6 兼容约束`
- 当前工作区：交接材料创建前应为干净状态；提交本文件后再次确认。
- 当前阶段：Phase 1 执行中，已完成 Kimi K2.6 兼容性改造；P0 画布与九类节点的完整 E2E 仍未完成。

## 已完成内容

### Kimi K2.6 兼容约束

- 专用模型元数据：256K 上下文、默认 `max_tokens=32768`、Base64 多模态输入能力。
- 仅在 Moonshot 中国站 `https://api.moonshot.cn/v1` 且模型为 `kimi-k2.6` 时启用专用请求逻辑。
- 请求不发送 `temperature`、`top_p`、`frequency_penalty`、`presence_penalty` 或 `max_completion_tokens`。
- 请求仅发送 `thinking.type=enabled|disabled`；默认启用思考，启用时将输出长度下限限制为 16000。
- Agent、RAG、意图识别及含工作流/插件/知识库工具的链路显式关闭思考。
- 历史 assistant 消息回放保留 `ReasoningContent`，保证多步工具调用可回传原始推理内容。
- 图片与视频仅允许 Base64 data URL 或 `ms://` 文件引用；外部 URL 转换失败时返回错误，不再静默透传。
- 创建与读取模型配置时均规范化到 Moonshot 中国站基础地址，并强制 Base64 多模态输入。

### 安全处理

- 认证与模型配置请求体已在访问日志中脱敏。
- 创建模型时不再打印连接凭据。
- API Key、账号密码、会话文件、临时响应及本地 `.env` 未进入 Git。

### 已有验证证据

- Moonshot `/models` 实测存在 `kimi-k2.6`，不含已退役模型标识。
- 非流式（思考关闭）与流式（默认思考）实测成功；流式先收到 `reasoning_content`，后收到正文。
- Coze 中模型配置已保存；隔离工作流 `phase1_e2e` 的模型节点单节点调试成功。
- 静态检查通过：`git diff --check`、模型元数据 JSON 解析、全仓旧模型/旧 Moonshot 标识检索、强制工具选择与 `$web_search` 检索。

## 关键文件

| 目的 | 文件 |
| --- | --- |
| 阶段验收、证据与未完成项 | `docs/coze-phase/phase-01-p0功能核验.md` |
| 阶段索引 | `docs/coze-phase/README.md` |
| Kimi 请求构建 | `backend/bizpkg/llm/modelbuilder/openai.go` |
| 模型配置创建与日志保护 | `backend/api/handler/coze/config_service.go`、`backend/api/middleware/log.go` |
| 已保存模型的兼容配置 | `backend/bizpkg/config/modelmgr/model_get.go` |
| 工作流/工具/RAG 思考策略 | `backend/domain/workflow/internal/nodes/llm/llm.go`、`backend/domain/workflow/internal/nodes/qa/question_answer.go`、`backend/domain/workflow/internal/nodes/intentdetector/intent_detector.go` |
| Agent 运行时多模态处理 | `backend/domain/agent/singleagent/internal/agentflow/agent_flow_runner.go` |
| 工具调用推理历史回放 | `backend/domain/conversation/agentrun/internal/message_builder.go` |

## 待办

1. 上传全部非敏感配置与说明文档；提交前检查其中没有 API Key、密码、Cookie、会话头或本地路径。
2. 既然 Node 已恢复，先运行仓库 Git 钩子依赖所需的 Node 环境检查，再执行以下验证并把真实输出补入 Phase 1 文档：

   ```bash
   go test ./bizpkg/llm/modelbuilder ./bizpkg/config/modelmgr ./api/handler/coze \
     ./domain/workflow/internal/nodes/llm ./domain/workflow/internal/nodes/intentdetector \
     ./domain/agent/singleagent/internal/agentflow ./domain/conversation/agentrun/internal
   go build ./...
   ```

3. 登录本地 Coze，修复/重建模拟知识库文档；验收条件是产生分段与检索命中，再完成 RAG 节点 E2E。
4. 完成画布拖拽、连线、删除的可重复计时证据，以及九类基础节点的最小全流程运行。
5. 构造确定性失败节点，验证 `NodeResult`、错误列表、节点错误态和错误定位一致。
6. 将真实命令输出、E2E 请求/响应摘要写入 `phase-01-p0功能核验.md` 的“九、执行记录”，再创建后续 Phase 1 提交。

## 接手步骤

```bash
git checkout main
git pull --ff-only origin main
git log --oneline -3
```

随后依次阅读：

1. `docs/coze-phase/00-任务简报.md`
2. `docs/coze-phase/phase-00-现状核验.md`
3. `docs/coze-phase/phase-01-p0功能核验.md`
4. 本交接文档

不要重新录入或提交任何测试凭据。模型测试继续使用已保存的本地配置；若需要变更配置，保持中国站基础地址并遵守 K2.6 的请求约束。

## Git 钩子说明

此前机器缺少 Node，pre-commit/post-commit 钩子未能运行，静态检查仍正常。现在该问题已由用户解决；下一次提交应让钩子正常执行。若钩子失败，先记录其真实输出并修复依赖或代码，不要默认跳过。
