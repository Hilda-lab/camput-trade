package main

import "os"

func main() {
html := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>用户登录 - {{ .title }}</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="/static/style.css" rel="stylesheet">
    <style>
        body { background-color: #f8f9fa; }
        .login-container { max-width: 400px; margin: 100px auto; }
        .card { border-radius: 12px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); border: none; }
        .card-header { background-color: #0d6efd; color: white; text-align: center; border-top-left-radius: 12px; border-top-right-radius: 12px; padding: 20px; }
        .btn-primary { width: 100%; border-radius: 8px; padding: 10px; font-weight: bold; }
        .form-control { border-radius: 8px; padding: 10px; }
    </style>
</head>
<body>

<nav class="navbar navbar-expand-lg navbar-dark bg-primary shadow-sm mb-4">
    <div class="container">
        <a class="navbar-brand" href="/">校园二手交易平台</a>
        <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav">
            <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbarNav">
            <ul class="navbar-nav me-auto">
                <li class="nav-item"><a class="nav-link" href="/">首页</a></li>
            </ul>
        </div>
    </div>
</nav>

<div class="container login-container">
    <div class="card">
        <div class="card-header">
            <h3 class="mb-0">欢迎回来</h3>
        </div>
        <div class="card-body p-4">
            {{if .error}}
            <div class="alert alert-danger">{{.error}}</div>
            {{end}}
            
            <form action="/login" method="POST">
                <div class="mb-3">
                    <label class="form-label text-muted">邮箱地址 / 账号</label>
                    <input type="text" class="form-control" name="email" placeholder="example@stu.edu.cn 或 u001" required>
                </div>
                <div class="mb-4">
                    <label class="form-label text-muted">密码</label>
                    <input type="password" class="form-control" name="password" placeholder="请输入密码" required>
                </div>
                <button type="submit" class="btn btn-primary mb-3">立 即 登 录</button>
                <div class="text-center">
                    <a href="/register" class="text-decoration-none text-muted">没有账号？点击注册</a>
                </div>
            </form>
        </div>
    </div>
</div>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/js/bootstrap.bundle.min.js"></script>
</body>
</html>`
os.WriteFile("templates/login.html", []byte(html), 0644)
}
