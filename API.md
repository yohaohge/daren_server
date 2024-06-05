# API文档
## 前提说明

### 返回消息结构（json)

| 字段名    | 类型 | 说明 |
| --------- | --------- | --------- |
|code| int | 错误码，详细定义见下|
|msg| string| 错误原因说明|
|data| json object|返回的具体数据|
- ps: 以下api的返回参数说明，对code及msg不再特别说明，只针对data展示说明，如果是空的就是data为{}

### 认证
除了登录本身不需要认证，其他的都需要认证，登录后会得到一个session，以后的每一个请求都需要将此sid（session）和uid（玩家id）传过来

### 错误码定义
| 错误码    | 说明 |
| --------- | ------- |
| 0 | 没有错误 |
| 1 | 参数错误 |
| 2 | 失败 |
| 3 | 需要登陆 |
| 4 | 签名错误 |
| 5 | 用户不存在 |
| 6 | 用户已经存在 |
| 7 | 密码错误 |
|其他||
- 后续可以定义在proto


## 登录
- url: /api/login
- method: POST
- 输入参数：
  
| 字段名    | 类型 | 说明 |
| --------- | --------- | --------- |
|open_id| string | 唯一id|
|password| string| 密码|

- 返回参数：

| 字段名    | 类型 | 说明 |
| --------- | --------- | --------- |
|uid| string | 用户id|
|name| string | 名字|
|vipTime| string | vip过期时间|
|createTime| string | 注册时间|
|sid|string| 会话id |


## 注册
- url: /api/regster
- method: POST
- 输入参数：
  
| 字段名    | 类型 | 说明 |
| --------- | --------- | --------- |
|open_id| string | 唯一id|
|password| string| 密码|

- 返回参数：

| 字段名    | 类型 | 说明 |
| --------- | --------- | --------- |


## 主页介绍等 (开发中)
- url: /api/index
- method: GET
- 输入参数：
  
| 字段名    | 类型 | 说明 |
| --------- | --------- | --------- |

- 返回参数：

| 字段名    | 类型 | 说明 |
| --------- | --------- | --------- |

## 获取配置
- url: /api/config
- method: GET
- 输入参数：

| 字段名    | 类型 | 说明 |
| --------- | --------- | --------- |

- 返回参数：

| 字段名    | 类型     | 说明     |
| --------- |--------|--------|
| pay_config      | obj    | 支付     |
| banner_config    | obj    | banner |
| notice    | string | 公告     |

## 视频列表信息
- url: /api/v_list
- method: GET
- 输入参数：
- 
| 字段名    | 类型  | 说明 |
|-----| --------- | --------- |
|page_no| 页码  | 唯一id|
|page_count| 数量  | 密码|

- 返回参数：

| 字段名     | 类型     | 说明    |
|---------| --------- |-------|
| id      | string | 视频标识  |
| name    | string | 视频名字  |
| data    | string | unuse |
| total   | string | 列表总数  |
| desc    | string | 描述    |
| label   | string | 标签    |


## 视频详细信息
- url: /api/v_detail
- method: POST
- 输入参数：
-
| 字段名        | 类型  | 说明   |
|------------|-----|------|
| id         | int | 视频标识 |

- 返回参数：

| 字段名     | 类型     | 说明      |
|---------|--------|---------|
| id      | string | 视频标识    |
| name    | string | 视频名字    |
| data    | string | unuse   |
| total   | string | 列表总数    |
| desc    | string | 描述      |
| label   | string | 标签      |
| episodes   | obj    | 剧每集信息   |
| w_min   | string | 用户可以看的集 |
| w_max   | string | 用户可以看的集      |


## 视频某集数据
- url: /api/v_episode
- method: POST
- 输入参数：
-
| 字段名   | 类型  | 说明   |
|-------|-----|------|
| id    | int | 视频标识 |
| index | int | 第几集  |

- 返回参数：

| 字段名     | 类型     | 说明   |
|---------| --------- |------|
| id      | string | 视频标识 |
| name    | string | 视频名字 |
| data    | string | url  |

## 搜索视频列表信息
- url: /api/search
- method: GET
- 输入参数：

| 字段名        | 类型     | 说明 |
|---------| --------- |------|
| name       | string | 名称 |
| page_no    | string | 页码 |
| page_count | string | 数量 |

- 返回参数：

| 字段名     | 类型     | 说明    |
|---------| --------- |-------|
| id      | string | 视频标识  |
| name    | string | 视频名字  |
| data    | string | unuse |
| total   | string | 列表总数  |
| desc    | string | 描述    |
| label   | string | 标签    |

## 搜索视频列表信息
- url: /user/refresh
- method: POST
- 输入参数：
-
| 字段名        | 类型     | 说明   |
|------------|--------|------|
| uid        | int    | uid  |
| sid        | string | 会话id |

- 返回参数：

| 字段名     | 类型     | 说明    |
|---------| --------- |-------|

## 一些配置信息
- url: /api/config
- method: GET
- 输入参数：
-
| 字段名        | 类型     | 说明   |
|------------|--------|------|

- 返回参数：

| 字段名        | 类型     | 说明   |
|------------| --------- |------|
| pay_config | obj | 支付配置 |
| ...        | obj | 其他配置 |

## 一些配置信息
- url: /user/use_cdkey
- method: POST
- 输入参数：
-
| 字段名        | 类型     | 说明   |
|------------|--------|------|
| uid        | int    | uid  |
| sid        | string | 会话id |

- 返回参数：

| 字段名  | 类型     | 说明   |
|------| --------- |------|
| data | []array | 奖励内容 |