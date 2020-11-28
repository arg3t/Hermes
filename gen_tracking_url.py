#!/bin/python

import sys
from base64 import b64decode
from hashlib import sha256
from urllib.parse import quote

# Change these variables depending on your setup.
userid = ""
key = ""
domain = ""

if len(sys.argv) > 2:
    subject = input("Please enter subject of email:")
    recipient = input("Please enter recipient of email: ")
else:
    message = b64decode(sys.argv[1]).decode()

    subject = "NIL"
    recipient = "NIL"

    for i in message.split("\n"):
        if i.split(" ")[0] == "Subject:":
            subject = i[9:]
        elif i.split(" ")[0] == "To:":
            recipient = i[4:]

identifier = subject + recipient + userid + key # Generate hash
identifier_hash = sha256(identifier.encode("utf-8")).hexdigest()

url = "{}/read/{}/{}/{}/{}".format(domain,
                                   userid,
                                   quote(subject),
                                   quote(recipient),
                                   identifier_hash) #Generate tracking url

print(url)
