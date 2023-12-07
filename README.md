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

## Making requests 
### Upload a file 
```bash
curl --location 'http://localhost:3000/upload' \
--form 'file=@"<video_file_location>"'
```
```json
{
    "filename": "<filename_from_server>"
}
```

### Get the preview images 
```bash
curl --location 'http://localhost:3000/images?filename=<filename_from_server>&height=480&width=640&interval=2'
```
This should allow you to download a zip file with all the preview images.

