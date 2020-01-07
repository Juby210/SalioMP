const fs = require('fs')
const admzip = require('adm-zip')
const fetch = require('node-fetch')
const { join } = require('path')
let lastx = -18, lasty = -18
const mod = join(__dirname, '..', 'mod')

module.exports = ({ saliopath, ip }) => {
  saliopath = join(saliopath, 'mp')
  const pathes = {
    p1:     join(saliopath, 'p1.txt'),
    p2:     join(saliopath, 'p2.txt'),
    data:   join(mod,       'data.salio'),
    levels: join(mod,       'levels'),
    lzip:   join(mod,       'levels.zip')
  }

  if(!fs.existsSync(saliopath)) fs.mkdirSync(saliopath)
  if(!fs.existsSync(pathes.p2)) fs.openSync(pathes.p2, 'w')
  fs.writeFileSync(pathes.p2, `-18\n-18`)

  const socket = require('socket.io-client')(ip)
  socket.on('connect', () => {
    console.log('Connected')
  })

  socket.on('config', async ({ syncLevels }) => {
    if(syncLevels) {
      console.log('Syncing levels..')
      const data = await (await fetch(ip + '/data.salio')).text()
      if(!fs.existsSync(pathes.data)) fs.openSync(pathes.data, 'w')
      const datao = fs.readFileSync(pathes.data).toString()
      if(data == datao) return console.log('Done')
      fs.writeFileSync(pathes.data, data)

      const res = await fetch(ip + '/levels.zip')
      const stream = fs.createWriteStream(pathes.lzip)
      res.body.pipe(stream)
      await new Promise(r => stream.on('finish', r))

      if(fs.existsSync(pathes.levels)) delFolder(pathes.levels)
      fs.mkdirSync(pathes.levels)
      new admzip(pathes.lzip).extractAllTo(pathes.levels, true)
      console.log('Done')
      setTimeout(() => fs.unlinkSync(pathes.lzip), 150)
    }
  })

  socket.on('playermove', data => {
    if(!fs.existsSync(saliopath + '/p2.txt')) fs.openSync(saliopath + '/p2.txt', 'w')
    fs.writeFileSync(saliopath + '/p2.txt', `${data.x}\n${data.y}`)
    console.log(`p2; ${data.x} | ${data.y}`)
  })

  if(!fs.existsSync(saliopath + '/p1.txt')) fs.openSync(saliopath + '/p1.txt', 'w')
  fs.watchFile(saliopath + '/p1.txt', { interval: 10 }, () => {
    let content = fs.readFileSync(saliopath + '/p1.txt', 'utf-8').toString()
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

function delFolder(path) {
  if(fs.existsSync(path)) {
    fs.readdirSync(path).forEach(file => {
      const curPath = path + '/' + file
      if(fs.lstatSync(curPath).isDirectory()) {
        delFolder(curPath)
      } else {
        fs.unlinkSync(curPath)
      }
    })
    fs.rmdirSync(path)
  }
}
