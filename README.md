# awsprofile
Small and simple tool to set the current profile values for aws many CLI's

![awsprofile](.github/Peek.gif)
## INSTALL
execute:
```
go build 
```
or get the binary from releases

copy awsprofile to somewhere in your $PATH
in my case is $HOME/.local/bin
then to change the active profile
```
$(awsprofile)
aws sts get-caller-identity
```

and you're ready !!!

## Contribute
All contributions are welcome open an issue or make a pr!

## License
This code is released under MIT license
