package main

import (
    "bytes"
    "time"
    "strings"

    "github.com/SREnity/epico/signers/aws_v4"
    utils "github.com/SREnity/epico/utils"
    generic_structs "github.com/SREnity/epico/structs"
    xj "github.com/basgys/goxml2json"
)

var PluginAuthFunction = PluginAuth
var PluginPostProcessFunction = PluginPostProcess
var PluginPagingPeekFunction = utils.DefaultXmlPagingPeek


// authParams expects the AWS ACCESS ID in slice slot [0] and the AWS SECRET KEY
//     in slice slot [1].
func PluginAuth( apiRequest generic_structs.ApiRequest, authParams []string ) generic_structs.ApiRequest {

    var err error
    signer := v4.NewSigner( generic_structs.ApiCredentials{
            Id: authParams[0],
            Key: authParams[1],
        }, func(v4Signer *v4.Signer) {
            v4Signer.DisableHeaderHoisting = false
            v4Signer.DisableRequestBodyOverwrite = true

    })

    if apiRequest.Settings.Vars["region"] != "{{region}}" {
        apiRequest.FullRequest.Header, err = signer.Presign( apiRequest.FullRequest, nil, apiRequest.Settings.Vars["service"], apiRequest.Settings.Vars["region"], 60 * time.Minute, time.Now() )
    } else {
        apiRequest.FullRequest.Header, err = signer.Presign( apiRequest.FullRequest, nil, apiRequest.Settings.Vars["service"], "us-east-1", 60 * time.Minute, time.Now() )
    }

    if err != nil {
        utils.LogFatal("AWS:PluginAuth", "Error presigning the AWS request", err)
    }

    return apiRequest

}


func PluginPostProcess( apiResponseMap map[generic_structs.ComparableApiRequest][]byte, jsonKeys []map[string]string, postParams []string ) []byte {

    parsedStructure := make(map[string]interface{})
    parsedErrorStructure := make(map[string]interface{})

    for response, apiResponse := range apiResponseMap {
        jsonBody, err := xj.Convert( bytes.NewReader( apiResponse ) )
        if err != nil {
            utils.LogFatal("AWS:PluginPostProcess", "Error converting XML API response", err)
            return nil
        }

        // This chunk reads in the list of responses and handles AWS' terrible
        //    XML return with infinite "item"/"member" tags.
        processedJson := jsonBody.Bytes()
        for _, v := range jsonKeys {
            if  v["api_call_name"] == response.Name {
                xmlKeys := strings.Split( v["xml_tags"], "," )
                for _, k := range xmlKeys {
                    processedJson = utils.RemoveXmlTagFromJson(
                        k, processedJson)
                }
            }
        }

        // TODO: Should I just rebuild the map and pass to the generic
        //    utils.DefaultJsonPostProcess function?
        utils.ParsePostProcessedJson( response, jsonKeys, processedJson, parsedStructure, parsedErrorStructure )

    }

    returnJson := utils.CollapseJson( parsedStructure, parsedErrorStructure )
    return returnJson
}
