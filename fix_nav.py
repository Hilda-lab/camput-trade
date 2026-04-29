import os
import re
import glob

nav_str = """    <div class="collapse navbar-collapse">
      <div class="navbar-nav me-auto">
        <a class="nav-link" href="/items">商品</a>
        <a class="nav-link" href="/users">用户</a>
        <a class="nav-link" href="/orders">订单</a>
        <a class="nav-link" href="/reports">查询统计</a>
      </div>
      <div class="navbar-nav ms-auto">
        {{if .loggedIn}}
        <span class="navbar-text text-light me-3">
          当前登录: {{.user_name}} ({{.user_id}})
        </span>
        <a class="nav-link btn btn-sm btn-outline-danger text-danger px-3 py-1" style="height: 32px; margin-top: 4px;" href="/logout">登出身份</a>
        {{else}}
        <a class="nav-link btn btn-sm btn-outline-success text-success px-3 py-1" style="height: 32px; margin-top: 4px;" href="/login">点击登录系统</a>
        {{end}}
      </div>
    </div>
  </div>
</nav>"""

rx = re.compile(r'<div class="navbar-nav.*?</nav>', re.DOTALL)

for f in glob.glob("templates/*.html"):
    with open(f, 'r', encoding='utf-8') as file:
        content = file.read()
    if 'navbar-nav' in content:
        content = rx.sub(nav_str, content)
        with open(f, 'w', encoding='utf-8') as file:
            file.write(content)
            print(f"Updated {f}")
