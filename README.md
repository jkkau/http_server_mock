build:
    cd server mock
    go build

run:
    mock_http_server listen_port delayTime[millisecond]