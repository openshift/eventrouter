/*
 * REST API
 *
 * Rockset's REST API allows for creating and managing all resources in Rockset. Each supported endpoint is documented below.  All requests must be authorized with a Rockset API key, which can be created in the [Rockset console](https://console.rockset.com). The API key must be provided as `ApiKey <api_key>` in the `Authorization` request header. For example: ``` Authorization: ApiKey aB35kDjg93J5nsf4GjwMeErAVd832F7ad4vhsW1S02kfZiab42sTsfW5Sxt25asT ```  All endpoints are only accessible via https.  Build something awesome!
 *
 * API version: v1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package rockset

import (
    "bufio"
    "bytes"
    "encoding/json"
	"io/ioutil"
    "log"
	"net/http"
	"net/url"
	"strings"
	"fmt"
)

type UsersApiService Service

/* 
UsersApiService Create User
Create a new user for an organization.
 * @param body JSON object

@return CreateUserResponse
*/
func (a *UsersApiService) Create(body CreateUserRequest) (CreateUserResponse, *http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Post")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		localVarReturnValue CreateUserResponse
	)

	// create path and map variables
	localVarPath := a.Client.cfg.BasePath + "/v1/orgs/self/users"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Api Key header
    localVarHttpHeaderAuthorization := ""
	localVarHttpHeaderApiKey := a.Client.selectHeaderAuthorization(localVarHttpHeaderAuthorization)
    if localVarHttpHeaderApiKey == "" {
        log.Fatal("missing required argument ApiKey")
    }
	localVarHeaderParams["authorization"] = "ApiKey " + localVarHttpHeaderApiKey

	// body params
	localVarPostBody = &body
	r, err := a.Client.prepareRequest(localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHttpResponse, err := a.Client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	if localVarHttpResponse.StatusCode < 300 {
		// If we succeed, return the data, otherwise pass on to decode error.
		err = a.Client.decode(&localVarReturnValue, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
		if err == nil { 
			return localVarReturnValue, localVarHttpResponse, err
		}
	}

	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}
		
		if localVarHttpResponse.StatusCode == 200 {
			var v CreateUserResponse
			err = a.Client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		
		return localVarReturnValue, localVarHttpResponse, newErr
	}

	return localVarReturnValue, localVarHttpResponse, nil
}

func (a *UsersApiService) CreateStream(body CreateUserRequest) (string, *http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Post")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
        localVarBody []byte
	)

	// create path and map variables
	localVarPath := a.Client.cfg.BasePath + "/v1/orgs/self/users"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Api Key header
    localVarHttpHeaderAuthorization := ""
	localVarHttpHeaderApiKey := a.Client.selectHeaderAuthorization(localVarHttpHeaderAuthorization)
    if localVarHttpHeaderApiKey == "" {
        log.Fatal("missing required argument ApiKey")
    }
	localVarHeaderParams["authorization"] = "ApiKey " + localVarHttpHeaderApiKey

	// body params
	localVarPostBody = &body
	r, err := a.Client.prepareRequest(localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return "", nil, err
	}

	localVarHttpResponse, err := a.Client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return "", localVarHttpResponse, err
	}

    reader := bufio.NewReader(localVarHttpResponse.Body)
    localVarBody, err= reader.ReadBytes('\n')

    var out bytes.Buffer
    err = json.Indent(&out, []byte(string(localVarBody)), "", "    ")
    if err != nil {
        return "", localVarHttpResponse, err
    }

    if localVarHttpResponse.StatusCode >= 300 {
		return out.String(), localVarHttpResponse, nil
	}

	return out.String(), localVarHttpResponse, nil
}

/* 
UsersApiService Delete User
Delete a user from an organization.
 * @param user user email

@return DeleteUserResponse
*/
func (a *UsersApiService) Delete(user string) (DeleteUserResponse, *http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Delete")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		localVarReturnValue DeleteUserResponse
	)

	// create path and map variables
	localVarPath := a.Client.cfg.BasePath + "/v1/orgs/self/users/{user}"
	localVarPath = strings.Replace(localVarPath, "{"+"user"+"}", fmt.Sprintf("%v", user), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Api Key header
    localVarHttpHeaderAuthorization := ""
	localVarHttpHeaderApiKey := a.Client.selectHeaderAuthorization(localVarHttpHeaderAuthorization)
    if localVarHttpHeaderApiKey == "" {
        log.Fatal("missing required argument ApiKey")
    }
	localVarHeaderParams["authorization"] = "ApiKey " + localVarHttpHeaderApiKey

	r, err := a.Client.prepareRequest(localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHttpResponse, err := a.Client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	if localVarHttpResponse.StatusCode < 300 {
		// If we succeed, return the data, otherwise pass on to decode error.
		err = a.Client.decode(&localVarReturnValue, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
		if err == nil { 
			return localVarReturnValue, localVarHttpResponse, err
		}
	}

	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}
		
		if localVarHttpResponse.StatusCode == 200 {
			var v DeleteUserResponse
			err = a.Client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		
		return localVarReturnValue, localVarHttpResponse, newErr
	}

	return localVarReturnValue, localVarHttpResponse, nil
}

func (a *UsersApiService) DeleteStream(user string) (string, *http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Delete")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
        localVarBody []byte
	)

	// create path and map variables
	localVarPath := a.Client.cfg.BasePath + "/v1/orgs/self/users/{user}"
	localVarPath = strings.Replace(localVarPath, "{"+"user"+"}", fmt.Sprintf("%v", user), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Api Key header
    localVarHttpHeaderAuthorization := ""
	localVarHttpHeaderApiKey := a.Client.selectHeaderAuthorization(localVarHttpHeaderAuthorization)
    if localVarHttpHeaderApiKey == "" {
        log.Fatal("missing required argument ApiKey")
    }
	localVarHeaderParams["authorization"] = "ApiKey " + localVarHttpHeaderApiKey

	r, err := a.Client.prepareRequest(localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return "", nil, err
	}

	localVarHttpResponse, err := a.Client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return "", localVarHttpResponse, err
	}

    reader := bufio.NewReader(localVarHttpResponse.Body)
    localVarBody, err= reader.ReadBytes('\n')

    var out bytes.Buffer
    err = json.Indent(&out, []byte(string(localVarBody)), "", "    ")
    if err != nil {
        return "", localVarHttpResponse, err
    }

    if localVarHttpResponse.StatusCode >= 300 {
		return out.String(), localVarHttpResponse, nil
	}

	return out.String(), localVarHttpResponse, nil
}

/* 
UsersApiService Get Current User
Retrieve currently active user.

@return User
*/
func (a *UsersApiService) Get() (User, *http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		localVarReturnValue User
	)

	// create path and map variables
	localVarPath := a.Client.cfg.BasePath + "/v1/orgs/self/users/self"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Api Key header
    localVarHttpHeaderAuthorization := ""
	localVarHttpHeaderApiKey := a.Client.selectHeaderAuthorization(localVarHttpHeaderAuthorization)
    if localVarHttpHeaderApiKey == "" {
        log.Fatal("missing required argument ApiKey")
    }
	localVarHeaderParams["authorization"] = "ApiKey " + localVarHttpHeaderApiKey

	r, err := a.Client.prepareRequest(localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHttpResponse, err := a.Client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	if localVarHttpResponse.StatusCode < 300 {
		// If we succeed, return the data, otherwise pass on to decode error.
		err = a.Client.decode(&localVarReturnValue, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
		if err == nil { 
			return localVarReturnValue, localVarHttpResponse, err
		}
	}

	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}
		
		if localVarHttpResponse.StatusCode == 200 {
			var v User
			err = a.Client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		
		return localVarReturnValue, localVarHttpResponse, newErr
	}

	return localVarReturnValue, localVarHttpResponse, nil
}

func (a *UsersApiService) GetStream() (string, *http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
        localVarBody []byte
	)

	// create path and map variables
	localVarPath := a.Client.cfg.BasePath + "/v1/orgs/self/users/self"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Api Key header
    localVarHttpHeaderAuthorization := ""
	localVarHttpHeaderApiKey := a.Client.selectHeaderAuthorization(localVarHttpHeaderAuthorization)
    if localVarHttpHeaderApiKey == "" {
        log.Fatal("missing required argument ApiKey")
    }
	localVarHeaderParams["authorization"] = "ApiKey " + localVarHttpHeaderApiKey

	r, err := a.Client.prepareRequest(localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return "", nil, err
	}

	localVarHttpResponse, err := a.Client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return "", localVarHttpResponse, err
	}

    reader := bufio.NewReader(localVarHttpResponse.Body)
    localVarBody, err= reader.ReadBytes('\n')

    var out bytes.Buffer
    err = json.Indent(&out, []byte(string(localVarBody)), "", "    ")
    if err != nil {
        return "", localVarHttpResponse, err
    }

    if localVarHttpResponse.StatusCode >= 300 {
		return out.String(), localVarHttpResponse, nil
	}

	return out.String(), localVarHttpResponse, nil
}

/* 
UsersApiService List Users
Retrieve all users for an organization.

@return ListUsersResponse
*/
func (a *UsersApiService) List() (ListUsersResponse, *http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		localVarReturnValue ListUsersResponse
	)

	// create path and map variables
	localVarPath := a.Client.cfg.BasePath + "/v1/orgs/self/users"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Api Key header
    localVarHttpHeaderAuthorization := ""
	localVarHttpHeaderApiKey := a.Client.selectHeaderAuthorization(localVarHttpHeaderAuthorization)
    if localVarHttpHeaderApiKey == "" {
        log.Fatal("missing required argument ApiKey")
    }
	localVarHeaderParams["authorization"] = "ApiKey " + localVarHttpHeaderApiKey

	r, err := a.Client.prepareRequest(localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHttpResponse, err := a.Client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	if localVarHttpResponse.StatusCode < 300 {
		// If we succeed, return the data, otherwise pass on to decode error.
		err = a.Client.decode(&localVarReturnValue, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
		if err == nil { 
			return localVarReturnValue, localVarHttpResponse, err
		}
	}

	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}
		
		if localVarHttpResponse.StatusCode == 200 {
			var v ListUsersResponse
			err = a.Client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		
		return localVarReturnValue, localVarHttpResponse, newErr
	}

	return localVarReturnValue, localVarHttpResponse, nil
}

func (a *UsersApiService) ListStream() (string, *http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
        localVarBody []byte
	)

	// create path and map variables
	localVarPath := a.Client.cfg.BasePath + "/v1/orgs/self/users"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Api Key header
    localVarHttpHeaderAuthorization := ""
	localVarHttpHeaderApiKey := a.Client.selectHeaderAuthorization(localVarHttpHeaderAuthorization)
    if localVarHttpHeaderApiKey == "" {
        log.Fatal("missing required argument ApiKey")
    }
	localVarHeaderParams["authorization"] = "ApiKey " + localVarHttpHeaderApiKey

	r, err := a.Client.prepareRequest(localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return "", nil, err
	}

	localVarHttpResponse, err := a.Client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return "", localVarHttpResponse, err
	}

    reader := bufio.NewReader(localVarHttpResponse.Body)
    localVarBody, err= reader.ReadBytes('\n')

    var out bytes.Buffer
    err = json.Indent(&out, []byte(string(localVarBody)), "", "    ")
    if err != nil {
        return "", localVarHttpResponse, err
    }

    if localVarHttpResponse.StatusCode >= 300 {
		return out.String(), localVarHttpResponse, nil
	}

	return out.String(), localVarHttpResponse, nil
}

