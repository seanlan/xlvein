CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go && \
rsync --progress main tx-tele-bot:/data/vein/main.new && \
ssh tx-tele-bot "supervisorctl stop vein && \
                \cp /data/vein/main /data/vein/main.last && \
                \cp /data/vein/main.new /data/vein/main && \
                supervisorctl start vein" && \
rsync --progress main tx-tele-bot-2:/data/vein/main.new && \
ssh tx-tele-bot-2 "supervisorctl stop vein && \
                \cp /data/vein/main /data/vein/main.last && \
                \cp /data/vein/main.new /data/vein/main && \
                supervisorctl start vein" && \
rm -rf main