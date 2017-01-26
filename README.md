brocadecli is a tool that enables you to enter configuration changes to Brocade Netiron devices (MLX, MLXE, CER) via ssh.

For example, if you want to commit the cloudflare.txt ip prefix lists, you can enter the command:

```bash 
brocadecli.linux -enable enablepassword  -hostname rt1 -password nocpassword -username noc  -readtimeout 10s -filename cloudflare.txt
```

```bash
./brocadecli.linux -h
Usage of ./brocadecli.linux:
  -debug
    	Enable debug for read / write
  -enable string
    	enable password (default "enablepassword")
  -filename string
    	Configuration file to insert
  -hostname string
    	Router hostname (default "rt1")
  -password string
    	user password (default "password")
  -readtimeout duration
    	timeout for reading poll on cli select (default 10s)
  -username string
    	username (default "username")
  -writetimeout duration
    	timeout to stall after a write to cli
```
