package config

import "time"

const DATEBASE = "postgres://user:password@localhost/metrics"
const ADDRSERV = "localhost:8080"
const ACCRUALSYSTEMADDRESS = "http://localhost:8081"

// hash
const PASSWORDHASH = "1234567"

// jwt
const TOKEN_EXP = time.Hour * 3
const SECRET_KEY = "supersecretkey"
const NAME_TOKEN = "token"

// chan
const ELEMENTS = 1
