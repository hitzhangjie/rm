# rm-safe
rm-safe is an alternative to shell /bin/rm, which supports `pin` or `unpin` operation to protect your data.

# commands

- `rm help`, display help info
- `rm pin`, pin directories to protect them   
    `rm pin`: pin current directory   
    `rm pin [target]...`: pin targets, if target is file, pin directory of target  
    `rm pin -r [target]...`: pin targets recursively
- `rm unpin`, unpin directories to unprotect them  
    `rm unpin`: like `rm pin`  
    `rm unpin [target]...`: like `rm unpin [target]`  
    `rm unpin -r [target]...`: like `rm unpin -r [target]...`
- `rm [options] [target]...`, absolutely works like shell `/bin/rm`, why? I use exec.Command('/bin/rm', args) to fullfill the deletion task.

# installation

```bash
go install github.com/hitzhangjie/rm
```

`go install` will install `rm` to `$GOPATH/bin`, please make sure `$GOPATH/bin` appears before `/bin, /usr/bin` in `$PATH`.


