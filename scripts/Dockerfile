
## NOT USED YET SOME DAY.......




## --- Build Go APP --- ##
FROM golang AS BE

# Set go bin which doesn't appear to be set already.
ENV GOBIN /go/bin

# Build directories
RUN mkdir /go/src/app.skyclerk.com

# Add SSH credentials to build
ARG GITHUB_PRIVATE_KEY
RUN mkdir /root/.ssh/
RUN echo "${GITHUB_PRIVATE_KEY}" >  /root/.ssh/id_rsa

# Make sure the github domain is accepted
RUN touch /root/.ssh/known_hosts
RUN chmod 600 /root/.ssh/id_rsa
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts

# Clone backend repo
RUN git clone git@github.com:cloudmanic/app.skyclerk.com.git /go/src/app.skyclerk.com
WORKDIR /go/src/app.skyclerk.com/backend

# Ensure vendor is up-to-date
RUN go get -u github.com/golang/dep/...
RUN dep ensure --vendor-only

# Build this puppy lean!
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o builds/app.skyclerk.com

# ## --- Build Front-End APP --- ###
FROM node:11.5.0 AS FE

# Install linux packages we need for the build
RUN apt-get update
RUN apt-get install gettext-base

# Copy FE repo
COPY --from=BE /go/src/app.skyclerk.com /usr/src/app.skyclerk.com
WORKDIR /usr/src/app.skyclerk.com/frontend

# Add `/usr/src/app.skyclerk.com/node_modules/.bin` to $PATH
ENV PATH /usr/src/app.skyclerk.com/node_modules/.bin:$PATH

# Install and cache app dependencies
RUN npm install
RUN npm install -g @angular/cli@7.2.3

# Build app
RUN ng build --prod --deploy-url=/

## --- Build Final Image --- ##
# Move to 'scratch' when the certificate issue has been resolved.
FROM golang:1.11-alpine

# Set working directory
WORKDIR /go/bin

# Copy Go app over
COPY --from=BE /go/src/app.skyclerk.com/backend/builds/app.skyclerk.com .

# Copy FE app over
COPY --from=FE /usr/src/app.skyclerk.com/frontend/dist/frontend/ /frontend/

CMD ["./app.skyclerk.com"]

EXPOSE 80
EXPOSE 443
