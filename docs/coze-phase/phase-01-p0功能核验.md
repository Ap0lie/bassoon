# Phase 1：P0 画布、九类基础节点与调试运行核验

## ① 阶段目标与范围

本阶段目标是在不替换 FlowGram、不重写已有功能的前提下，完成三项 P0 的真实产品核验：画布的拖拽/连线/删除响应，九类节点的可拖出/配置/最小运行，以及包含失败节点的调试运行、逐节点输入输出与错误定位。若真实核验明确发现小且局部的实现缺口，才在本阶段补齐；未证实缺口前不改产品代码。

In-scope：登录后进入工作流编辑器；用最小工作流逐项验证大模型对话、知识库检索、条件分支、变量赋值、HTTP 请求、代码执行、回复输出、循环、子流程；运行全流程与单节点调试；检查日志、输入、输出、错误信息和错误视觉状态；记录实际响应时间与证据。

本阶段不做什么：不替换画布引擎；不开发 AI 改节点、模板市场或 MCP；不新建节点类型、IDL、数据库迁移或第三方依赖；不为绕过外部资源而伪造模型响应；不在未收到“确认 Phase 1”前写任何功能代码。

## ② 现状核验证据

### 画布与编辑操作

- `frontend/packages/workflow/playground/src/components/workflow-container/index.tsx` 持有 `Sidebar`、`AddNodeModalProvider` 与 FlowGram 编辑区；`onDragOver` 计算坐标并交给 `addNodeRef.current?.handleAddNode`，证明主编辑器有节点拖入入口。
- `frontend/packages/workflow/playground/src/services/workflow-edit-service.ts` 的 `addNode`、`deleteNode` 与 `recreateNodeJSON` 是节点编辑服务；`form-extensions/components/node-header/index.tsx` 调用 `editService.deleteNode(node, true)`。
- 连线拖动由 `frontend/packages/workflow/playground/src/services/workflow-line-service.ts` 的 `onDragLineEnd` 接入 drag service；删除快捷键逻辑位于 `shortcuts/contributions/delete/index.ts` 的 `removeNode`。
- 这些源码事实证明能力链存在，但没有可认证编辑器会话可测量拖拽、连线、删除的“≤1 秒”真实响应；不得以构建通过替代性能/E2E 结论。

### 九类节点

`backend/domain/workflow/entity/node_meta.go` 同时声明并配置以下节点元数据：`NodeTypeLLM`、`NodeTypeKnowledgeRetriever`、`NodeTypeSelector`、`NodeTypeVariableAssigner`、`NodeTypeHTTPRequester`、`NodeTypeCodeRunner`、`NodeTypeOutputEmitter`、`NodeTypeLoop`、`NodeTypeSubWorkflow`。前端注册与配置落点已在 Phase 0 核验，Phase 1 将以运行态复验：

| 节点 | 最小运行例与依赖 | 前端配置证据 | 运行态状态 |
| --- | --- | --- | --- |
| 大模型对话 | Start → LLM → Output；需可用模型 | `nodes-v2/llm/llm-form-meta.tsx` | 待认证 E2E |
| 知识库检索 | Start → Dataset → Output；需已索引知识库 | `node-registries/dataset/dataset-search/` | 待认证 E2E |
| 条件分支 | Start → If → 两个 Output；输入变量 | `node-registries/if/form-meta.tsx` | 待认证 E2E |
| 变量赋值 | Start → Set-variable → Output | `node-registries/set-variable/form-meta.tsx` | 待认证 E2E |
| HTTP 请求 | Start → HTTP → Output；需可访问测试 URL | `node-registries/http/form-meta.tsx` | 待认证 E2E |
| 代码执行 | Start → Code → Output；使用确定性脚本 | `node-registries/code/form-meta.tsx` | 待认证 E2E |
| 回复输出 | Start → Output | `node-registries/output/form-meta.tsx` | 待认证 E2E |
| 循环 | Start → Loop（内部变量处理）→ Output | `node-registries/loop/form-meta.tsx` | 待认证 E2E |
| 子流程 | Start → Sub-workflow → Output；需可运行子流程 | `node-registries/sub-workflow/form-meta.tsx` | 待认证 E2E |

### 调试运行与错误定位

- IDL `idl/workflow/workflow_svc.thrift` 暴露 `WorkFlowTestRun`（`/api/workflow_api/test_run`）和 `WorkflowNodeDebugV2`（`/api/workflow_api/nodeDebug`）；`frontend/packages/workflow/playground/src/services/workflow-operation-service.ts` 分别由 `testRun` 与 `testOneNode` 调用。
- `idl/workflow/workflow.thrift:NodeResult` 含 `NodeId`、`Input`、`Output`、`ErrorInfo`、`ErrorLevel`；后端 `backend/application/workflow/workflow.go:TestRun`、`NodeDebug` 与 `convertNodeExecution` 生产这些结果。
- 前端错误呈现链明确存在：`components/test-run/log-navigation-v2/index.tsx` 按 `errorInfo` 筛选，`page-selector.tsx` 区分 `Error` 与 warning，`execute-result-side-sheet-v2/error-item.tsx` 展示 `errorInfo`；旧 `workflow/test-run` 也在 `features/problem` 与 `log-detail/pagination` 消费同一字段。
- 本地 `http://localhost:8888/` 可打开；用户提供的本地测试会话已登录到 `Personal Space`。真实应用 `app` 中的已有工作流 `teat1` 可打开；其大模型节点显示“未配置模型”。
- 已在同一应用创建隔离工作流 `phase1_e2e`（ID `7680756530848727040`）。默认代码节点、大模型、知识库检索、选择器、变量赋值、HTTP 请求、输出、循环和子流程 `teat1` 均可成功加入；首次代码节点添加为 615 ms，其余八项为 476–645 ms，均小于 1 秒。子流程须从“添加工作流”对话框选择已有流程，这是实际配置语义而非节点缺失。
- 对未配置的 `teat1` 点击“试运行”，页面显示“错误列表 / 结束 / 引用变量不存在”。对 `phase1_e2e` 的代码节点点击单节点调试，空输入被拦截并弹出“有校验失败的节点，请检查配置好后，再试运行”；填入确定性字符串后输入字段已有效，但单节点运行按钮仍禁用，表明该临时画布还存在其他调试前置条件。尚未观察到已执行的 `NodeResult.input/output` 或运行时红色节点态。
- 对完整的 `phase1_e2e` 点击“试运行”后，约 1.94 秒内出现运行输入面板和“错误列表 / 结束 / 引用变量不存在”。这是保存前配置预校验，非后端执行结果；它与 `teat1` 的同类错误一致，进一步确认当前结束节点尚未绑定有效返回变量。
- 删除核验的目标是临时代码节点 `194128`。当前 UI 的节点操作菜单由悬停触发；自动化浏览器未能稳定暴露其“删除”条目。一次画布快捷键尝试误选并删除了临时“循环”节点，已立即通过“添加节点”恢复，九类节点现仍齐全；`teat1` 与用户其他工作流未被修改。该次尝试不计入删除通过证据，后续须用可定位的菜单或人工操作复验。
- 用户已授权使用本地测试资源和模拟数据。已创建自定义文本知识库 `phase1_mock_knowledge`（ID `7680767047764017152`）；首份模拟文档的嵌入任务失败并保持 0 分段，原因是知识库默认嵌入服务未配置，不是对话模型连通性失败。模型新增的实际后端入口是 `POST /api/admin/config/model/create`，会先发送一次最小生成请求验证连通性；工作区 UI 只提供会话管理/变量设置，没有模型管理页面。
- 已用 Moonshot `/models` 真实响应确认 `kimi-k2.6` 在售、此前的退役标识不在列表中；`https://api.moonshot.cn/v1` 上的非流式（关闭思考）和流式（默认思考）调用均成功。流式响应先出现推理内容、后出现正文；模型配置已成功保存，隔离工作流的大模型节点单节点调试也已成功。临时会话和凭据文件不进入文档或仓库。

### 差异修正

- 任务简报将 `fabric-canvas` 作为画布定位提示；实际主编辑器是 FlowGram 的 `workflow/playground`。`fabric-canvas` 保留为独立包，Phase 1 不以它替换或重写主编辑器。
- 当前并非“能力缺失”的证据，而是“运行态尚未认证核验”。九类节点和调试通路已有源码与后端路由证据；是否对业务开发者可完整跑通，必须待真实 E2E 判定。

## ③ 方案设计

确认后，先使用现有账号进入本地工作空间，创建一个专用临时工作流并逐步添加节点；每次只增加一个最小可运行节点或分支，保存后执行并记录节点 ID、配置要点、输入、输出、耗时和截图/日志。画布操作使用浏览器计时，从拖拽/连线/删除动作开始到 DOM/视觉稳定为止；若超过 1 秒，先用浏览器控制台与 React/FlowGram 事件确认是渲染、网络还是数据规模问题。

调试用例分两类：全流程用 `WorkFlowTestRun` 取得执行 ID 并读取节点结果；单节点用 `WorkflowNodeDebugV2`。失败用例优先采用确定性的 HTTP 4xx/5xx 或代码节点显式异常，避免消耗模型配额；应验证节点侧错误标记、错误列表、输入输出和 `ErrorInfo` 一致。所有临时业务数据仅在用户已授权的本地工作空间内创建；在验收完成后由用户决定是否保留或删除。

如果 E2E 发现明确的小修缺口，先把根因、精确文件、改动方案和回归面补写到本文件第④、⑦，再实施。Moonshot K2.6 的请求、工具链、推理历史和多模态输入兼容性缺口已由真实响应与源码确认，最小改动范围已在第④登记；其余 P0 能力仍不凭源码推断进行重复开发。

## ④ 改动文件清单

在收到确认前，本阶段仅改文档：

| 文件 | 动作 | 内容 |
| --- | --- | --- |
| `docs/coze-phase/README.md` | 修改 | 链接本阶段文档并将状态标为“待确认”。 |
| `docs/coze-phase/phase-01-p0功能核验.md` | 新增 | 记录本阶段范围、源码/页面证据、E2E 方案、验收表与外部依赖。 |

确认后的首轮执行仍只做 E2E 核验和文档更新，不预设产品代码文件。若缺口被证实，必须先在本表逐个补列精确文件及改动内容，再写代码，确保不扩大范围。

Moonshot k2.6 兼容性缺口已经由真实验证和源码确认，以下为本次确认后追加的最小改动范围：

| 文件 | 动作 | 内容 |
| --- | --- | --- |
| `backend/conf/model/model_meta.json` | 修改 | 为 `GPT/kimi-k2.6` 增加专用元数据：256K 能力描述、默认 `max_tokens` 32768，并启用 Base64 多模态输入；不声明 `temperature`、`top_p`、`frequency_penalty` 或 `presence_penalty`。 |
| `backend/bizpkg/llm/modelbuilder/openai.go` | 修改 | 当配置为 Moonshot 中国站的 `kimi-k2.6` 时，使用 `MaxTokens`（而非 `MaxCompletionTokens`）32768，清除四个受限采样字段，并通过 Eino `ExtraFields` 仅写入 `thinking.type=enabled|disabled`；其他 OpenAI 兼容模型保持原有行为。 |
| `backend/bizpkg/llm/modelbuilder/openai_test.go` | 新增 | 单元测试 k2.6 专用请求配置：base URL、模型 ID、32768、无受限采样字段，以及启用/禁用思考；覆盖非 Kimi OpenAI 配置不回归。 |
| `backend/bizpkg/llm/modelbuilder/llm_params.go` | 修改 | 增加只覆盖内部思考开关的参数辅助方法，避免重新写入或污染用户设置的受限采样字段。 |
| `backend/bizpkg/llm/modelbuilder/model_builder.go` | 修改 | 增加按模型配置或应用模型设置构建且可显式关闭思考的窄接口，供 Agent 工具链复用，避免把该策略扩散到其他兼容模型。 |
| `backend/bizpkg/llm/modelbuilder/model_builder_test.go` | 新增 | 回归验证显式思考开关辅助方法的空值与启用/禁用行为。 |
| `backend/bizpkg/llm/modelbuilder/builtin.go` | 修改 | 内置知识库召回模型经显式关闭思考的构建路径创建，覆盖不经过工作流节点参数的即时 RAG 链路。 |
| `backend/domain/workflow/entity/vo/modelmgr.go` | 修改 | 将工作流节点的内部 `EnableThinking` 选项映射到模型构建参数，供 RAG 与工具调用链显式关闭思考。 |
| `backend/bizpkg/config/modelmgr/model_get.go` | 修改 | 读取已保存的 `kimi-k2.6` 配置时归一化到 Moonshot 中国站基础地址并强制启用 Base64 多模态输入，覆盖旧记录中可能遗留的错误地址或关闭值。 |
| `backend/bizpkg/config/modelmgr/model_get_test.go` | 新增 | 回归验证已保存的 K2.6 模型记录在读取时获得中国站基础地址和 Base64 输入约束。 |
| `backend/api/handler/coze/config_service.go` | 修改 | 创建模型时对 Moonshot 中国站 `kimi-k2.6` 强制固定基础地址并启用 Base64 输入；移除可能包含连接凭据的调试日志。 |
| `backend/api/handler/coze/config_service_test.go` | 新增 | 回归验证创建 K2.6 模型时规范化允许的基础地址、拒绝错误资源路径；模型配置路径的凭据日志由中间件脱敏实现覆盖。 |
| `backend/domain/workflow/internal/nodes/qa/question_answer.go` | 修改 | 问答/RAG 的内部即时 LLM 调用显式设置 `EnableThinking=false`，避免检索链路被推理内容占满输出。 |
| `backend/domain/workflow/internal/nodes/intentdetector/intent_detector.go` | 修改 | 意图识别的即时 LLM 调用同样显式关闭思考，避免分类链路生成不参与判定的推理输出。 |
| `backend/domain/workflow/internal/nodes/intentdetector/intent_detector_test.go` | 新增 | 回归验证意图识别构建参数始终显式关闭思考。 |
| `backend/domain/workflow/internal/nodes/llm/llm.go` | 修改 | 在工作流、插件或知识库工具存在时，于建模前显式关闭思考；工具选择维持默认 `auto`/`none`，不引入强制指定工具。 |
| `backend/domain/agent/singleagent/internal/agentflow/agent_flow_builder.go` | 修改 | Agent/ReAct 建模使用显式关闭思考的构建路径，满足 Agent 与工具即时链路约束。 |
| `backend/domain/agent/singleagent/service/single_agent_impl.go` | 修改 | 将 Agent 输入多模态转换失败的可读错误返回给调用方，避免继续发出不合规请求。 |
| `backend/domain/conversation/agentrun/internal/message_builder.go` | 修改 | 保留历史 assistant 消息的原始 `ReasoningContent`，使多步工具调用能够原样回传上一轮推理内容。 |
| `backend/domain/workflow/internal/nodes/llm/prompt.go` | 修改 | Base64 转换失败时返回可读错误，不再把 `http(s)` 图片或视频 URL 原样传给 K2.6；仅已合规的 `data:*;base64,` 与 `ms://` 引用直通。 |
| `backend/domain/agent/singleagent/internal/agentflow/agent_flow_runner.go` | 修改 | Agent 输入预处理同样拒绝转换失败后的外部图片/视频直链，仅保留已合规的 `data:*;base64,` 与 `ms://` 引用，避免与工作流节点行为不一致。 |
| `backend/domain/workflow/internal/nodes/llm/prompt_test.go` | 修改 | 覆盖 Base64 成功、已合规引用直通和外部直链转换失败三种结果。 |
| `backend/domain/workflow/internal/nodes/llm/thinking_test.go` | 新增 | 回归验证工作流、插件和知识库工具存在时主/备用模型参数均显式关闭思考。 |
| `backend/domain/conversation/agentrun/internal/message_builder_test.go` | 修改 | 回归验证历史消息回放不清空推理内容。 |
| `backend/domain/agent/singleagent/internal/agentflow/agent_flow_runner_test.go` | 新增 | 回归验证 Agent 多模态输入不将转换失败的图片/视频直链送往 K2.6。 |
| `backend/api/middleware/log.go` | 修改 | 对认证和模型配置路径的请求体统一脱敏，避免密码或 API Key 进入服务日志。 |
| `docs/coze-phase/phase-01-p0功能核验.md` | 修改 | 记录 k2.6 迁移、真实 API 验证、grep 证据与执行结果；不记录密钥。 |

`OpenAI SDK` 的 Python `hasattr/getattr` 要求不适用于本仓库：全仓不存在调用该 SDK 的 Python 生产代码；现有 Go 流式节点读取 `schema.Message.ReasoningContent`。本次不引入 Python 文件或伪造该分支。

## ⑤ 风险与回滚

主要风险是模型、知识库、HTTP 目标或子流程不可用造成的外部失败，被误判为编辑器缺陷；每次运行须记录依赖资源和原始 `ErrorInfo`。另一个风险是临时工作流/调试调用写入用户空间或产生模型费用，创建、运行、删除均须在用户提供可用资源和授权后进行。

当前文档改动可用 Git 恢复；确认后的 E2E 不改产品代码。若实施小修，必须只经既有 FlowGram 文档/HistoryService 路径改动节点，验证 Ctrl/Cmd+Z 后再保存；回滚先用历史栈撤销，必要时恢复该阶段独立 commit，绝不重置用户其他工作流或宽泛清理数据库。

## ⑥ 验收标准与验证方法

1. 登录后打开新/临时工作流，分别执行拖拽节点、从端口连线、删除节点；每项以浏览器计时与截图记录，预期界面稳定时间 ≤1 秒且无控制台错误。
2. 按上表九类节点各建一个最小工作流；保存、全流程运行，预期节点可选、可配置、可保存并产出与例子一致的结果。知识库、模型、HTTP、子流程的外部配置失败必须记录为依赖失败，不能伪称通过。
3. 对可调试节点点击单节点运行；预期请求经 `WorkflowNodeDebugV2` 成功返回，界面显示当前节点输入/输出。
4. 构造确定性失败节点并全流程执行；预期 `NodeResult.errorInfo` 非空、错误列表可导航、节点/分页显示错误态且原因可读。
5. 对 Phase 1 任何代码改动（若经确认后产生）运行 glibc Go 1.24 的 `go build ./...`、相关单元测试、前端 `rush build --to @coze-studio/app` 与 `rush lint`；临时容器的 Node 堆参数和所有 warning 均写入“九、执行记录”。
6. 对 K2.6 专用构建请求做单元级断言：基础地址精确为 `https://api.moonshot.cn/v1`，请求中无受限采样字段，`max_tokens` 默认 32768，且思考仅为 `enabled` 或 `disabled`；启用思考的请求不得低于 16000。
7. 对 Agent、RAG 和含工具的工作流断言思考关闭；确认没有强制指定工具的调用，并验证多步工具调用回放 assistant 历史时保留原始推理内容。
8. 对图片/视频输入断言只会向 K2.6 发送 Base64 `data:` 或 `ms://` 引用；外部直链或非 Base64 data URL 转换失败必须中止并显示错误，不能静默透传。

## ⑦ 任务拆解 checklist

- [x] 回读 Phase 0、任务简报和工作流/调试源码，确认既有能力边界。
- [x] 确认本地服务 HTTP 可用，并用浏览器检查工作空间入口。
- [x] 记录认证门槛，未在无账户情况下绕过登录或伪造业务数据。
- [x] 写入本 Phase 的 E2E 用例矩阵、错误用例与验收标准。
- [x] 收到用户口令“确认 Phase 1”。
- [x] 获得可登录账号及所需模型、知识库、HTTP、子流程资源的使用授权；以用户提供的本地测试会话及实际可用资源为准。
- [x] 创建隔离工作流 `phase1_e2e`，并核验九类节点均可添加、可显示基础配置；添加延迟均小于 1 秒。
- [ ] 执行画布拖拽/连线/删除计时核验；已获临时代码节点 `194128` 的单次删除确认，但自动化尚未可靠定位悬停菜单，且一次误删的临时循环节点已恢复，不计为通过证据。
- [ ] 按矩阵核验九类节点的配置、连线、全流程运行与输出；当前模型、知识库、HTTP、循环配置未就绪。
- [ ] 执行全流程、单节点和失败节点调试核验；目前仅获得前端预校验错误，尚无后端节点执行结果。
- [ ] 如发现小修缺口，先更新第④文件清单和本 checklist，再做最小改动与回归。
- [x] 以 Moonshot `/models` 确认 `kimi-k2.6` 存在且旧模型 ID 不在列表；完成 k2.6 非流式（思考关闭）和流式（默认思考）真实调用。
- [x] 创建并保存 K2.6 模型配置；隔离工作流的大模型节点单节点调试成功，已观察到推理内容与正文。
- [x] 按第④追加清单实现 K2.6 专用请求与模型配置读写约束、RAG/Agent/工具链思考关闭、推理历史回放、多模态 URL 限制，并新增针对这些约束的回归用例；待有 Go 运行环境后实际执行。
- [x] 为模型构建、内置知识库召回、Agent/ReAct 及工作流工具链增加思考关闭路径；非 K2.6 的 OpenAI 参数保留已有构建行为。
- [x] 为多步工具调用的推理内容原样回放，以及图片/视频 Base64 `data:`、`ms://` 直通和外部直链转换失败拒绝添加回归用例；待有 Go 运行环境后实际执行。
- [x] 为已保存模型与创建模型两条配置路径添加中国站基础地址、Base64 输入强制与凭据脱敏实现及回归用例；待有 Go 运行环境后实际执行。
- [x] 复查全仓无已下线模型 ID 与旧版 Moonshot 系列引用；确认没有强制工具选择或内置网页检索引用。
- [ ] 在可用 Go 运行环境中执行新增/相关单测与后端构建，记录真实输出。
- [ ] 修复或重建模拟知识库文档，确认索引产生分段与命中后执行 RAG 节点 E2E。
- [ ] 追加“九、执行记录”、全量自验、更新 README 交接摘要并创建独立 Phase 1 commit。

## ⑧ 遗留问题与外部依赖

1. 本地会话与 K2.6 模型配置均已验证；模拟知识库的首份文档因默认嵌入服务配置失败而没有分段或命中。需在不写入凭据的前提下修复本地模拟嵌入配置并重建文档，才能完成 RAG 节点与包含 RAG 的全流程运行。
2. 子流程 `teat1` 已可选入临时画布，但自身存在“引用变量不存在”的预校验错误，不能作为成功子流程运行证据。
3. 运行验证可能消耗模型/外部请求配额；当前仅有使用本地测试空间的授权，临时工作流是否在阶段结束后删除仍待用户决定。
4. 当前社区版工作区 UI 未暴露模型管理；`/api/admin/config/model/create` 需要管理会话，现有业务会话不能直接使用。已通过用户授权的本地测试登录创建独立临时会话，不读取浏览器 Cookie；后续只复用已验证的 K2.6 配置，不更换密钥或切换站点。
5. 节点菜单的悬停交互与当前自动化控制存在兼容性限制；删除与连线仍须在真实 UI 中取得可重复的实测证据，不能用源码存在替代。
6. 已确认的 K2.6 兼容性小修必须按第④逐项完成并通过回归；画布和其余节点的产品能力仍须以 E2E 实测为准，不能由静态源码替代。
