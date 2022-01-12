<p align="center">
    <img width="196" height="196" src="https://avatars.githubusercontent.com/u/96201133">
    <h3 align="center">Minecraft-Server</h3>
    <p align="center">A Minecraft server implementation written in pure Go.</p>
    <p align="center">
        <a href="https://github.com/GoLangMinecraft/Minecraft-Server/issues/new">Report an Issue</a> &bullet; <a href="https://github.com/GoLangMinecraft/Minecraft-Server/compare">Open a Pull Request</a> &bullet; <a href="https://discord.gg/8G5hDECBWk">Discord Server</a>
    </p>
</p>
<hr/>
<h2 id="about">About</h2>
<p>This project was started to better understand the interaction between the server and the client, but has recently changed into a desire to create a working Minecraft server from scratch in Go. You can learn more about our development and interact with the community by joining our <a href="#discord">Discord Server</a>.</p>
<h2 id="state-of-project">Features</h2>
<ul>
    <li>
        <input type="checkbox" checked>
        Server status
    </li>
    <li>
        <input type="checkbox" checked>
        Login
    </li>
    <li>
        <input type="checkbox" checked>
        Two-way encryption
    </li>
    <li>
        <input type="checkbox">
        World generation (work-in-progress)
    </li>
    <li>
        <input type="checkbox">
        Physics
    </li>
    <li>
        <input type="checkbox">
        Update blocks
    </li>
    <li>
        <input type="checkbox">
        Dropped entity interaction
    </li>
    <li>
        <input type="checkbox">
        NPCs (animals, monsters, etc.)
    </li>
    <li>
        <input type="checkbox">
        Plugins
    </li>
</ul>
<h2 id="installation">Installation</h2>
<p>There will be no distributable binaries until there is a working release with most of the features checklist completed. In the mean time, you are free to build from source if you would like to.</p>
<pre><code>git clone https://github.com/GoLangMinecraft/Minecraft-Server.git
cd Minecraft-Server
./scripts/build # or .\scripts\build.bat for Windows</code></pre>
<p>The executable will now be located in the <code>bin</code> folder for you to run. It is recommended that you create a test folder and copy this executable to there so the server does not create a bunch of directories during startup within the source folder.</p>
<h2 id="credit">Credit</h2>
<p>The goal was to create this Minecraft server implementation from scratch, but it became clear that this was not going to happen by myself. I would like to give a big thanks to Tnze, the creator of the <a href="https://github.com/Tnze/go-mc">go-mc</a> package, for their work on the amazing utilities for interacting with the Minecraft protocol.</p>
<h2 id="discord">Discord Server</h2>
<p><a href="https://discord.gg/8G5hDECBWk">https://discord.gg/8G5hDECBWk</a></p>
<h2 id="license">License</h2>
<p><a href="https://github.com/GoLangMinecraft/Minecraft-Server/blob/main/LICENSE">GNU General Public License v3.0</a></p>