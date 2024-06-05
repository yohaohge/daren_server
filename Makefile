all:LittleVideo

LittleVideo:
	cd app/ && go build -o ../LittleVideo

clean :
	-rm LittleVideo
