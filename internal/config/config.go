package config

import "time"

const DATEBASE = "postgres://user:password@localhost/metrics"
const ADDRSERV = "localhost:8080"
const ACCRUALSYSTEMADDRESS = "http://localhost:8081"

// hash
const PASSWORDHASH = "1234567"

// jwt
const TOKENEXP = time.Hour * 3
const SECRETKEY = "supersecretkey"
const NAMETOKEN = "token"

// chan
const ELEMENTS = 1
