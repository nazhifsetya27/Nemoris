#!/usr/bin/env node
/**
 * NEMORIS Docs Server
 * Serves docs so you can read them on your phone.
 * Run: node serve.js
 * Then open the printed URL on your phone (same WiFi).
 */

const http = require('http');
const fs = require('fs');
const path = require('path');
const os = require('os');

const PORT = process.env.PORT || 3847;
const DOCS_DIR = __dirname;

const MIME = {
  '.html': 'text/html',
  '.md': 'text/markdown',
  '.css': 'text/css',
  '.js': 'application/javascript',
  '.json': 'application/json',
  '.ico': 'image/x-icon',
};

function getLocalIP() {
  try {
    const nets = os.networkInterfaces();
    for (const name of Object.keys(nets)) {
      for (const net of nets[name]) {
        if (net.family === 'IPv4' && !net.internal) return net.address;
      }
    }
  } catch (_) {}
  return null;
}

const server = http.createServer((req, res) => {
  let url = req.url === '/' ? '/index.html' : req.url;
  url = url.split('?')[0];
  const filePath = path.join(DOCS_DIR, path.normalize(url));

  if (!filePath.startsWith(DOCS_DIR)) {
    res.writeHead(403);
    res.end('Forbidden');
    return;
  }

  fs.readFile(filePath, (err, data) => {
    if (err) {
      if (err.code === 'ENOENT') {
        res.writeHead(404);
        res.end('Not found');
      } else {
        res.writeHead(500);
        res.end('Server error');
      }
      return;
    }
    const ext = path.extname(filePath);
    res.setHeader('Content-Type', MIME[ext] || 'application/octet-stream');
    res.end(data);
  });
});

server.listen(PORT, '0.0.0.0', () => {
  const local = getLocalIP();
  console.log('\n  NEMORIS Docs\n');
  console.log('  On this computer:  http://localhost:' + PORT);
  if (local) {
    console.log('  On your phone:     http://' + local + ':' + PORT);
    console.log('\n  Make sure phone and computer are on the same WiFi.\n');
  } else {
    console.log('  On your phone:     Use your Mac\'s IP + :' + PORT + ' (run: ipconfig getifaddr en0)\n');
  }
});
