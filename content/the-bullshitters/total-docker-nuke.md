---
title: The Total Docker Nuke
tags: Docker, Catharsis
description: When diplomacy with Docker fails.
icon: RM
reading_time: 3 min
date: 2026-05-06
---

Every Docker user eventually meets the moment: a container that won't start, logs that read like ransom notes, and volumes that may or may not contain a small demon. You've tried turning it off and on again. You've tried it twice. You've re-read the logs just in case they changed their minds. They didn't.

Diplomacy has failed. We nuke from orbit — it's the only way to be sure.

## The Nuke

Run these one at a time, in order. Containers first, then images, then volumes, then custom networks, then a final sweep:

```bash
docker rm -f $(docker ps -aq) 2>/dev/null
docker rmi -f $(docker images -aq) 2>/dev/null
docker volume rm $(docker volume ls -q) 2>/dev/null
docker network rm $(docker network ls -q --filter type=custom) 2>/dev/null
docker system prune -a --volumes -f
```

The order matters — containers hold images, images hold layers, volumes and networks hold references. Going top-down means no "resource in use" errors. The final system prune is technically redundant with the four above it, but it's a belt-and-braces move and it clears the build cache while it's at it.

## The Disclaimer

This is not a project reset. It is a Docker reset. Everything goes:

- Every other project's containers, images, and volumes on this machine — gone
- Every base image you've ever pulled (node, postgres, nginx) — re-downloading on next build
- Build cache across every project — cold from scratch
- Custom networks from compose files you forgot you ran — vapourised

What it doesn't touch: buildx builders (`docker buildx rm --all-inactive -f` for those) and swarm state (`docker swarm leave --force` if you ever ran swarm init). Most people don't have either.

If you only wanted to nuke one project, you reached for the wrong tool. `docker compose down -v --rmi all --remove-orphans` from the project directory does the targeted version. But if you're reading this, you probably knew that and chose violence anyway. Respect.
