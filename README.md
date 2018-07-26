### Why?

There are not a lot of options in Kubernetes cluster to pass small blobs of data from automation task creator to the Kubernetes job. Environment variables are not meant for that, nor configmaps.

### API

`POST /` or `PUT /name` will create a blob and return `Location: /name-or-random`. `?ttl=<seconds>` will set lifetime of the blob, default to 300 (5min). `GET /name` is to retrieve the blob.
