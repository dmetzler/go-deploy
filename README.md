# GO Deploy

This simple Go tool allows to deploy a static single page app in a volume or in a S3 bucket. It also evaluates variables present in the source `.env` file and generates a `env-config.js` with those variables at the root of the destination folder which allows to configure the app at runtime.
The regular way to use this tool, is to use the published Docker image as the base to build your own image.




## How to build your distroless image

If for instance you have a base JS app using yarn, you can write this Dockerfile:


```Dockerfile
FROM node:alpine as builder
WORKDIR /app
COPY package.json .
COPY yarn.lock .
RUN yarn
COPY . .
RUN yarn build

FROM dmetzler/go-deploy
ENV SRC_DIR=/src
COPY --from=builder /app/build $SRC_DIR

```

and build it with:

```console
docker build -t myuser/mystaticapp .
```

## How To Deploy My App?

You have various way to use this image depending on you deployment target

### With Docker Compose

```yaml
version: '3'
services:

  html:
    image: dmetzler/static-html
    command: volume /html_dir
    environment:
      API_URL: https://jsonplaceholder.typicode.com/users
    volumes:
    - html:/html_dir/
  nginx:
    image: nginx
    ports:
      - "8080:80"
    volumes:
    - html:/usr/share/nginx/html/:ro

volumes:
  html:
```

### With Kubernetes

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: myapp
  labels:
    app: myapp
spec:
  initContainers:
  - name: html
    image: dmetzler/static-html
    imagePullPolicy: Always
    command:
    - "/entrypoint.sh"
    - "vol"
    - "/html_dir"
    env:
    - name: API_URL
      value: https://jsonplaceholder.typicode.com/users
    volumeMounts:
    - name: htmldir
      mountPath: /html_dir

  containers:
  - name: nginx
    image: twalter/openshift-nginx:stable
    ports:
    - containerPort: 8081
    volumeMounts:
    - name: htmldir
      mountPath: /usr/share/nginx/html/
      readOnly: true
  volumes:
    - name: htmldir
      emptyDir: {}
```

### On S3

```console
# docker run --rm \
     -e API_URL=https://jsonplaceholder.typicode.com/users \
     -e AWS_ACCESS_KEY_ID \
     -e AWS_SECRET_ACCESS_KEY \
     -e AWS_DEFAULT_REGION \
     -e AWS_SESSION_TOKEN \
     -it dmetzler/static-html s3 s3://mysamplestaticapp.com

# open http://mysamplestaticapp.com.s3-website-us-east-1.amazonaws.com
```


## Environment Variables

All environment variables references in the `.env` file are evaluated at runtime and rendered in a `env-config.js` file that can be included in index.html.

The environment variable values are then available under the window._env_ variable. For instance:

```
API_URL=https://jsonplaceholder.typicode.com/users
DEFAULT_VAR=default
```
would give the JS file:

```javascript
window._env_ = {
  API_URL: "https://jsonplaceholder.typicode.com/users",
  DEFAULT_VAR: "default",
}
```

and if the `API_URL` variable is set to something else (`https://jsonplaceholder.typicode.com/photos`) when running `go-deploy`, then it would generate:
```javascript
window._env_ = {
  API_URL: "https://jsonplaceholder.typicode.com/photos",
  DEFAULT_VAR: "default",
}
```
