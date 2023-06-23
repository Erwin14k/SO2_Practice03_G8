//Allows access to functions and structures necessary for the creation and administration of kernel modules. (Mandatory)
#include <linux/module.h>
//Header para usar KERN_INFO
#include <linux/kernel.h>
//Header para las funciones module_init y module_exit
#include <linux/init.h>
//Header necesary to use proc_ops/file_operations
#include <linux/proc_fs.h>
// Header para usar la lib seq_file y manejar el archivo en /proc*
#include <linux/seq_file.h>
//Header to use the memory info pages
#include <linux/mm.h>
// Header to use the task_struct
#include <linux/sched.h> 
#include <linux/sched/signal.h> //to recognize for_each_procces 
// Use for filp_open and filp_close
#include <linux/fs.h> 
// Reading process info 0_RDONLY
#include <linux/stat.h> 
// library for using string functions like mmemset ( to clear the buffer)
#include <linux/string.h>  
// functionalities for secure copi for copy_from_user (user space to kernel space)
#include <linux/uaccess.h> 
// Information about memory / for si_meminfo
#include <linux/sysinfo.h> 



MODULE_LICENSE("GPL");
MODULE_DESCRIPTION("Modulo CPU, Laboratorio Sistemas Operativos 2");
MODULE_AUTHOR("GRUPO 8");

// calcular_porcentaje_cpu_total calculates the total CPU usage percentage.
// It reads the CPU times from the "/proc/stat" file, calculates the total CPU time,
// and then calculates the CPU usage percentage based on the idle time.
// It returns the CPU usage percentage as an integer.
static int calcular_porcentaje_cpu_total(void)
{
    // create a file pointer
    struct file *archivo;
    char lectura[256]; // buffer
    //char *etiqueta;
    int usuario, niced, sistema, idle, iowait, irq, suaveirq, steal, guest, guest_nice; // varibles to store the CPU times
    int total;
    int porcentaje;

    archivo = filp_open("/proc/stat", O_RDONLY, 0); // open the file for reading
    if (IS_ERR(archivo)) {
        printk(KERN_ALERT "No se pudo abrir el archivo /proc/stat\n");
        return -1;
    }
    
    memset(lectura, 0, sizeof(lectura)); // clear the buffer
    kernel_read(archivo, lectura, sizeof(lectura) - 1, &archivo->f_pos); // read the file

    sscanf(lectura, "cpu %d %d %d %d %d %d %d %d %d %d",
           &usuario, &niced, &sistema, &idle, &iowait, &irq, &suaveirq, &steal, &guest, &guest_nice); // storte the CPU times in the variables

    total = usuario + niced + sistema + idle + iowait + irq + suaveirq + steal + guest + guest_nice; // calculate the total CPU time

    
    porcentaje = 100 - (idle * 100 / total); // calculate the CPU usage percentage
    filp_close(archivo, NULL); //   close the file
    
    return porcentaje;
}

// obtain_state is a function that returns the state of a process based on the state number.
// It takes the state number as input and returns the corresponding state as a string.
static const char* obtain_state(int estado)
{
    const char* estado_str;

    switch (estado) {
        case 0:
            estado_str = "ejecucion";
            break;
        case 1:
        case 1026:
            estado_str = "suspendido";
            break;
        case 128:
            estado_str = "detenido";
            break;
        case 260:
            estado_str = "zombie";
            break;
        default:
            estado_str = "desconocido";
            break;
    }

    return estado_str;
}

// escribir_archivo is a function that writes information to a file.
// It takes a file pointer and a void pointer as input and returns an integer.
static int escribir_archivo(struct seq_file *archivo, void *v)
{
    int porcentaje;
    //seq_printf(archivo, "================\n");


    porcentaje = calcular_porcentaje_cpu_total();
    
    if (porcentaje == -1) {
        seq_printf(archivo, "No se pudo calcular el porcentaje de CPU total\n");
    } else {
        //seq_printf(archivo, "Porcentaje de CPU total: %d\n", porcentaje);

        seq_printf(archivo, "{\n");
        seq_printf(archivo, "\"totalcpu\":%d,\n", porcentaje);
        seq_printf(archivo, "\"tasks\":\n");

        int count_running = 0, count_sleeping = 0, count_stopped = 0, count_zombie = 0, count_total = 0;
        struct task_struct* cpu;
        int ram, separator, childseparator;
        separator = 0;
        childseparator = 0;
        seq_printf(archivo, "[");
        for_each_process(cpu){
            if(separator){
                seq_printf(archivo, ",");
            }
            seq_printf(archivo, "{\"pid\":");
            seq_printf(archivo, "%d", cpu->pid);
            seq_printf(archivo, ",\"nombre\":");
            seq_printf(archivo, "\"%s\"", cpu->comm);
            seq_printf(archivo, ",\"usuario\": \"");
            seq_printf(archivo, "%d", cpu->real_cred->uid);
            seq_printf(archivo, "\",\"estado\": \"");
            seq_printf(archivo, "%s", obtain_state(cpu->state));
	    seq_printf(archivo, "\"");
            if (cpu->mm) {
                ram = (get_mm_rss(cpu->mm)<<PAGE_SHIFT)/(1024*1024); // MB
                seq_printf(archivo, ",\"ram\":");
                seq_printf(archivo, "%d", ram);
            }
            seq_printf(archivo, ",\"padre\":");
            seq_printf(archivo, "%d",  cpu->parent->pid);
            seq_printf(archivo, "}\n");
            separator = 1;

            //contar 
            count_total++;
            switch(cpu->state) {
                case TASK_RUNNING:
                    count_running++;
                    break;
                case TASK_INTERRUPTIBLE:
                case TASK_UNINTERRUPTIBLE:
                    count_sleeping++;
                    break;
                case TASK_STOPPED:
                    count_stopped++;
                    break;
                case EXIT_ZOMBIE:
                    count_zombie++;
                    break;
                default:
		    count_sleeping++;
                    break;
            }
        }

        seq_printf(archivo, "],\n");
        seq_printf(archivo, "\"running\": %d,\n", count_running);
        seq_printf(archivo, "\"sleeping\": %d,\n", count_sleeping);
        seq_printf(archivo, "\"stopped\": %d,\n", count_stopped);
        seq_printf(archivo, "\"zombie\": %d,\n", count_zombie);
        seq_printf(archivo, "\"total\": %d\n", count_total);
        seq_printf(archivo, "}\n");
    }

    return 0;
}

// al_abrir is a function that is executed every time with a cat command.
// inode - pointer to the inode structure
// file - pointer to the file structure
static int al_abrir(struct inode *inode, struct file *file)
{
    return single_open(file, escribir_archivo, NULL);
}

// operaciones is a structure of file operations.
// If the kernel version is lower than 5.6:
// - open function is set to al_abrir
// - read function is set to seq_read
static struct file_operations operaciones =
{
    .open = al_abrir,
    .read = seq_read
};

// _insert is a function that is called when the module is inserted using insmod.
// It performs the following tasks:
// - Creates a proc entry named "cpu_grupo8" with the file operations specified by operaciones.
// - Displays a message indicating the successful creation of the module and the group information.
// Returns 0 to indicate success.
static int _insert(void)
{
    proc_create("cpu_grupo8", 0, NULL, &operaciones);
    printk(KERN_INFO "Hola mundo, somos el grupo 8 y este es el monitor de CPU\n");
    return 0;
}

// _remove is a function that is called when the module is removed using rmmod.
// It performs the following tasks:
// - Removes the proc entry named "cpu_grupo8".
// - Displays a farewell message indicating the removal of the module and the group information.
// This function does not return a value.
static void _remove(void)
{
    remove_proc_entry("cpu_grupo8", NULL);
    printk(KERN_INFO "Sayonara mundo, somos el grupo 8 y este fue el monitor de CPU\n");
}

module_init(_insert);
module_exit(_remove);