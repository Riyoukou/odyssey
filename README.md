# odyssey

tips:当前主包实力有限，只能整合cicd工具使用，当前功能实现依赖:Argocd,Jenkins,Git这三个传统工具，还不能脱离传统工具流，后续会逐步完善

GO + VUE3编写的运维平台

后端:使用GOGIN框架

VUE3:使用Vue3+TypeScript+Vite

本地调试指南:
```
##数据库准备
#启动一个mysql和redis(如下配置无需更改dev配置文件)
#mysql中创建odyssey数据库
#导入initsql中的表结构
vim /etc/host/
127.0.0.1 odyssey.mysql ##127.0.0.1改为自己数据库的地址
127.0.0.1 odyssey.redis ##127.0.0.1改为自己数据库的地址
##后端调试
go mod tidy
GO_ENV=dev go run cmd/server/main.go
##启动后端后尝试
curl 127.0.0.1:8000/man
##前端调试
#默认用户名/密码:admin/123,user/123
pnpm install
cd web && npm run dev
```

预期实现功能CI/CD、资产管理、工单系统

目前已有版本写在公司私有仓库，没有脱敏，正在慢慢迁移到github

代码初学者，有什么代码上的意见可以直接提出来，听劝

觉得有点东西的希望家人们可以给个🌟

功能迁移进度
| 名称         | 描述         |
|--------------|--------------|
| 登陆/注册     | ✅ |
| 统一身份验证  |  |
| 用户管理 | ✅ |
| 集群管理 | ✅ |
| 凭证管理 | ✅ |
| 工具管理 | ✅ |
| 环境管理 | ✅ |
| 项目管理 | ✅ |
| 服务管理 | ✅ |
| CI | ✅ |
| CD |  |
