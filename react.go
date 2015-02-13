package photoshare

import (
	"encoding/json"
	"github.com/juju/errgo"
	"io/ioutil"
	"net/http"
	"net/url"
)

var reactClient *http.Client

func popular(ctx *context, w http.ResponseWriter, r *http.Request) error {

	photos, err := ctx.datamapper.getPhotos(newPage(1), "votes")

	if err != nil {
		return errgo.Mask(err)
	}

	return renderToReact(w, r, "", photos)
}

func latest(ctx *context, w http.ResponseWriter, r *http.Request) error {

	photos, err := ctx.datamapper.getPhotos(newPage(1), "")

	if err != nil {
		return errgo.Mask(err)
	}

	return renderToReact(w, r, "latest", photos)
}

func renderToReact(w http.ResponseWriter, r *http.Request, route string, props interface{}) error {

	if reactClient == nil {
		t := &http.Transport{}
		reactClient = &http.Client{Transport: t}
	}

	propsJSON, err := json.Marshal(props)
	if err != nil {
		return errgo.Mask(err)
	}

	values := url.Values{}
	values.Set("route", route)
	values.Set("props", string(propsJSON))

	resp, err := reactClient.PostForm("http://127.0.0.1:3000/react/", values)
	if err != nil {
		return errgo.Mask(err)
	}

	content, err := ioutil.ReadAll(resp.Body)
	return writeBody(w, content, http.StatusOK, "text/html")

}
