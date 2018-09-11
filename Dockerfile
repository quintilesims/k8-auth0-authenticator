FROM alpine
RUN apk add --no-cache ca-certificates
ADD ./static /static
ADD ./templates /templates
ADD ./k8-auth0-authenticator /
CMD ["/k8-auth0-authenticator"]
