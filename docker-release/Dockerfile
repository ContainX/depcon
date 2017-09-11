FROM alpine
RUN apk update && apk add ca-certificates

ADD depcon /usr/bin

ENV DEPCON_MODE="" \
	MARATHON_HOST="http://localhost:8080" \
	MARATHON_USER="" \
	MARATHON_PASS="" 

ENTRYPOINT ["/usr/bin/depcon"]
