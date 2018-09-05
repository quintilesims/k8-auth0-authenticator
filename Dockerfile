FROM alpine
RUN apk add --no-cache ca-certificates
ADD ./views /views
ADD ./k8-auth0-authenticator /
CMD ["/k8-auth0-authenticator"]
