import * as api from '../api';
import ActionTypes from '../actionTypes/changePassword';

const {
  CHANGE_PASSWORD_RESET,
  CHANGE_PASSWORD_FORM_INVALID,
  CHANGE_PASSWORD_PENDING,
  CHANGE_PASSWORD_SUCCESS,
  CHANGE_PASSWORD_FAILURE
} = ActionTypes;

function validate(password, passwordConfirm) {
  const errors = new Map();

  if (!password) {
    errors.set("password", "Password is required");
  }

  if (!passwordConfirm) {
    errors.set("passwordConfirm", "Please confirm your new password");
  }

  if (password && passwordConfirm && password !== passwordConfirm) {
    errors.set("passwordConfirm", "The passwords do not match");
  }

  return errors;

}


export function resetForm() {
  return { type: CHANGE_PASSWORD_RESET };
}

export function submitForm(password, passwordConfirm, code) {

  const errors = validate(password, passwordConfirm);

  if (errors.size > 0) {
    return {
      type: CHANGE_PASSWORD_FORM_INVALID,
      errors: errors
    }
  }

  return {
    types: [
      CHANGE_PASSWORD_PENDING,
      CHANGE_PASSWORD_SUCCESS,
      CHANGE_PASSWORD_FAILURE
    ],
    payload: {
      promise: api.changePassword(password, code)
    }
  }
}


