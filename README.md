[![Go Report Card](https://goreportcard.com/badge/github.com/ipcjk/mlxsh)](https://goreportcard.com/report/github.com/ipcjk/mlxsh)
[![Docker Repository on Quay](https://quay.io/repository/ipcjk/mlxsh/status?token=ee40e2b7-dc86-4ee2-8269-184501ca09a2 "Docker Repository on Quay")](https://quay.io/repository/ipcjk/mlxsh)
[![Build Status](https://travis-ci.org/ipcjk/mlxsh.svg?branch=master)](https://travis-ci.org/ipcjk/mlxsh)

# mlxsh

mlxsh is the missing power command-line that enables you to enter configuration changes to groups of Brocade / Extreme Networks Netiron devices (
MLX, MLXE, CER, XMR), other Ironware style devices like Turboiron, ICX and also SLX/VDX switches and new (since 0.3) also for Juniper switches.

## Junos Support
In version 0.3 I have added basic JunOS support. To use your device as Juniper-router you need to add "DeviceType: juniper" to your YAML-configuration file. 

## modes
 
 mlxsh has two different modes
 
 - exec mode 
 - config mode
 
 and two different sources of origin hosts:
 
 - cli  (command line arguments)
 - yaml - file
 
 In cli origin source mlxsh reads all params for a **single router** directly from the command line arguments. It is good for one-shots, one-liners or testing connectivity.
  
 In YAML mode mlxsh reads **records of routers** from a YAML-file. Therefore it is possible to work on groups of routers by calling out user-defined labels.
  It also allows to overwrite certain params from the command line to calling out  scripts or config-commands without re-editing the YAML configuration.
   
    
### exec vs config mode 
   If you pass a file or a command with the -script command, the router will drop into the
   _exec or privileged mode_. If you pass in the file with the -config parameter, the router will be inserting configuration into the devices _configuration mode_.
    
   E.g. if you want to run commands in the executable mode, be sure to set the script-parameter at start, else it will drop into config mode:
    
   ```bash
   crontab -l
    0 4 * * *  mlxsh -hostname rt1 -password nocpassword -username noc -enable enablepassword\
     -script "show ip bgp summary"  
   ```

 
## cli source examples

For example, if you want to quickly commit the cloudflare.txt ip prefix lists, you can enter the command:

```bash 
mlxsh -enable enablepassword  -hostname rt1 -password nocpassword -username noc \
 -config cloudflare.txt 
```

Also this is very handy for daily maintenance tasks or cronjobs:

```bash
crontab -l
 0 4 * * *  mlxsh -hostname rt1 -password nocpassword -username noc -enable enablepassword\
  -script /home/noc/brocade/shutdown_bgp
```


## YAML source examples

Routers can be configured in a YAML file and it is possible to execute commands or configuration settings
on a group of routers by calling user-defined labels or connect to a single router by setting the hostname parameter.
   
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
   ```

Now from the command line it is only necessary to specify a hostname for the connection to your favourite router. If there is no script set (ScriptFile) for configuration or executable mode set,
you can still give this parameters from the command line. Lets run a command for rt2:
 
 ```bash
mlxsh -hostname rt2 -script "show ip cache" 
2017/06/25 15:01:32 sh ip cache
Total IP and IPVPN Cache Entry Usage on LPs:
 Module        Host    Network       Free      Total
      1          24     640960     559016    1200000
2017/06/25 15:01:32 sh ipv6 cache
Total IPv6 and IPv6 VPN Cache Entry Usage on LPs:
 Module        Host    Network       Free      Total
      1           7      38339      81654     120000
 ```
 
 If you want to execute the command on several routers, you can call a label, that
 is user-defined in the YAML-file. For example to read the ip cache command from a file and execute it on any router
 that is located in the location in Frankfurt you can enter the command line:
  
  ```bash
 mlxsh -label "location=frankfurt" -script scripts/ip_caches 
  ```
  
  If you only want to execute on all production devices in Frakfurt, you can just add a label and also set a command-one
   liner directly on the prompt: 
```bash
   mlxsh -label "location=frankfurt,environment=production" -script "show ip bgp summary"
```
 
 - chain commands 
 ```bash
    mlxsh -label "location=frankfurt,environment=production" -script "show ip bgp summary; show ip cache; show uptime"
 ```
 
 - parallel execution in background on router-groups with the -c flag, defaults to ten
 ```bash
 mlxsh -c20 -label "location=munich" -script "show ip bgp 8.8.8.8"
 ```
 
- other cool examples ro run mlxsh:
```bash
mlxsh -hostname frankfurt-rt1 -script "show uptime"
mlxsh -hostname frankfurt-rt1 -username operator -password foo -enable foo -script "show ip bgp sum"

```
- grep-able output:
```bash
mlxsh -hostname frankfurt-rt1  -script "show uptime" | grep MP
```

- label-based execution and configuration on router-groups. Great for scheduled maintenance within cron, 
reloading IX-configs at night, reload the router for testing HA, ….
```bash
mlxsh -label "location=frankfurt,type=mlx" -script 'show ip cache'
mlxsh -label "location=munich" -config scripts/bgp_neighbor
mlxsh -label "mission=DECIX" -routerdb='/home/mlxsh/mlxsh.yaml' -config /home/ixgen/decix
```

### docker

mlxsh is container ready, joerg/mlxsh is the name of the docker image available at hub.docker.com.
```bash
docker run -ti joerg/mlxsh /bin/sh
./mlxsh.linux -h
```

### full list of command line parameters
 
 Command line arguments:
 
 ```bash
 Usage of ./mlxsh:
  -c int
    	concurrent working threads \(default 20\)
  -clitype string
    	Router type \(default mlxe\)
  -config string
    	Configuration file to insert, its used as a direct command
  -debug
    	Enable debug for read / write
  -enable string
    	enable password
  -hostname string
    	Router hostname
  -i string
    	Path to a ssh private key \(in openssh2-format\) that will be used for connections 
  -label string
    	label-selection for run commands on a group of routers, e.g. 'location=munich,environment=prod'
  -password string
    	user password
  -q	quiet mode, no output except error on connecting & co
  -readtimeout duration
    	timeout for reading poll on cli select \(default 30s\)
  -routerdb string
    	Input file in yaml for username,password and host configuration if not specified on command-line \(default "mlxsh.yaml"\)
  -s	Enable strict hostkey checking for ssh connections
  -script string
    	script file to to execute, if no file is found, its used as a direct command
  -sf string
    	Path to the known-hosts-file \(in openssh2-format\) that will be used for validating hostkeys, defaults to .ssh/known_hosts 
  -speedmode
    	Enable speed mode write, will ignore any output from the cli while writing
  -username string
    	username
  -version
    	prints version and exit
  -writetimeout duration
    	timeout to stall after a write to cli
 ```
 
 ### full list of possible host parameters in YAML
 
 - ConfigFile: File with configuration statements
 - DeviceType: Type of Device, possible: MLX,CER,MLXE,XMR,IRON,TurboIron,ICX,FCS,SLX,VDX,MX,EX 
 - EnablePassword: Password that may be needed for privileged mode
 - ExecMode (internal): True or false, if its necessary to execute commands or configure
 - FileName (internal): Filename with config or command statements
 - HostName: Hostname to connect to
 - KeyFile: SSH private key that is needed for auth
 - KnownHosts: SSH Hostkeys for host-auth and to prevent MitM
 - Labels: Map of labels to group devices for command execution (see example yaml-file)
 - Password: SSH password for the initial connection
 - ReadTimeout: Timeout waiting for output from the device, tune for slow devices
 - ScriptFile: File with execution statements
 - SpeedMode: true or false: wait for prompt to return after execution
 - SSHIP: IP to connect to, will overwrite Hostname if set
 - SSHPort: SSH Port to connect to, default is 22
 - StrictHostCheck: yes/no or true/false, on true/yes we will scan the known_hosts_file 
 - Username: User for the initial ssh connection
 - WriteTimeout: time to wait after a command statement, tune for slow devices 
 
