const fs = require('fs')
let lastx = -18
let lasty = -18

module.exports = ({saliopath, ip}) => {
  if(!fs.existsSync(saliopath + "/p2.txt")) fs.openSync(saliopath + "/p2.txt", "w")
  fs.writeFileSync(saliopath + "/p2.txt", `-18\n-18`)

  const socket = require('socket.io-client')(ip)
  socket.on('connect', () => {
    console.log('Connected')
  })
  socket.on('playermove', data => {
    if(!fs.existsSync(saliopath + "/p2.txt")) fs.openSync(saliopath + "/p2.txt", "w")
    fs.writeFileSync(saliopath + "/p2.txt", `${data.x}\n${data.y}`)
    console.log(`p2; ${data.x} | ${data.y}`)
  })

  if(!fs.existsSync(saliopath + "/p1.txt")) fs.openSync(saliopath + "/p1.txt", "w")
  fs.watchFile(saliopath + "/p1.txt", { interval: 10 }, () => {
    let content = fs.readFileSync(saliopath + "/p1.txt", "utf-8").toString()
    let x = content.split('\n')[0]
    if(!x) return
    x = x.replace('\r', '')
    let y = content.split('\n')[1]
    if(!y) return
    y = y.replace('\r', '')
    if(!x || !y) x = -18, y = -18;
    if(lastx != x || lasty != y) {
      socket.emit('playermove', { x, y })
      lastx = x
      lasty = y
    }
  })
}
