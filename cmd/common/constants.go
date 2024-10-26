package common

const ServerHost string = "http://localhost"
const ServerPort string = ":8080"
const ServerUrl = ServerHost + ServerPort

const GetEndpoint string = "/value"
const GetFullEndpoint = ServerUrl + GetEndpoint + "/"

const UpdateEndpoint string = "/update"
const UpdateFullEndpoint = ServerUrl + UpdateEndpoint + "/"
