.PHONY: start run_api run_racing run_sport stop_all

start: run_api run_racing run_sport

run_api:
	cd api && go build && ./api &

run_racing:
	cd racing && go build && ./racing &

run_sport:
	cd sport && go build && ./sport &

stop_all: 
	chmod +x stop_all.sh && ./stop_all.sh