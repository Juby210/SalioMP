# SalioMP
[Salio](https://store.steampowered.com/app/875810/Salio) Multiplayer Mod

This is alpha version, only player position & levels synced and support only 2 players!


### How to use?
- Run `npm install` in server directory
- Run `node server/server.js`
- Run `npm install` in client directory
- Change `saliopath` for salio installation directory in `client/config.json`
- Run `node client/start.js`
- Move `mod/data.salio` to salio in workshop menu

On your friend pc:
- Run `npm install` in client directory
- Change `saliopath` for salio installation directory in `client/config.json`
- Change `ip` to your ip in `client/config.json`
- Run `node client/start.js`
- Move `mod/data.salio` to salio in workshop menu

Optional step:
Change `server/data.salio` and `server/levels.zip` to your levels
