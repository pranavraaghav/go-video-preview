# Go Video Preview 

A golang REST Server that takes a video file and generates preview images for every n seconds. 

## How to run (Docker)

Build the docker container 
```bash
docker build -t go-video-preview:latest .
```

Run the docker container
```bash
docker run -p 3000:3000 go-video-preview:latest
```

