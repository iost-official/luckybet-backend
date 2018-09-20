run : build
	nohup ./luckybet -w
drop :
	mongo --eval 'db=connect("localhost:27017/lucky_bet"); db.dropDatabase()'
build :
	go build -o luckybet
stop :
	killall luckybet
set :
	git commit -a -m "contract" && git push
