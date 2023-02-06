build:
	go build 
pi:
	xgo --targets=linux/arm github.com/kaepa3/getters
gcp:
	xgo --targets=linux/amd64 github.com/kaepa3/getters
