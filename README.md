<p align="center">
    <img width="196" height="196" src="https://avatars.githubusercontent.com/u/96201133">
    <h3 align="center">Anchor</h3>
    <p align="center">A Minecraft server implementation written in Go.</p>
    <p align="center">
        <a href="https://github.com/AnchorMC/Anchor/issues/new">Report an Issue</a> &bullet; <a href="https://github.com/AnchorMC/Anchor/compare">Open a Pull Request</a> &bullet; <a href="https://discord.gg/8G5hDECBWk">Discord Server</a>
    </p>
</p>

## About

This project was started to better understand the interaction between the server and the client, but has recently changed into a desire to create a working Minecraft server from scratch in Go. You can learn more about our development and interact with the community by joining our [Discord Server](https://discord.gg/8G5hDECBWk).

**Please note** that this project recently underwent a complete rewrite to better work around the restrictive packaging structure that Go puts in place. All code is being converted into a single base Go package, except for a few utilities that are better fit into their own package.

## Features

- [ ] Server status (with favicon)
- [ ] Login
- [ ] Authentication
- [ ] Two-way encryption
- [ ] Query
- [ ] World generation (only flat available)
- [ ] Region store
- [ ] Player movement and rotation
- [ ] Chunk updates
- [ ] Physics
- [ ] Update blocks
- [ ] Dropped entity interaction
- [ ] NPCs (animals, mobs, etc.)
- [ ] Plugins
- [ ] Console input
- [ ] RCON
- [ ] Chat
- [ ] Crafting
- [ ] ... and so many utility APIs

## Installation

There will be no distributable binaries until there is a working release with most of the features checklist completed. In the mean time, you are free to build from source if you would like to.

```
git clone https://github.com/AnchorMC/Anchor.git
cd Anchor
make
```

The executable will now be located in the `bin` folder for you to run. It is recommended that you create a test folder and copy this executable to there so the server does not create a bunch of directories during startup within the source folder.

## Compatibility

Please note that this is not meant to be a drop-in replacement for existing servers. While the goal of this software is to match the format of standard Java Edition configuration and data storage formats, it may not be perfect and will be missing quite a few features. Over time, more features will be available and functional but the ETA is not known.

## The Goal

The goal of this project is to create a fully functional custom Minecraft server software that is highly efficient and outperforms existing Java-based server architecture. This is a slow and highly involved endeavour that will not be completed quickly. This project may forever be in a non-complete state, meaning that it will not perfectly match Java Edition servers, but we can always work towards providing maximum compatibility.

## Contributors

<table>
    <tr>
        <td align="center"><a href="https://github.com/PassTheMayo"><img src="https://avatars.githubusercontent.com/u/16949253?v=4&s=100" width="100px;" alt=""/><br /><b>Jacob Gunther</b></a><br/><sub>Main Developer</sub></td>
    </tr>
</table>

## Discord Server

https://discord.gg/8G5hDECBWk

## License

[GNU General Public License v3.0](https://github.com/AnchorMC/Anchor/blob/main/LICENSE)