# AGENTS.md — eGZ-Board 项目规范

教学信息壁纸软件（电子班牌）。前端展示高考倒计时、公告栏、课程表、天气，支持多班级管理与班牌（Android）展示。

## 技术栈

| 端 | 技术 |
| --- | --- |
| 后端 | Go 1.18 / GoFrame v2.7（`ghttp` 路由）/ GORM + SQLite（`resource/database/data.db`）/ excelize 解析课程表 xlsx |
| 前端 | Vue 3 + Vite + Vuetify 3 + UnoCSS + Pinia + wangEditor（git 子模块 `eGZ-board-front`） |
| 构建 | `make`（内含 `hack/hack.mk`），Dockerfile 在 `manifest/docker/` |

> 注意：前端 `eGZ-board-front` 是 git submodule，改动需在子模块内单独提交。

## 目录结构

```
├── api/                  # API 定义层（GoFrame 约定）：每个模块一个包
│   └── <module>/v1/      # 请求/响应结构体（Req/Res），含 g.Meta 路由元信息
│   └── <module>/<module>.go   # ISet<V>V1 接口（由 GoFrame CLI 生成，勿手改）
├── internal/
│   ├── cmd/              # cmd.go：HTTP server 装配 + 路由分组注册
│   ├── controller/       # 控制器：实现 api 层接口，处理请求、调用 dao
│   ├── service/          # 业务服务（天气代理、定时器），目前逻辑较薄
│   ├── logic/            # 业务逻辑层（GoFrame 分层约定，当前为空）
│   ├── dao/              # 数据访问层：封装 GORM 查询
│   └── model/
│       ├── data.go       # GORM 初始化 + AutoMigrate + 初始数据种子
│       └── entity/       # 数据库实体结构体
├── utility/              # 通用工具：HTTP 请求、按 IP 取班级等
├── manifest/             # 配置、Docker、K8s 部署文件
├── public/               # 前端构建产物（server root）
├── resource/             # 运行时资源：database/ 数据库、upload/ 上传文件
└── shared/api/           # 前后端共享 API 规则文档（见下）
```

## 后端架构与分层

请求链路：`HTTP -> ghttp 路由 -> api/<module>/v1 的 Req/Res -> controller/<module> -> dao -> model(数据库)`

- **路由注册**：全部集中在 `internal/cmd/cmd.go`，以 `/api/<module>` 分组绑定各模块 `NewV1()`。
- **路由命名**：GoFrame 由控制器方法名自动生成 URL，如 `GetSettingValue -> /setting/get-setting-value`（蛇形）。
- **响应包裹**：`/api` 分组挂载 `ghttp.MiddlewareHandlerResponse`，所有响应统一为 `{code, message, data}`，前端读取 `res.data.data.<...>`。
- **多班级隔离**：无登录态。通过客户端 IP 在 `device` 表查班级（`utility.GetClassID(ctx)`），设备绑定/解绑即增删 `device` 记录。
- **API 定义层**：`api/<module>/v1/*.go` 内的 Req/Res 是唯一契约来源。新增接口 = 加 Req/Res 结构体 + 加接口方法 + controller 实现。顶层 `api/<module>/<module>.go` 与 `internal/controller/<module>/*_new.go` 为 GoFrame CLI 生成，`DO NOT EDIT` 头可保留但要与新方法保持同步。
- **service 层**：当前承载天气代理（`service/WeatherProxy` + `gcron` 每 3 分钟刷新，`service/proxy_timer.go`）。

## 数据库约定

- 单文件 SQLite：`resource/database/data.db`（`resource/database` 与 `resource/upload` 由 `main.go initPath()` 自动创建）。
- 实体放 `internal/model/entity/`，全部用 GORM tag（`gorm:"primary_key;auto_increment;not null"`），主键为 `Id`。
- 表结构变更：修改实体后依赖 `AutoMigrate` 自动迁移（`model.InitData()`）。新表需加入 `AutoMigrate` 调用。
- 初始数据种子也在 `model.InitData()`，注意只首次（id=1 用户不存在时）写入。
- 新增数据访问方法统一封装进 `internal/dao/`，controller 不要直接写 SQL。

## 前后端 API 规则

前后端共享的接口契约与规范见 `shared/api/README.md`，改动接口时必须同步更新该文档。

## 前端约定

- 无统一 axios 封装/api 层：各组件 `import axios from "axios"` 直接调 `/api/...`，取数统一 `res.data.data.<字段>`。
- 开发时经 Vite proxy 转发到 `http://192.168.2.47:5020`（`vite.config.mjs`）；生产由 Go 后端同源托管 `public/`。
- 路由由 unplugin-vue-router 按 `src/pages/**` 自动生成（文件路由，勿手写路由）。
- 组件命名 kebab-case，文件自动路由。Vuetify 组件自动导入。
- 业务接口一律加统一响应解包；改后端字段名时同步排查前端所有调用处。

## 开发命令

```bash
# 后端
go run main.go                 # 启动 :5020，自动建 resource/ 目录
go build -o main .             # 编译
# 前端（在 eGZ-board-front/ 子模块内）
yarn dev                       # 开发服务器 :3000
yarn build                     # 构建到 dist/，由后端 public/ 托管
yarn lint                      # ESLint --fix
```

接口文档：启动后访问 `/swagger` 与 `/api.json`（见 `manifest/config/config.yaml`）。

## 提交工作流（重要）

本仓库遵守「edits.md 驱动提交」规则：

1. Agent 每完成一处改动，必须在仓库根目录 `edits.md`（已 gitignore）追加一条改动记录（文件、改动内容、影响）。
2. 用户说「推上去」/「提交」时，**只读 `edits.md`**，按其清单提交对应改动，**不要重新通读代码**。
3. 提交成功后**清空 `edits.md` 正文**（保留标题与本说明）。

`edits.md` 模板见文件头部说明，违反此规则视为不规范。

## 其他约定

- 除 `edits.md` 外禁止向仓库内添加与任务无关的文件；不主动创建 README/*.md 文档，除非用户要求。
- 服务端配置：`manifest/config/config.yaml`（监听 `:5020`）；`.gitignore` 已忽略 `**/config/config.yaml`，勿提交本地配置。
