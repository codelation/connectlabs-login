FROM alpine:3.4

# TODO: The 'git' dependency should be updated with fixed revisions
ENV BUILD_DEPS 'go=1.6.3-r0 git alpine-sdk'
ENV DEL_BUILD_DEPS 'alpine-sdk expat libcurl libssh2 pcre git go'
ENV APP_NAME=connectlabs-login APP_SRC=github.com/ryanhatfield/connectlabs-login

WORKDIR /opt/build/src

RUN apk --update --no-cache add openssl ca-certificates

ADD . /opt/build/src/$APP_SRC

# This runs as one command/layer, otherwise deleting and
# cleaning up files wouldn't reduce the server file size.
RUN echo $APP_NAME $APP_SRC && \
    apk add --update $BUILD_DEPS && \
    export GOPATH=/opt/build/ CGO_ENABLED=0 && \
    cd /opt/build/src/$APP_SRC && \
    go get ./... && \
    cd /opt/build && \
    go build -o /opt/static/app $APP_SRC && \
    cp -r /opt/build/src/$APP_SRC/www /opt/static/www && \
    apk del $DEL_BUILD_DEPS && \
    rm -rf /opt/build /var/cache/apk/*

# The old work directory has been deleted, change to avoid errors
# in some Docker hosting systems (heroku for one)
WORKDIR /opt/static

CMD /opt/static/app
