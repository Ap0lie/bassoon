# Phase 0：全局只读现状核验

## ① 阶段目标与范围

本阶段只核验本仓库当前 `main` 工作树中“智能体编排”相关能力的真实实现，并保存可复核的源码证据，为后续 Phase 1–5 定边界。本阶段已按任务简报启动本地环境：`make web` 成功，`http://localhost:8888` 返回 HTTP 200。

In-scope：工作流 DSL/画布、九类节点、调试运行、模板能力、MCP、AI 改节点六个落点，以及其保存、撤销和对话链路。

本阶段不做什么：不修改任何产品源代码、IDL、数据库、配置或测试；不创建工作流样例；不对画布交互作性能压测；不登录产品界面执行需账号和模型配置的端到端用例；不开始 Phase 1 或其后的设计/实现。

## ② 现状核验证据

### A. 五项功能现状表

| 项目 | 结论 | 真实证据与代码事实 |
| --- | --- | --- |
| 可视化画布 | 部分具备 | `frontend/packages/workflow/playground/src/workflow-playground.tsx` 使用 `WorkflowRenderProvider`，并装配 `WorkflowPageContainerModule`、`WorkflowNodesContainerModule` 与 `WorkflowHistoryContainerModule`；`src/container/workflow-page-contribution.ts` 将 `NodeRender` 注册给 `FlowRendererKey.NODE_RENDER`。`frontend/packages/workflow/fabric-canvas/src/components/fabric-editor/fabric-editor.tsx` 的确有右键菜单监听，但它是独立 Fabric 画布包，不是该工作流编辑器的主渲染路径。拖拽、连线、删除已有实现与单测目录，但本阶段没有真实交互性能基准，不能证明“≤1 秒”。 |
| 九类基础节点 | 已具备（源码层面）；运行态待 Phase 1 逐例验证 | 后端 `backend/domain/workflow/entity/node_meta.go` 定义 `NodeTypeLLM`、`NodeTypeKnowledgeRetriever`、`NodeTypeSelector`、`NodeTypeVariableAssigner`、`NodeTypeHTTPRequester`、`NodeTypeCodeRunner`、`NodeTypeOutputEmitter`、`NodeTypeLoop`、`NodeTypeSubWorkflow`，并在 `NodeTypeMetas` 配置可用元数据。前端 `frontend/packages/workflow/playground/src/nodes-v2/constants.ts` 的 `NODES_V2` 注册 LLM、Dataset、If、Set-variable、HTTP、Code、Output、Loop、Sub-workflow；相应节点目录均有 `form-meta.tsx` 和/或 `data-transformer.ts`，LLM 使用 `nodes-v2/llm/llm-form-meta.tsx` 与 `llm-node-registry.ts`。 |
| 调试运行 | 已具备（单节点和全流程）；失败显示待 Phase 1 UI 验收 | `idl/workflow/workflow_svc.thrift` 暴露 `WorkFlowTestRun`、`GetWorkFlowProcess`、`GetNodeExecuteHistory`、`WorkflowNodeDebugV2`。`frontend/packages/workflow/playground/src/services/workflow-operation-service.ts` 调用 `workflowApi.WorkFlowTestRun` 与 `WorkflowNodeDebugV2`；`components/test-run/test-run-button/single-node.tsx` 提供单节点运行，排除 Start、End、If。`workflow/test-run/`、`workflow/test-run-next/` 均存在；`components/test-run/log-navigation-v2/index.tsx` 以 `NodeResult.errorInfo` 过滤错误项。`workflow/nodes/src/setting-on-error/` 也存在。尚未登录并构造失败工作流，因此“红色定位”和逐节点输入输出为源码证据，非已观察的界面证据。 |
| 模板市场 | 部分具备 | `backend/domain/template/` 只有 `entity/template.go`、`internal/dal/`、`repository/repository.go`，没有 service 层。不存在独立“模板市场”路由页，但工作流编辑器内有模板面板：`frontend/packages/workflow/playground/src/components/template-panel/`。`use-workflow-template-list.ts` 调 `workflowApi.GetExampleWorkFlowList`；`template-card.tsx` 取得模板 schema 后以 `WorkflowSaveService.reloadDocument` 写入当前编辑器并 `highPrioritySave`。后端现有 `backend/application/workflow/workflow.go:GetExampleWorkFlowList` 从 `consts.TemplateSpaceID` 的工作流草稿读取；这不是 `domain/template` 服务化市场，也不是“复制为独立副本”的完整市场页。 |
| MCP 工具接入 | 部分具备（骨架，不能实际调用） | `backend/crossdomain/plugin/consts/consts.go` 定义 `PluginTypeOfMCP`；`backend/domain/plugin/service/tool/invocation_mcp.go` 存在 `NewMcpCallImpl`，但 `mcpCallImpl.Do` 直接返回 `errors.New("mcp call not implemented")`。前端生成代码含 MCP 管理 API（`frontend/packages/arch/idl/src/auto-generated/prompt_api/index.ts`），但未找到与本开源后端对应的完整管理服务。`backend/domain/connector/entity/connector.go` 只映射发布 Connector 到 `developer_api.ConnectorInfo`，`service/connector.go` 只列出/查询 Connector，职责不是 MCP。 |

九类节点的后端—前端映射如下，供 Phase 1 按一类一个最小运行例逐项验收：

| 需求节点 | 后端 NodeType | 前端注册/配置落点 |
| --- | --- | --- |
| 大模型对话 | `NodeTypeLLM` | `nodes-v2/llm/llm-node-registry.ts`、`llm-form-meta.tsx` |
| 知识库检索 | `NodeTypeKnowledgeRetriever` | `node-registries/dataset/dataset-search/` |
| 条件分支 | `NodeTypeSelector` | `node-registries/if/form-meta.tsx`、`data-transformer.ts` |
| 变量赋值 | `NodeTypeVariableAssigner` | `node-registries/set-variable/form-meta.tsx`、`data-transformer.ts` |
| HTTP 请求 | `NodeTypeHTTPRequester` | `node-registries/http/form-meta.tsx`、`data-transformer.ts` |
| 代码执行 | `NodeTypeCodeRunner` | `node-registries/code/form-meta.tsx`、`data-transformer.ts` |
| 回复输出 | `NodeTypeOutputEmitter` | `node-registries/output/form-meta.tsx`、`data-transformer.ts` |
| 循环 | `NodeTypeLoop` | `node-registries/loop/form-meta.tsx`、`data-transformer.ts` |
| 子流程 | `NodeTypeSubWorkflow` | `node-registries/sub-workflow/form-meta.tsx`、`data-transformer.ts` |

### B. AI 改节点六条落点问题

1. **编辑器布局与左侧面板挂点**：`frontend/packages/workflow/playground/src/use-workflow-playground.tsx` 是一个用于预览/运行的 hook，实例化 `WorkflowPlayground`；实际总装配在 `src/workflow-playground.tsx`，其中的 `WorkflowRenderProvider` 装配容器。`src/container/workflow-page-container-module.ts` 向 Inversify 容器绑定 `WorkflowSaveService`、`WorkflowOperationService`、`WorkflowRunService`、`WorkflowFloatLayoutService` 等服务；`src/components/workflow-container/index.tsx` 是页面区域和浮层的实际 React 容器。Phase 2 左侧“编程助手”应由 `WorkflowContainer` 的布局层新增可控区域/浮层，并通过新 container binding 获取文档与保存服务；不应把它塞进只读的 `useWorkflowPlayground` hook。

2. **节点菜单、右键与选中节点**：主工作流节点的“…”菜单真实位于 `frontend/packages/workflow/playground/src/form-extensions/components/node-header/index.tsx`，其 `Dropdown.Menu` 已有重命名、复制、删除和 `extraOperation` 扩展，且组件已持有 `node`（`WorkflowNodeEntity`）。这比改 `fabric-canvas` 更符合实际编辑器。`fabric-canvas/src/components/fabric-editor/fabric-editor.tsx` 的 `contextmenu` 是独立 Fabric 包的实现，不能直接作为 FlowGram 工作流节点右键扩展点。当前未在工作流主渲染链路找到一个现成的 FlowGram `onContextMenu` 节点回调；Phase 2 应在 `NodeHeader` 增加“…”入口，并在对应 FlowGram 节点渲染组件中补右键入口，两个入口复用同一 action。当前节点实体可由 `NodeHeader` 的 `node` 参数直接取得。

3. **读取/写回 DSL、撤销与保存**：后端 `backend/domain/workflow/entity/vo/canvas.go:Canvas` 存储 `Nodes`、`Edges`；`Node` 含 `ID`、`Type`、`Meta`、`Data`、`Blocks`、`Edges`。`WorkflowSaveService.loadDocument` 把 `WorkflowJSON` 载入 `WorkflowDocument`，`save` 则执行 `workflowDocument.toJSON()`。前端保存动作是 `WorkflowOperationService.save`（`src/services/workflow-operation-service.ts`），调用生成的 `workflowApi.SaveWorkflow`；IDL 名为 `WorkflowService.SaveWorkflow`，路由 `/api/workflow_api/save`（`idl/workflow/workflow_svc.thrift`）。撤销/重做由 `workflow/history/src/workflow-history-shortcuts-contribution.ts` 的 `HistoryService.undo/redo` 和 `ctrl/meta z` 提供。`WorkflowNodeData.updateNodeData` 只更新节点的前端附加实体数据，不能用于持久化 DSL 的 `data`。Phase 2 必须经 FlowGram 表单/命令式模型更新真实节点表单数据，再验证该更新触发内容变更并进入 HistoryService；具体 API 名不能在本阶段凭猜测写死，应以已安装适配器类型与最小原型验证为准。

4. **复用聊天 UI 的 adapter 方案**：真实导出文件为 `frontend/packages/common/chat-area/chat-area/src/index.tsx`（不是任务简报所写的 `.ts`），导出 `ChatArea`、`ChatAreaProvider`、`useChatArea`、`useSendTextMessage` 等。`frontend/packages/agent-ide/chat-debug-area/src/index.tsx` 的 `BotDebugChatArea` 以 `<ChatArea componentTypes={...}>` 注入 `ReceiveMessageBox`、`MessageBoxActionBarAdapter`、工作流消息渲染和输入插槽；其 provider adapter 在 `agent-ide/chat-area-provider-adapter/src/provider/index.tsx`，负责 `requestToInit`、botId、插件注册和消息列表初始化。Phase 2 应新建“助手外壳 + 专用 Provider adapter”：外壳管理开闭、当前节点与提示上下文；adapter 提供独立会话初始化和发送实现；内核复用 `ChatArea` 的消息列表、输入框、生命周期和 `componentTypes`，不重写聊天 UI。

5. **模型返回 `data` 的校验/转换**：`frontend/packages/workflow/nodes/src/workflow-json-format.ts:WorkflowJSONFormat` 在初始化/提交阶段全局处理变量、批处理兼容和节点 `nodeMeta`，其中 `formatNodeOnSubmit` 是工作流 JSON 提交通路。每个节点的 `formMeta` 决定本节点校验和 DTO 转换；例如 `node-registries/code/form-meta.tsx:CODE_FORM_META` 配置 `validate`、`formatOnInit`、`formatOnSubmit`，其 `data-transformer.ts:transformOnSubmit` 输出后端 `data`。LLM 在 `nodes-v2/llm/llm-node-registry.ts:LLM_NODE_REGISTRY` 注册 `formMeta: LLM_FORM_META`。Phase 2 应把严格 JSON patch 合并到当前节点 data 的候选副本，按节点 registry 执行 `formatOnInit`/form 校验/`formatOnSubmit`，并让 `WorkflowJSONFormat` 参与全局格式化；任一步失败即只显示错误，不写文档。`nodes-v2/constants.ts` 仅注册 `START_NODE_REGISTRY`/`END_NODE_REGISTRY`，没有单独系统节点判定函数；实际判定应复用 `shortcuts/contributions/copy/is-system-nodes.ts` 对 `StandardNodeType.Start` 和 `StandardNodeType.End` 的比较，并同时尊重 start/end registry 的 `deleteDisable`、`copyDisable` 和 `headerReadonly`。

6. **模型对话/SSE 链路及 JSON 强约束可行性**：`common/chat-area/chat-core/src/request-manager/request-config.ts` 的 `RequestScene.SendMessage` 指向 `POST /api/conversation/chat`。后端 `backend/api/handler/coze/agent_run_service.go:AgentRun` 用 Hertz SSE 设置 `X-Accel-Buffering: no`；`backend/application/conversation/openapi_agent_run.go:pullStream` 转发 agent 流事件。IDL 为 `idl/conversation/agentrun_service.thrift:AgentRun`。此链路面向“已有 bot 的 AgentRun”，请求需要 bot、会话和 scene；`arch/bot-api` 本身主要包装常规/IDL API，`arch/bot-http` 是 Axios 设施，均未暴露一个可直接指定已配置模型、JSON mode 或 function calling schema 的通用前端补全接口。因此可复用聊天 UI 和已有 SSE 基础设施，但“P1 仅前端、直接调用现成对话并强制 JSON mode/function calling”尚无已验证的现成接口。Phase 2 需要先用真实已配置 bot 做技术验证；若不能稳定约束 JSON，必须在 Phase 2 文档中作为阻塞/改后端 API 的授权点，而不是伪造能力。

### 差异修正

- 仓库实际路径是 `/Users/hangwang/Documents/ChatGPT/agent/coze-studio`，不是任务简报所述的 `/Users/hangwang/code26`；当前分支为 `main`。
- `frontend/packages/workflow/fabric-canvas` 存在，但工作流主编辑器的渲染链路是 `WorkflowRenderProvider` + `@flowgram-adapter/free-layout-editor`；不能把 Fabric 包当作主工作流节点菜单的唯一落点。
- `common/chat-area` 的入口为 `chat-area/src/index.tsx`，不是 `src/index.ts`。
- 前端存在编辑器内模板面板和“示例工作流”接口，不存在独立模板市场页；`domain/template` 虽无 service，但当前模板展示主要不由该领域模块提供。
- MCP 不是完全不存在：插件域有 MCP 类型和调用入口，但调用实现明确未完成，因此结论为“部分具备、不可实际调用”。

## ③ 方案设计

Phase 0 的交付仅是事实基线文档，没有产品数据流或接口改动。后续阶段必须以如下已验证边界设计：工作流编辑态始终是 `WorkflowDocument` 中的结构化 JSON；节点配置改动从节点表单/文档模型产生内容变更，`WorkflowSaveService` 将 `toJSON()` 结果经 `WorkflowOperationService.save` 调 `SaveWorkflow` 持久化；HistoryService 管理撤销/重做。任何 AI 能力只能针对该 DSL，不得引入磁盘 Python 文件、文件监听或 Python 工程假设。

预定设计约束：Phase 1 先以端到端样例证实已有能力；Phase 2 仅增加单节点 `data` 编辑，使用 NodeHeader/主节点渲染路径、ChatArea adapter、节点 formMeta 与 FlowGram history；Phase 3 应复用现有示例工作流读取/复制能力，但补齐独立市场语义；Phase 4 才允许图结构变更；Phase 5 只能扩展 plugin 域，不能改 connector 域。

## ④ 改动文件清单

本阶段只新增/修改文档：

| 文件 | 动作 | 内容 |
| --- | --- | --- |
| `docs/coze-phase/00-任务简报.md` | 新增 | 用户提供任务简报的逐字副本，作为跨会话权威来源。 |
| `docs/coze-phase/README.md` | 新增、修改 | Phase 总索引、状态、文档链接与交接摘要骨架。 |
| `docs/coze-phase/phase-00-现状核验.md` | 新增 | 本 Phase 的只读核验设计文档。 |

没有产品代码、IDL、DDL、配置、测试或锁文件改动。

## ⑤ 风险与回滚

影响面仅为仓库内 Markdown 文档；运行 `make web` 创建/重建了本地 Docker 容器，但没有改变应用源码或数据库 schema。文档提交前的回滚方式是移除本 Phase 新增文档或用 Git 恢复文档版本；不得对用户已有改动执行宽泛重置。

对功能的主要风险是错误地把独立 `fabric-canvas` 当成 FlowGram 工作流主画布，或把 MCP/模板误判为完全缺失。后续设计已在本文件的“差异修正”中限定：节点编辑应以 FlowGram 文档和 `NodeHeader` 为主；模板能力为部分具备；MCP 调用为未实现骨架。HistoryService 的既有 Ctrl/Cmd+Z 栈不在本阶段改动。

## ⑥ 验收标准与验证方法

本阶段验收的是证据完整性与文档落盘，不是功能验收。

1. 运行 `cmp -s <原始粘贴文本> docs/coze-phase/00-任务简报.md`；预期退出码为 0，说明简报逐字保存。
2. 运行 `make web`；预期 Docker Compose 成功启动/重建服务。已执行，退出码 0。
3. 运行 `curl -I http://localhost:8888`；预期 HTTP 200。已执行，返回 `HTTP/1.1 200 OK`、`Server: nginx/1.25.5`。
4. 运行 `docker compose -f docker/docker-compose.yml ps`；预期 `coze-web`、`coze-server` 及依赖均为 Up，数据库、缓存和向量库显示 healthy。已执行并满足。
5. 运行 `rg -n 'NodeTypeLLM|NodeTypeKnowledgeRetriever|NodeTypeSelector|NodeTypeVariableAssigner|NodeTypeHTTPRequester|NodeTypeCodeRunner|NodeTypeOutputEmitter|NodeTypeLoop|NodeTypeSubWorkflow' backend/domain/workflow/entity/node_meta.go`；预期九类枚举全部命中。
6. 运行 `rg -n 'SaveWorkflow|WorkFlowTestRun|WorkflowNodeDebugV2|GetExampleWorkFlowList' idl/workflow/workflow_svc.thrift`；预期保存、调试、示例工作流接口全部命中。
7. 运行 `rg -n 'mcp call not implemented' backend/domain/plugin/service/tool/invocation_mcp.go`；预期命中，证明 MCP 尚不可调用。

画布响应 ≤1 秒、每类节点真实执行、失败节点视觉表现、模板复制隔离和 MCP 调用均需要认证账户、模型/外部资源或后续实现，留给相应 Phase 的 E2E 验收。

## ⑦ 任务拆解 checklist

- [x] 将任务简报逐字保存为 `docs/coze-phase/00-任务简报.md`。
- [x] 建立并更新 `docs/coze-phase/README.md` 索引骨架。
- [x] 运行 `make web`，并用 HTTP 与 Docker 状态确认本地服务可用。
- [x] 核验结构化 DSL、FlowGram 主渲染路径与保存/撤销链路。
- [x] 核验九类基础节点的后端枚举与前端注册/表单落点。
- [x] 核验全流程/单节点调试、日志与错误筛选证据。
- [x] 核验模板领域、编辑器模板面板和示例工作流读取链路。
- [x] 核验 MCP/connector 的真实职责与未实现边界。
- [x] 回答 AI 改节点六条落点并记录差异修正。
- [x] 写入本 Phase 文档并同步 README 状态。
- [x] 已收到用户口令“确认 Phase 0”。
- [x] 完成全局构建/lint 自验；宿主机未安装 Go、Node 或 Rush，已使用经校验的临时容器工具链完成。
- [x] 自验通过后写入“九、执行记录”、更新交接摘要并提交独立 Git commit。

## ⑧ 遗留问题与外部依赖

1. 画布 ≤1 秒只能在登录后的真实浏览器交互中测量；本阶段服务健康不等于性能达标。
2. 九类节点的“可执行”受模型配置、知识库、HTTP 目标和子流程样例影响；Phase 1 必须逐类构造最小运行用例。
3. 当前通用聊天接口是 bot AgentRun SSE，未发现可直接指定已配置模型并要求 JSON mode/function calling 的前端 API。Phase 2 的“纯前端快路”需真实验证；不稳定时需用户授权改为正规的后端/IDL 方案或调整范围。
4. `WorkflowNodeData.updateNodeData` 并非持久化节点 DSL 更新接口。Phase 2 必须先确认 FlowGram/节点表单的正式写入 API 与 HistoryService 记录方式。
5. 编辑器内模板面板会把模板 DSL 重载入当前工作流并保存，不能替代“独立市场页 + 复制出隔离副本”的产品验收。Phase 3 应先判断可复用的 `CopyWkTemplateApi`/`CopyWorkflow` 是否满足隔离语义。
6. MCP 调用器目前明确返回未实现；Phase 5 需要独立会话，并在 plugin 域设计 Server 配置、tools 拉取、注册、超时/错误隔离和运行时调用。

## 九、执行记录

### 实际改动文件

- `docs/coze-phase/00-任务简报.md`：保存用户任务简报的逐字副本。
- `docs/coze-phase/README.md`：维护 Phase 状态、索引和阻塞交接信息。
- `docs/coze-phase/phase-00-现状核验.md`：记录核验结果、本阶段 checklist 与本执行记录。

未修改产品代码、IDL、DDL、配置、测试或运行中的镜像。

### 实现说明

按确认后的 Phase 0 范围完成文档交付与环境自验。运行镜像是无源码、无 Go/Node/Rush 的运行时镜像，因此使用临时官方 Go 1.24 容器；前端 Node 22.23.2 ARM64 发行包从 nodejs.org 下载并以其官方 `SHASUMS256.txt` 校验后运行 Rush。当前运行的 `coze-server`、`coze-web` 均未被替换。

### 构建、测试与界面/日志验收证据

| 项目 | 命令/证据 | 结果 |
| --- | --- | --- |
| 本地服务 | `make web`；随后 `curl -I http://localhost:8888` | 通过：Docker Compose 启动成功；HTTP 200，Nginx 响应。 |
| 后端构建（Alpine） | `golang:1.24-alpine` 中运行 `go build ./...` | 不通过：`github.com/mattn/go-sqlite3` 的 CGO 代码在 musl 下缺少 `pread64`、`pwrite64`、`off64_t`。该临时平台不代表源码改动错误。 |
| 后端构建（glibc） | `golang:1.24` 中运行 `go build ./...` | 通过：退出码 0。 |
| 前端 build | `rush build --to @coze-studio/app`（Node 22.23.2、Rush 5.147.1） | 通过：256 个操作成功、1 个无操作；应用包成功但有既有 warning，见下。 |
| 前端 lint | `rush lint`；最终采用单并发与 `NODE_OPTIONS=--max-old-space-size=6144` | 通过：退出码 0；最终汇总为 254 个缓存命中、5 个实际成功（含 `@coze-workflow/playground` 与 `@coze-studio/app`）。 |

前端 build 的既有 warning：Browserslist/caniuse-lite 数据已过期；`frontend/config/tailwind-config/package.json` 未声明 `type`，导致 `design-token.ts` 以 ES module 重解析。临时 npm 安装还报告多个上游弃用包及 1 个审计高危项；它们未由本 Phase 改动引入，且不在只读核验范围内，因此仅记录、未改依赖或配置。

首次全仓 lint 在默认资源下出现两次非代码失败：8 并发时 `@coze-studio/bot-detail-store` 被系统以 137 终止，2 并发时 `@coze-workflow/playground` 在约 2 GiB Node 堆触发 134。后者在 Docker 可用约 7.75 GiB 内存、单并发、6 GiB Node 堆下成功；最终同一全仓 lint 命令退出码为 0。未产生本阶段业务代码，因此没有新增单元测试；全局构建/lint 验收已满足。

### 与计划的偏差

原计划将使用临时 Go/Node 镜像完成全部验证。实际后端必须从 Alpine 改用 glibc Go 1.24 镜像，原因是项目全量包构建包含 CGO SQLite 依赖；Docker Hub 拉取 `node:22-alpine` 连续 EOF 后，改为校验 nodejs.org 官方 Node 22.23.2 发行包。全仓 lint 的默认并发/堆配置在 Docker Desktop 内存限制下两次失败，最终仅调整临时容器资源与运行时参数取得通过。没有为绕过这些环境问题修改 Dockerfile、依赖版本、项目配置或运行镜像。

### 遗留问题

1. Phase 1 仍需认证账户、可用模型、知识库、HTTP 目标和子流程样例，逐项验证九类节点的真实交互、运行和错误视觉反馈；这不是 Phase 0 的阻塞项。
2. Phase 2 仍需先验证既有 AgentRun SSE 是否能稳定获得严格 JSON；若不能，须按本文件“八、遗留问题”取得正规的后端/IDL 方案授权。
