package elastico

import (
	"bytes"
	json "elastico/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
)

type Elastico struct {
	Client  *http.Client
	BaseURL *url.URL
}

func (c *Elastico) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	var buf io.Reader
	if body == nil {
	} else if v, ok := body.(io.Reader); ok {
		buf = v
	} else if v, ok := body.(json.M); ok {
		buf = new(bytes.Buffer)
		if err := json.NewEncoder(buf.(io.ReadWriter)).Encode(v); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("not supported type: %s", reflect.TypeOf(body))
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "text/json; charset=UTF-8")
	req.Header.Add("Accept", "text/json")
	return req, nil
}

func New(u string) (*Elastico, error) {
	baseURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	return &Elastico{
		Client:  http.DefaultClient,
		BaseURL: baseURL,
	}, nil
}

type elasticsearchResponse struct {
}

/*
{"acknowledged":true}map[string]interface {}{
    "acknowledged": bool(true),
}⏎
*/
/*
   "error": json.M{
       "root_cause": []interface {}{
           json.M{
               "type":   "illegal_argument_exception",
               "reason": "Malformed action/metadata line [1], expected START_OBJECT or END_OBJECT but found [null]",
           },
       },
       "type":   "illegal_argument_exception",
       "reason": "Malformed action/metadata line [1], expected START_OBJECT or END_OBJECT but found [null]",
   },


   {
  "acknowledged": true
}⏎

*/

type Error struct {
	State  string `json:"state"`
	Status int64  `json:"status"`
}

func (e Error) Error() string {
	return ""
}

func (wd *Elastico) do(req *http.Request, v interface{}) (*elasticsearchResponse, error) {
	resp, err := wd.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var r io.Reader = resp.Body

	// r = io.TeeReader(r, os.Stdout)
	err = json.NewDecoder(r).Decode(&v)
	if err != nil {
		return nil, err
	}

	/*
		if dr.Status != 0 {
			return &dr, &Error{
				Status: dr.Status,
				State:  dr.State,
			}
		}*/

	return nil, nil
}

func (wd *Elastico) Do(req *http.Request, v interface{}) error {
	_, err := wd.do(req, v)
	if err != nil {
		return err
	}

	/*
		if dr.Value == nil {
			return nil
		}

		switch v := v.(type) {
		case io.Writer:
			value := ""
			if err = json.Unmarshal(*dr.Value, &value); err != nil {
				return err
			}

			v.Write([]byte(value))
		case interface{}:
			return json.Unmarshal(*dr.Value, &v)
		}
	*/

	return nil
}
