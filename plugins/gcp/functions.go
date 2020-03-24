package main

import (
	generic_structs "github.com/SREnity/epico/structs"
	utils "github.com/SREnity/epico/utils"
)

var PluginAuthFunction = PluginAuth
var PluginResponseToJsonFunction = PluginResponseToJson
var PluginPostProcessFunction = PluginPostProcess
var PluginPagingPeekFunction = PluginPagingPeek

// authParams expects the Auth function in slice slot [0] followed by a []string
//    of auth parameters
func PluginAuth(apiRequest generic_structs.ApiRequest, authParams []string) generic_structs.ApiRequest {

	if len(authParams) < 1 {
		authParams = []string{""}
	}

	// In the JSON plugin, these should all be JSON requests.
	apiRequest.FullRequest.Header.Set("Content-Type", "application/json")
	_, ok := apiRequest.FullRequest.Header["Accept"]
	_, ok1 := apiRequest.FullRequest.Header["accept"]
	if !ok && !ok1 {
		apiRequest.FullRequest.Header.Set("Accept", "application/json")
	}

	// TODO: I would like to be able to dynamically call whatever function name
	//    is passed in authParams[0], but this is difficult in Go - particularly
	//    if they are going outside the standard utils library. At this point,
	//    we'll just manage the duplication here between core utils and these
	//    generic files.  If folks want to use their own functions, they should
	//    create their own plugins.
	switch authParams[0] {
	case "BasicAuth":
		return utils.BasicAuth(apiRequest, authParams[1:])
	case "CustomHeaderAndBasicAuth":
		return utils.CustomHeaderAndBasicAuth(apiRequest, authParams[1:])
	case "CustomQuerystringAuth":
		return utils.CustomQuerystringAuth(apiRequest, authParams[1:])
	case "CustomHeaderAuth":
		return utils.CustomHeaderAuth(apiRequest, authParams[1:])
	case "JwtAuth":
		// Expects authParams[1..5] to be:
		// email, private key, private key id, scopes (csv), token url
		return utils.JwtAuth(apiRequest, authParams[1:])
	case "Oauth2TwoLegAuth":
		return utils.Oauth2TwoLegAuth(apiRequest, authParams[1:])
	case "SessionTokenAuth":
		return utils.SessionTokenAuth(apiRequest, authParams[1:])
	default:
		return apiRequest
	}

}

func PluginResponseToJson(vars map[string]string, response []byte) []byte {
	return response
}

func PluginPagingPeek(response []byte, responseKeys []string, oldPageValue interface{}, peekParams []string) (interface{}, bool) {

	if len(peekParams) < 1 {
		peekParams = []string{""}
	}

	switch peekParams[0] {
	case "CalculatePagingPeek":
		return utils.CalculatePagingPeek(response, responseKeys,
			oldPageValue, peekParams[1:])
	case "RegexJsonPagingPeek":
		return utils.RegexJsonPagingPeek(response, responseKeys,
			oldPageValue, peekParams[1:])
	default:
		return utils.DefaultJsonPagingPeek(response, responseKeys,
			oldPageValue, peekParams[1:])
	}

}

func PluginPostProcess(apiResponseMap map[generic_structs.ComparableApiRequest][]byte, jsonKeys []map[string]string, postParams []string) []byte {

	if len(postParams) < 1 {
		postParams = []string{""}
	}

	switch postParams[0] {
	default:
		// No jsonKeys required because we aren't expecting any special vars
		//    or configurations.
		return utils.DefaultJsonPostProcess(apiResponseMap, jsonKeys)
	}

}
