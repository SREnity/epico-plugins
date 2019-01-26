package main

import (
    "bytes"
    "time"
    "strings"

    utils "github.com/SREnity/epico/utils"
    generic_structs "github.com/SREnity/epico/structs"
    xj "github.com/basgys/goxml2json"
)

var PluginAuthFunction = PluginAuth
var PluginPostProcessFunction = PluginPostProcess
var PluginPagingPeekFunction = PluginPagingPeek


// authParams expects the Auth function in slice slot [0] followed by a []string
//    of auth parameters
func PluginAuth( apiRequest generic_structs.ApiRequest, authParams []string ) generic_structs.ApiRequest {

    if (authParams) < 1 {
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
            return utils.JwtAuth( apiRequest, authParams[1:] )
        default:
            return apiRequest
    }

}


func PluginPagingPeek( response []byte, responseKeys []string, oldPageValue interface{}, peekParams []string ) ( interface{}, bool ) {

    if (peekParams) < 1 {
        peekParams = []string{""}
    }

    switch peekParams[0] {
        default:
            return utils.DefaultJsonPagingPeek( response, responseKeys,
                oldPageValue, peekParams )
    }

}


func PluginPostProcess( apiResponseMap map[generic_structs.ComparableApiRequest][]byte, jsonKeys []map[string]string, postParams []string ) []byte {

    if (postParams) < 1 {
        postParams = []string{""}
    }

    switch postParams[0] {
        default:
            parsedStructure := make(map[string]interface{})
            parsedErrorStructure := make(map[string]interface{})

            for response, apiResponse := range apiResponseMap {
                utils.ParsePostProcessedJson( response, apiResponse,
                    parsedStructure, parsedErrorStructure )
            }

            returnJson := utils.CollapseJson(
                parsedStructure, parsedErrorStructure )
            return returnJson
    }

}
