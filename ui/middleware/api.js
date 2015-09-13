export const CALL_API = Symbol('CALL API');

function callAPI(endpoint, method, data) {

  const args = { method: method };
  const token = getToken();
  let headers = {};

  if (token) {
    headers[AUTH_TOKEN] = token;
  }

  if (data) {
      // check if window.FormData
      if (data instanceof window.FormData) {
        args.body= data;
      } else {
        args.body  = JSON.stringify(data);
        headers = Object.assign({}, headers, {
          "Accept": "application/json",
          "Content-Type": "application/json"
        });
      }
  }

  if (headers) {
    args.headers = headers;
  }

  return fetch(API_URI + endpoint, args)
    .then(response => {
      if (response.headers.has(AUTH_TOKEN)) {
        const token = response.headers.get(AUTH_TOKEN);
        if (token) {
          setToken(token);
        }
      }
      return response.json();
    });

}


export default store => next => action => {
  const callApi = action[CALL_API];

  if (typeof callAPI === 'undefined') {
    return next(action);
  }

  // example:
  //
  // {
  //  [CALL_API]: {
  //  endpoint: '/auth/',
  //    types: [LOGIN_REQUEST, LOGIN_SUCCESS, LOGIN_FAILURE],
  //  method: 'POST',
   //  body: {
    //    identifier,
    //    password
    //  }
  //  }
  // }
  let { types, args } = callAPI;

  const [ requestType, successType, failureType ] = types;

  function actionWith(data) {
    const finalAction = Object.assign({}, action, data);
    delete finalAction[CALL_API];
    return finalAction;
  }

  next(actionWith({ type: requestType }));

  api(requestType, args)
    .then(
        response => next(actionWith({ response, type: successType })),
        error => next(actionWith({ error.message, type: failureType }))
    );

}
