# 2024-new-blog-server
放空心声 后端服务,采用go+ echo + gorm  等技术栈开发.  


## Model 设计

### 1. 用户表
- id
- nickname
- avatar
- sex
- email
- password
- phone
- created_at
- updated_at

### 2. 漂流瓶表 bottle
- id
- user_id
- content
- image_url
- audio_url
- mood
- topic_id
- is_public
- views
- resonance_value  // 共鸣值
- created_at
- updated_at

### 3. 用户打开的漂流瓶
- id
- bottle_id
- user_id
- created_at
- updated_at

### 3. 话题表 topic
- id
- title
- cover_url
- description
- content          // 话题的详细内容
- views           // 话题的浏览量
- participants    // 参与人数（发布相关漂流瓶的用户数）
- bottle_count    // 该话题下的漂流瓶总数
- status         // 话题状态：0-草稿 1-发布 2-结束 3-删除
- start_time     // 话题开始时间
- end_time       // 话题结束时间（可选）
- creator_id     // 话题发起人ID
- creator_name   // 话题发起人昵称（冗余字段，方便查询展示）
- created_at
- updated_at

> 在这种场景下，使用冗余字段是比较好的选择，因为：
> 存储空间成本相对较低（昵称字段不会占用太多空间）
>   话题列表的查询频率远高于用户修改昵称的频率
> 用户体验更重要，快速加载话题列表是关键
> 这是一个典型的"以空间换时间"的设计策略，在类似的社交内容展示场景中很常见。
> 如果后续发现维护数据一致性成本过高，我们也可以考虑：
> 使用视图（View）
> 使用数据库触发器自动同步
> 或者去掉冗余字段，改用联表查询
> 但就目前的需求来说，保留冗余字段是合理的选择。

### 4. 话题参与记录表 topic_participant
- id
- topic_id       // 话题ID
- user_id        // 参与用户ID
- bottle_id      // 参与漂流瓶ID
- created_at     // 首次参与时间
- updated_at     // 最后参与时间

统计分析场景
统计每个话题的真实参与人数
分析用户在话题中的活跃度（通过 bottle_count）
追踪用户首次参与时间和最后参与时间，了解话题的生命周期
用户个人中心场景
展示"我参与过的话题"
显示用户在每个话题下发布的漂流瓶数量
记录用户的参与历史
3. 话题热度分析场景
通过参与人数和漂流瓶数量计算话题热度
识别最活跃的参与者
分析话题的参与趋势（通过时间戳）
推荐系统场景
基于用户参与历史推荐相似话题
发现用户感兴趣的话题类型
识别话题间的关联性
成就系统场景
追踪用户参与话题的数量
设置参与奖励（如"话题活跃者"徽章）
记录用户的参与成就
运营分析场景
分析话题的参与度
识别最受欢迎的话题类型
评估话题的运营效果



## API 接口设计

### 用户资源 (Users)
#### 认证相关
POST   /api/auth/register          # 用户注册
POST   /api/auth/login            # 用户登录
POST   /api/auth/logout           # 用户登出
POST   /api/auth/refresh-token    # 刷新token

#### 用户信息
GET    /api/users/me              # 获取当前用户信息
PUT    /api/users/me              # 更新当前用户信息
GET    /api/users/:id             # 获取指定用户信息
GET    /api/users/:id/bottles     # 获取用户的漂流瓶列表
GET    /api/users/:id/topics      # 获取用户创建的话题列表

### 漂流瓶资源 (Bottles)
POST   /api/bottles              # 创建漂流瓶
GET    /api/bottles             # 获取漂流瓶列表
GET    /api/bottles/random      # 随机获取漂流瓶
GET    /api/bottles/:id         # 获取漂流瓶详情
PUT    /api/bottles/:id         # 更新漂流瓶
DELETE /api/bottles/:id         # 删除漂流瓶
PUT    /api/bottles/:id/visibility  # 修改漂流瓶可见性

### 话题资源 (Topics)
POST   /api/topics              # 创建话题
GET    /api/topics             # 获取话题列表
GET    /api/topics/trending    # 获取热门话题
GET    /api/topics/:id         # 获取话题详情
PUT    /api/topics/:id         # 更新话题
DELETE /api/topics/:id         # 删除话题
GET    /api/topics/:id/stats   # 获取话题统计数据

### 话题参与记录资源 (Topic Participants)
GET    /api/topics/:id/participants        # 获取话题参与者列表
GET    /api/topics/:id/participants/stats  # 获取话题参与统计
GET    /api/users/:id/participated-topics  # 获取用户参与的话题列表

### 查询参数规范
所有列表接口通用查询参数：
- page: 页码（默认1）
- page_size: 每页数量（默认10）
- sort: 排序方式（字段名_asc/字段名_desc）

### 响应格式规范
- 所有成功响应都使用200状态码
- 所有错误响应都使用400状态码
- 所有未认证响应都使用401状态码
- 所有未授权响应都使用403状态码
- 所有资源不存在响应都使用404状态码
- 所有服务器错误响应都使用500状态码

约束
1. 接口需要写在handle 里面, 返回值需要使用tools/response.go 里面的函数.
2. 接口采用RESTful 风格, 使用GET, POST, PUT, DELETE 等HTTP 方法.
3. 接口需要进行权限控制, 需要使用jwt 进行token 认证.
4. 接口需要进行参数校验, 需要使用validator 进行参数校验.
5. 接口需要进行日志记录, 需要使用logrus 进行日志记录.
6. 接口需要进行错误处理, 需要使用echo 的错误处理函数.
7. 对于需要返回分页数据, 需要使用tools/response.go 里面的PagedOkResponse 函数.
8. 对于一些通用的工具函数 统一写在tools 目录下.
9. 对于一些通用的配置 统一写在config 目录下.
10. 对于一些通用的中间件 统一写在middleware 目录下.
10. 对于一些通用的模型 统一写在model 目录下.
11. 对于一些通用的常量 统一写在constant 目录下.
12. 对于某些接口 需要编写大量代码, 可以写在service 目录下.

13. 路由统一写在router.go 里面.
14. 主函数统一写在main.go 里面.

## 项目规范

### 目录结构

├── main.go         # 主程序入口
├── router.go       # 路由配置
├── config/         # 配置文件目录
│ ├── config.go     # 配置结构定义
│ └── config.yaml # 配置文件
├── constant/       # 常量定义
│ ├── error.go      # 错误码常量
│ └── common.go     # 通用常量
├── handler/        # 接口处理器
│ ├── auth.go       # 认证相关处理器
│ ├── user.go # 用户相关处理器
│ ├── bottle.go     # 漂流瓶相关处理器
│ └── topic.go      # 话题相关处理器
├── middleware/     # 中间件
│ ├── auth.go       # JWT认证中间件
│ ├── logger.go     # 日志中间件
│ └── validator.go  # 参数校验中间件
├── model/          # 数据模型
│ ├── user.go       # 用户模型
│ ├── bottle.go     # 漂流瓶模型
│ └── topic.go      # 话题模型
├── service/        # 业务逻辑层
│ ├── auth.go       # 认证相关服务
│ ├── user.go       # 用户相关服务
│ ├── bottle.go     # 漂流瓶相关服务
│ └── topic.go      # 话题相关服务
└── tools/          # 工具函数
│ ├── response.go   # 统一响应处理
│ ├── jwt.go        # JWT工具
│ └── validator.go  # 参数校验工具

### 开发规范

#### 1. 接口处理（Handler）
- 所有接口处理器必须位于 `handler` 目录下
- 处理器方法命名规则：`Handle{Action}{Resource}`
- 必须使用 `tools/response.go` 中的响应函数返回数据

#### 2. RESTful API规范
- 使用标准HTTP方法（GET, POST, PUT, DELETE）
- URL使用小写字母，单词间用连字符"-"分隔
- 资源名称使用复数形式

#### 3. 认证与授权
- 使用JWT进行身份认证
- 在 `middleware/auth.go` 中实现认证中间件
- 需要认证的路由必须使用认证中间件

#### 4. 参数校验
- 使用 `validator` 包进行参数校验
- 在模型中定义验证规则

#### 5. 日志记录
- 使用 `logrus` 进行日志记录
- 在中间件中统一记录请求日志
- 重要操作必须记录日志

#### 6. 错误处理
- 使用 Echo 的错误处理机制
- 统一错误响应格式

#### 7. 分页处理
- 使用 `tools/response.go` 中的 `PagedOkResponse` 函数

#### 8. 配置管理
- 配置文件放在 `config` 目录下
- 使用 `viper` 包加载配置

#### 9. 业务逻辑
- 复杂业务逻辑放在 `service` 层
- Service 层方法命名规则：`{Action}{Resource}`

#### 10. 数据模型
- 模型定义放在 `model` 目录下
- 使用 `gorm` tag 定义数据库映射
- 使用 `json` tag 定义JSON序列化规则

#### 11. 常量定义
- 错误码常量放在 `constant/error.go`
- 其他常量放在 `constant/common.go`

#### 12. 工具函数
- 通用工具函数放在 `tools` 目录下
- 工具函数应该是无状态的

#### 13. 路由配置
- 所有路由配置统一在 `router.go` 中管理
- 按资源类型组织路由

#### 14. 主程序入口
- 在 `main.go` 中初始化各组件
- 启动服务器


