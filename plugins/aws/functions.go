package main

import (
    "bytes"
    "time"
    "strings"
    "strconv"

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

    requestCount := 0
    for request, apiResponse := range apiResponseMap {
        requestCount = requestCount + 1
        jsonBody, err := xj.Convert( bytes.NewReader( apiResponse ) )
        if err != nil {
            utils.LogFatal("AWS:PluginPostProcess", "Error converting XML API response", err)
            return nil
        }

        // This chunk reads in the list of responses and handles AWS' terrible
        //    XML return with infinite "item"/"member" tags.
        count := 0
        processedJson := jsonBody.Bytes()
        for _, v := range jsonKeys {
            if  v["api_call_name"] == request.Name {
                count = count + 1
                utils.LogWarn("POST", "(" + strconv.Itoa(requestCount) + ") KEY COUNT: " + strconv.Itoa(count), nil)
                for kz, vz := range v {
                    utils.LogWarn("POST", kz + ": " + vz, nil)
                }
                xmlKeys := strings.Split( v["xml_tags"], "," )
                for _, k := range xmlKeys {
                    processedJson = utils.RemoveXmlTagFromJson(
                        k, processedJson)
                }
                break
            }
        }
        utils.LogWarn("POST", "Continuing...", nil)
        utils.LogWarn("POST", string(processedJson), nil)

        // TODO: Should I just rebuild the map and pass to the generic
        //    utils.DefaultJsonPostProcess function?
        structureVar, errorVar := utils.ParsePostProcessedJson( request,
            jsonKeys, processedJson, parsedStructure, parsedErrorStructure )
        parsedStructure = structureVar
        parsedErrorStructure = errorVar
        //for kz, vz := range parsedStructure {
        //    for kkz, _ := range vz.(map[string]interface{}) {
        //        utils.LogWarn("POST PS", "(" + kz + ") " + kkz, nil)
        //    }
        //}

    }

    returnJson := utils.CollapseJson( parsedStructure, parsedErrorStructure )
    return returnJson
}
