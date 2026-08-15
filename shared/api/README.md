# shared/api — 前后端共享 API 规则

> 本目录是前后端唯一的接口契约文档，位于仓库根（可同时被 Go 后端与前端子模块引用）。
> **改动任何接口时，必须同步更新本文件**，否则视为破坏前后端契约。

## 一、通用约定

- 所有业务接口挂在 `/api/<module>/...` 下，通过 GoFrame `ghttp` 自动注册（controller 方法名 → 蛇形 URL，如 `GetSettingValue` → `/api/setting/get-setting-value`）。
- 响应统一包裹：`{ code, message, data }`（`ghttp.MiddlewareHandlerResponse`）。
  - 前端取数一律 `res.data.data.<字段>`（axios 请求时）。
- 请求参数校验：后端 Req 结构体用 `v:"required"` tag；前端提交前做基本校验。
- **多班级隔离**：后端按客户端 IP 从 `device` 表推导班级（`utility.GetClassID(ctx)`）；未绑定设备的班级接口返回错误，不再回落到班级 `0`。
- 上传接口用 `multipart/form-data`，字段名为 `file`；返回 `{ fileName }`，前端拼 `/resource/upload/<fileName>` 访问。服务端仅公开该上传目录且禁止目录列表，`resource/database` 不可通过 HTTP 访问。
- 当前仍无认证或授权机制；IP 绑定只提供班级隔离，不构成安全身份认证。这是已知限制，待产品明确凭据和权限模型后处理。

## 二、模块与路由总览

| 模块 | 路由前缀 | 接口 | 方法 | 说明 |
| --- | --- | --- | --- | --- |
| device | `/api/device` | `is-bind` | GET | 当前 IP 是否已绑定班级，返回 `isBind`（DeviceTable，未绑定各字段为空） |
| device | `/api/device` | `bind-device` | POST | 绑定：body `{ grade, class }`（整数，grade 1-3、class 1-24），按当前 IP 新增或更新映射，返回 `{ msg }` |
| device | `/api/device` | `un-bind` | GET | 解绑：删除当前 IP 的 device 记录 |
| device | `/api/device` | `get-class` | GET | 返回当前绑定 `{ grade, class }` |
| setting | `/api/setting` | `get-setting-value` | GET | 传 `key`（如 `wall_url`/`notice`/`notice_switch`/`slogan_time`），返回 `Res`（SettingTable） |
| setting | `/api/setting` | `set-setting-value` | POST | body `{ key, value, class? }`，写班级设置 |
| setting | `/api/setting` | `upload-wall-paper` | POST | 上传 jpg/jpeg/png/webp 壁纸；扩展名和文件内容必须匹配，返回 `{ fileName }` |
| schedule | `/api/schedule` | `get-schedule` | GET | 返回 `Res.Schedule`（JSON 字符串，键 `1-7` 为每日课程、`8/9` 为节次起止时间） |
| schedule | `/api/schedule` | `upload-schedule` | POST | 仅接受 `.xlsx` 课程表，解析后按当前班级覆盖写入；临时文件会被删除 |
| count-down | `/api/count-down` | `get-event` | GET | 返回 `Res[]`（当前班级倒计时事件列表） |
| count-down | `/api/count-down` | `add-event` | POST | body `{ event, date, type, during }` 新增事件 |
| count-down | `/api/count-down` | `del-event` | POST | body `{ id }`，仅删除当前绑定班级的事件 |
| proxy | `/api/proxy` | `get-weather` | GET | 传 `key`（`now`/`3d`/`rain`），返回 `Data`（WeatherTable，`Value` 为天气 JSON 字符串） |

## 三、请求示例（curl）

```bash
# 查询设置
curl 'http://localhost:5020/api/setting/get-setting-value?key=notice'

# 绑定设备
curl -X POST 'http://localhost:5020/api/device/bind-device' \
  -H 'Content-Type: application/json' \
  -d '{"grade":1,"class":3}'

# 上传壁纸
curl -X POST 'http://localhost:5020/api/setting/upload-wall-paper' \
  -F 'file=@./bg.jpg'
```

## 四、响应包裹示例

```json
{
  "code": 0,
  "message": "",
  "data": {
    "isBind": { "Id": 1, "ClassID": 3, "DeviceID": "192.168.1.10" }
  }
}
```

## 五、改动规范

1. 后端新增/改字段：同步修改 `api/<module>/v1/*.go`（契约）与 controller/dao 实现。
2. 前端调用处与 `shared/api` 本文件同步更新（字段名、路由、方法）。
3. 改字段名时用 grep 排查前端全部调用处，防止遗漏 `res.data.data.<字段>`。

## 六、天气代理配置

- `QWEATHER_KEY`：必填；未配置时跳过刷新并保持服务运行，仓库不保存天气 API 密钥。
- `QWEATHER_LOCATION`：可选，默认 `119.15,34.81`。
- `QWEATHER_BASE_URL`：可选，默认 `https://devapi.qweather.com`。
- 定时刷新每 3 分钟执行一次，使用带超时和状态检查的 HTTP 客户端；刷新不会重叠，错误只记录而不会终止服务。
