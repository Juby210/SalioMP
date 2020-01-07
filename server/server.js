let p1, p2

const { readFileSync, createReadStream } = require('fs')
const config = require('./config.json')

const server = require('http').createServer((req, res) => {
  let end = true
  if(req.url == '/data.salio') {
    res.write(readFileSync('data.salio').toString())
  } else if(req.url == '/levels.zip') {
    res.writeHead(200, { 'Content-Type': 'application/zip' })
    createReadStream('levels.zip').pipe(res)
    end = false
  } else {
    res.writeHead(302, { Location: 'https://github.com/juby210-PL/SalioMP' })
  }
  if(end) res.end()
})
const io = require('socket.io')(server)

io.on('connection', socket => {
  let p11
  if(!p1) {
    p1 = socket
    p11 = true
  } else p2 = socket;

  socket.emit('config', config)

  socket.on('playermove', data => {
    console.log(`${p11 ? 'p1' : 'p2'}; ${data.x} | ${data.y}`)
    if(p11) {
      if(p2 != null) p2.emit('playermove', data)
    } else {
      if(p2 != null) p1.emit('playermove', data)
    }
  })

  socket.on('disconnect', () => {
    if(p11) {
      p1 = null
      if(p2 != null) p2.emit('playermove', { x: -18, y: -18 })
    } else {
      p2 = null
      if(p1 != null) p1.emit('playermove', { x: -18, y: -18 })
    }
  })
})
server.listen(config.port)
