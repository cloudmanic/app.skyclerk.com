# {{ ansible_managed }}

FROM golang:1.9.2-alpine3.7

MAINTAINER Spicer Mathtews <spicer@cloudmanic.com>

ARG P_UID=0
ARG P_GID=0

# Gots to go Root to make some changes
USER root

# Reset Password and group files
RUN set -x \
      && echo "root:x:0:0:root:/root:/bin/ash" > /etc/passwd \
      && echo "root:::0:::::" > /etc/shadow \
      && echo "root:x:0:root" > /etc/group 
      
# Ensure www-data user exists (as our local user)
RUN set -x \ 
  && addgroup -g ${P_GID} -S backup \
  && adduser -u ${P_UID} -D -S -G backup backup 
   
# Build App
RUN apk --no-cache add --virtual build-dependencies go git \
  && go get github.com/cloudmanic/cloudmanic-backup \
  && apk del --purge build-dependencies

# Install Helpful Programs
RUN apk --no-cache add bash mysql-client

# This needs to be at the bottom. (back from root)
USER backup

WORKDIR /go/bin

CMD ["cloudmanic-backup"]