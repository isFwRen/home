/**
 * The preload script runs before. It has access to web APIs
 * as well as Electron's renderer process modules and some
 * polyfilled Node.js functions.
 *
 * https://www.electronjs.org/docs/latest/tutorial/sandbox
 */

window.addEventListener('DOMContentLoaded', () => {
  // const { io } = require("socket.io-client");
  // const socket = io("http://127.0.0.1:30080");
  // const element = document.getElementById("precent")
  // socket.on('connect', () => {
  //   console.log("connected renderer");
  // });
  // console.log(socket.connected);
  // socket.on("update", (data) => {
  //   console.log("update ", data);
  //   if (element) element.innerText = (data.current * 100 / data.total).toFixed(2) + "%"
  //   if (data.action == "结束更新") window.close()
  // });
  // socket.on('error', (e) => {
  //   console.log(e);
  // });
})