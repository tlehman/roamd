# `roamd(1)` roam dump (and roam daemon)

Command to quickly dump thoughts into Roam for later editing and review.

The `roamd` command appends to a queue. The `roamd.service' uses systemd to 
keep a running process that watches the queue and uploads the contents to 
roam in the background.

## How to use

Set the environment variables:
- `ROAM_API_GRAPH`
- `ROAM_API_EMAIL`
- `ROAM_API_PASSWORD`

Then, assume it's 02:11pm, type:
```
roamd "the linux16 keyword is used in grub.cfg for legacy MBR boot configurations"
```

The `roamd` command will append the string "14:11 the linux16 keyword is used in grub.cfg for legacy MBR boot configurations" to the queue.

A background process will check the '~/.roamqueue', and if there are any entries, it will pass them to the following command:

```
dapper roam-api create "12:11 'the linux16 keyword is used in grub.cfg for legacy MBR boot configurations'"
```

This `roam-api create` command is slow (~10-15s) when I tried. My solution is to 
queue it up and then post it in the background.

# More links
For more about [dapper](https://tobilehman.com/posts/dapper-build-consistent/) 

And the [private Roam Research API](https://github.com/artpi/roam-research-private-api) 

