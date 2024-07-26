FROM alpine:latest

ENV TOR_SOCKS_PORT=9050
ENV TOR_CONTROL_PORT=9051
ENV TOR_MAX_CIRCUIT_DIRTINESS=1
ENV TOR_PASSWORD="proxy"

RUN apk update && \
    apk add --no-cache tor && \
    rm -rf /var/cache/apk/*

# Generate hashed password
RUN apk add --no-cache --virtual .build-deps make gcc g++ && \
    tor --hash-password $TOR_PASSWORD > /etc/tor/hashed_password && \
    apk del .build-deps

# Extract the hashed password and configure torrc
RUN mkdir -p /etc/tor && \
    HASHED_PASSWORD=$(cat /etc/tor/hashed_password | grep 16:) && \
    echo "SocksPort 0.0.0.0:$TOR_SOCKS_PORT" > /etc/tor/torrc && \
    echo "ControlPort 0.0.0.0:$TOR_CONTROL_PORT" >> /etc/tor/torrc && \
    echo "CookieAuthentication 0" >> /etc/tor/torrc && \
    echo "HashedControlPassword $HASHED_PASSWORD" >> /etc/tor/torrc && \
    echo "MaxCircuitDirtiness $TOR_MAX_CIRCUIT_DIRTINESS" >> /etc/tor/torrc

EXPOSE $TOR_SOCKS_PORT $TOR_CONTROL_PORT

CMD ["tor", "-f", "/etc/tor/torrc"]
