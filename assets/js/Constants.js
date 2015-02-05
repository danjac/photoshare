var keyMirror = require('keymirror');

module.exports = {

    Events: keyMirror({
        CHANGE_EVENT: null
    }),

    ActionSources: keyMirror({
        VIEW_ACTION: null,
        SERVER_ACTION: null
    }),

    Alerts: keyMirror({
        ALERT_SUCCESS: null,
        ALERT_INFO: null,
        ALERT_WARNING: null,
        ALERT_DANGER: null
    }),

    Actions: keyMirror({
        GET_TAGS: null,
        GET_PHOTOS: null,
        GET_PHOTO_DETAIL: null,
        NEW_PHOTO: null,
        NEW_PHOTO_PREVIEW: null,
        UPLOAD_STARTED: null,
        PHOTO_DELETED: null,
        PHOTO_EDIT_MODE: null,
        PHOTO_EDIT_DONE: null,
        LOGIN_SUCCESSFUL: null,
        LOGOUT: null,
        NEW_ALERT_MESSAGE: null
    })
    
};
