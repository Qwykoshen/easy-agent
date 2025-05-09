# easy-agent

大模型代理服务Web系统设计文档
技术栈：Gin (Go) + Vue3 + DeepSeek API + MySQL + Redis
系统架构：前后端分离模式，采用RESTful API交互

一、系统功能模块设计
1. 用户管理模块
   功能清单：

用户注册（邮箱/手机号+密码）

JWT鉴权登录/登出

Token自动续期

用户信息修改（头像、昵称）

2. 多模态交互模块
   功能	输入	输出	调用API
   文本对话	用户文本消息	大模型生成文本回复	DeepSeek Chat API
   主题图像生成	图像主题描述文本	图像URL（PNG）	DeepSeek DALL-E API
   文本转语音（TTS）	待合成文本	音频文件URL（MP3）	DeepSeek TTS API

二、技术架构设计
1. 后端架构（Gin框架）
   go
   // 项目结构  
   ├── config              // 配置管理  
   ├── controllers         // 控制器层  
   │   ├── auth.go         // 鉴权逻辑  
   │   ├── chat.go         // 对话处理  
   │   ├── image.go        // 图像生成  
   │   └── tts.go          // 语音合成  
   ├── middleware          // 中间件  
   │   ├── jwt.go          // JWT校验  
   │   └── ratelimit.go    // 接口限流  
   ├── models              // 数据模型  
   │   ├── user.go         // 用户ORM  
   │   └── task.go         // 异步任务记录  
   ├── services            // 业务逻辑  
   │   ├── deepseek.go     // API调用封装  
   │   └── queue.go        // 异步任务队列  
   └── main.go             // 入口文件
2. 前端架构（Vue3 + Pinia）
   javascript
   // 功能模块划分  
   src/  
   ├── api/                // Axios接口封装  
   │   ├── auth.js         // 用户相关API  
   │   └── model.js        // 大模型交互API  
   ├── stores/             // Pinia状态管理  
   │   ├── user.js         // 用户状态  
   │   └── chat.js         // 对话历史  
   ├── views/              // 页面组件  
   │   ├── ChatView.vue     // 文本对话界面  
   │   ├── ImageGen.vue    // 图像生成界面  
   │   └── TTSView.vue     // 语音合成界面  
   三、核心接口设计
1. 用户管理接口
   端点	方法	参数	响应
   /api/v1/register	POST	username, password	{user_id, token}
   /api/v1/login	POST	username, password	{token, expires_in}
   /api/v1/userinfo	GET	Header: Authorization	{username, avatar, created_at}
2. 大模型服务接口
   端点	方法	参数	响应
   /api/v1/chat	POST	{content: string}	{response: string, timestamp}
   /api/v1/gen-image	POST	{prompt: string}	{task_id, status_url}
   /api/v1/gen-tts	POST	{text: string}	{audio_url: string}
   四、关键技术实现
1. 流式对话响应（SSE）
   javascript
   // 前端建立SSE连接  
   const eventSource = new EventSource('/api/v1/chat-stream');  
   eventSource.onmessage = (e) => {  
   this.chatHistory.push(JSON.parse(e.data));  
   };
