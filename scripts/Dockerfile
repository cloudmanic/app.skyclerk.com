# NOT USED YET.....

############################
# STEP 1 build executable binary
############################

FROM golang:1.14-alpine AS builder

# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache ca-certificates
RUN update-ca-certificates

############################
# STEP 2 build a small image
############################
FROM scratch

# Set working directory
WORKDIR /go/bin

# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy Go app over
COPY build/app.skyclerk.com .

# Copy FE app over
COPY build/frontend/ /frontend/

CMD ["./app.skyclerk.com"]

EXPOSE 80
EXPOSE 443
