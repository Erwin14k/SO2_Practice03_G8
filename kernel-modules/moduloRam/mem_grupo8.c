//Header to use the memory info
#include <linux/mm.h> 
//Allows access to functions and structures necessary for the creation and administration of kernel modules. (Mandatory)
#include <linux/module.h>
// To use kern info
#include <linux/kernel.h>
//Header to use the init and exit macros
#include <linux/init.h>
//Header necesary to use proc_ops/file_operations
#include <linux/proc_fs.h>
/* for copy_from_user */
#include <asm/uaccess.h>	
/* Header to use the seq_file library and manage the file on  /proc*/
#include <linux/seq_file.h>


// Module info
MODULE_LICENSE("GPL");
MODULE_DESCRIPTION("Modulo RAM, Laboratorio Sistemas Operativos 2");
MODULE_AUTHOR("GRUPO8");

// escribir_archivo is a function that is responsible for writing the file.
// It performs the following tasks:
// - Retrieves system memory information using the sysinfo function.
// - Writes the total RAM on the system, free RAM, and occupied RAM to the file in JSON format.
// - Closes the file by writing a closing curly brace.
// This function returns 0.
static int escribir_archivo(struct seq_file *archivo, void *v)
{
    // struct sysinfo
    struct sysinfo info;
    si_meminfo(&info);

    //Init the file with a {
    seq_printf(archivo, "{\n");
    //The total ram on the system
    seq_printf(archivo, "\"totalram\":%lu,\n", info.totalram * info.mem_unit / 1024 / 1024);
    //The free ram on the system
    seq_printf(archivo, "\"ramlibre\":%lu,\n", info.freeram * info.mem_unit / 1024 / 1024);
    //The occupied ram on the system
    seq_printf(archivo, "\"ramocupada\":%lu\n", (info.totalram -  info.freeram )* info.mem_unit / 1024 / 1024);
    //Closing the file with a }
    seq_printf(archivo, "}\n");

    return 0;
}

// al_abrir is a function that is executed every time a cat command is used on the file.
// It performs the following tasks:
// - It opens the file for reading using the single_open function and passes escribir_archivo as the callback function.
// - Returns the result of single_open.
// This function returns the result of single_open.
static int al_abrir(struct inode *inode, struct file *file)
{
    return single_open(file, escribir_archivo, NULL);
}

// operaciones is a structure of type struct file_operations that defines the file operations for the module.
// If the kernel version is lower than 5.6, the structure is initialized with the following operations:
// - The .open field is set to the al_abrir function, which is called when the file is opened.
// - The .read field is set to the seq_read function, which is called when the file is read.
// Note that other file operations can be added to this structure if needed.
static struct file_operations operaciones =
{
    .open = al_abrir,
    .read = seq_read
};

// _insert is a function that is executed when the module is loaded using insmod.
// It performs the following tasks:
// - Creates a proc entry named "mem_grupo8" using the proc_create function. It associates the file operations defined in the operaciones structure with this entry.
// - Prints a message to the kernel log using the printk function.
// - Returns 0 to indicate success.
static int _insert(void)
{
    proc_create("mem_grupo8", 0, NULL, &operaciones);
    printk(KERN_INFO "Hola mundo, somos el grupo 8 y este es el monitor de memoria\n");
    return 0;
}

// _remove is a function that is executed when the module is unloaded using rmmod.
// It performs the following tasks:
// - Removes the proc entry named "mem_grupo8" using the remove_proc_entry function.
// - Prints a message to the kernel log using the printk function.
// - Does not return a value (void).
static void _remove(void)
{
    remove_proc_entry("mem_grupo8", NULL);
    printk(KERN_INFO "Sayonara mundo, somos el grupo 8 y este fue el monitor de memoria\n");
}

module_init(_insert);
module_exit(_remove);