## golang 框架工具
### 打包
```
go mod init github.com/ljzh152535/go-framework
go mod tidy
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
