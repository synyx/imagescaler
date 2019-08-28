# Message Driven Image Scaler
The imagescaler is a program that uses Go's native image processing capabilities to scale down images in an object storage (currently min.io) system.

## Motivation
Traditionally image processing is done from within a given program using native bindings to the overly popular ImageMagick (or similar programs). This has drawbacks:

* Operational complexity: The host or container needs to provide a recent version of the image processing tool. The lifecycle of the tool needs to managed, too.

* Synchronous calls of potentially expensive operations.

* ...

## Flow

The imagescaler consumes image update message, and reads them from the provided URL. This URL is not required to reside in the configured object storage, BTW. There are 2 different target scales:
* WEB for usage on web sites
* THUMBNAIL for icons or thumbnails (also on websites)

The scaled versions are stored (put) to a configured object storage (currently min.io) and new image update events for the scaled versions and their URLs are published.

## Message Format
Currently consuming and emitting messages formatted like this:

```
{
"userUUID": "schnubbli",
"imageUUID": "yolo",
"url": "http://localhost:9000/testbucket/vivian_robin.jpg",
"imageScale": "ORIGINAL"
}
```

## Configuration

tbw.
