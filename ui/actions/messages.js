import ActionTypes from '../actionTypes';

const {
  DELETE_MESSAGE
} = ActionTypes;


export function deleteMessage(key) {
  return {
    type: DELETE_MESSAGE,
    key: key
  }
}
