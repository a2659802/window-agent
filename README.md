## window agent
该项目用于windows自动改密

### 功能
1. agent发起并维护与堡垒机的TCP连接
2. 从连接中接收指令并执行
3. 将指令执行结果返回给堡垒机

### 认证方案
1. RSA, agent持有私钥, 将公钥在初始化时发给堡垒机. 每个agent一套密钥
2. RSA, 堡垒机持有私钥，agent用公钥加密请求，所有agent共用一套密钥
3. Basic, 堡垒机硬编码agent专用账号，agent使用密码认证(不安全)
4. ApiKey, 管理页面提供生成ApiKey页面，每个agent申请一个apiKey
5. Basic+RSA, 初始化时使用Basic, 然后生成RSA, 将pubKey保存在堡垒机,private保存在本地. 后续使用RSA

### 整体流程
1. agent向FP发起SSL连接
2. 建立连接
3. 服务器ping agent
4. 服务器发送command
5. agent通知服务器command执行结果

### Layer
App      应用
Message  消息切割
Security 认证+加解密
TCP      


### 启动流程
1. 非服务模式启动，检查配置文件
2. 检查windows service是否已注册
3. 注册到service, 类型为自动启动
4. 以服务模式启动，读取配置文件，解析出堡垒机IP,端口,认证凭证
5. 连接堡垒机，协议握手、认证
6. 定时发送心跳，等待服务端消息

### 参考资料
