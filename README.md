Jukebox
=======
A minimal web-based MPlayer controller written in Go.

Can be cross-compiled for Raspberry Pi, Dockstar etc.

For example, on your workstation:

```
GOOS=linux GOARCH=arm GOARM=5 go get github.com/mrw34/jukebox
scp $GOPATH/bin/linux_arm/jukebox raspberrypi:
```

or via Docker:

```
docker run --rm -v "$PWD":/go -e GOOS=linux -e GOARCH=arm -e GOARM=5 golang go get github.com/mrw34/jukebox
scp bin/linux_arm/jukebox raspberrypi:
```

Then on your Pi, assuming you have MPlayer installed and your music is /mnt/music/artist/album/*.mp3:

```
~/jukebox -root /mnt/music -port 8000
```

You can then visit:

```
http://raspberrypi:8000/
```
