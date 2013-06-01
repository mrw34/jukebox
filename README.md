Jukebox
=======
A minimal web-based MPlayer controller written in Go.

Can be [cross-compiled](https://github.com/davecheney/golang-crosscompile) for Raspberry Pi, Dockstar etc.

For example, on your workstation:

```
GOARM=5 go-linux-arm get github.com/mrw34/jukebox
scp $GOPATH/bin/linux_arm/jukebox raspberrypi:
```

Then on your Pi, assuming your music is organised as /mnt/music/artist/album/*.mp3:

```
mkfifo /tmp/mplayer
mplayer -really-quiet -cache 64 -slave -input file=/tmp/mplayer -idle &
~/jukebox -root /mnt/music -port 8000
```

And back on your workstation:

```
google-chrome http://raspberrypi:8000/
```
