<p align="center">
    <img width="196" height="196" src="https://avatars.githubusercontent.com/u/96201133">
    <h3 align="center">Minecraft-Server</h3>
    <p align="center">A Minecraft server implementation written in pure Go.</p>
    <p align="center">
        <a href="https://github.com/GoLangMinecraft/Minecraft-Server/issues/new">Report an Issue</a> &bullet; <a href="https://github.com/GoLangMinecraft/Minecraft-Server/compare">Open a Pull Request</a> &bullet; <a href="https://discord.gg/8G5hDECBWk">Discord Server</a>
    </p>
</p>

## About

This project was started to better understand the interaction between the server and the client, but has recently changed into a desire to create a working Minecraft server from scratch in Go. You can learn more about our development and interact with the community by joining our [Discord Server](https://discord.gg/8G5hDECBWk).

## Features

- [x] Server status (with favicon)
- [x] Login
- [x] Authentication
- [x] Two-way encryption
- [x] Query
- [ ] World generation (work-in-progress)
- [ ] Physics
- [ ] Update blocks
- [ ] Dropped entity interaction
- [ ] NPCs (animals, mobs, etc.)
- [ ] Plugins
- [x] Console input
- [ ] RCON
- [ ] Chat
- [ ] Crafting
- [ ] ... and so many utility APIs

## Installation

There will be no distributable binaries until there is a working release with most of the features checklist completed. In the mean time, you are free to build from source if you would like to.

```
git clone https://github.com/GoLangMinecraft/Minecraft-Server.git
cd Minecraft-Server
./scripts/build # or .\scripts\build.bat for Windows
```

The executable will now be located in the `bin` folder for you to run. It is recommended that you create a test folder and copy this executable to there so the server does not create a bunch of directories during startup within the source folder.

## Credit

The goal was to create this Minecraft server implementation from scratch, but it became clear that this was not going to happen by myself. I would like to give a big thanks to Tnze, the creator of the [go-mc](https://github.com/Tnze/go-mc) package, for their work on the amazing utilities for interacting with the Minecraft protocol.

## Contributors

<table>
    <tr>
        <td align="center"><a href="https://github.com/PassTheMayo"><img src="https://avatars.githubusercontent.com/u/16949253?v=4&s=100" width="100px;" alt=""/><br /><sub><b>Jacob Gunther</b></sub></a><br/> <p>Main Developer</p></td>
    </tr>
</table>

## Discord Server

https://discord.gg/8G5hDECBWk

## License

[GNU General Public License v3.0](https://github.com/GoLangMinecraft/Minecraft-Server/blob/main/LICENSE)