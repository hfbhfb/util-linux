

# 逻辑梳理(比较笨的办法)
# 1、代理pod :  获取 /proc 下 【所有进程id 与pod对应起来：ip】      》》      [   为 nsenter 使用 进程命名空间   ]
# 2、代理pod :  运行相应命令（ 可以进入到 进程网络命名空间，但无法进入到cgroup,这足够了）
#                                           f6
make docker

/mytools1/proxy-prod
dive swr.cn-north-4.myhuaweicloud.com/hfbbg4/proxy-prod:v0.1



# test  --------------------------------------------------------------------------------
# test  --------------------------------------------------------------------------------
# test  --------------------------------------------------------------------------------
# test  --------------------------------------------------------------------------------
ps -e -o pid,ppid,cmd --forest | grep -C 5 "sleep 12"
ps aux --forest |grep -C 5 "sleep 12"


nsenter -t 7161 -n -p ps aux
nsenter -t 7161 -n -p -m  ls /bin
nsenter -t 7161 -n -p -m  which ls 
nsenter -t 7161 -n -p -m  which pwd 
nsenter -t 7161 -n -p -m  pwd 
nsenter -t 7161 -n -p -m   ls -lh /usr/bin/ls
nsenter -t 7161 -n -p -m --wd=/usr/bin ./ls pwd
nsenter -t 7161 -n -p -m --wd=/usr/bin ./pwd
nsenter -t 7161 -n -p -m --root=/usr/bin ./ls




cat /sys/fs/cgroup/cpuset/cpuset.mem_exclusive
nsenter -t 7161 -n -p cat /sys/fs/cgroup/cpuset/cpuset.mem_exclusive
ls /sys/fs/cgroup/cpuset

nsenter -t 7161 -n -C bash



挂载命名空间的相关文件：
/proc/[pid]/mounts：列出指定进程的挂载点。
/proc/[pid]/mountinfo：提供更详细的挂载信息。
/proc/[pid]/ns/mnt：指定进程的挂载命名空间文件。


-t, --target<pid>：指定目标进程的进程ID，表示要进入其所属的命名空间。
-a, --all：进入所有命名空间，而不仅仅是进程命名空间。
-m, --mount：进入挂载命名空间。
-u, --uts：进入UTS命名空间（主机名和域名）。
-i, --ipc：进入IPC命名空间（进程间通信）。
-n, --net：进入网络命名空间。
-p, --pid：进入PID命名空间。
-C, --cgroup<file>：进入指定的cgroup命名空间。
例如，如果想进入系统的全局命名空间，允许用户在该命名空间内执行操作，以实现对系统资源的管理或调试，可使用以下命令：
nsenter --target <PID> --mount --uts --ipc --net --pid -- bash



Run a program with namespaces of other processes.

Options:
 -t, --target <pid>     target process to get namespaces from
 -m, --mount[=<file>]   enter mount namespace
 -u, --uts[=<file>]     enter UTS namespace (hostname etc)
 -i, --ipc[=<file>]     enter System V IPC namespace
 -n, --net[=<file>]     enter network namespace
 -p, --pid[=<file>]     enter pid namespace
 -U, --user[=<file>]    enter user namespace
 -S, --setuid <uid>     set uid in entered namespace
 -G, --setgid <gid>     set gid in entered namespace
     --preserve-credentials do not touch uids or gids
 -r, --root[=<dir>]     set the root directory
 -w, --wd[=<dir>]       set the working directory
 -F, --no-fork          do not fork before exec'ing <program>
 -Z, --follow-context   set SELinux context according to --target PID

 -h, --help     display this help and exit
 -V, --version  output version information and exit




