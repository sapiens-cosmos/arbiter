FROM ubuntu:20.04

RUN apt-get update && apt-get install -y git build-essential wget
RUN wget https://github.com/mikefarah/yq/releases/download/v4.9.3/yq_linux_amd64 -O /usr/bin/yq && chmod +x /usr/bin/yq

RUN wget https://dl.google.com/go/go1.17.3.linux-amd64.tar.gz
RUN tar -xvf go1.17.3.linux-amd64.tar.gz
RUN mv go /usr/local

ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV GOBIN /go/bin
ENV PATH $GOBIN:$GOROOT/bin:$PATH

COPY ./ /app

WORKDIR /app

RUN make build

WORKDIR /app/build

RUN ./arbiterd init localnet --chain-id arbiter-localnet-1

RUN yq eval '.app_state.bond.params.policies[0].bond_denom = "ugreen"' ~/.arbiter/config/genesis.json \
    | yq eval '.app_state.bond.base_denom = "uarb"' - \
    | yq eval '.app_state.staking.params.bond_denom = "uarb"' - \
    | yq eval '.app_state.crisis.constant_fee.denom = "uarb"' - \
    | yq eval '' -j - \
    > ~/.arbiter/config/genesis.tmp.json
RUN cp ~/.arbiter/config/genesis.tmp.json ~/.arbiter/config/genesis.json

RUN sed -i 's#tcp://127.0.0.1:26657#tcp://0.0.0.0:26657#g' ~/.arbiter/config/config.toml
RUN sed -i 's/cors_allowed_origins = \[\]/cors_allowed_origins = \["*"\]/g' ~/.arbiter/config/config.toml
RUN sed -i 's/pruning = "default"/pruning = "nothing"/g' ~/.arbiter/config/app.toml
RUN sed -i 's/enable = false/enable = true/g' ~/.arbiter/config/app.toml
RUN sed -i 's/swagger = false/swagger = true/g' ~/.arbiter/config/app.toml
RUN sed -i 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/g' ~/.arbiter/config/app.toml

RUN echo "high gain deposit chuckle hundred regular exist approve peanut enjoy comfort ride" | ./arbiterd keys add val --recover --keyring-backend test
RUN echo "health nest provide snow total tissue intact loyal cargo must credit wrist" | ./arbiterd keys add local1 --recover --keyring-backend test
RUN echo "canyon stone next tenant trial ugly slim luggage ski govern outside comfort" | ./arbiterd keys add local2 --recover --keyring-backend test
RUN echo "travel renew first fiction trick fly disease advance hunt famous absurd region" | ./arbiterd keys add local3 --recover --keyring-backend test

RUN ./arbiterd add-genesis-account $(./arbiterd keys show val -a --keyring-backend test) 10000000000uarb,10000000000ugreen
RUN ./arbiterd add-genesis-account $(./arbiterd keys show local1 -a --keyring-backend test) 10000000000uarb,10000000000ugreen
RUN ./arbiterd add-genesis-account $(./arbiterd keys show local2 -a --keyring-backend test) 10000000000uarb,10000000000ugreen
RUN ./arbiterd add-genesis-account $(./arbiterd keys show local3 -a --keyring-backend test) 10000000000uarb,10000000000ugreen

RUN ./arbiterd gentx val 10000000000uarb --chain-id arbiter-localnet-1 --keyring-backend test
RUN ./arbiterd collect-gentxs

EXPOSE 26657
EXPOSE 1317
EXPOSE 9090
CMD ./arbiterd start
