# BrocadeCLI

## Overview
Brocadecli is a tool that enables you to enter configuration changes to Brocade Netiron devices (
Brocade MLX, Brocade MLXE, Brocade CER-series) via Secure Shell (ssh).

## modi 

Brocadecli can take a list of full parameters on the command line or can read configuration parameters from a yaml-style configuration file.
 
### cli mode

For example, if you want to quickly commit the cloudflare.txt ip prefix lists, you can enter the command:

```bash 
brocadecli.linux -enable enablepassword  -hostname rt1 -password nocpassword -username noc \
-readtimeout 10s -config cloudflare.txt -speedmode
```

Also it is a handy tool for daily maintenance tasks or cronjobs:

```bash
crontab -l
 0 4 * * *  brocadecli.linux -hostname rt1 -password nocpassword -username noc -enable enablepassword\
  -script /home/noc/brocade/shutdown_bgp
```

There is a fine and subtle difference in both commands. If you pass a file with the -script command, the router will drop into the
exec or privileged mode. If you pass in the file with the -config parameter, the router will be inserting configuration in the configuration-mode.
 
E.g. if you want to run commands in the executable mode, be sure to set the script-parameter at start, else the tool\
 will drop into config mode:
 
```bash
crontab -l
 0 4 * * *  brocadecli.linux -hostname rt1 -password nocpassword -username noc -enable enablepassword\
  -script /home/noc/brocade_scripts/bgp_sum  
```

Command line arguments:

```bash
Usage of ./brocadecli:
   -config string
     	Configuration file to insert
   -debug
     	Enable debug for read / write
   -enable string
     	enable password
   -hostname string
     	Router hostname
   -logdir string
     	Record session into logDir, automatically gzip
   -outputfile string
     	Output file, else stdout
   -password string
     	user password
   -readtimeout duration
     	timeout for reading poll on cli select (default 15s)
   -routerdb string
     	Input file in yaml for username,password and host configuration if not specified on command-line (default "broconfig.yaml")
   -script string
     	script file to to execute
   -speedmode
     	Enable speed mode write, will ignore any output from the cli while writing
   -username string
     	username
   -version
     	prints version and exit
   -writetimeout duration
     	timeout to stall after a write to cli
 exit status 2
```

### configfile mode

When brocadecli is reading a yaml-configuration file, it will overwrite any given parameter from the commandline. You can specify
  the configuration file with the configfile-parameter, by default it is looking for a file named config.yaml in the current working directory.
   
   A typical config.yaml is included in the distribution file and could look like this:
   ```yaml
   -   Hostname: rt1
       Username: noc
       Password: noc
       EnablePassword: noc
       DeviceType: CER
       KeyFile: key
       StrictHostCheck: False
       SpeedMode: False
       FileName: scripts
       ExecMode: True
   -   Hostname: rt2
       Username: mucuser
       Password: mucpass
       EnablePassword: enablePass
       DeviceType: MLX
   ```

Now from the command line it is only necessary to specify a hostname for the connection to your favourite router. If there is no script set (FileName) for configuration or executable mode set,
you can still give this parameters from the command line. Lets run a command for rt2:
 
 ```bash
brocadecli -hostname rt2 -script brocade_scripts/ip_caches 
2017/06/25 15:01:32 sh ip cache
Total IP and IPVPN Cache Entry Usage on LPs:
 Module        Host    Network       Free      Total
      1          24     640960     559016    1200000
2017/06/25 15:01:32 sh ipv6 cache
Total IPv6 and IPv6 VPN Cache Entry Usage on LPs:
 Module        Host    Network       Free      Total
      1           7      38339      81654     120000
 ```
 
 Great!

### docker

Run brocadecli with the help of docker, joerg/brocadecli is the name of the docker image available at hub.docker.com.
```bash
docker run -ti joerg/brocadecli /bin/sh
./brocadecli.linux -h
```