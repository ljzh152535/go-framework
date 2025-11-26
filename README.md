## golang 框架工具
### 打包
```
go mod init github.com/ljzh152535/go-framework
go mod tidy
```

### git全局配置

```bash
# 设置你的用户名（替换成你的名字或昵称）
git config --global user.name "Your Name"

# 设置你的邮箱地址（替换成你常用的邮箱）
git config --global user.email "your.email@example.com"


git config --global user.name "ljzh152535"
git config --global user.email "ljzh152535@163.com"
```





### 初始化推送到github

```
git init
git add .
git commit -m "first commit"
git branch -M main
git remote add origin git@github.com:ljzh152535/go-framework.git
git push -u origin main
```

### 更新代码推送github
```
git add .
git commit -am "first commit"
git push -u origin main
```


### 下载代码使用
```
# 清除 Go 模块缓存
go clean -modcache
go mod tidy -modcache -module=github.com/ljzh152535/go-framework
go clean -modcache -module=github.com/ljzh152535/go-framework
go get -u github.com/ljzh152535/go-framework
```

### 添加tag
```
git tag v0.0.1
git push origin v0.0.2
```



### 查看当前远程仓库地址

```bash
# git remote -v
origin  git@github.com:YourUsername/your-repository-name.git (fetch)
origin  git@github.com:YourUsername/your-repository-name.git (push)

# 修改为 HTTPS 格式
执行以下命令，将远程仓库地址改为 HTTPS 格式：
git remote set-url origin https://github.com/YourUsername/your-repository-name.git

git remote set-url origin https://github.com/ljzh152535/go-framework.git

验证修改是否成功
git remote -v
```

