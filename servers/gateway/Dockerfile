FROM alpine
COPY gateway /gateway 
RUN apk update && apk add ca-certificates
EXPOSE 443
ENTRYPOINT ["/gateway"]
