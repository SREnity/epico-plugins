package main

import (
	"bytes"
	"strings"
	"time"

	xj "github.com/basgys/goxml2json"

	"github.com/SREnity/epico/signers/aws_v4"
	generic_structs "github.com/SREnity/epico/structs"
	utils "github.com/SREnity/epico/utils"
)

var PluginAuthFunction = PluginAuth
var PluginResponseToJsonFunction = PluginResponseToJson
var PluginPostProcessFunction = PluginPostProcess
var PluginPagingPeekFunction = utils.DefaultXmlPagingPeek

// authParams expects the AWS ACCESS ID in slice slot [0] and the AWS SECRET KEY
//     in slice slot [1].
func PluginAuth(apiRequest generic_structs.ApiRequest, authParams []string) generic_structs.ApiRequest {

	var err error
	signer := v4.NewSigner(generic_structs.ApiCredentials{
		Id:  authParams[0],
		Key: authParams[1],
	}, func(v4Signer *v4.Signer) {
		v4Signer.DisableHeaderHoisting = false
		v4Signer.DisableRequestBodyOverwrite = true

	})

	if apiRequest.Settings.Vars["region"] != "{{region}}" {
		apiRequest.FullRequest.Header, err = signer.Presign(apiRequest.FullRequest, nil, apiRequest.Settings.Vars["service"], apiRequest.Settings.Vars["region"], 60*time.Minute, time.Now())
	} else {
		apiRequest.FullRequest.Header, err = signer.Presign(apiRequest.FullRequest, nil, apiRequest.Settings.Vars["service"], "us-east-1", 60*time.Minute, time.Now())
	}

	if err != nil {
		utils.LogError("AWS:PluginAuth", "Error presigning the AWS request", err)
	}

	return apiRequest
}

func PluginResponseToJson(vars map[string]string, response []byte) []byte {
	jsonBody, err := xj.Convert(bytes.NewReader(response))
	if err != nil {
		utils.LogError("AWS:PluginResponseToJson", "Error converting XML API response", err)
		return nil
	}

	processedJson := jsonBody.Bytes()
	xmlKeys := strings.Split(vars["xml_tags"], ",")
	for _, k := range xmlKeys {
		processedJson = utils.RemoveXmlTagFromJson(k, processedJson)
	}

	return processedJson
}

func PluginPostProcess(apiResponseMap map[generic_structs.ComparableApiRequest][]byte, jsonKeys []map[string]string, postParams []string) []byte {

	parsedStructure := make(map[string]interface{})
	parsedErrorStructure := make(map[string]interface{})

	for request, apiResponse := range apiResponseMap {

		processedJson := []byte(nil)
		for _, v := range jsonKeys {
			if v["api_call_uuid"] == request.Uuid {
				processedJson = PluginResponseToJson(v, apiResponse)
			}
		}

		// TODO: Should I just rebuild the map and pass to the generic
		//    utils.DefaultJsonPostProcess function?
		structureVar, errorVar := utils.ParsePostProcessedJson(request,
			jsonKeys, processedJson, parsedStructure, parsedErrorStructure)
		parsedStructure = structureVar
		parsedErrorStructure = errorVar

	}

	returnJson := utils.CollapseJson(parsedStructure, parsedErrorStructure)

	return returnJson
}
