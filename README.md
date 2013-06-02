Jukebox
=======
A minimal web-based MPlayer controller written in Go.

Can be [cross-compiled](https://github.com/davecheney/golang-crosscompile) for Raspberry Pi, Dockstar etc.

For example, on your workstation:

```
GOARM=5 go-linux-arm get github.com/mrw34/jukebox
scp $GOPATH/bin/linux_arm/jukebox raspberrypi:
```

Then on your Pi, assuming you have MPlayer installed and your music is /mnt/music/artist/album/*.mp3:

```
~/jukebox -root /mnt/music -port 8000
```

You can then visit:

```
http://raspberrypi:8000/
```
