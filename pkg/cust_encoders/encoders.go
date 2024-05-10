package cust_encoders

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"forum/internal/exceptions"
	"net/url"
)

func EncodeParams(data error) string {
	var code int
	switch data.(type) {
	case exceptions.AuthenticationError:
		code = data.(exceptions.AuthenticationError).StatusCode
	case exceptions.ForbiddenError:
		code = data.(exceptions.ForbiddenError).StatusCode
	case exceptions.ResourceNotFoundError:
		code = data.(exceptions.ResourceNotFoundError).StatusCode
	case exceptions.StatusMethodNotAllowed:
		code = data.(exceptions.StatusMethodNotAllowed).StatusCode
	case exceptions.StatusConflictError:
		code = data.(exceptions.StatusConflictError).StatusCode
	case exceptions.ValidationError:
		code = data.(exceptions.ValidationError).StatusCode
	case exceptions.InternalServerError:
		code = data.(exceptions.InternalServerError).StatusCode
	case exceptions.BadRequestError:
		code = data.(exceptions.BadRequestError).StatusCode
	}
	jsonData, _ := json.Marshal(data)
	base64EncodedData := base64.StdEncoding.EncodeToString(jsonData)
	base64Code := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v", code)))
	params := url.Values{}
	params.Set("params", base64EncodedData)
	params.Set("code", base64Code)
	encodedParams := params.Encode()
	return encodedParams
}

func DecodeParams(params string) (error, error) {
	decodedParams, err := url.QueryUnescape(params)
	if err != nil {
		return nil, err
	}

	queryParams, err := url.ParseQuery(decodedParams)
	if err != nil {
		return nil, err
	}

	base64EncodedData := queryParams.Get("params")
	codeData := queryParams.Get("code")

	jsonData, err := base64.StdEncoding.DecodeString(base64EncodedData)
	if err != nil {
		return nil, err
	}
	code, err := base64.StdEncoding.DecodeString(codeData)
	if err != nil {
		return nil, err
	}
	fmt.Println("here")
	switch string(code) {
	case "422":
		var errData exceptions.ValidationError
		err = json.Unmarshal(jsonData, &errData)
		if err != nil {
			fmt.Println("Here")
			return nil, err
		}
		return errData, nil
	case "401":
		var errData exceptions.AuthenticationError
		err = json.Unmarshal(jsonData, &errData)
		if err != nil {
			fmt.Println("Here")
			return nil, err
		}

		return errData, nil
	case "404":
		var errData exceptions.ResourceNotFoundError
		err = json.Unmarshal(jsonData, &errData)
		if err != nil {
			fmt.Println("Here")
			return nil, err
		}

		return errData, nil
	case "405":
		var errData exceptions.StatusMethodNotAllowed
		err = json.Unmarshal(jsonData, &errData)
		if err != nil {
			fmt.Println("Here")
			return nil, err
		}

		return errData, nil
	case "403":
		fmt.Println("here2")
		var errData exceptions.ForbiddenError
		err = json.Unmarshal(jsonData, &errData)
		if err != nil {
			fmt.Println("Here")
			return nil, err
		}

		return errData, nil
	case "500":
		var errData exceptions.InternalServerError
		err = json.Unmarshal(jsonData, &errData)
		if err != nil {
			fmt.Println("Here")
			return nil, err
		}

		return errData, nil
	default:
		var errData exceptions.ResourceNotFoundError
		err = json.Unmarshal(jsonData, &errData)
		if err != nil {
			fmt.Println("Here")
			return nil, err
		}

		return errData, nil
	}
}
