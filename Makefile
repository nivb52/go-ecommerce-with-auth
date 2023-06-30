# WIN
SHELL=cmd
GOSTRIPE_PORT=4000
API_PORT=4001
#DSN="root@(localhost:3306)/widgets?parseTime=true&tls=false"

########### SERVE / DEV  => WATCH ALL  ###########
dev : dev_front dev_back db_up 
serve : dev

########### BUILD ALL  ###########
build: clean build_front build_back
	@echo All binaries built!

## clean: cleans all binaries and runs go clean
clean:
	@echo Cleaning...
	@echo y | DEL /S dist
	@go clean
	@echo Cleaned and deleted binaries

########### WATCH / DEV ###########
watch_front:
	@echo Watch front end...
	@air -c ./web.air.toml
	@echo Front end built!
dev_front: watch_front

watch_back:
	@echo Watch back end...
	@air -c ./api.air.toml
	@echo Front end built!
dev_back: watch_back

debug_front:
	@echo Debug front end...
	@ D:\d-dev\goworkspace\bin\dlv.exe dap --listen=127.0.0.1:62712 from ${pwd}\cmd\web
	@echo Front end built!

debug_back:
	@echo Debug back end...
	@ D:\d-dev\goworkspace\bin\dlv.exe dap --listen=127.0.0.1:62712 from ${pwd}\cmd\api
	@echo Front end built!

########### START  ###########
start : start_front start_back

start_front: build_front
	@echo Starting the front end...
	start /B .\dist\gostripe.exe 
	@echo Front end running!

build_front:
	@echo Building front end...
	@go build -o dist/gostripe.exe ./cmd/web
	@echo Front end built!



start_back: build_back
	@echo Starting the back end...
	start /B .\dist\gostripe_api.exe 
	@echo Back end running!

build_back:
	@echo Building back end...
	@go build -o dist/gostripe_api.exe ./cmd/api
	@echo Back end built!

########### STOP ###########
stop: stop_front stop_back db_down
	@echo All applications stopped

stop_front:
	@echo Stopping the front end...
	@taskkill /IM gostripe.exe /F
	@echo Stopped front end

stop_back:
	@echo Stopping the back end...
	@taskkill /IM gostripe_api.exe /F
	@echo Stopped back end


########### DB ###########
db_up: db_down
	@echo "Docker compose up: db image..."
	@docker-compose --env-file .env --env-file local.env -p widgets_db  up -d  
	@echo "Docker db up!"

db_down:
	@echo "Docker compose down: db image..."
	@docker-compose  -p widgets_db down
	@echo "Docker db down!"

########### UTILS  ###########
copy_env:
	@echo "Copying env file..."
	@cat .env | head -1
	@cp .env ./cmd/api/.env
	@cp .env ./cmd/web/.env

########### TEST ###########
