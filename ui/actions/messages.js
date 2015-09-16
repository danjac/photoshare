import ActionTypes from '../actionTypes/messages';

const {
  DELETE_MESSAGE
} = ActionTypes;


export function deleteMessage(key) {
  return {
    type: DELETE_MESSAGE,
    key: key
  }
}
