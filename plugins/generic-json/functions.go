package main

import (
    utils "github.com/SREnity/epico/utils"
    generic_structs "github.com/SREnity/epico/structs"
)

var PluginAuthFunction = PluginAuth
var PluginPostProcessFunction = PluginPostProcess
var PluginPagingPeekFunction = PluginPagingPeek


// authParams expects the Auth function in slice slot [0] followed by a []string
//    of auth parameters
func PluginAuth( apiRequest generic_structs.ApiRequest, authParams []string ) generic_structs.ApiRequest {

    if len(authParams) < 1 {
        authParams = []string{""}
    }

    // TODO: I would like to be able to dynamically call whatever function name
    //    is passed in authParams[0], but this is difficult in Go - particularly
    //    if they are going outside the standard utils library. At this point,
    //    we'll just manage the duplication here between core utils and these
    //    generic files.  If folks want to use their own functions, they should
    //    create their own plugins.
    switch authParams[0] {
        case "JwtAuth":
            // Expects authParams[1..5] to be:
            // email, private key, private key id, scopes (csv), token url
            return utils.JwtAuth( apiRequest, authParams[1:] )
        default:
            return apiRequest
    }

}


func PluginPagingPeek( response []byte, responseKeys []string, oldPageValue interface{}, peekParams []string ) ( interface{}, bool ) {

    if len(peekParams) < 1 {
        peekParams = []string{""}
    }

    switch peekParams[0] {
        default:
            return utils.DefaultJsonPagingPeek( response, responseKeys,
                oldPageValue )
    }

}


func PluginPostProcess( apiResponseMap map[generic_structs.ComparableApiRequest][]byte, jsonKeys []map[string]string, postParams []string ) []byte {

    if len(postParams) < 1 {
        postParams = []string{""}
    }

    switch postParams[0] {
        default:
            // No jsonKeys required because we aren't expecting any special vars
            //    or configurations.
            return utils.DefaultJsonPostProcess( apiResponseMap, jsonKeys )
    }

}
