#include <linux/module.h>
#include <linux/device.h>
#include <linux/uaccess.h>
#include <linux/init.h>
#include <linux/fs.h>
#include <linux/cdev.h>
#include <linux/errno.h>
	
	// start of definitions
	#define BUF_SIZE 100
	struct testlllmodule_dev {
		char buf[BUF_SIZE];
		struct cdev cdev;
		int buf_size;
		struct device *devicep;
	};
	unsigned int  testlllmodule_major;
	struct testlllmodule_dev testlllmodule_devs[4];
	struct class * classp=NULL;
	// end of definitions
	static int testlllmodule_open(struct inode * inode, struct file *filp)
	{
		int minor=MINOR(inode->i_rdev);
		if (minor > 4-1) {
			return -ENODEV;
		}

		filp->private_data=&testlllmodule_devs[4];
		return 0;
	}
	
	static int testlllmodule_release(struct inode * inode, struct file *filp)
	{
		return 0;
	}
	static int testlllmodule_ioctl(struct file *filp, unsigned int cmd, unsigned long arg)
	{
	struct testlllmodule_dev * testlllmodule_devp = filp->private_data;
	switch(cmd)
	{
	default:
		return -EINVAL;
	}
	return 0;
}
static ssize_t testlllmodule_read(struct file *filp, char __user *user_buf,size_t size, loff_t *ppos)
{
	int ret=0;
	int count=size;
	unsigned long p=*ppos;
	struct testlllmodule_dev *testlllmodule_devp=filp->private_data;

	if(p>=testlllmodule_devp->buf_size){
		return count? -ENXIO:0;
	}
	if( count>testlllmodule_devp->buf_size-p) {
		count=testlllmodule_devp->buf_size-p;
	}
	if( copy_to_user(user_buf,(void *)(testlllmodule_devp->buf+p),count))
	{
		ret=-EFAULT;
	} else {
		*ppos+=count;
		return count;
		printk(KERN_INFO"testlllmodule: read %d bytes form %d\n",count,p);
	}
	return ret;
}
static ssize_t testlllmodule_write(struct file *filp, const char __user *user_buf, size_t size, loff_t *ppos)
{
	unsigned long p =*ppos;
	unsigned int count=size;
	int ret=0;
	struct testlllmodule_dev *testlllmodule_devp=filp->private_data;

	if(p>=testlllmodule_devp->buf_size){
		return count? -ENXIO:0;
	}
	if( count>testlllmodule_devp->buf_size-p) {
		count=testlllmodule_devp->buf_size-p;
	}

	if(copy_from_user(testlllmodule_devp->buf+p,user_buf,count)) {
		return -EFAULT;
	}else {
		*ppos+=count;
		ret=count;
		printk(KERN_INFO"testlllmodule: written %d bytes at %d.\n",count,p);
	}
	return ret;
}
	// start of llseek
	static loff_t testlllmodule_llseek(struct file * filp, loff_t offset, int orig)
	{
		struct testlllmodule_dev *testlllmodule_devp=filp->private_data;
		loff_t ret=0;
		switch( orig) {
		case 0:
			if(offset <0) {
				ret=-EINVAL;
				break;
			}
			if((unsigned int)offset>testlllmodule_devp->buf_size)
			{
				ret=-EINVAL;
				break;
			}
			filp->f_pos=(unsigned int)offset;
			ret=filp->f_pos;
			break;
		case 1:
			if((filp->f_pos+offset) >testlllmodule_devp->buf_size)
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
	
	// start of fops definition
	static struct file_operations testlllmodule_fops=
	{
		.owner=THIS_MODULE,
		.llseek=testlllmodule_llseek,
		.read=testlllmodule_read,
		.write=testlllmodule_write,
		.open=testlllmodule_open,
		.unlocked_ioctl=testlllmodule_ioctl,
		.release=testlllmodule_release,

	};
	// end of fops definition
	
	static void testlllmodule_setup_cdev(struct testlllmodule_dev *testlllmodule_dev, int index)
	{
		int err,devno =MKDEV(testlllmodule_major,index);

		testlllmodule_dev->buf_size=BUF_SIZE;
		cdev_init(&testlllmodule_dev->cdev,&testlllmodule_fops);
		testlllmodule_dev->cdev.owner = THIS_MODULE;
		testlllmodule_dev->cdev.ops = &testlllmodule_fops;
		err =cdev_add(&testlllmodule_dev->cdev, devno,1);
		if (err)
		printk(KERN_NOTICE "Error %d adding testlllmodule%d",err, index);
	}
	
	int __init testlllmodule_init(void)
	{
		int result;
		int i;
		dev_t devno =MKDEV(testlllmodule_major,0);

		if (testlllmodule_major)
		{
			result = register_chrdev_region(devno,4,"testlllmodule");
		}
		else
		{
			result = alloc_chrdev_region(&devno, 0, 4, "testlllmodule");
			testlllmodule_major =MAJOR(devno);
		}
		if (result < 0)
		return result;

		printk("chardev major:%d, number of minors:%d\n",testlllmodule_major,4);


		for(i=0;i<4;i++) {
			testlllmodule_setup_cdev(&testlllmodule_devs[i],0);
		}

		classp=class_create(THIS_MODULE,"testlllmodule");
		if( IS_ERR(classp)) {
			printk("Error registering class.\n");
			return -ENOMEM;
		}

		for(i=0;i<4;i++) {
			testlllmodule_devs[i].devicep=device_create(classp,NULL,MKDEV(testlllmodule_major,i),NULL,"testlllmodule%d",i);
		}
		return 0;
	}
	
	void  __exit testlllmodule_exit(void)
	{
		int i;
		for(i=0;i<4;i++) {
			device_destroy(classp, MKDEV(testlllmodule_major,i));
		}
		class_destroy(classp);
		for(i=0;i<4;i++) {
			cdev_del(&testlllmodule_devs[i].cdev);
		}
		unregister_chrdev_region(MKDEV(testlllmodule_major, 0), 4);
	}

	MODULE_AUTHOR("Dean Zhang");
	MODULE_LICENSE("Dual BSD/GPL");

	module_param(testlllmodule_major,int, S_IRUGO);

	module_init(testlllmodule_init);
	module_exit(testlllmodule_exit);
	