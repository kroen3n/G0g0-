<b> Building Docker image (basics)  </b>


In this tutorial we will focus on writing Golang code for creating a Docker image

For building a Docker image, it is necesarry to create a Dockerfile.

<b> Dockerfile </b>

A Dockerfile is a text document that contains specific instructions.

For our (very simple example), the Dockerfile contains the following lines:

```
root@kr03nen:/home/gog0# more Dockerfile
FROM python:3.4
RUN apt-get update -y
RUN echo "hielau!"
```
Now, for each line:

<b><i>FROM</i></b> instruction - this one is a must for every Dockerfile when building an image. This is the foundation on which the other layers will be built.
This instruction will extract from Docker repository the image type we need for our containers (in our case, we want our containers to be extended from 
an image that has python3)
<br></br>
<b><i>RUN</i></b> instruction - this will allow you to "run" commands you have passed. 
In our case you want to run an update (-y is necessary to pass, since you cannot interract directly when building the image)


Let's build our image - we shall call it "hielaunewimage". <br>
Suppose you are under same path where the Dockerfile is located (hence the dot ( . ) at the end of the command)

```
root@kr03nen:/home/gog0# docker build -t hielaunewimage .
Sending build context to Docker daemon  15.36kB
Step 1/3 : FROM python:3.4
 ---> 8c62b065252f
Step 2/3 : RUN apt-get update -y
 ---> Using cache
 ---> 26bc7a14085f
Step 3/3 : RUN echo "hielau!"
 ---> Running in f61916bfe400
hielau!
Removing intermediate container f61916bfe400
 ---> 8054fe4c907b
Successfully built 8054fe4c907b
Successfully tagged hielaunewimage:latest
```
Let's explain what is going on here a bit.

Each line of instruction represents a layer. 
This layer can be treated as an "ephemeral" container. Once the instruction has been successfully implemented in its specific container,
that specific container is removed.
Or as the Docker engine mentions in the output "Removing intermediate container" <br>

Notice the "Using cache" line at Step 2/3 ... on my machine, that means there is an existent image with dependent child images. 
For this step there was no intermediate container created.

However, if I want to create intermediate containers for every step, I can use the option "--no-cache". 
<br>
<i>
Let's remove the image, just for phun:

```
root@kr03nen:/home/gog0# docker rmi hielaunewimage
Untagged: hielaunewimage:latest
Deleted: sha256:eacb7f78036fcefadcbc5b021ddc210c31aeacb7bb436f00dcb5179a251f5c28
```
</i>

Observe the output for step 2/3 when creating the image with "--no-cache" option:

```
root@kr03nen:/home/gog0# docker build --no-cache -t hielaunewimage . 
Sending build context to Docker daemon  18.43kB
Step 1/3 : FROM python:3.4
 ---> 8c62b065252f
Step 2/3 : RUN apt-get update -y
 ---> Running in 6f8ebdf9cd1b
Ign:1 http://deb.debian.org/debian stretch InRelease
Get:2 http://security.debian.org/debian-security stretch/updates InRelease [53.0 kB]
Get:3 http://deb.debian.org/debian stretch-updates InRelease [93.6 kB]
Get:4 http://deb.debian.org/debian stretch Release [118 kB]
Get:5 http://deb.debian.org/debian stretch Release.gpg [2410 B]
Get:6 http://security.debian.org/debian-security stretch/updates/main amd64 Packages [551 kB]
Get:7 http://deb.debian.org/debian stretch-updates/main amd64 Packages [2596 B]
Get:8 http://deb.debian.org/debian stretch/main amd64 Packages [7080 kB]
Fetched 7901 kB in 3s (2375 kB/s)
Reading package lists...
Removing intermediate container 6f8ebdf9cd1b
 ---> fabe6b22ee4d
Step 3/3 : RUN echo "hielau!"
 ---> Running in 3e93a2cb1eb6
hielau!
Removing intermediate container 3e93a2cb1eb6
 ---> 0417c8d4ddc0
Successfully built 0417c8d4ddc0
Successfully tagged hielaunewimage:latest

```

<i>Consider this as a good practice, since using cache is quite insecure for containers -  this is caused by the lack of updates on that specific image.</i> 


You can check the creation of these intermediate containers by either applying below command, 
and watch (lol!) the created/removed containers for each instruction...

```
watch -n 1 'docker ps'
```

... or you can play with this small python code (no-cache included)

```
import docker
import os, subprocess
import sys, six

def BuildImage():
    dockerfile = '''
    FROM python:3.4
    RUN apt-get update -y
    RUN echo "hielau!"
    '''

    f = six.BytesIO(dockerfile.encode('utf-8'))
    cli = docker.APIClient(base_url='unix://var/run/docker.sock')

    # if rm is set to True, the 'ephemeral'/intermediate containers are removed
    # if rm is set to False or not even mentioned, the 'ephemeral' containers are not deleted
    # uncomment/comment the response line as you want to
    # documentation: https://docker-py.readthedocs.io/en/stable/images.html
    #####
    
    #response = [line for line in cli.build(fileobj=f, rm=True,  nocache=True, tag='hielaunewimage:latest')]
    response = [line for line in cli.build(fileobj=f, rm=False,  nocache=True, tag='hielaunewimage:latest')]

    print(response)
    print("      ")
    print("The small image has been built...")


def main():
   BuildImage()

if __name__ == "__main__":
    main()
```

Let's run it, with rm=False line as uncommented:
```
root@kr03nen:/home/gog0# python3 build_image.py 

[b'{"stream":"Step 1/3 : FROM python:3.4"}\r\n', 
b'{"stream":"\\n"}\r\n', 
b'{"stream":" ---\\u003e 8c62b065252f\\n"}\r\n', 
b'{"stream":"Step 2/3 : RUN apt-get update -y"}\r\n', 
b'{"stream":"\\n"}\r\n', b'{"stream":" ---\\u003e Running in 729fe571d784\\n"}\r\n', 
b'{"stream":"Get:1 http://security.debian.org/debian-security stretch/updates InRelease [53.0 kB]\\n"}\r\n', 
b'{"stream":"Ign:2 http://deb.debian.org/debian stretch InRelease\\n"}\r\n', 
b'{"stream":"Get:3 http://deb.debian.org/debian stretch-updates InRelease [93.6 kB]\\n"}\r\n', 
b'{"stream":"Get:4 http://deb.debian.org/debian stretch Release [118 kB]\\n"}\r\n', 
b'{"stream":"Get:5 http://deb.debian.org/debian stretch Release.gpg [2410 B]\\n"}\r\n', 
b'{"stream":"Get:6 http://security.debian.org/debian-security stretch/updates/main amd64 Packages [551 kB]\\n"}\r\n', 
b'{"stream":"Get:7 http://deb.debian.org/debian stretch-updates/main amd64 Packages [2596 B]\\n"}\r\n', 
b'{"stream":"Get:8 http://deb.debian.org/debian stretch/main amd64 Packages [7080 kB]\\n"}\r\n', 
b'{"stream":"Fetched 7901 kB in 3s (2448 kB/s)\\nReading package lists..."}\r\n', 
b'{"stream":"\\n"}\r\n', b'{"stream":" ---\\u003e b4f869d8bc3c\\n"}\r\n', 
b'{"stream":"Step 3/3 : RUN echo \\"hielau!\\""}\r\n', 
b'{"stream":"\\n"}\r\n', b'{"stream":" ---\\u003e Running in e01c4dab395c\\n"}\r\n', 
b'{"stream":"hielau!\\n"}\r\n', b'{"stream":" ---\\u003e 13d7f085fda3\\n"}\r\n', 
b'{"aux":{"ID":"sha256:13d7f085fda367758c717c9d5852e98cbc595bf27c198a033b3cd546aa45e429"}}\r\n', 
b'{"stream":"Successfully built 13d7f085fda3\\n"}\r\n', 
b'{"stream":"Successfully tagged hielaunewimage:latest\\n"}\r\n']

The small image has been built...
```

And let's check if the intermediate containers were /not/ removed (check them from bottom-up)

```
root@kr03nen:/home/gog0#  docker ps -a | head -4
CONTAINER ID        IMAGE                                 COMMAND                   CREATED              STATUS                          PORTS                NAMES
b395cc4dab5c        b4f869d8bc3c                          "/bin/sh -c 'echo \"h…"   About a minute ago   Exited (0) About a minute ago                        jolly_torvalds
71dfe571d784        8c62b065252f                          "/bin/sh -c 'apt-get…"    About a minute ago   Exited (0) About a minute ago     
```
Use command "docker rm <container ID>" to remove them. 

```
root@kr03nen:/home/gog0# docker rm e01c4dab395c 729fe571d784
e01c4dab395c
729fe571d784
root@kr03nen:/home/gog0# 
```

<b> Let's find our image </b>

You can either use the "docker images" command, with grep or sed/awk:

```
root@kr03nen:/home/gog0# docker images | grep hielau
hielaunewimage                  latest              13d7f085fda3        3 minutes ago       941MB
root@kr03nen:/home/gog0#
root@kr03nen:/home/gog0# # some sed just for the chills
root@kr03nen:/home/gog0#
root@kr03nen:/home/gog0# docker images | sed -n 's/\(hielau\)/\1/p'
hielaunewimage                  latest              13d7f085fda3        13 minutes ago      941MB
root@kr03nen:/home/gog0#
```


... or you could just list it with "docker image" command:

```
root@kr03nen:/home/gog0# docker image ls hielaunewimage
REPOSITORY          TAG                 IMAGE ID            CREATED              SIZE
hielaunewimage      latest              13d7f085fda3        About a minute ago   941MB
root@kr03nen:/home/gog0#
```

Let's check a bit the history of creating this docker image (use name or ID):

```
root@kr03nen:/home/gog0# docker history hielaunewimage
IMAGE               CREATED             CREATED BY                                      SIZE                COMMENT
13d7f085fda3        17 minutes ago      /bin/sh -c echo "hielau!"                       0B                  
b4f869d8bc3c        17 minutes ago      /bin/sh -c apt-get update -y                    16.4MB              
8c62b065252f        17 months ago       /bin/sh -c #(nop)  CMD ["python3"]              0B                  
<missing>           17 months ago       /bin/sh -c set -ex;   wget -O get-pip.py 'h8054fe4c907bt…   6.04MB              
<missing>           17 months ago       /bin/sh -c #(nop)  ENV PYTHON_PIP_VERSION=19…   0B                  
<missing>           17 months ago       /bin/sh -c cd /usr/local/bin  && ln -s idle3…   32B                 
<missing>           17 months ago       /bin/sh -c set -ex   && wget -O python.tar.x…   58.2MB              
<missing>           17 months ago       /bin/sh -c #(nop)  ENV PYTHON_VERSION=3.4.10    0B                  
<missing>           17 months ago       /bin/sh -c #(nop)  ENV GPG_KEY=97FC712E4C024…   0B                  
<missing>           17 months ago       /bin/sh -c apt-get update && apt-get install…   24.3MB              
<missing>           17 months ago       /bin/sh -c #(nop)  ENV LANG=C.UTF-8             0B                  
<missing>           17 months ago       /bin/sh -c #(nop)  ENV PATH=/usr/local/bin:/…   0B                  
<missing>           17 months ago       /bin/sh -c set -ex;  apt-get update;  apt-ge…   562MB               
<missing>           17 months ago       /bin/sh -c apt-get update && apt-get install…   142MB               
<missing>           17 months ago       /bin/sh -c set -ex;  if ! command -v gpg > /…   7.81MB              
<missing>           17 months ago       /bin/sh -c apt-get update && apt-get install…   23.2MB              
<missing>           17 months ago       /bin/sh -c #(nop)  CMD ["bash"]                 0B                  
<missing>           17 months ago       /bin/sh -c #(nop) ADD file:e4bdc12117ee95eaa…   101MB       
```


<b> Let's Go ... </b>
<br>
And now, let's try to implement same thing we did previously in Python, only in Golang ...

You can find the code in same repository, 
or just check out the raw form, by clicking <a href="https://raw.githubusercontent.com/kroen3n/G0g0-/master/related/Dockerfile/code/build_image.go">build_image.go</a>


The output you'd be obtaining should be looking as below:

```
root@kr03nen:/home/gog0# go run build_image.go 
{"stream":"Step 1/3 : FROM python:3.4"}
{"stream":"\n"}
{"stream":" ---\u003e 8c62b065252f\n"}
{"stream":"Step 2/3 : RUN apt-get update -y"}
{"stream":"\n"}
{"stream":" ---\u003e Using cache\n"}
{"stream":" ---\u003e 26bc7a14085f\n"}
{"stream":"Step 3/3 : RUN echo \"hielau!\""}
{"stream":"\n"}
{"stream":" ---\u003e Running in 62e19e8ff996\n"}
{"stream":"hielau!\n"}
{"stream":" ---\u003e c25823fc340f\n"}
{"aux":{"ID":"sha256:c25823fc340f7b0da6c10c404b6471446515d878112255bd7387a9da66790162"}}
{"stream":"Successfully built c25823fc340f\n"}
{"stream":"Successfully tagged hielaunewimage:latest\n"}
```
