# eGZ-Board 全新开发计划书（Greenfield）

> **重要说明**：本仓库的旧代码已全部删除，从零开始实现。本文档是唯一权威的需求与设计文档，
> 后续所有开发（无论人类还是 AI）都以此为准。文档要求「事无巨细」——所有功能、数据、接口、
> 布局、权限都必须描述到可直接落地实现的程度，不许存在"沿用旧版""保持原样"的模糊表述。

---

## 0. 项目定位

教学信息壁纸软件（电子班牌）。班级大屏（Windows 10 班牌 / 安卓屏）开机即全屏显示信息壁纸，
后台网页进行内容管理。设计目标为**通用化**：可服务多所学校、多年级、多班级、多地区，而非单一学校。

- **展示端**：教室/校门口的大屏，纯展示。
- **管理端**：网页后台，分角色管理内容。

---

## 1. 技术栈（定死，不再讨论）

### 1.1 后端

| 组件 | 选型 | 说明 |
| --- | --- | --- |
| 语言 | Go 1.22+ | |
| 路由 | Gin | HTTP 框架 |
| 依赖注入 | google/wire | 编译期 DI |
| ORM | GORM | + `github.com/glebarez/sqlite`（纯 Go 的 SQLite 驱动，单文件库，无 CGO） |
| 数据库 | SQLite | `resource/database/data.db`，WAL，busy_timeout 5s，foreign_keys on |
| 密码 | bcrypt | 数据库只存哈希 |
| Token | 自定义 opaque | 32 随机字节 base64url；DB 只存 SHA-256；见 §3.5 |
| 定时 | `robfig/cron` 或 `time.Ticker` | 天气刷新 |
| 配置 | 启动参数/环境变量仅承载机密与端口 | **业务配置全部后台化**，见 §9 |

### 1.2 前端

| 组件 | 选型 |
| --- | --- |
| 框架 | Vue 3（`<script setup>`，组合式 API） |
| 构建 | Vite |
| UI | Vuetify 3 |
| 状态 | Pinia |
| 路由 | Vue Router + unplugin-vue-router（文件路由，`web/src/pages/**`） |
| 样式 | UnoCSS |
| HTTP | axios |
| 富文本 | wangEditor（`@wangeditor/editor-for-vue`） |

### 1.3 仓库结构

```
eGZ-Board/
├── plan.md                  # 本文档（唯一权威）
├── cmd/
│   └── server/main.go       # 入口：wire 组装 + 启动
├── internal/
│   ├── config/              # 启动配置（端口、DB 路径、天气 key 等）
│   ├── server/              # Gin 装配、路由注册、中间件挂载
│   ├── handler/             # HTTP 处理器（按模块拆分文件）
│   ├── service/             # 业务逻辑（权限解析、scope 聚合、课程表、天气代理）
│   ├── repository/          # GORM 数据访问（每表一文件）
│   ├── model/               # entity（GORM 模型）+ migrate 函数 + seed
│   ├── auth/                # token 生成/哈希、bcrypt
│   └── middleware/          # 认证中间件、角色/范围校验中间件
├── utility/                 # 通用工具（上下文主体、scope 计算）
├── pkg/                     # 纯工具包（可选）
├── web/                     # 前端（同一仓库，不再用子模块）
│   ├── src/
│   │   ├── api/             # axios 封装 + 各模块 API
│   │   ├── pages/           # 文件路由页面
│   │   ├── components/      # 通用组件
│   │   ├── stores/          # Pinia
│   │   └── ...
│   └── dist/                # 构建产物
├── public/                  # 前端构建产物托管目录（= web/dist 拷贝，server root）
├── resource/                # 运行时：database/、upload/（gitignore）
├── .gitignore
└── edits.md                 # 提交工作流（gitignore，见 §10）
```

构建产物流程：`web/dist` → 拷贝到 `public/`，由 Gin 静态托管；`/resource/upload` 单独静态挂载（关闭目录列表）。

### 1.4 响应包裹约定

所有 `/api` 响应统一为：

```json
{ "code": 0, "message": "", "data": { ... } }
```

- `code=0` 成功；非 0 失败，`message` 为人类可读错误。
- 前端取数统一 `res.data.data.<字段>`。
- 错误码约定：`401` 未认证/过期；`403` 主体无权；`404` 资源不存在；`400` 参数错误。
- 所有路由路径在定义处显式写死（无自动推导）。

---

## 2. 功能需求清单（全部需求）

### 2.1 展示端

1. **`/` 班牌主屏**（Windows 10 班级大屏，**左侧是桌面图标区，红线约束见 §6.2**）
   - 全屏壁纸背景（可滚动，铺满）。
   - 顶部：日期 / 时间 / 星期 / 周数条。
   - 课程表条：显示"今天"的课程，当前节次高亮，前后节次弱化。
   - 公告栏：右侧列表（标题+来源+置顶标签），点击展开详情（右侧抽屉，**不得遮挡左侧**）。
   - 倒计时 / 正计时卡片。
   - 天气卡片（当前天气 + 3 天预报 + 降雨）。
   - 进度条（今日已过 %）。
   - Slogan（hitokoto 一言 + 出处）。
   - **左侧图标区留白，宽度后台可配，默认 0 不占用左侧。**

2. **`/outside` 室外屏**
   - 全屏壁纸 + 课程表条 + 倒计时列表 + 公告 + 时间/进度 + 天气，右对齐布局。
   - **无左侧遮挡约束**，公告详情可用居中弹窗。

3. **`/pair` 设备配对页**
   - 输入 8 位配对码 + 设备名称，成功保存设备 token，跳转 `/` 或 `/outside`。

### 2.2 内容模块

4. **公告栏**：富文本公告，支持发布到全校/某年级/选择班级，置顶、上下线、定时生效。
5. **倒计时**：事件 + 目标日期，显示"距 X 还有 N 天"。
6. **正计时**：事件 + 起始日期，显示"X 已进行 N 天"。
7. **课程表**：作息时间配置（多份）+ 课程表（多份、周模式/单双日模式、可增删改、可设使用、上游下发）。
8. **天气**：和风天气数据，展示端定时刷新。
9. **壁纸**：链接或上传图片。
10. **Slogan**：hitokoto，刷新秒数可配。

### 2.3 后台管理

11. 登录 / 退出；角色感知导航与班级选择器。
12. 机构管理：学校、年级、班级 CRUD；**班级不默认创建**，支持批量创建。
13. 账户管理：创建子账户、禁用/启用、改密码、删除；账户↔班级双向绑定。
14. 公告管理：新建/编辑/删除/置顶/上下线/定时；发布范围选择。
15. 倒计时 / 正计时管理：增删改 + 分发范围。
16. 课程表管理：作息配置 CRUD、课程表在线编辑（周/单双日）、设使用、复制、删除、下发、查看下游使用情况。
17. 设备管理：按班级生成/删除配对码；设备列表/改名/调班/撤销。
18. 基础设置：壁纸、slogan 秒数、公告开关、**左侧图标区宽度**。
19. 系统配置（仅 superadmin）：天气地区、默认壁纸、默认作息、默认图标区宽度、默认 slogan 秒数。
20. 审计日志：所有变更操作留痕（谁、何时、做了什么）。

### 2.4 权限

21. 四级角色 + 设备主体；高权限可授权低权限；高权限发布的内容低权限可见不可管。
22. 内容 scope 覆盖模型（all / school / grade / classes）闭环。

---

## 3. 权限体系（核心，实现必须闭环）

### 3.1 角色与层级

| 角色 `role` | 层级 `level` | 管辖范围 | 可授权给 |
| --- | --- | --- | --- |
| `superadmin` 超级管理员 | 0 | 全部学校/班级 | school、grade、class |
| `school` 学校管理员 | 1 | 一个/多个学校 | 其校内 grade、class |
| `grade` 年级管理员 | 2 | 一个年级（多班） | 其年级内 class |
| `class` 班级管理员 | 3 | 一个班级 | 无 |
| `device` 设备主体 | - | 绑定班级只读 | 无 |

### 3.2 账户绑定模型（双向）

- 一张绑定表：`user_bindings(user_id, target_type, target_id)`，其中 `target_type ∈ {school, grade, class}`。
- 「账户分配班级」= 为一个账户添加多个 `class` 绑定。
- 「每个班级分配子账户」= 一个 `class` 可被多个账户绑定（也即班级管理员）。
- **管辖集合**：一个账户有效管辖的班级 = 其所有绑定展开到班级的并集：
  - 绑 `school` → 该校全部班级；
  - 绑 `grade` → 该年级全部班级；
  - 绑 `class` → 仅该班。
- 唯一索引：`(user_id, target_type, target_id)`；删除账户时级联删除绑定，不影响班级数据。

### 3.3 授权规则

- `level` 更小 = 权限更高。授权只能**向下**：
  - superadmin 可创建 school/grade/class 账户；
  - school 可创建其管辖范围内的 grade/class 账户；
  - grade 可创建其年级内的 class 账户；
  - class 不能创建账户。
- 创建账户时必须指定绑定范围，且**范围必须落在授权者管辖集合内**，否则 403。
- 账户可禁用/启用、重置密码、删除；删除即失去所有绑定。

### 3.4 内容可见性与可管理性

每条「内容」（公告 / 倒计时 / 正计时 / 课程表 / 作息）都带发布信息：

- `publisher_user_id`：发布者账户。
- `publisher_level`：发布时快照的 level（用于"低权限不可管理高权限内容"判定）。
- `scope_type`：`all` | `school` | `grade` | `classes`。
- `scope_targets`：JSON 数组。
  - `all`：空数组 = 发布者管辖范围的全部班级（由发布者 level 动态展开）；
  - `school`：`[schoolId, ...]`；
  - `grade`：`[gradeId, ...]`；
  - `classes`：`[classId, ...]`。
- 服务端提供 `ResolveScopeClassIDs(scope_type, targets, publisher)`：把 scope 展开为**目标班级 id 集合**。

**读取（设备展示）**：`device.class_id` 在目标班级集合中即显示。天然实现"高权限发布 → 低权限可见"。

**管理（后台）**，对某条内容 `c` 和操作者 `u`：
- `c.publisher_user_id == u.id` → 可改可删；
- `u.level < c.publisher_level` 且 `c` 的 scope 全部落在 `u` 管辖集合内 → 可改可删（高管低）；
- 否则 → 只读（列表可见，前端标"只读"，隐藏编辑/删除按钮；后端同样拒绝写操作）。

### 3.5 认证与 Token

- 密码：bcrypt。不存在明文。
- Token：32 随机字节，无填充 base64url（`-` `_` 无 `=`）。DB 只存 `sha256(token)`。
- 两种主体：
  - **管理员会话**：`admin_sessions`，12 小时过期，可撤销，登录/退出。
  - **设备**：`devices`，无过期，配对时签发，`unpair` 或管理员撤销后失效。
- 中间件每请求解析一次 Bearer：先查管理员会话，再查设备，写入上下文主体；无效返回 401，主体能力不匹配返回 403。
- 配对码：8 位大写无歧义字母数字（去除 I/O/0/1），按班级生成，**无失效日期**、一用、可硬删除；DB 只存哈希。

---

## 4. 数据模型（字段级，GORM）

> 全部表主键为 `Id`（`gorm:"primary_key;auto_increment;not null"`）。时间用 UTC 存储，展示层转本地。

### 4.1 机构

**schools**
| 字段 | 类型 | 约束 |
| --- | --- | --- |
| Id | int | PK |
| Name | string | not null（学校名） |
| Region | string | 可空（地区名，用于泛化） |
| CreatedAt / UpdatedAt | time | |

**grades**
| 字段 | 类型 | 约束 |
| --- | --- | --- |
| Id | int | PK |
| SchoolId | int | not null，index |
| GradeNo | int | not null（1/2/3…） |
| Name | string | not null（高一/高二/高三…） |
| Sort | int | 排序 |

**classes**
| 字段 | 类型 | 约束 |
| --- | --- | --- |
| Id | int | PK |
| SchoolId | int | not null，index |
| GradeId | int | not null，index |
| GradeNo | int | not null（冗余，便于展示） |
| ClassNo | int | not null（班号） |
| Name | string | not null（班级名，如"1班"） |
| IconInset | int | 左侧图标区宽度（px），0 = 不占左侧 |

唯一索引：`(SchoolId, GradeNo, ClassNo)`。

### 4.2 账户

**users**
| 字段 | 类型 | 约束 |
| --- | --- | --- |
| Id | int | PK |
| Username | string | not null，unique |
| PasswordHash | string | not null |
| DisplayName | string | 可空 |
| Role | string | not null（superadmin/school/grade/class） |
| Enabled | bool | default true |
| LastLoginAt | *time | |
| CreatedAt / UpdatedAt | time | |

**user_bindings**
| 字段 | 类型 | 约束 |
| --- | --- | --- |
| Id | int | PK |
| UserId | int | not null，index |
| TargetType | string | not null（school/grade/class） |
| TargetId | int | not null |

唯一索引：`(UserId, TargetType, TargetId)`。

**admin_sessions**
| 字段 | 类型 | 约束 |
| --- | --- | --- |
| Id | int | PK |
| UserId | int | not null，index |
| TokenHash | string | not null，unique |
| ExpiresAt | time | not null |
| CreatedAt | time | not null |
| RevokedAt | *time | index |

**audit_logs**
| 字段 | 类型 | 约束 |
| --- | --- | --- |
| Id | int | PK |
| UserId | int | not null |
| Username | string | not null |
| Action | string | not null（如 announcement.create） |
| Detail | string | 可空（JSON 摘要） |
| CreatedAt | time | index |

### 4.3 设备与配对

**devices**
| 字段 | 类型 | 约束 |
| --- | --- | --- |
| Id | int | PK |
| ClassId | int | not null，index |
| Name | string | not null |
| TokenHash | string | not null，unique |
| LastIp | string | not null（审计字段） |
| LastSeenAt | *time | |
| CreatedAt / UpdatedAt | time | |
| RevokedAt | *time | index |

**pairing_codes**
| 字段 | 类型 | 约束 |
| --- | --- | --- |
| Id | int | PK |
| ClassId | int | not null，index |
| CodeHash | string | not null，unique |
| CreatedAt | time | not null |
| CreatedBy | int | not null |
| RedeemedAt | *time | |

### 4.4 内容

**announcements**
| 字段 | 类型 | 约束 |
| --- | --- | --- |
| Id | int | PK |
| SchoolId | int | not null，index |
| Title | string | not null |
| ContentHtml | text | not null |
| PublisherUserId | int | not null |
| PublisherLevel | int | not null |
| ScopeType | string | not null |
| ScopeTargets | text | not null（JSON） |
| Sticky | bool | 置顶 |
| Enabled | bool | 上下线 |
| PublishStartAt | *time | 定时生效 |
| PublishEndAt | *time | 定时下线 |
| CreatedAt / UpdatedAt | time | |
| DeletedAt | *time | 软删除，index |

**countdowns**
| 字段 | 类型 | 约束 |
| --- | --- | --- |
| Id | int | PK |
| SchoolId | int | not null，index |
| Title | string | not null |
| Type | string | not null（`count_down` 倒计时 / `count_up` 正计时） |
| TargetDate | string | not null（YYYY-MM-DD；倒计时=目标日，正计时=起始日） |
| PublisherUserId | int | not null |
| PublisherLevel | int | not null |
| ScopeType | string | not null |
| ScopeTargets | text | not null（JSON） |
| CreatedAt / UpdatedAt | time | |
| DeletedAt | *time | 软删除 |

### 4.5 课程表

**time_tables**（作息时间配置，可多份）
| 字段 | 类型 | 约束 |
| --- | --- | --- |
| Id | int | PK |
| SchoolId | int | not null，index |
| Name | string | not null（如"夏季作息"） |
| Periods | text | not null（JSON：`[{index,label,start,end}]`，`start/end` 为 "HH:MM"） |
| CreatedBy | int | not null |
| CreatedAt / UpdatedAt | time | |
| DeletedAt | *time | 软删除 |

**schedules**（课程表，多份）
| 字段 | 类型 | 约束 |
| --- | --- | --- |
| Id | int | PK |
| SchoolId | int | not null，index |
| ClassId | int | not null，index |
| Name | string | not null（如"2026 春季"） |
| Mode | string | not null（`weekly` 周模式 / `odd_even` 单双日模式） |
| Content | text | not null（JSON，格式见 §7） |
| IsActive | bool | 当前使用，每班级至多一份 true |
| CreatedBy | int | not null |
| CreatedAt / UpdatedAt | time | |
| DeletedAt | *time | 软删除 |

### 4.6 设置

**system_settings**（全局，仅 superadmin 可写）
| 字段 | 类型 | 约束 |
| --- | --- | --- |
| Id | int | PK |
| Key | string | not null，unique |
| Value | text | not null |

预设 Key：`weather.location`（天气地区，如 "119.15,34.81"）、`weather.key`（和风 key，不回显）、`default.wall_url`、`default.icon_inset`、`default.slogan_seconds`、`default.time_table_id`、`default.wallpaper_id`。

**class_settings**（班级级）
| 字段 | 类型 | 约束 |
| --- | --- | --- |
| Id | int | PK |
| ClassId | int | not null，index |
| Key | string | not null |
| Value | text | not null |

唯一索引：`(ClassId, Key)`。预设 Key：`wall_url`、`slogan_seconds`、`notice_switch`、`icon_inset`、`time_table_id`（使用中的作息）、`slogan_time`。

**weather_cache**
| 字段 | 类型 | 约束 |
| --- | --- | --- |
| Id | int | PK |
| Key | string | not null，unique（now/3d/rain） |
| Value | text | not null（和风原始 JSON） |
| UpdatedAt | time | |

### 4.7 迁移与种子

- 启动时 `AutoMigrate` 全部表。
- **新装不创建任何班级**。仅当 `users` 为空时创建 `superadmin`：
  - 密码取 `ADMIN_INITIAL_PASSWORD` 环境变量；未设置则随机生成并打印一次到日志。
- 提供 API 由 superadmin 建立首个学校/年级/班级（也可批量建班）。

---

## 5. API 设计（完整清单）

> 全部前缀 `/api`，响应包裹见 §1.4。`{id}` 路径参数。权限列表示最低要求。

### 5.1 认证 `/api/auth`

| 权限 | 方法 | 路径 | 请求 | 响应 data |
| --- | --- | --- | --- | --- |
| public | POST | `/login` | `{username,password}` | `{token, expiresAt, user:{id,username,displayName,role}}` |
| 管理员 | GET | `/me` | - | `{user, bindings:[{targetType,targetId}], managedClasses:[{id,gradeNo,classNo,name}]}` |
| 管理员 | POST | `/logout` | - | `{}` |

### 5.2 机构 `/api/org`

| 权限 | 方法 | 路径 | 请求 | 响应 data |
| --- | --- | --- | --- | --- |
| superadmin | GET | `/schools` | - | `{schools:[...]}` |
| superadmin | POST | `/schools` | `{name,region?}` | `{school}` |
| superadmin | PUT | `/schools/{id}` | `{name?,region?}` | `{}` |
| superadmin | DELETE | `/schools/{id}` | - | `{}` |
| super/school | GET | `/grades?schoolId=` | - | `{grades:[...]}` |
| super/school | POST | `/grades` | `{schoolId,gradeNo,name,sort?}` | `{grade}` |
| super/school | PUT | `/grades/{id}` | `{name?,gradeNo?,sort?}` | `{}` |
| super/school | DELETE | `/grades/{id}` | - | `{}` |
| 任一管理员 | GET | `/classes?schoolId=&gradeId=` | - | `{classes:[...]}`（仅返回权限范围内） |
| super/school/grade | POST | `/classes/batch` | `{schoolId,gradeNo,classNoFrom,classNoTo}` | `{created:[...]}` |
| super/school | PUT | `/classes/{id}` | `{name?,iconInset?}` | `{}` |
| super/school | DELETE | `/classes/{id}` | - | `{}` |

说明：`grade` 管理员批量建班仅限其年级；`school` 限其学校。班级不默认创建。

### 5.3 账户 `/api/accounts`

| 权限 | 方法 | 路径 | 请求 | 响应 data |
| --- | --- | --- | --- | --- |
| super/school/grade | GET | `/accounts` | `?role=` | `{accounts:[{id,username,displayName,role,enabled,bindings}]}` |
| super/school/grade | POST | `/accounts` | `{username,password,displayName?,role,bindings:[{targetType,targetId}]}` | `{account}` |
| super/school/grade | PUT | `/accounts/{id}` | `{displayName?,enabled?,bindings?}` | `{}` |
| super/school/grade | POST | `/accounts/{id}/reset-password` | `{password}` | `{}` |
| super/school/grade | DELETE | `/accounts/{id}` | - | `{}` |

校验：目标角色层级必须低于操作者；绑定范围必须落入操作者管辖集合；`bindings` 至少一个。

### 5.4 公告 `/api/announcements`

| 权限 | 方法 | 路径 | 请求 | 响应 data |
| --- | --- | --- | --- | --- |
| 管理员 | GET | `/announcements` | `?scope=own&classId=` | `{announcements:[...]}` |
| 管理员 | POST | `/announcements` | `{title,contentHtml,scopeType,scopeTargets,sticky?,enabled?,publishStartAt?,publishEndAt?}` | `{announcement}` |
| 管理员 | PUT | `/announcements/{id}` | 同上可改字段 | `{}` |
| 管理员 | DELETE | `/announcements/{id}` | - | `{}` |
| 管理员 | POST | `/announcements/{id}/toggle` | `{enabled}` | `{}` |
| 设备 | GET | `/api/announcement/list` | - | `{announcements:[{id,title,summary,sticky,source}]}` |
| 设备 | GET | `/api/announcement/detail/{id}` | - | `{id,title,contentHtml}` |

写入前校验可管理性（§3.4）。`/list` 仅返回"设备所属班级"范围内、已上线且在生效期内的公告。

### 5.5 倒计时 / 正计时 `/api/countdowns`

| 权限 | 方法 | 路径 | 请求 | 响应 data |
| --- | --- | --- | --- | --- |
| 管理员 | GET | `/countdowns` | `?classId=` | `{countdowns:[...]}` |
| 管理员 | POST | `/countdowns` | `{title,type,targetDate,scopeType,scopeTargets}` | `{countdown}` |
| 管理员 | PUT | `/countdowns/{id}` | 可改字段 | `{}` |
| 管理员 | DELETE | `/countdowns/{id}` | - | `{}` |
| 设备 | GET | `/api/count-down/get-event` | - | `{events:[{title,type,targetDate,days}]}` |

`days` 由后端计算：倒计时 = 目标日 - 今天（正数）；正计时 = 今天 - 起始日（非负）。设备接口保持旧路径 `/api/count-down/get-event` 兼容。

### 5.6 课程表 `/api/schedules` 与 `/api/time-tables`

| 权限 | 方法 | 路径 | 请求 | 响应 data |
| --- | --- | --- | --- | --- |
| 管理员 | GET | `/time-tables` | `?schoolId=` | `{timeTables:[...]}` |
| 管理员 | POST | `/time-tables` | `{name,periods:[{index,label,start,end}]}` | `{timeTable}` |
| 管理员 | PUT | `/time-tables/{id}` | 可改字段 | `{}` |
| 管理员 | DELETE | `/time-tables/{id}` | - | `{}` |
| 管理员 | POST | `/time-tables/{id}/push` | `{classIds:[]}` | `{pushed: n}`（下发作息：写目标班级 `class_settings.time_table_id`） |
| 管理员 | GET | `/schedules?classId=` | - | `{schedules:[...]}`（含 active 标记） |
| 管理员 | POST | `/schedules` | `{classId,name,mode,content}` | `{schedule}` |
| 管理员 | PUT | `/schedules/{id}` | 可改字段 | `{}` |
| 管理员 | DELETE | `/schedules/{id}` | - | `{}` |
| 管理员 | POST | `/schedules/{id}/activate` | - | `{}`（设为使用，同班其他置 false） |
| 管理员 | POST | `/schedules/{id}/push` | `{classIds:[]}` | `{pushed: n}`（下发副本到目标班级并设为使用） |
| 管理员 | GET | `/schedules/status?gradeId=` | - | `{rows:[{classId,className,timeTable:{id,name,source},schedule:{id,name,mode,source}}]}` |
| 设备 | GET | `/api/schedule/get-schedule` | - | `{mode, content, periods}` |

`source` 标记 `own`（班级自配）/ `pushed`（上游下发）/ `default`（系统默认），供上游查看下游配置情况。

### 5.7 设备 `/api/device` 与 `/api/admin/devices`

| 权限 | 方法 | 路径 | 请求 | 响应 data |
| --- | --- | --- | --- | --- |
| public | POST | `/api/device/pair` | `{code,name}` | `{token, device:{id,classId,name}, class:{id,gradeNo,classNo,name,schoolName}}` |
| 设备 | GET | `/api/device/me` | - | `{device, class}` |
| 设备 | POST | `/api/device/unpair` | - | `{}` |
| 管理员 | GET | `/api/admin/pairing-codes` | `?classId=` | `{pairingCodes:[...]}` |
| 管理员 | POST | `/api/admin/pairing-codes` | `{classId}` | `{code, pairingCode}`（code 仅返回一次） |
| 管理员 | POST | `/api/admin/delete-pairing-code` | `{id}` | `{}`（硬删除） |
| 管理员 | GET | `/api/admin/devices` | `?classId=&gradeId=` | `{devices:[...]}` |
| 管理员 | POST | `/api/admin/update-device` | `{id,classId,name}` | `{}` |
| 管理员 | POST | `/api/admin/revoke-device` | `{id}` | `{}` |

### 5.8 设置 `/api/setting` 与 `/api/admin/system-settings`

| 权限 | 方法 | 路径 | 请求 | 响应 data |
| --- | --- | --- | --- | --- |
| 设备/管理员 | GET | `/api/setting/get-setting-value` | `key`；管理员另传 `classId` | `{value}` |
| 管理员 | POST | `/api/setting/set-setting-value` | `{key,value,classId}` | `{}` |
| 管理员 | POST | `/api/setting/upload-wall-paper` | multipart `file` | `{fileName}` |
| superadmin | GET | `/api/admin/system-settings` | - | `{settings:[{key,value}]}` |
| superadmin | PUT | `/api/admin/system-settings` | `{settings:[{key,value}]}` | `{}` |

### 5.9 天气 `/api/proxy`

| 权限 | 方法 | 路径 | 请求 | 响应 data |
| --- | --- | --- | --- | --- |
| public | GET | `/api/proxy/get-weather` | `key=now\|3d\|rain` | `{data}` |

服务端每 3 分钟拉取和风并写 `weather_cache`；无 key 时返回空且不报错。

---

## 6. 前端设计

### 6.1 路由（文件路由）

| 路径 | 页面 | 守卫 |
| --- | --- | --- |
| `/` | 班牌主屏（board） | 需设备 token |
| `/outside` | 室外屏 | 需设备 token |
| `/pair` | 设备配对 | public |
| `/admin/login` | 后台登录 | public |
| `/admin` | 后台壳（重定向到首个标签） | 需管理员 |
| `/admin/announcements` | 公告管理 | 管理员 |
| `/admin/countdowns` | 倒计时/正计时 | 管理员 |
| `/admin/schedules` | 课程表管理 | 管理员 |
| `/admin/devices` | 设备管理 | 管理员 |
| `/admin/settings` | 基础设置 | 管理员 |
| `/admin/org` | 机构管理 | super/school |
| `/admin/accounts` | 账户管理 | super/school/grade |
| `/admin/system` | 系统配置 | superadmin |
| `/admin/audit` | 审计日志 | superadmin |

路由守卫逻辑：无 token → `/pair` 或 `/admin/login`；token 校验失败同样跳转；后台按角色控制可访问菜单（隐藏无权限入口，后端兜底拒绝）。

### 6.2 `/` 班牌主屏布局（红线）

```
┌──────────────────────────────────────────────┐
│  [日期][时间][星期][周数]        (顶部时间条)   │
│  [ 课 程 表 条  （今天课程 + 当前节高亮）  ]     │
│ ┌──────┐ ┌─────────────────────────────────┐ │
│ │ 左侧  │ │  公告栏（列表，点击右侧抽屉详情）   │ │
│ │ 留白  │ │  倒计时/正计时卡片               │ │
│ │ 图标区│ │  天气卡片                       │ │
│ │ 不遮挡│ │  进度条 + Slogan                │ │
│ └──────┘ └─────────────────────────────────┘ │
└──────────────────────────────────────────────┘
```

- 整屏壁纸背景。
- 顶部时间条与课程表条为**全宽**，跨整个屏（含左侧图标区上方）——允许，因为桌面图标在下方。
- 下方内容：**左侧固定留白 = 图标区**，宽度取 `class_settings.icon_inset`（px），默认 0（不留白，内容占满）。所有内容组件（公告/倒计时/天气/进度/slogan）都放在右侧容器内。
- **公告详情为右侧抽屉**（`v-navigation-drawer location="right"` 或底部滑出），**绝不从左侧弹出、绝不覆盖左侧**。
- 组件结构建议：
  - `web/src/pages/index.vue`（board 容器，读 `icon_inset`，计算内容区 `margin-left`）
  - `components/board/TopBar.vue`（日期/时间/星期/周数）
  - `components/board/ScheduleBar.vue`（课程表条）
  - `components/board/AnnouncementBoard.vue`（公告列表 + 右侧详情抽屉）
  - `components/board/CountdownCards.vue`（倒计时/正计时）
  - `components/board/WeatherCard.vue`
  - `components/board/ProgressCard.vue`（今日已过 %）
  - `components/board/SloganCard.vue`
- 数据：`usePolling` 组合式（每 60s 刷新内容，天气 60s，时间每秒）。

### 6.3 `/outside` 布局

全屏壁纸 + 顶部课程表条 + 右侧 `ToolBox`（倒计时列表、公告、时间/进度、天气），右对齐，无左侧约束，公告详情居中弹窗。

### 6.4 后台壳

- 顶部栏：登录账户 + 角色标签 + 班级选择器（`managedClasses`）+ 退出。
- 左侧导航：按角色渲染（见 §6.1 路由表权限列）。
- 各管理页统一使用：表格（分页/搜索）、表单弹窗、确认框、Snackbar；只读内容显示"只读（来源：XX管理员）"。
- 班级选择器仅在需要指定班级的操作（如基础设置、设备、课程表）使用，其他页面按账号管辖范围聚合。

### 6.5 页面组件清单

- 公告管理：`AnnouncementManager.vue`（列表 + 富文本编辑弹窗 + 发布范围选择 + 置顶/上下线/定时）
- 倒计时管理：`CountdownManager.vue`（列表 + 新建弹窗，含 type 选择、日期、发布范围）
- 课程表管理：`ScheduleManager.vue`（作息配置 Tab + 课程表 Tab）
  - 作息配置：`TimeTableEditor.vue`（节次列表：index/label/start/end 增删）
  - 课程表：`ScheduleGrid.vue`（周模式：7 列 × 节次行；单双日模式：单/双两列 × 节次行；单元格文本输入）+ 多份列表（设使用/复制/删除/下发）
  - 下发：`PushDialog.vue`（选目标班级集合）
  - 下游状态：`ScheduleStatusTable.vue`
- 设备管理：`DevicesManager.vue`（配对码生成/列表/删除 + 设备列表/改名/调班/撤销）
- 机构管理：`OrgManager.vue`（学校/年级/班级树 + 批量建班对话框）
- 账户管理：`AccountManager.vue`（账户列表 + 创建/编辑弹窗 + 绑定范围选择）
- 基础设置：`BasicSettings.vue`（壁纸、slogan 秒数、公告开关、图标区宽度）
- 系统配置：`SystemSettings.vue`
- 审计日志：`AuditLog.vue`

---

## 7. 课程表模块详细设计（重点）

### 7.1 作息时间配置（time_table）

- `periods = [{index:1, label:"第一节", start:"08:00", end:"08:45"}, ...]`。
- 一个学校可有多份；班级通过 `class_settings.time_table_id` 指定使用哪份（或系统默认）。
- 下游班级管理员可「复制上游配置后修改」成专属配置；上游下发 = 把配置 id 写入目标班级 `time_table_id`（可加 `source` 字段标识下发来源）。
- 展示端课程表条的节次 = 该班级使用中的作息配置。

### 7.2 课程表内容格式

**weekly 周模式**（兼容经典 JSON）：
```json
{
  "1": ["语文","数学","英语"],
  "2": ["数学","物理","语文"],
  "3": [], "4": [], "5": [], "6": [], "7": [],
  "8":  ["08:00","08:50","09:40"],
  "9":  ["08:45","09:35","10:25"]
}
```
- key `1`~`7` 对应周一~周日；每个 key 的值数组长度 = 节次数。
- `8` = 每节开始时间，`9` = 每节结束时间（由作息配置节次生成，可留空让前端取作息配置）。

**odd_even 单双日模式**：
```json
{
  "odd":  ["语文","数学","英语"],
  "even": ["数学","物理","语文"],
  "8":  ["08:00","08:50","09:40"],
  "9":  ["08:45","09:35","10:25"]
}
```
- 单双日按**自然日奇偶**判断：日期号（day of month）为奇数用 `odd`，偶数用 `even`。
- 也可在后台加开关「按周次奇偶」（周一为单周起点）——实现按需求取舍，默认按自然日奇偶。

### 7.3 展示端渲染

- 设备 GET `/api/schedule/get-schedule` → `{mode, content, periods}`，服务端已按当天解析好：
  - weekly：返回今天的课程数组 + 节次；
  - odd_even：按今天奇偶返回对应数组 + 节次。
- 前端根据当前时间与节次区间高亮当前节（沿用"当前节高亮、前后弱化"样式）。
- 无课程表/无作息 → 显示"暂无课程数据"。

### 7.4 上游下发与查看

- 上游（school/grade）创建作息或课程表 → `push` 到目标班级集合：写入/新建副本并设为使用。
- `GET /schedules/status?gradeId=` 返回该年级每班当前使用的作息与课程表，`source` 标记 own/pushed/default，供上游审查下游配置。

---

## 8. 非功能需求

- **安全**：token/密码/配对码/天气 key 永不进 DTO 与响应；HTML 内容前端渲染前 `DOMPurify` 消毒；文件上传白名单扩展名与 MIME；上传目录关闭目录列表；SQLite 文件不通过 HTTP 暴露。
- **性能**：展示端轮询聚合（60s 一次，天气 60s）；倒计时天数为后端计算避免时区误差；静态资源走 Gin 静态托管 + 浏览器缓存。
- **可运维**：启动日志输出端口、数据库路径；`/swagger` 或 `/api.json`（可选，Gin 可接 swag 或手工维护接口表）；审计日志记录关键操作。
- **兼容**：设备端旧路径 `/api/count-down/get-event`、`/api/schedule/get-schedule`、`/api/setting/get-setting-value`、`/api/device/{pair,me,unpair}` 保持路径与响应结构，减少展示端迁移成本。

---

## 9. 配置后台化（泛化支撑）

- 业务配置全部走 `system_settings`（全局，superadmin 管理）+ `class_settings`（班级级）：
  - 天气地区、天气 key、默认壁纸、默认 icon_inset、默认 slogan 秒数、默认作息。
- 环境变量仅保留：`ADMIN_INITIAL_PASSWORD`（首启）、天气 key（可选，也可后台录入且不回显）、监听端口、DB 路径。
- 多学校/多地区 = 机构树支撑 + 每个学校独立配置，不做区域硬编码。

---

## 10. 开发与提交约定

1. **edits.md 驱动提交**：每次完成一处改动，在仓库根目录 `edits.md`（已 gitignore）追加记录（文件、改动内容、影响）。
2. 用户说「推上去/提交」时：只读 `edits.md`，按清单提交对应改动，**不要重新通读代码**；提交成功后清空 `edits.md` 正文（保留标题与本说明）。
3. 除 `edits.md` 外不向仓库添加与任务无关的文件；不主动创建 README/*.md，除非用户要求。
4. 改动接口必须同步更新本文档（如果涉及 API/数据结构）。
5. 提交信息简洁，按仓库风格。

---

## 11. 实施阶段

- **P0 地基**：Go 项目骨架（gin + wire + gorm + sqlite）、配置加载、响应包裹、认证中间件、token/bcrypt、迁移与 superadmin 种子、健康检查。
- **P1 权限**：机构（school/grade/class）CRUD + 批量建班、账户 CRUD + 绑定、scope 解析、可管理性校验、/auth/me、后台壳。
- **P2 内容基础**：公告栏（CRUD + 分发 + 设备列表/详情 + board 右侧公告组件）、倒计时/正计时（CRUD + 分发 + 设备接口 + 展示卡片）。
- **P3 课程表**：作息配置 + 课程表在线编辑（周/单双日）+ 多份管理 + 使用/下发/状态查看 + 展示端渲染。
- **P4 设备与设置**：配对码、设备管理、基础设置、系统配置、天气代理。
- **P5 打磨**：审计日志、UI 完善、`/outside`、`/pair`、测试、文档同步。

---

## 12. 验收红线（实现时必须满足）

1. `/` 左侧图标区**绝不遮挡**，`icon_inset` 可配、默认 0。
2. 新装不自动建班；批量建班可用。
3. 账户双向绑定；授权只能向下；越权操作后端拒绝（403）。
4. 高权限发布的内容低权限可见但**不可管理**（前端隐藏 + 后端拦截）。
5. 课程表：多份配置、可增删改、可设使用、周与单双日两模式、上游下发、下游可覆盖、上游可查状态。
6. 公告栏列表 + 详情；倒计时与正计时都支持。
7. 旧设备端接口路径保持兼容。
8. 机密（密码/token/配对码/天气 key）不进入响应与日志。
