# Changelog

## 1.0.0 (2025-08-07)


### Features

* add ability to post message (does not update messages live yet) ([16b43ff](https://github.com/TheQueenIsDead/go-chat/commit/16b43ffa28d8bd4ed2e5671afaa3e3472422d5e8))
* add hard coded demo data and layout for rooms and chats ([e86b47f](https://github.com/TheQueenIsDead/go-chat/commit/e86b47f23227dc75019e20be9e4e397d462cae53))
* add user middleware to generate username cookie and clear input after posting a new message ([472c577](https://github.com/TheQueenIsDead/go-chat/commit/472c577cc00149f5e417e6fac7b38792e844708f))
* adjust templates to render based on demo input data ([4103849](https://github.com/TheQueenIsDead/go-chat/commit/4103849b832843fcfa9d922415a19382306146d3))
* create a persistent SSE connection in order to add messages in real time ([e084a99](https://github.com/TheQueenIsDead/go-chat/commit/e084a9973c9d2da585318e49bc0d43d3cc4b4f7b))
* highlight active room based on data signal ([7d378fc](https://github.com/TheQueenIsDead/go-chat/commit/7d378fc00161b138fcd76c1f3f89bc6a0ad83864))
* initialise embedded NATS server ([078d90c](https://github.com/TheQueenIsDead/go-chat/commit/078d90c8c52bfccdaff3f2180b7c2d62eea62101))
* switch between rooms by clicking on the side-bar ([33d30b6](https://github.com/TheQueenIsDead/go-chat/commit/33d30b66c0c0f4edad178f027b996ee1f92cc5e2))
* use nats to subscribe to SSE via JS and refactor code ([7b6a761](https://github.com/TheQueenIsDead/go-chat/commit/7b6a7617777c1b30e4e670a5a59cd453ee76f0a5))
* wire up templ ui index with datastar loading empty rooms and chats ([758baa6](https://github.com/TheQueenIsDead/go-chat/commit/758baa62dff3f67f6019bc4f01fed2d925323ace))


### Bug Fixes

* render message id correctly and scroll into view on room change ([983bd61](https://github.com/TheQueenIsDead/go-chat/commit/983bd6134a73cd5e862d8eeb2f3070da89c5d116))
