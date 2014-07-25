DIR=/lib/modules/$(shell uname -r)/build
SRC =$(shell pwd)
obj-m:=testlllmodule.o
all:
	make -C $(DIR) M=$(SRC) modules
clean: 
	rm -f *.o *.ko
	