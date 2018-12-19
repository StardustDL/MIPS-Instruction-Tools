.PHONY : rund build run

SHELL = powershell.exe
ARGS = -help

rund : # run directly
	cd src ; go run .

build :
	cd src ; go build -o ../bin/mip.exe

run :
	cd bin ; ./mip $(ARGS)
