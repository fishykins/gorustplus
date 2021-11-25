# GoRustPlus
<div id="top"></div>

This is an **unofficial** websocket api for the PC game [Rust](https://store.steampowered.com/app/252490/Rust/) and its companion app, Rusplus.

Websocket interactions were first documented in [Rustplus](https://github.com/liamcottle/rustplus.js), for more details hop on over to Liam's repo. 

## About The Project
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
