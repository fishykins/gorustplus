<div id="top"></div>

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/fishykins/gorustplus/assets/smartAlarm.png">
    <img src="assets/smartAlarm.png" alt="Logo" width="120" height="120">
  </a>

  <h3 align="center">GoRustPlus</h3>

  <p align="center">
    An unofficial websocket api for the PC game <strong>Rust</strong> and its companion app, <strong>Rustplus</strong>.
    <br />
    <a href="https://github.com/fishykins/gorustplus/">Documentation</a>
    ·
    <a href="https://store.steampowered.com/app/252490/Rust/">Rust</a>
    ·
    <a href="https://github.com/liamcottle/rustplus.js">Rustplus.js</a>

  </p>
</div>



## About The Project
Websocket interactions were first documented in [Rustplus.js](https://github.com/liamcottle/rustplus.js), for more details hop on over to Liam's repo. 
GoRustPlus was developed to help facilitate more advanced discord/rust interactions, as well as providing the community with an alternative option to JS when programing bots.
As apps written to handle such interactions rely heavily on concurrency and scalability, Golang seems like the perfect fit for the job!

## Design Philosophy

This project aims to handle basic websocket interactions for the end user without cutting them off from any of the data. Where possible, an effort has been made to keep
interactions minimalistic, with any abstractions kept to the bare minimum. Methods for reading/writing have been implemented, but it is down to the end user how and when 
these are used.

## Quirks and Usage

While it is possible to handle device states using raw requests, a caching system is available to help prevent excessive api calls. By registering devices
with the Client, it is possible to add more complex behavior, such as event handling for any broadcasts, and ensure that all device data is fully validated and kept up to date.
Additionally, there is support for event callbacks, allowing for asynchronous handling of read/write calls.

While devices are more tightly managed, chat messages and team updates are palmed off to the api caller via dedicated channels. By default these are left empty and will 
not be used, but once they are set all messages will be relayed via their dedicated chanel, causing the execution to block until the channel is read. handle with care!
