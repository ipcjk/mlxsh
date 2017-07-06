[![Go Report Card](https://goreportcard.com/badge/github.com/ipcjk/mlxsh)](https://goreportcard.com/report/github.com/ipcjk/mlxsh)
[![Docker Repository on Quay](https://quay.io/repository/ipcjk/mlxsh/status?token=ee40e2b7-dc86-4ee2-8269-184501ca09a2 "Docker Repository on Quay")](https://quay.io/repository/ipcjk/mlxsh)

# mlxsh

mlxsh is the missing power command-line that enables you to enter configuration changes to groups of Brocade / Extreme Networks Netiron devices (
MLX, MLXE, CER, XMR) via Secure Shell (ssh).

## modi 

mlxsh can take a list of full parameters on the command line or can read configuration parameters from a yaml-style configuration file.
 
### cli mode

For example, if you want to quickly commit the cloudflare.txt ip prefix lists, you can enter the command:

```bash 
mlxsh -enable enablepassword  -hostname rt1 -password nocpassword -username noc \
 -config cloudflare.txt 
```

Also it is a handy tool for daily maintenance tasks or cronjobs:

```bash
crontab -l
 0 4 * * *  mlxsh -hostname rt1 -password nocpassword -username noc -enable enablepassword\
  -script /home/noc/brocade/shutdown_bgp
```

There is a fine and subtle difference in both commands. If you pass a file with the -script command, the router will drop into the
exec or privileged mode. If you pass in the file with the -config parameter, the router will be inserting configuration in the configuration-mode.
 
E.g. if you want to run commands in the executable mode, be sure to set the script-parameter at start, else the tool\
 will drop into config mode:
 
```bash
crontab -l
 0 4 * * *  mlxsh -hostname rt1 -password nocpassword -username noc -enable enablepassword\
  -script /home/noc/brocade_scripts/bgp_sum  
```

Command line arguments:

```bash
Usage of ./mlxsh:
  -c int
    	concurrent working threads / connections to the routers default 2
  -config string
    	Configuration file to insert, its used as a direct command
  -debug
    	Enable debug for read / write
  -enable string
    	enable password
  -hostname string
    	Router hostname
  -label string
    	label-selection for run commands on a group of routers, e.g. 'location=munich,environment=prod'
  -password string
    	user password
  -readtimeout duration
    	timeout for reading poll on cli select default 15s
  -routerdb string
    	Input file in yaml for username,password and host configuration if not specified on command-line default "mlxsh.yaml"
  -script string
    	script file to to execute, if no file is found, its used as a direct command
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

When mlx is reading a yaml-configuration file, it will overwrite any given parameter from the commandline. You can specify
  the configuration file with the configfile-parameter, by default it is looking for a file named config.yaml in the current working directory.
   
   A typical config.yaml is included in the distribution file and could look like this:
   ```yaml
- Hostname: rt2
  Username: mucuser
  Password: mucpass
  SSHPort: 22
  EnablePassword: enablePass
  StrictHostCheck: False
  SpeedMode: False
  ScriptFile: scripts/bgp_sum
  Labels:
    location: dus
    environment: stage
    type: cer

   ```

Now from the command line it is only necessary to specify a hostname for the connection to your favourite router. If there is no script set (FileName) for configuration or executable mode set,
you can still give this parameters from the command line. Lets run a command for rt2:
 
 ```bash
mlxsh -hostname rt2 -script scripts/ip_caches 
2017/06/25 15:01:32 sh ip cache
Total IP and IPVPN Cache Entry Usage on LPs:
 Module        Host    Network       Free      Total
      1          24     640960     559016    1200000
2017/06/25 15:01:32 sh ipv6 cache
Total IPv6 and IPv6 VPN Cache Entry Usage on LPs:
 Module        Host    Network       Free      Total
      1           7      38339      81654     120000
 ```
 
 If you want to execute the command on several routers, you can send a label, that
 is user-defined in the yaml configuration. For example to run the ip cache command on any router
 that is located in the location in Frankfurt you can enter the command line:
  
  ```bash
 mlxsh -label "location=frankfurt" -script scripts/ip_caches 
  ```
  
  If you only want to execute on any production device in Frakfurt, you can just add a label and also set a command-one liner directly on the prompt: 
```bash
   mlxsh -label "location=frankfurt,environment=production" -script "show ip bgp summary"
```
 
 Great!

Other cool examples ro run mlxsh:
```bash
mlxsh -hostname frankfurt-rt1 -script "show uptime"
mlxsh -hostname frankfurt-rt1 -username operator -password foo -enable foo -script "show ip bgp sum"
```

- grep-able output:

```bash

mlxsh -hostname frankfurt-rt1  -script "show uptime" | grep MP
```

- label-based execution and configuration on router-groups, also great for scheduled maintenance within cron, reloading IX-configs at night, reload the router for testing HA, ….

```bash
mlxsh -label "location=frankfurt,type=mlx" -script 'show ip cache'
mlxsh -label "location=munich" -config scripts/bgp_neighbor
mlxsh -label "mission=DECIX" -routerdb='/home/mlxsh/mlxsh.yaml' -config /home/ixgen/decix
```

- parallel execution in background on router-groups with the -c flag, defaults to two
```bash
mlxsh -c10 -label "location=munich" -script "show ip bgp 8.8.8.8"
````



### docker

Run mlxsh with the help of docker, joerg/mlxsh is the name of the docker image available at hub.docker.com.
```bash
docker run -ti joerg/mlxsh /bin/sh
./mlxsh.linux -h
```
