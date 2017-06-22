### BrocadeCLI

brocadecli is a tool that enables you to enter configuration changes to Brocade Netiron devices (MLX, MLXE, CER) via ssh.

For example, if you want to commit the cloudflare.txt ip prefix lists, you can enter the command:

```bash 
brocadecli.linux -enable enablepassword  -hostname rt1 -password nocpassword -username noc \
-readtimeout 10s -filename cloudflare.txt -speedmode
```

Also it is a handy tool for daily maintenance tasks or cronjobs:

```bash
crontab -l
 0 4 * * *  brocadecli.linux -hostname rt1 -password nocpassword -username noc -enable enablepassword\
  -filename /home/noc/brocade/shutdown_bgp
```

If you want to run commands in the executable mode, be sure to set the parameter at start, else the tool\
 will drop into config mode:
 
 
```bash
crontab -l
 0 4 * * *  brocadecli.linux -hostname rt1 -password nocpassword -username noc -enable enablepassword\
  -filename /home/noc/brocade_scripts/bgp_sum -execmode 
```

Command line arguments:

```bash
Usage of ./brocadecli:
  -debug
    	Enable debug for read / write
  -enable string
    	enable password (default "enablepassword")
  -execmode
    	Exec commands / input from filename instead of paste configuration
  -filename string
    	Configuration file to insert
  -hostname string
    	Router hostname (default "rt1")
  -password string
    	user password (default "password")
  -readtimeout duration
    	timeout for reading poll on cli select (default 15s)
  -speedmode
    	Enable speed mode write, will ignore any output/feedback from the cli while writing
  -username string
    	username (default "username")
  -writetimeout duration
    	timeout to stall after a write to cli
```

Run brocadecli with the help of docker:
```bash
docker run -ti joerg/brocadecli /bin/sh
./brocadecli.linux -h
```