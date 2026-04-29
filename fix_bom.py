import os
import glob
import re

for f in glob.glob("templates/*.html"):
    if os.path.basename(f) in ["base.html", "login.html"]: continue
    with open(f, 'r', encoding='utf-8-sig') as file:
        s = file.read()
    s = re.sub(r'^\s*\{\{define\s+"[^"]+"\}\}\s*', '', s)
    s = re.sub(r'\s*\{\{end\}\}\s*$', '', s)
    with open(f, 'w', encoding='utf-8') as file:
        file.write(s)
