package main

import "os"

func main() {
	var g chardevgen
	g.setName("testlllmodule")
	g.setNDevs("4")
	g.setBufSize("100")

	os.Mkdir(g.name, 0777)

	f, _ := os.OpenFile(g.name+"/"+g.name+".c", os.O_RDWR|os.O_CREATE, 0644)
	f.WriteString(g.genHeader() +
		g.genDefinitions() +
		g.genOpen() +
		g.genRelease() +
		g.genIoctl() +
		g.genRead() +
		g.genWrite() +
		g.genLlseek() +
		g.genFops() +
		g.genSetup() +
		g.genInit() +
		g.genExit())
	f.Close()

	f, _ = os.OpenFile(g.name+"/Makefile", os.O_RDWR|os.O_CREATE, 0644)
	f.WriteString(g.genMakefile())
	f.Close()

}

type chardevgen struct {
	name     string
	n_devs   string
	buf_size string
}

func (g *chardevgen) setBufSize(n string) {
	g.buf_size = n
}

func (g *chardevgen) setName(name string) {
	g.name = name
}
func (g *chardevgen) setNDevs(n string) {
	g.n_devs = n
}

func (g *chardevgen) genHeader() string {
	return `#include <linux/module.h>
#include <linux/device.h>
#include <linux/uaccess.h>
#include <linux/init.h>
#include <linux/fs.h>
#include <linux/cdev.h>
#include <linux/errno.h>
	`
}

func (g *chardevgen) genDefinitions() string {
	return `
	// start of definitions
	#define BUF_SIZE ` + g.buf_size + `
	struct ` + g.name + `_dev {
		char buf[BUF_SIZE];
		struct cdev cdev;
		int buf_size;
		struct device *devicep;
	};
	unsigned int  ` + g.name + `_major;
	struct ` + g.name + `_dev ` + g.name + `_devs[` + g.n_devs + `];
	struct class * classp=NULL;
	// end of definitions
	`
}

func (g *chardevgen) genOpen() string {
	return `static int ` + g.name + `_open(struct inode * inode, struct file *filp)
	{
		int minor=MINOR(inode->i_rdev);
		if (minor > ` + g.n_devs + `-1) {
			return -ENODEV;
		}

		filp->private_data=&` + g.name + `_devs[` + g.n_devs + `];
		return 0;
	}
	`
}

func (g *chardevgen) genRelease() string {
	return `
	static int ` + g.name + `_release(struct inode * inode, struct file *filp)
	{
		return 0;
	}`
}

func (g *chardevgen) genIoctl() string {
	return `
	static int ` + g.name + `_ioctl(struct file *filp, unsigned int cmd, unsigned long arg)
	{
	struct ` + g.name + `_dev * ` + g.name + `_devp = filp->private_data;
	switch(cmd)
	{
	default:
		return -EINVAL;
	}
	return 0;
}`
}

func (g *chardevgen) genRead() string {
	return `
static ssize_t ` + g.name + `_read(struct file *filp, char __user *user_buf,size_t size, loff_t *ppos)
{
	int ret=0;
	int count=size;
	unsigned long p=*ppos;
	struct ` + g.name + `_dev *` + g.name + `_devp=filp->private_data;

	if(p>=` + g.name + `_devp->buf_size){
		return count? -ENXIO:0;
	}
	if( count>` + g.name + `_devp->buf_size-p) {
		count=` + g.name + `_devp->buf_size-p;
	}
	if( copy_to_user(user_buf,(void *)(` + g.name + `_devp->buf+p),count))
	{
		ret=-EFAULT;
	} else {
		*ppos+=count;
		return count;
		printk(KERN_INFO"` + g.name + `: read %d bytes form %d\n",count,p);
	}
	return ret;
}`
}
func (g *chardevgen) genWrite() string {
	return `
static ssize_t ` + g.name + `_write(struct file *filp, const char __user *user_buf, size_t size, loff_t *ppos)
{
	unsigned long p =*ppos;
	unsigned int count=size;
	int ret=0;
	struct ` + g.name + `_dev *` + g.name + `_devp=filp->private_data;

	if(p>=` + g.name + `_devp->buf_size){
		return count? -ENXIO:0;
	}
	if( count>` + g.name + `_devp->buf_size-p) {
		count=` + g.name + `_devp->buf_size-p;
	}

	if(copy_from_user(` + g.name + `_devp->buf+p,user_buf,count)) {
		return -EFAULT;
	}else {
		*ppos+=count;
		ret=count;
		printk(KERN_INFO"` + g.name + `: written %d bytes at %d.\n",count,p);
	}
	return ret;
}`
}
func (g *chardevgen) genLlseek() string {
	return `
	// start of llseek
	static loff_t ` + g.name + `_llseek(struct file * filp, loff_t offset, int orig)
	{
		struct ` + g.name + `_dev *` + g.name + `_devp=filp->private_data;
		loff_t ret=0;
		switch( orig) {
		case 0:
			if(offset <0) {
				ret=-EINVAL;
				break;
			}
			if((unsigned int)offset>` + g.name + `_devp->buf_size)
			{
				ret=-EINVAL;
				break;
			}
			filp->f_pos=(unsigned int)offset;
			ret=filp->f_pos;
			break;
		case 1:
			if((filp->f_pos+offset) >` + g.name + `_devp->buf_size)
			{
				ret=-EINVAL;
				break;
			}
			if((filp->f_pos+offset)<0)
			{
				ret=-EINVAL;
				break;
			}
			filp->f_pos+=offset;
			ret=filp->f_pos;
			break;
		default:
			ret=-EINVAL;
			break;
		}
		return ret;
	}
	//	end of llseek
	`
}
func (g *chardevgen) genFops() string {
	return `
	// start of fops definition
	static struct file_operations ` + g.name + `_fops=
	{
		.owner=THIS_MODULE,
		.llseek=` + g.name + `_llseek,
		.read=` + g.name + `_read,
		.write=` + g.name + `_write,
		.open=` + g.name + `_open,
		.unlocked_ioctl=` + g.name + `_ioctl,
		.release=` + g.name + `_release,

	};
	// end of fops definition
	`
}

func (g *chardevgen) genSetup() string {
	return `
	static void ` + g.name + `_setup_cdev(struct ` + g.name + `_dev *` + g.name + `_dev, int index)
	{
		int err,devno =MKDEV(` + g.name + `_major,index);

		` + g.name + `_dev->buf_size=BUF_SIZE;
		cdev_init(&` + g.name + `_dev->cdev,&` + g.name + `_fops);
		` + g.name + `_dev->cdev.owner = THIS_MODULE;
		` + g.name + `_dev->cdev.ops = &` + g.name + `_fops;
		err =cdev_add(&` + g.name + `_dev->cdev, devno,1);
		if (err)
		printk(KERN_NOTICE "Error %d adding ` + g.name + `%d",err, index);
	}
	`
}

func (g *chardevgen) genInit() string {
	return `
	int __init ` + g.name + `_init(void)
	{
		int result;
		int i;
		dev_t devno =MKDEV(` + g.name + `_major,0);

		if (` + g.name + `_major)
		{
			result = register_chrdev_region(devno,` + g.n_devs + `,"` + g.name + `");
		}
		else
		{
			result = alloc_chrdev_region(&devno, 0, ` + g.n_devs + `, "` + g.name + `");
			` + g.name + `_major =MAJOR(devno);
		}
		if (result < 0)
		return result;

		printk("chardev major:%d, number of minors:%d\n",` + g.name + `_major,` + g.n_devs + `);


		for(i=0;i<` + g.n_devs + `;i++) {
			` + g.name + `_setup_cdev(&` + g.name + `_devs[i],0);
		}

		classp=class_create(THIS_MODULE,"` + g.name + `");
		if( IS_ERR(classp)) {
			printk("Error registering class.\n");
			return -ENOMEM;
		}

		for(i=0;i<` + g.n_devs + `;i++) {
			` + g.name + `_devs[i].devicep=device_create(classp,NULL,MKDEV(` + g.name + `_major,i),NULL,"` + g.name + `%d",i);
		}
		return 0;
	}
	`
}

func (g *chardevgen) genExit() string {
	return `
	void  __exit ` + g.name + `_exit(void)
	{
		int i;
		for(i=0;i<` + g.n_devs + `;i++) {
			device_destroy(classp, MKDEV(` + g.name + `_major,i));
		}
		class_destroy(classp);
		for(i=0;i<` + g.n_devs + `;i++) {
			cdev_del(&` + g.name + `_devs[i].cdev);
		}
		unregister_chrdev_region(MKDEV(` + g.name + `_major, 0), ` + g.n_devs + `);
	}

	MODULE_AUTHOR("Dean Zhang");
	MODULE_LICENSE("Dual BSD/GPL");

	module_param(` + g.name + `_major,int, S_IRUGO);

	module_init(` + g.name + `_init);
	module_exit(` + g.name + `_exit);
	`
}

func (g *chardevgen) genMakefile() string {
	return `DIR=/lib/modules/$(shell uname -r)/build
SRC =$(shell pwd)
obj-m:=` + g.name + `.o
all:
	make -C $(DIR) M=$(SRC) modules
clean: 
	rm -f *.o *.ko
	`
}
