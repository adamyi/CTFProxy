# hacky script to start gunicorn

import re
import sys
import os
import logging

from gunicorn.app.wsgiapp import run

logging.basicConfig()

os.system("python -m compileall /app/infra/ctfd/image.binary.runfiles/")
os.system("python /manage.py db upgrade")

sys.argv[0] = re.sub(r'(-script\.pyw?|\.exe)?$', '', sys.argv[0])

sys.argv.append("CTFd:create_app()")
sys.argv.append("--workers=1")
sys.argv.append("--worker-class=gevent")
sys.argv.append("--bind")
sys.argv.append("0.0.0.0:80")

os.chdir("../ctfd/CTFd/")

sys.exit(run())
