# Installation
Then compile it with:
```
go build -o ansem 
```

![](https://github.com/PeterParser/Ansem/workflows/Go/badge.svg)

# Configuration
The configuration is stored in yaml file with the following (mandatory) attributes:
* **exploits_dir**: Path to the directory that contains the exploits.
* **tick**: The tick of the gameserver in minutes.
* **team_file**: Path to a plain text file with the ip of the teams.
* **gameserver**: IP address and port of the gameserver in the following format "IP:port".
* **workers**: The number of workers
* **submission_type**: The type of submission required by the game server, at the moment only "TCP" is supported.
* **flag_regex**: The regex of the flag to be submitted.
* **flag_accepted**: The string sent by the server that tells us if the flag was right or not.
* **token**: The token which identifies your team to submit to the checking system.
Example:
```yaml
exploits_dir: "/path/to/exploits"
tick: 5
team_file : "/path/to/teams"
gameserver: "127.0.0.1:31337"
workers: 8
submission_type: "TCP"
flag_regex: "\\d"
flag_accepted: "ACCEPTED"
token: "some_kind_of_token"
```
