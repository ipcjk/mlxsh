# BrocadeCLI

## Overview
Brocadecli is a tool that enables you to enter configuration changes to Brocade Netiron devices (
Brocade MLX, Brocade MLXE, Brocade CER-series) via Secure Shell (ssh).

## modi 

Brocadecli can takes a list of full parameters on the command line or can read configuration parameters from a yaml-style configuration file.
 
### cli mode

For example, if you want to quickly commit the cloudflare.txt ip prefix lists, you can enter the command:

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
  -configfile string
    	Input file in yaml for username,password and host configuration if not specified on command-line (default "broconfig.yaml")
  -debug
    	Enable debug for read / write
  -enable string
    	enable password (default "enablepassword")
  -execmode
    	Exec commands / input from filename instead of paste configuration
  -filename string
    	Configuration file to insert
  -hostname string
    	Router hostname (default "dus-rt1.premium-datacenter.eu")
  -logdir string
    	Record session into logDir, automatically gzip
  -outputfile string
    	Output file, else stdout
  -password string
    	user password (default "password")
  -readtimeout duration
    	timeout for reading poll on cli select (default 15s)
  -speedmode
    	Enable speed mode write, will ignore any output from the cli while writing
  -username string
    	username (default "username")
  -writetimeout duration
    	timeout to stall after a write to cli
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
       FileName: brocade_scripts/bgp_sum
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
brocadecli -hostname rt2 -filename brocade_scripts/ip_caches -execmode 
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

Run brocadecli with the help of docker:
```bash
docker run -ti joerg/brocadecli /bin/sh
./brocadecli.linux -h
```