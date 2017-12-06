# simplechain
A simple blockchain implementation based on [lhartikk/naivechain](https://github.com/lhartikk/naivechain).
## API

### GET `/blocks`
```console
$ curl -X GET http://localhost:3001/blocks
[{"index":1,"previousHash":"0","timestamp":1465154705,"data":"my genesis block!!","hash":"816534932c2b7154836da6afc367695e6337db8a921823784c14378abed4f7d7"}]
```

### POST `/mineBlock`
```console
$ curl -X POST --data '{"data":"some data"}' http://127.0.0.1:3001/mineBlock
```

### GET `/peers`
```console
$ curl -X GET http://127.0.0.1:3002/peers
["[::1]:6001"]
```

### GET `/addPeer`
```console
$ curl -X POST --data '{"peer":"ws://localhost:6001"}' http://localhost:3002/addPeer
```

## License
MIT License
