# Dev Container Router

The core idea around this project is that:

1. I have a set of devcontainers running on a machine in my local network
2. I want to have long running sessions that attach to those containers on that machine
3. A proxy machine which I can ssh into and then tunnels into the container on the other machine
4. A CLI that allows me to ssh into a container through that tunnel.

So the essence of the project is that I am able to say on any machine that has an internet connection:

dcr music

And it ssh's into my router machine directly, and then bounces into the devcontainer running on a different machine.

We can assume that I have permissions to ssh directly to the router machine, and from the router to devcontainer machine.