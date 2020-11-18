# Use the official go docker image built on debian.
FROM golang

# Grab the source code and add it to the workspace.
ADD . /Hermes
# Install revel and the revel CLI.
RUN go get github.com/revel/revel
RUN go get github.com/revel/cmd/revel

# Use the revel CLI to start up our application.
WORKDIR /Hermes
# ENTRYPOINT bash
ENTRYPOINT revel run -a Hermes prod

# Open up the port where the app is running.
EXPOSE 9000

