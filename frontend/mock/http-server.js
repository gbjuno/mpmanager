const express = require('express')
const http = require('http')
const fs = require('fs')
const path = require('path');

const chokidar = require('chokidar');

const watcher = chokidar.watch(path.join(__dirname, '/mock'));

watcher.on('ready', () => {

        watcher.on('change', (path) => {

                console.log('<---- watched file change, do something ---->');

        });

        watcher.on('add', (path) => {

                console.log('<---- watched new file add, do something ---->');

        });

        watcher.on('unlink', (path) => {

                console.log('<---- watched file remove, do something ---->');

        });

});

const app = new express()

//allow custom header and CORS
app.all('*',function (req, res, next) {
  res.header('Access-Control-Allow-Origin', '*');
  res.header('Access-Control-Allow-Headers', 'Content-Type, Content-Length, Authorization, Accept, X-Requested-With , yourHeaderFeild');
  res.header('Access-Control-Allow-Methods', 'PUT, POST, GET, DELETE, OPTIONS');

  if (req.method == 'OPTIONS') {
    res.send(200); /让options请求快速返回/
  }
  else {
    next();
  }
});

app.get('/pic', function(req, res){
    client(req, function(info){
        res.send(info)
    })
});

const server = app.listen(8081, function(){
    const host = server.address().address
    const port = server.address().port

    console.log("应用实例，访问的地址为 http://%s:%s", host, port)
})

// function client(req, callback){
//     http.get('http://localhost:3006/pic-example.json', function(res){
//         let buffers = []
//         console.log('req ssss', req)
//         res.on('data', function(chunk){
//             buffers.push(chunk)
//         })

//         res.on('end', function(chunk){
//             const wholeData = Buffer.concat(buffers);
//             const dataStr = wholeData.toString('utf8')
//             callback(dataStr)
//         })
//     }).on('error', function(e){
//         callback(e.message)
//     })
// }

function client(req, callback){
    const prefix = __dirname
    let picName = req.query.picName
    const file = path.join(__dirname, "/store/pic-example.json")
    const result = JSON.parse(fs.readFileSync(file))
    callback(filter(picName, result)) 
}

function filter(picName, result){
    let filterRes = []
    for(let item of result){
        if(item.name.indexOf(picName) > -1){
            filterRes.push(item)
        }
    }
    return filterRes
}