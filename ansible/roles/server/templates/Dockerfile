# {{ ansible_managed }}

FROM golang:1.8.1-alpine

MAINTAINER Spicer Matthews <spicer@cloudmanic.com>

ARG P_UID=0
ARG P_GID=0

ENV GOBIN=/go/bin

# Gots to go Root to make some changes
USER root

# Install bash
RUN apk add --no-cache bash git

# Reset Password and group files
RUN set -x \
      && echo "root:x:0:0:root:/root:/bin/ash" > /etc/passwd \
      && echo "root:::0:::::" > /etc/shadow \
      && echo "root:x:0:root" > /etc/group 
      
# Ensure deploy user exists (as our local user)
RUN set -x \ 
  && addgroup -g ${P_GID} -S deploy \
  && adduser -u ${P_UID} -D -S -G deploy deploy 
   
# Setup directory and Set perms for app root
RUN set -x \
  && mkdir /work \
  && chmod 755 /work \    
  && chown -R deploy:deploy /work  
  
# Copy binary into place.
COPY app.options.cafe /work/app

# Copy frontend into place.
COPY frontend /frontend

# Expose http and https ports (we do not use 80 & 443 as we run the docker container as non-root)
EXPOSE 7080
EXPOSE 7043  

# This needs to be at the bottom. (back from root)
USER deploy

# Working directory
WORKDIR /work

# Start the server
CMD ["./app"]