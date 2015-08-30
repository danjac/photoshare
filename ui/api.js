import fetch from 'isomorphic-fetch';
//import { Schema, arrayOf, normalize } from 'normalizr';

const API_URI = '/api';

function makeURI(uri) {
  return `${API_URI}${uri}`;
}

export function getPhotos(page, orderBy) {
  return fetch(`${makeURI('/photos/')}?page=${page}&orderBy=${orderBy}`)
  .then(response => response.json());
}

export function getPhotoDetail(id) {
  const url  = `${makeURI('/photos/' + id)}`;
  return fetch(url).then(response => response.json());
}
