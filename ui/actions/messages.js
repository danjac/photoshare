import { ActionTypes } from '../constants';

const {
  NEW_MESSAGE,
  DELETE_MESSAGE
} = ActionTypes;


export function newMessage(msg, level) {
  return {
    type: NEW_MESSAGE,
    message: {
      level: level,
      msg: msg
    }
  }
}

export function deleteMessage(key) {
  return {
    type: DELETE_MESSAGE,
    key: key
  }
}
